package reporter

import (
	"time"

	"github.com/luisangelf11/frr-go/pkg/charts"
)

// ReportBuilder permite construir documentos mediante una API fluida y ergonómica.
type ReportBuilder struct {
	doc Document
}

// New inicia la construcción de un nuevo reporte con su título principal.
func New(title string) *ReportBuilder {
	return &ReportBuilder{
		doc: Document{
			Title:        title,
			Date:         time.Now().Format("2006-01-02"),
			Theme:        ThemeCorporateBlue,
			PageSize:     PageSizeA4,
			Orientation:  OrientationPortrait,
			ShowHeader:   true,
			ShowFooter:   true,
			MarginLeft:   15.0, // 15mm margen estándar
			MarginRight:  15.0,
			MarginTop:    15.0,
			MarginBottom: 18.0,
			Blocks:       make([]Block, 0),
		},
	}
}

// SetSubtitle define un subtítulo descriptivo bajo el título del reporte.
func (b *ReportBuilder) SetSubtitle(sub string) *ReportBuilder {
	b.doc.Subtitle = sub
	return b
}

// SetAuthor define el autor o departamento emisor del reporte.
func (b *ReportBuilder) SetAuthor(author string) *ReportBuilder {
	b.doc.Author = author
	return b
}

// SetDate define una fecha explícita para el reporte.
func (b *ReportBuilder) SetDate(date string) *ReportBuilder {
	b.doc.Date = date
	return b
}

// SetLogo establece la ruta de un archivo de imagen (PNG/JPEG) para el logo corporativo.
func (b *ReportBuilder) SetLogo(path string) *ReportBuilder {
	b.doc.LogoPath = path
	return b
}

// SetLogoBytes establece los bytes de la imagen del logo directamente en memoria.
func (b *ReportBuilder) SetLogoBytes(data []byte) *ReportBuilder {
	b.doc.LogoBytes = data
	return b
}

// SetTheme aplica una paleta y tema corporativo visual.
func (b *ReportBuilder) SetTheme(theme Theme) *ReportBuilder {
	b.doc.Theme = theme
	return b
}

// SetPageSize configura el tamaño físico de hoja (A4, Letter, Legal).
func (b *ReportBuilder) SetPageSize(size PageSize) *ReportBuilder {
	b.doc.PageSize = size
	return b
}

// SetOrientation configura la orientación de la hoja (Vertical u Horizontal).
func (b *ReportBuilder) SetOrientation(orient Orientation) *ReportBuilder {
	b.doc.Orientation = orient
	return b
}

// SetMargins define los márgenes en milímetros.
func (b *ReportBuilder) SetMargins(left, top, right, bottom float64) *ReportBuilder {
	b.doc.MarginLeft = left
	b.doc.MarginTop = top
	b.doc.MarginRight = right
	b.doc.MarginBottom = bottom
	return b
}

// SetShowHeader activa o desactiva la cabecera en cada página.
func (b *ReportBuilder) SetShowHeader(show bool) *ReportBuilder {
	b.doc.ShowHeader = show
	return b
}

// SetShowFooter activa o desactiva el pie de página con numeración.
func (b *ReportBuilder) SetShowFooter(show bool) *ReportBuilder {
	b.doc.ShowFooter = show
	return b
}

// AddHeading añade un título jerárquico (H1, H2, H3).
func (b *ReportBuilder) AddHeading(text string, level HeadingLevel) *ReportBuilder {
	b.doc.Blocks = append(b.doc.Blocks, Block{
		Type: BlockTypeHeading,
		Heading: &HeadingBlock{
			Text:  text,
			Level: level,
		},
	})
	return b
}

// ParagraphOption define modificadores funcionales para párrafos.
type ParagraphOption func(*ParagraphBlock)

// WithBold aplica negrita al párrafo.
func WithBold(bold bool) ParagraphOption {
	return func(p *ParagraphBlock) {
		p.Bold = bold
	}
}

// WithItalic aplica cursiva al párrafo.
func WithItalic(italic bool) ParagraphOption {
	return func(p *ParagraphBlock) {
		p.Italic = italic
	}
}

// WithAlign define la alineación del texto.
func WithAlign(align TextAlign) ParagraphOption {
	return func(p *ParagraphBlock) {
		p.Align = align
	}
}

// WithColor define un color de texto personalizado.
func WithColor(color charts.Color) ParagraphOption {
	return func(p *ParagraphBlock) {
		p.TextColor = &color
	}
}

// WithFontSize define un tamaño de fuente en puntos.
func WithFontSize(size float64) ParagraphOption {
	return func(p *ParagraphBlock) {
		p.FontSize = size
	}
}

// AddParagraph añade un párrafo de texto normal o formateado.
func (b *ReportBuilder) AddParagraph(text string, opts ...ParagraphOption) *ReportBuilder {
	p := &ParagraphBlock{
		Text:     text,
		Align:    AlignLeft,
		FontSize: 10.0,
	}
	for _, opt := range opts {
		opt(p)
	}
	b.doc.Blocks = append(b.doc.Blocks, Block{
		Type:      BlockTypeParagraph,
		Paragraph: p,
	})
	return b
}

// AddTable añade una tabla de datos estructurada.
func (b *ReportBuilder) AddTable(table TableBlock) *ReportBuilder {
	b.doc.Blocks = append(b.doc.Blocks, Block{
		Type:  BlockTypeTable,
		Table: &table,
	})
	return b
}

// AddChart añade un gráfico estadístico generado con pkg/charts.
func (b *ReportBuilder) AddChart(cfg charts.Config, dims ...float64) *ReportBuilder {
	w := 180.0 // Por defecto ancho completo en página A4 (210 - 30mm)
	h := 90.0  // Altura estándar balanceada
	if len(dims) > 0 && dims[0] > 0 {
		w = dims[0]
	}
	if len(dims) > 1 && dims[1] > 0 {
		h = dims[1]
	}

	b.doc.Blocks = append(b.doc.Blocks, Block{
		Type: BlockTypeChart,
		Chart: &ChartBlock{
			Config: cfg,
			Width:  w,
			Height: h,
		},
	})
	return b
}

// AddImage añade una imagen desde un archivo en disco.
func (b *ReportBuilder) AddImage(path string, width, height float64, align ...TextAlign) *ReportBuilder {
	a := AlignCenter
	if len(align) > 0 {
		a = align[0]
	}
	b.doc.Blocks = append(b.doc.Blocks, Block{
		Type: BlockTypeImage,
		Image: &ImageBlock{
			Path:   path,
			Width:  width,
			Height: height,
			Align:  a,
		},
	})
	return b
}

// AddMetricCards añade un grupo de tarjetas de indicadores clave (KPIs).
func (b *ReportBuilder) AddMetricCards(cards ...MetricCard) *ReportBuilder {
	b.doc.Blocks = append(b.doc.Blocks, Block{
		Type: BlockTypeMetricCards,
		MetricCards: &MetricCardsBlock{
			Cards: cards,
		},
	})
	return b
}

// AddPageBreak fuerza un salto de página limpio.
func (b *ReportBuilder) AddPageBreak() *ReportBuilder {
	b.doc.Blocks = append(b.doc.Blocks, Block{
		Type: BlockTypePageBreak,
	})
	return b
}

// AddDivider inserta una línea horizontal divisoria sutil.
func (b *ReportBuilder) AddDivider() *ReportBuilder {
	b.doc.Blocks = append(b.doc.Blocks, Block{
		Type: BlockTypeDivider,
	})
	return b
}

// Build finaliza la construcción y retorna la instancia del documento.
func (b *ReportBuilder) Build() *Document {
	return &b.doc
}

// --- Table Builder ---

// TableBuilder facilita la construcción de tablas con API fluida.
type TableBuilder struct {
	table TableBlock
}

// NewTable inicializa un nuevo constructor de tabla.
func NewTable() *TableBuilder {
	return &TableBuilder{
		table: TableBlock{
			Headers:          make([]string, 0),
			Rows:             make([][]string, 0),
			ColumnWidths:     make([]float64, 0),
			ColumnAlignments: make([]TextAlign, 0),
			ZebraStriping:    true,
		},
	}
}

func (tb *TableBuilder) SetID(id string) *TableBuilder {
	tb.table.ID = id
	return tb
}

func (tb *TableBuilder) SetTitle(title string) *TableBuilder {
	tb.table.Title = title
	return tb
}

func (tb *TableBuilder) SetHeaders(headers ...string) *TableBuilder {
	tb.table.Headers = headers
	return tb
}

func (tb *TableBuilder) AddRow(cells ...string) *TableBuilder {
	tb.table.Rows = append(tb.table.Rows, cells)
	return tb
}

func (tb *TableBuilder) AddRows(rows [][]string) *TableBuilder {
	tb.table.Rows = append(tb.table.Rows, rows...)
	return tb
}

func (tb *TableBuilder) SetColumnWidths(widths ...float64) *TableBuilder {
	tb.table.ColumnWidths = widths
	return tb
}

func (tb *TableBuilder) SetAlignments(aligns ...TextAlign) *TableBuilder {
	tb.table.ColumnAlignments = aligns
	return tb
}

func (tb *TableBuilder) SetZebraStriping(zebra bool) *TableBuilder {
	tb.table.ZebraStriping = zebra
	return tb
}

func (tb *TableBuilder) Build() TableBlock {
	return tb.table
}
