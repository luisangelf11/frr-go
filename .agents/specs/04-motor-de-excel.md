# Especificación 04: Motor de Excel (.xlsx)

## 1. Objetivo
Desarrollar un generador de hojas de cálculo de formato nativo **Microsoft Excel (.xlsx / OpenXML)** de alto rendimiento y **cero dependencias externas**, aprovechando los paquetes `archive/zip` y `encoding/xml` de la biblioteca estándar de Go.

---

## 2. Anatomía de un Archivo .xlsx Nativo
Un archivo `.xlsx` es un contenedor comprimido (ZIP) que alberga archivos XML estandarizados (Open Packaging Conventions - ISO/IEC 29500):

```
reporte.xlsx (ZIP)
├── [Content_Types].xml              # Define tipos MIME de las partes del archivo
├── _rels/
│   └── .rels                        # Relaciones raíz del paquete
└── xl/
    ├── workbook.xml                 # Definición de hojas (sheets)
    ├── _rels/
    │   └── workbook.xml.rels        # Relaciones del workbook con las hojas y estilos
    ├── styles.xml                   # Paleta de fuentes, colores, bordes y formatos
    └── worksheets/
        ├── sheet1.xml               # Datos, filas, celdas y fórmulas de la hoja 1
        └── sheet2.xml               # Hoja 2 (si aplica)
```

Al utilizar `inlineStr` (cadenas en línea) o una tabla compartida compacta y escribir los XML directamente al compresor ZIP en modo streaming:
- Se elimina por completo el uso de memoria residual.
- Se generan hojas de cálculo con decenas de miles de celdas en **menos de 5 milisegundos**.
- **Cero dependencias externas**.

---

## 3. Características Soportadas

### 3.1 Estilos y Formato Visual
- **Fila de Cabecera Corporativa**:
  - Color de fondo (*fill*) según el tema elegido (ej. Azul Marino `#1E3A8A`, Esmeralda `#065F46`, etc.).
  - Tipografía en negrita con color de texto en blanco o contraste alto.
  - Alineación centrada o justificada.
- **Filas de Datos con *Zebra Striping***:
  - Alternancia automática de filas con fondo blanco y gris claro tenue (`#F9FAFB`).
- **Bordes de Celda**:
  - Bordes finos elegantes para delimitar la cuadrícula de datos.
- **Tipos de Datos y Formatos Numéricos**:
  - Texto general.
  - Números enteros y decimales con separador de miles.
  - Moneda (`$#,##0.00`).
  - Porcentajes (`0.0%`).
  - Fechas (`YYYY-MM-DD`).

### 3.2 Usabilidad de la Hoja
- **Auto-Ajuste de Columnas (*Auto-width*)**:
  - Cálculo automático del ancho óptimo de cada columna según la longitud del contenido para evitar texto truncado o columnas estrechas.
- **Inmovilización de Paneles (*Freeze Panes*)**:
  - La fila de cabecera queda fija al hacer scroll hacia abajo para no perder de vista los nombres de las columnas.
- **Soporte Multi-Hoja (*Multi-sheet*)**:
  - Capacidad de crear hojas separadas (ej. "Resumen Ejecutivo", "Detalle de Transacciones", "Métricas").

---

## 4. API de Uso en Go

```go
excelExp := excel.NewExporter(excel.Config{
    Theme:       styles.ThemeCorporateBlue,
    FreezeHeader: true,
    AutoWidth:    true,
})

err := excelExp.Render(doc, "ventas.xlsx")
```

---

## 5. Criterios de Aceptación
1. **Zero Dependencies**: Ningún paquete externo en `go.mod`.
2. **Compatibilidad Total**: El archivo `.xlsx` debe abrirse directamente en **Microsoft Excel** (Windows y Mac), **Google Sheets** y **LibreOffice Calc** sin ninguna advertencia de reparación o formato inválido.
3. **Velocidad**: Generar un archivo con 5 hojas y 5,000 celdas con estilos en **menos de 4 milisegundos**.
