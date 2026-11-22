package charts

import (
	"fmt"
)

// Generate procesa la configuración y devuelve las operaciones vectoriales elementales.
func Generate(cfg Config) ([]DrawOp, error) {
	if cfg.Width <= 0 {
		cfg.Width = 500
	}
	if cfg.Height <= 0 {
		cfg.Height = 300
	}

	switch cfg.Type {
	case ChartTypeBar:
		return renderBarChart(cfg)
	case ChartTypeLine:
		return renderLineChart(cfg)
	case ChartTypePie:
		return renderPieChart(cfg, false)
	case ChartTypeDonut:
		return renderPieChart(cfg, true)
	default:
		return nil, fmt.Errorf("tipo de gráfico no soportado: %s", cfg.Type)
	}
}

// GenerateSVG procesa la configuración y devuelve el gráfico directamente como SVG estándar.
func GenerateSVG(cfg Config) (string, error) {
	ops, err := Generate(cfg)
	if err != nil {
		return "", err
	}
	return RenderSVG(ops, cfg.Width, cfg.Height), nil
}

// --- Fluent Builders para Gráficos en Go ---

// BarChartBuilder construye gráficos de barras con API fluida.
type BarChartBuilder struct {
	cfg Config
}

func NewBarChart(title string) *BarChartBuilder {
	return &BarChartBuilder{
		cfg: Config{
			Type:       ChartTypeBar,
			Title:      title,
			ShowGrid:   true,
			ShowLegend: true,
			Width:      500,
			Height:     300,
		},
	}
}

func (b *BarChartBuilder) SetSubtitle(sub string) *BarChartBuilder {
	b.cfg.Subtitle = sub
	return b
}

func (b *BarChartBuilder) SetSize(width, height float64) *BarChartBuilder {
	b.cfg.Width = width
	b.cfg.Height = height
	return b
}

func (b *BarChartBuilder) SetCategories(cats ...string) *BarChartBuilder {
	b.cfg.Categories = cats
	return b
}

func (b *BarChartBuilder) AddSeries(name string, values []float64, colorHex ...string) *BarChartBuilder {
	col := ""
	if len(colorHex) > 0 {
		col = colorHex[0]
	}
	b.cfg.Series = append(b.cfg.Series, Series{
		Name:   name,
		Color:  col,
		Values: values,
	})
	return b
}

func (b *BarChartBuilder) SetShowLegend(show bool) *BarChartBuilder {
	b.cfg.ShowLegend = show
	return b
}

func (b *BarChartBuilder) SetShowValues(show bool) *BarChartBuilder {
	b.cfg.ShowValues = show
	return b
}

func (b *BarChartBuilder) SetValueFormat(prefix, suffix string) *BarChartBuilder {
	b.cfg.ValuePrefix = prefix
	b.cfg.ValueSuffix = suffix
	return b
}

func (b *BarChartBuilder) Config() Config {
	return b.cfg
}

// LineChartBuilder construye gráficos de líneas con API fluida.
type LineChartBuilder struct {
	cfg Config
}

func NewLineChart(title string) *LineChartBuilder {
	return &LineChartBuilder{
		cfg: Config{
			Type:       ChartTypeLine,
			Title:      title,
			ShowGrid:   true,
			ShowLegend: true,
			Width:      500,
			Height:     300,
		},
	}
}

func (b *LineChartBuilder) SetSubtitle(sub string) *LineChartBuilder {
	b.cfg.Subtitle = sub
	return b
}

func (b *LineChartBuilder) SetSize(width, height float64) *LineChartBuilder {
	b.cfg.Width = width
	b.cfg.Height = height
	return b
}

func (b *LineChartBuilder) SetCategories(cats ...string) *LineChartBuilder {
	b.cfg.Categories = cats
	return b
}

func (b *LineChartBuilder) AddSeries(name string, values []float64, colorHex ...string) *LineChartBuilder {
	col := ""
	if len(colorHex) > 0 {
		col = colorHex[0]
	}
	b.cfg.Series = append(b.cfg.Series, Series{
		Name:   name,
		Color:  col,
		Values: values,
	})
	return b
}

func (b *LineChartBuilder) SetShowLegend(show bool) *LineChartBuilder {
	b.cfg.ShowLegend = show
	return b
}

func (b *LineChartBuilder) SetShowValues(show bool) *LineChartBuilder {
	b.cfg.ShowValues = show
	return b
}

func (b *LineChartBuilder) SetValueFormat(prefix, suffix string) *LineChartBuilder {
	b.cfg.ValuePrefix = prefix
	b.cfg.ValueSuffix = suffix
	return b
}

func (b *LineChartBuilder) Config() Config {
	return b.cfg
}

// PieChartBuilder construye gráficos de pastel o dona con API fluida.
type PieChartBuilder struct {
	cfg Config
}

func NewPieChart(title string) *PieChartBuilder {
	return &PieChartBuilder{
		cfg: Config{
			Type:       ChartTypePie,
			Title:      title,
			ShowLegend: true,
			ShowValues: true,
			Width:      500,
			Height:     300,
		},
	}
}

func NewDonutChart(title string) *PieChartBuilder {
	return &PieChartBuilder{
		cfg: Config{
			Type:       ChartTypeDonut,
			Title:      title,
			ShowLegend: true,
			ShowValues: true,
			DonutRatio: 0.55,
			Width:      500,
			Height:     300,
		},
	}
}

func (b *PieChartBuilder) SetSubtitle(sub string) *PieChartBuilder {
	b.cfg.Subtitle = sub
	return b
}

func (b *PieChartBuilder) SetSize(width, height float64) *PieChartBuilder {
	b.cfg.Width = width
	b.cfg.Height = height
	return b
}

func (b *PieChartBuilder) SetDonutRatio(ratio float64) *PieChartBuilder {
	b.cfg.DonutRatio = ratio
	return b
}

func (b *PieChartBuilder) SetData(items []string, values []float64) *PieChartBuilder {
	b.cfg.Categories = items
	b.cfg.Series = []Series{
		{
			Name:   "Valores",
			Values: values,
		},
	}
	return b
}

func (b *PieChartBuilder) SetValueFormat(prefix, suffix string) *PieChartBuilder {
	b.cfg.ValuePrefix = prefix
	b.cfg.ValueSuffix = suffix
	return b
}

func (b *PieChartBuilder) Config() Config {
	return b.cfg
}
