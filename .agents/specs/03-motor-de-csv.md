# Especificación 03: Motor de CSV

## 1. Objetivo
Implementar un exportador a **CSV (Comma-Separated Values)** y **TSV (Tab-Separated Values)** ultra-rápido y con **cero dependencias externas**, utilizando exclusivamente el paquete `encoding/csv` de la biblioteca estándar de Go.

---

## 2. Requerimientos Funcionales

### 2.1 Extracción de Datos
- El documento del reporte puede contener múltiples elementos (textos, tablas, gráficos).
- El motor de CSV debe permitir:
  1. **Exportar una tabla específica** seleccionada por ID o título.
  2. **Exportar todas las tablas** presentes en el reporte, separadas por bloques informativos legibles (título de sección y cabeceras).
  3. **Exportar datos de gráficos**: Exportar las series y categorías de los gráficos en formato tabular de datos brutos.

### 2.2 Compatibilidad Internacional con Excel (BOM)
- **UTF-8 Byte Order Mark (BOM)**: Incluir la opción de anteponer el prefijo `\xEF\xBB\xBF`. Esto garantiza que Microsoft Excel (especialmente en Windows y en idiomas como español) abra el CSV reconociendo automáticamente tildes, caracteres especiales (`ñ`, `€`, `á`, etc.) sin que el usuario deba importar manualmente el archivo.

### 2.3 Delimitadores Configurables
- Coma estándar (`,`) - Formato anglosajón por defecto.
- Punto y coma (`;`) - Estándar para países de habla hispana y Europa donde la coma es el separador decimal.
- Tabulador (`\t`) - Para archivos TSV.
- Salto de línea configurable (`\n` o `\r\n`).

### 2.4 Streaming de Alto Rendimiento
- Operar directamente sobre `io.Writer` o `bytes.Buffer` para evitar acumulación innecesaria en memoria.

---

## 3. Modelo de Configuración (`pkg/exporters/csv`)

```go
package csvexp

type CSVConfig struct {
    Delimiter   rune   // ',' o ';' o '\t'
    UseCRLF     bool   // true para Windows (\r\n)
    IncludeBOM  bool   // true para que Excel reconozca UTF-8 automáticamente
    TableID     string // Si está vacío, exporta todas las tablas
}
```

---

## 4. Criterios de Aceptación
1. **Zero Dependencies**: Uso exclusivo de `encoding/csv`.
2. **Velocidad Extrema**: Exportación de 10,000 filas en **menos de 5 milisegundos**.
3. **Compatibilidad Excel**: Al hacer doble clic en el archivo `.csv` en Windows con Excel, los caracteres en español deben leerse perfectamente.
