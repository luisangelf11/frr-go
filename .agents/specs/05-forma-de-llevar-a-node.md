# Especificación 05: Forma de Llevar a Node.js

## 1. Objetivo y Análisis de Enfoques
El propósito es permitir que desarrolladores de Node.js (TypeScript y JavaScript) utilicen el motor de Go con una experiencia nativa (**Fluent Builder**), sin necesidad de instalar Go en sus servidores de producción, y con un rendimiento de **1 a 3 ms por reporte**.

### 1.1 Comparativa de Enfoques Técnicos

| Enfoque | Compilador C (CGO) | Seguridad Node | Rendimiento | Experiencia en Producción |
| :--- | :--- | :--- | :--- | :--- |
| **FFI / C-Shared (.dll / .so)** | ⚠️ Obligatorio (MinGW / GCC) | 🔴 Riesgo de *segmentation fault* si Go hace panic | < 1 ms | Frágil en Windows/Alpine |
| **WASM (WebAssembly)** | ❌ No | 🟢 Seguro | 5-10 ms (overhead V8 WASM) | Limitaciones de fuentes e I/O |
| **Persistent Worker Stdio IPC (Elegido)** | ❌ No requiere CGO | 🟢 Aislado (Node nunca se cae) | **1 a 3 ms** | **El estándar moderno (esbuild, Prisma, Biome)** |

---

## 2. Protocolo de Comunicación (Stdio Binary Framing)

Para soportar buffers binarios arbitrarios (PDFs o archivos Excel pueden contener bytes `0x0A`, `0x0D` o caracteres nulos que romperían protocolos basados en texto puro), definimos un protocolo de enmarcado por longitud (**Length-Prefixed Binary Protocol**):

```
[ Petición de Node.js a Go (Stdin) ]
┌───────────────────────────┬──────────────────────────────────────────────┐
│ Longitud JSON (4 bytes)   │ Payload JSON (UTF-8)                         │
│ uint32 Big-Endian         │ {"id":"1","format":"pdf","spec":{...}}       │
└───────────────────────────┴──────────────────────────────────────────────┘

[ Respuesta de Go a Node.js (Stdout) ]
┌───────────────────────────┬──────────────────┬───────────────────────────┐
│ Longitud Datos (4 bytes)  │ Estado (1 byte)  │ Payload de Salida         │
│ uint32 Big-Endian         │ 0x00=OK, 0x01=ERR│ Bytes binarios (PDF/XLSX) │
└───────────────────────────┴──────────────────┴───────────────────────────┘
```

### Ventajas del protocolo:
1. **Zero-copy**: Node.js lee el encabezado de 4 bytes, sabe exactamente cuántos bytes esperar, y reserva el `Buffer` de forma inmediata.
2. **Sin serialización Base64**: Ahorra un 33% de ancho de banda y evita la recolección de basura de strings pesados en V8.

---

## 3. Arquitectura del SDK en TypeScript (`nodejs-wrapper`)

```
nodejs-wrapper/
├── package.json
├── tsconfig.json
├── src/
│   ├── index.ts                # Exporta ReportBuilder, ChartType, Theme, etc.
│   ├── builder.ts              # Fluent Builder idéntico al de Go
│   ├── types.ts                # Interfaces TypeScript del Spec
│   ├── worker.ts               # Administrador del proceso Go en segundo plano
│   └── binary-resolver.ts      # Localizador del binario compilado por plataforma
└── test/
    └── worker.test.ts          # Pruebas de integración
```

### 3.1 Ciclo de Vida del Worker en Node.js
1. **Lazy Initialization**: El proceso Go no se inicia hasta que se genera el primer reporte.
2. **Persistencia**: El proceso se mantiene vivo en segundo plano escuchando en `stdin`.
3. **Auto-reinicio (*Self-healing*)**: Si por alguna razón externa el proceso termina, el cliente de Node lo detecta y lo levanta transparentemente en la siguiente solicitud.
4. **Cierre limpio**: Escucha de señales `SIGINT`, `SIGTERM` y evento `beforeExit` de Node para terminar el worker sin dejar procesos huérfanos.

---

## 4. API de Uso en Node.js

```typescript
import { ReportBuilder, ChartType, HeadingLevel } from 'reporteria-node';

const report = new ReportBuilder("Reporte Financiero")
  .setAuthor("Finanzas Corp")
  .addHeading("Resumen de Ventas", HeadingLevel.H1)
  .addParagraph("Consolidado mensual de operaciones...")
  .addChart({
    type: ChartType.Bar,
    title: "Ventas Q1",
    categories: ["Ene", "Feb", "Mar"],
    series: [{ name: "2026", values: [15000, 18500, 22000] }]
  })
  .addTable({
    headers: ["Mes", "Meta", "Alcanzado"],
    rows: [
      ["Enero", "$12,000", "$15,000"],
      ["Febrero", "$16,000", "$18,500"],
      ["Marzo", "$20,000", "$22,000"],
    ]
  });

// Generar a Buffer en memoria (para responder en Express/Fastify)
const pdfBuffer = await report.exportToBuffer("pdf");

// O guardar directamente a disco
await report.exportToFile("excel", "./reporte.xlsx");
```

---

## 5. Criterios de Aceptación
1. **Latencia Round-Trip**: Generación de un reporte PDF completo desde Node.js en **menos de 3 milisegundos**.
2. **Estabilidad**: Capacidad de procesar 500 solicitudes consecutivas sin fugas de memoria en Node ni bloqueos en el worker.
3. **Ergonomía TypeScript**: Tipado estricto con autocompletado para todas las funciones del builder.
