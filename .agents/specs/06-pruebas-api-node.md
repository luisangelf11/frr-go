# Especificación 06: Pruebas en Proyecto API de Node.js

## 1. Objetivo
Validar el rendimiento, la estabilidad y la experiencia de desarrollo en un entorno real creando un servidor API en **Node.js** que exponga endpoints para generar y descargar reportes al vuelo en **PDF**, **Excel** y **CSV**.

---

## 2. Estructura del Proyecto de Pruebas (`examples/node-api`)

```
examples/node-api/
├── package.json
├── server.ts                   # Servidor HTTP (Express / Node Nativo)
├── test-client.ts              # Script de prueba funcional y descarga
└── benchmark.ts                # Test de estrés y medición de latencia
```

---

## 3. Endpoints a Implementar

### 3.1 `GET /api/reports/sales/pdf`
- Genera dinámicamente un reporte de ventas con encabezados, gráfico de barras y tabla de detalle.
- Responde con cabeceras HTTP:
  - `Content-Type: application/pdf`
  - `Content-Disposition: inline; filename="ventas.pdf"`
- Permite abrirlo directamente en el visor PDF del navegador.

### 3.2 `GET /api/reports/sales/excel`
- Genera la hoja de cálculo con estilos de cabecera, autofit y datos numéricos formateados.
- Responde con cabeceras HTTP:
  - `Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`
  - `Content-Disposition: attachment; filename="ventas.xlsx"`

### 3.3 `GET /api/reports/sales/csv`
- Extrae la tabla de datos en formato CSV con cabecera UTF-8 BOM.
- Responde con:
  - `Content-Type: text/csv; charset=utf-8`
  - `Content-Disposition: attachment; filename="ventas.csv"`

### 3.4 `POST /api/reports/custom`
- Recibe un cuerpo JSON con datos dinámicos enviados por el cliente y el formato deseado (`format: "pdf" | "excel" | "csv"`).
- Devuelve el binario correspondiente en tiempo real.

---

## 4. Pruebas de Carga y Benchmarking
Se ejecutará un script de medición concurrente para verificar:
1. **Latencia bajo carga**:
   - Evaluar tiempo de respuesta para 100 solicitudes concurrentes.
   - Meta: Latencia p95 < **15 ms** para PDF y < **5 ms** para Excel/CSV a través de HTTP.
2. **Consumo de Memoria**:
   - Monitorear el RSS (*Resident Set Size*) del proceso Node.js y del proceso Go Worker tras 1,000 generaciones consecutivas.
   - Meta: Retorno a niveles base de memoria sin fugas (*memory leaks*).
3. **Resiliencia**:
   - Simular la terminación abrupta del proceso Go (`kill -9`) y comprobar que la API de Node.js se recupera automáticamente en la siguiente llamada sin necesidad de reiniciar el servidor Node.
