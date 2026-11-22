# Especificación 02: Motor de PDF

## 1. Objetivo
Implementar el motor de generación de documentos **PDF** utilizando la biblioteca `github.com/go-pdf/fpdf` (1 única dependencia externa sin dependencias transitivas). 

El motor proporcionará un sistema de **maquetación fluida por bloques** con soporte completo para encabezados corporativos, tablas avanzadas con repetición de cabeceras, gráficos vectoriales incrustados e imágenes.

---

## 2. Características del Motor PDF

### 2.1 Configuración de Página y Plantilla Base
- Formatos estándar: **A4**, **Letter** y **Legal**.
- Orientación: **Vertical (Portrait)** y **Horizontal (Landscape)**.
- Márgenes configurables (superior, inferior, izquierdo, derecho).
- **Encabezado y Pie de página automáticos**:
  - Logo corporativo a la izquierda o derecha.
  - Título del reporte y fecha de generación.
  - Numeración dinámica: *"Página {P} de {Total}"* usando el mecanismo de alias de páginas.
  - Línea divisoria sutil.

### 2.2 Bloques de Contenido
- **Títulos y Encabezados**: Niveles H1, H2, H3 con tipografías proporcionales, colores de acento y espaciado automático.
- **Párrafos y Textos**: Ajuste automático de texto (*word wrap*), interlineado, estilos (negrita, cursiva, subrayado) y alineación (izquierda, centro, derecha, justificado).
- **Cajas de Resumen / Métricas (*KPI Cards*)**: Contenedores redondeados con valor destacado, etiqueta descriptiva e indicador de cambio (+15%, -3%).
- **Líneas divisorias y saltos de página controlados**.

### 2.3 Motor de Tablas Avanzadas
- **Cálculo de Ancho de Columnas**:
  - Automático (proporcional al contenido) o porcentual explícito.
- **Estilos Visuales**:
  - Cabecera con fondo de color corporativo y texto en blanco o de alto contraste.
  - Filas alternadas con sombreado tenue (*zebra striping*).
  - Bordes de celda configurables (completos, solo horizontales o sin bordes).
  - Alineación independiente por columna (texto a la izquierda, números a la derecha, códigos al centro).
- **Paginación Inteligente**:
  - Cálculo de altura para celdas con texto multilínea.
  - Si una fila excede el espacio de la página actual, realiza un salto de página limpio y **repite automáticamente la cabecera de la tabla** en la nueva página.

### 2.4 Integración de Gráficos e Imágenes
- **Gráficos Vectoriales**: Consumo de las primitivas de `pkg/charts` dibujadas directamente sobre el canvas del PDF (`pdf.Rect`, `pdf.Polygon`, etc.).
- **Imágenes**: Inserción de imágenes PNG y JPEG desde archivo en disco, URL o buffer en memoria (Base64), con escalado proporcional y ajuste de resolución DPI.

---

## 3. Arquitectura del Explotador (`pkg/exporters/pdf`)

```
pkg/exporters/pdf/
├── pdf.go              # Generador principal y ciclo de vida del documento
├── header_footer.go    # Gestión de cabeceras, pies y numeración
├── layout.go           # Control de flujo de página, márgenes y saltos
├── table.go            # Renderizado de tablas avanzadas y cálculo de celdas
├── chart_renderer.go   # Traductor de primitivas vectoriales a fpdf
└── image.go            # Registro y decodificación de imágenes
```

---

## 4. API de Uso en Go

```go
pdfExp := pdf.NewExporter(pdf.Config{
    PageSize:    pdf.PageSizeA4,
    Orientation: pdf.OrientationPortrait,
    Theme:       styles.ThemeModernBlue,
})

err := pdfExp.Render(doc, "salida.pdf") // O hacia un io.Writer / bytes.Buffer
```

---

## 5. Criterios de Aceptación
1. **Rendimiento**: Un reporte de 5 páginas con 2 gráficos y 1 tabla de 50 filas debe generarse en **menos de 3 milisegundos**.
2. **Integridad Visual**: No deben existir solapamientos de texto al saltar de página ni cortes abruptos en filas de tabla.
3. **Validación**: El PDF generado debe ser 100% compatible con Adobe Acrobat, Chrome PDF Viewer, macOS Preview y visores móviles.
