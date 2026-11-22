package reporter

import (
	"github.com/luisangelf11/reporteria-go/pkg/charts"
)

// BlockType identifica el tipo de bloque de contenido dentro del reporte.
type BlockType string

const (
	BlockTypeHeading     BlockType = "heading"
	BlockTypeParagraph   BlockType = "paragraph"
	BlockTypeTable       BlockType = "table"
	BlockTypeChart       BlockType = "chart"
	BlockTypeImage       BlockType = "image"
	BlockTypeMetricCards BlockType = "metric_cards"
	BlockTypePageBreak   BlockType = "page_break"
	BlockTypeDivider     BlockType = "divider"
)

// Block representa un elemento individual de contenido en el flujo del reporte.
type Block struct {
	Type        BlockType
	Heading     *HeadingBlock
	Paragraph   *ParagraphBlock
	Table       *TableBlock
	Chart       *ChartBlock
	Image       *ImageBlock
	MetricCards *MetricCardsBlock
}

// HeadingBlock representa un título o subtítulo jerárquico.
type HeadingBlock struct {
	Text  string
	Level HeadingLevel // 1 = H1, 2 = H2, 3 = H3
}

// ParagraphBlock representa un párrafo de texto con estilos opcionales.
type ParagraphBlock struct {
	Text      string
	Bold      bool
	Italic    bool
	Align     TextAlign
	TextColor *charts.Color
	FontSize  float64
}

// TableBlock representa una tabla estructurada de datos.
type TableBlock struct {
	ID               string
	Title            string
	Headers          []string
	Rows             [][]string
	ColumnWidths     []float64   // En mm o porcentual; si está vacío se calcula automáticamente
	ColumnAlignments []TextAlign // Alineación por columna
	ZebraStriping    bool        // Alternancia de color de fila
	HeaderColor      *charts.Color
}

// ChartBlock representa un gráfico estadístico integrado.
type ChartBlock struct {
	Config charts.Config
	Width  float64 // en mm (ej: 180mm para ancho completo A4)
	Height float64 // en mm (ej: 90mm)
}

// ImageBlock representa una imagen (PNG o JPEG) incrustada.
type ImageBlock struct {
	Path   string
	Data   []byte
	Width  float64
	Height float64
	Align  TextAlign
}

// MetricCard representa un indicador individual de resumen (KPI).
type MetricCard struct {
	Title       string
	Value       string
	Change      string // ej: "+15.2%"
	IsPositive  bool
	Description string
}

// MetricCardsBlock representa una fila de tarjetas de métricas/KPIs.
type MetricCardsBlock struct {
	Cards []MetricCard
}

// Document representa la estructura canónica en memoria de un reporte completo.
type Document struct {
	Title        string
	Subtitle     string
	Author       string
	Date         string
	LogoPath     string
	LogoBytes    []byte
	Theme        Theme
	PageSize     PageSize
	Orientation  Orientation
	ShowHeader   bool
	ShowFooter   bool
	MarginLeft   float64
	MarginRight  float64
	MarginTop    float64
	MarginBottom float64
	Blocks       []Block
}
