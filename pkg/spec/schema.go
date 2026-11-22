package spec

// ReportSpec es la estructura JSON interoperable entre Go y Node.js.
type ReportSpec struct {
	Title       string      `json:"title"`
	Subtitle    string      `json:"subtitle,omitempty"`
	Author      string      `json:"author,omitempty"`
	Date        string      `json:"date,omitempty"`
	Theme       string      `json:"theme,omitempty"`       // "corporate_blue", "emerald", "dark_slate"
	PageSize    string      `json:"page_size,omitempty"`   // "A4", "Letter", "Legal"
	Orientation string      `json:"orientation,omitempty"` // "portrait", "landscape"
	Blocks      []BlockSpec `json:"blocks"`
}

// BlockSpec representa un elemento individual del reporte en formato JSON.
type BlockSpec struct {
	Type        string           `json:"type"` // "heading", "paragraph", "table", "chart", "metric_cards", "divider", "page_break"
	Heading     *HeadingSpec     `json:"heading,omitempty"`
	Paragraph   *ParagraphSpec   `json:"paragraph,omitempty"`
	Table       *TableSpec       `json:"table,omitempty"`
	Chart       *ChartSpec       `json:"chart,omitempty"`
	MetricCards *MetricCardsSpec `json:"metric_cards,omitempty"`
}

type HeadingSpec struct {
	Text  string `json:"text"`
	Level int    `json:"level"` // 1, 2, 3
}

type ParagraphSpec struct {
	Text      string  `json:"text"`
	Bold      bool    `json:"bold,omitempty"`
	Italic    bool    `json:"italic,omitempty"`
	Align     string  `json:"align,omitempty"` // "left", "center", "right", "justify"
	FontSize  float64 `json:"font_size,omitempty"`
	TextColor string  `json:"text_color,omitempty"` // "#1F2937"
}

type TableSpec struct {
	ID            string     `json:"id,omitempty"`
	Title         string     `json:"title,omitempty"`
	Headers       []string   `json:"headers"`
	Rows          [][]string `json:"rows"`
	ColumnWidths  []float64  `json:"column_widths,omitempty"`
	Alignments    []string   `json:"alignments,omitempty"` // ["left", "center", "right"]
	ZebraStriping *bool      `json:"zebra_striping,omitempty"`
}

type ChartSpec struct {
	Type        string        `json:"type"` // "bar", "line", "pie", "donut"
	Title       string        `json:"title,omitempty"`
	Subtitle    string        `json:"subtitle,omitempty"`
	Width       float64       `json:"width,omitempty"`
	Height      float64       `json:"height,omitempty"`
	Categories  []string      `json:"categories,omitempty"`
	Series      []SeriesSpec  `json:"series"`
	ShowLegend  *bool         `json:"show_legend,omitempty"`
	ShowValues  *bool         `json:"show_values,omitempty"`
	ShowGrid    *bool         `json:"show_grid,omitempty"`
	Colors      []string      `json:"colors,omitempty"`
	DonutRatio  float64       `json:"donut_ratio,omitempty"`
	ValuePrefix string        `json:"value_prefix,omitempty"`
	ValueSuffix string        `json:"value_suffix,omitempty"`
}

type SeriesSpec struct {
	Name   string    `json:"name"`
	Color  string    `json:"color,omitempty"`
	Values []float64 `json:"values"`
}

type MetricCardsSpec struct {
	Cards []MetricCardSpec `json:"cards"`
}

type MetricCardSpec struct {
	Title       string `json:"title"`
	Value       string `json:"value"`
	Change      string `json:"change,omitempty"`
	IsPositive  bool   `json:"is_positive"`
	Description string `json:"description,omitempty"`
}

// RequestSpec define la petición recibida por el ejecutable de Go (formato + especificación).
type RequestSpec struct {
	ID        string      `json:"id,omitempty"`
	Format    string      `json:"format"` // "pdf", "excel", "csv"
	Delimiter string      `json:"delimiter,omitempty"` // para CSV: ",", ";", "\t"
	Spec      ReportSpec  `json:"spec"`
}
