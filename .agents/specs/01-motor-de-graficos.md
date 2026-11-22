# Especificación 01: Motor de Gráficos (Zero Dependencies)

## 1. Objetivo
Desarrollar un motor de generación de gráficos estadísticos en **Go puro** sin dependencias externas, utilizando exclusivamente el paquete matemático (`math`) y la biblioteca estándar de Go.

El motor debe producir descripciones geométricas vectoriales que puedan:
1. **Dibujarse directamente en el canvas vectorial del PDF** (`pdf.Rect`, `pdf.Polygon`, `pdf.Line`, etc.) para máxima nitidez y peso mínimo.
2. **Generar código SVG estándar** (`<svg>...</svg>`) para vistas previas en navegador o exportación independiente.

---

## 2. Tipos de Gráficos Soportados

### 2.1 Gráfico de Barras (`BarChart`)
- **Variantes**: Barras verticales agrupadas y apiladas.
- **Componentes**:
  - Ejes X e Y con marcas de graduación (*ticks*) y líneas de cuadrícula tenue.
  - Etiquetas de categorías en el eje X.
  - Escala automática inteligente en el eje Y (máximo redondeado a múltiplos de 5, 10, 50, etc.).
  - Barras coloreadas con paleta corporativa.
  - Valores numéricos opcionales sobre cada barra.
  - Leyenda superior o lateral.

### 2.2 Gráfico de Líneas (`LineChart`)
- **Variantes**: Líneas simples, multiserie y áreas sombreadas bajo la curva (*Area Chart*).
- **Componentes**:
  - Cuadrícula milimétrica o suave.
  - Puntos de datos resaltados con círculos vectoriales.
  - Conexión de puntos mediante segmentos de línea continua.
  - Leyenda de series.

### 2.3 Gráfico de Pastel y Dona (`PieChart` / `DonutChart`)
- **Componentes**:
  - Cálculo trigonométrico de ángulos y sectores (`math.Cos`, `math.Sin`).
  - Agrupación por porcentaje o valor absoluto.
  - Efecto Dona con radio interno configurable.
  - Etiquetas de porcentaje colocadas radialmente o mediante líneas guía.
  - Cuadro de leyenda con color, nombre y porcentaje.

---

## 3. Modelo de Datos y Estructuras (`pkg/charts`)

```go
package charts

type ChartType string

const (
    ChartTypeBar   ChartType = "bar"
    ChartTypeLine  ChartType = "line"
    ChartTypePie   ChartType = "pie"
    ChartTypeDonut ChartType = "donut"
)

type Series struct {
    Name   string    `json:"name"`
    Color  string    `json:"color,omitempty"`
    Values []float64 `json:"values"`
}

type ChartConfig struct {
    Type       ChartType `json:"type"`
    Title      string    `json:"title"`
    Width      float64   `json:"width"`      // en mm o px
    Height     float64   `json:"height"`     // en mm o px
    Categories []string  `json:"categories"` // Etiquetas eje X
    Series     []Series  `json:"series"`
    ShowLegend bool      `json:"show_legend"`
    ShowValues bool      `json:"show_values"`
    Colors     []string  `json:"colors,omitempty"`
}
```

---

## 4. Primitivas de Dibujo Vectorial

Para abstraer el renderizador (PDF vs SVG), el motor calculará una lista de primitivas vectoriales:

```go
type DrawOpType int

const (
    OpRect DrawOpType = iota
    OpLine
    OpPolygon
    OpCircle
    OpText
)

type DrawOp struct {
    Type      DrawOpType
    X, Y      float64
    Width     float64
    Height    float64
    Points    [][2]float64
    Color     [3]int   // RGB
    FillColor [3]int   // RGB
    Text      string
    FontSize  float64
    Align     string
}
```

---

## 5. Criterios de Aceptación y Pruebas
1. **Zero Dependencies**: `go list -m all` no debe incluir ningún paquete externo para el módulo `pkg/charts`.
2. **Generación Instantánea**: Un gráfico con 10 series y 50 puntos debe calcularse en menos de **0.5 milisegundos**.
3. **Salida SVG Válida**: Comprobar que el string SVG generado abre correctamente en navegadores y visores.
4. **Dibujo en PDF**: El exportador de PDF debe consumir las primitivas y renderizar el gráfico sin artefactos visuales.
