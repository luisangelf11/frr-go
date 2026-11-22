package spec

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/luisangelf11/reporteria-go/pkg/charts"
	"github.com/luisangelf11/reporteria-go/pkg/reporter"
)

// ParseSpecJSON decodifica un JSON crudo y construye la instancia de reporter.Document.
func ParseSpecJSON(data []byte) (*reporter.Document, error) {
	var s ReportSpec
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("error decodificando spec JSON: %w", err)
	}
	return ConvertSpecToDoc(s), nil
}

// ConvertSpecToDoc transforma una estructura ReportSpec en un reporter.Document listo para renderizar.
func ConvertSpecToDoc(s ReportSpec) *reporter.Document {
	b := reporter.New(s.Title)

	if s.Subtitle != "" {
		b.SetSubtitle(s.Subtitle)
	}
	if s.Author != "" {
		b.SetAuthor(s.Author)
	}
	if s.Date != "" {
		b.SetDate(s.Date)
	}

	// Tema
	switch strings.ToLower(s.Theme) {
	case "emerald", "verde":
		b.SetTheme(reporter.ThemeEmerald)
	case "dark_slate", "slate", "oscuro":
		b.SetTheme(reporter.ThemeDarkSlate)
	default:
		b.SetTheme(reporter.ThemeCorporateBlue)
	}

	// Tamaño de página
	switch strings.ToUpper(s.PageSize) {
	case "LETTER", "CARTA":
		b.SetPageSize(reporter.PageSizeLetter)
	case "LEGAL", "OFICIO":
		b.SetPageSize(reporter.PageSizeLegal)
	default:
		b.SetPageSize(reporter.PageSizeA4)
	}

	// Orientación
	if strings.EqualFold(s.Orientation, "landscape") || strings.EqualFold(s.Orientation, "horizontal") {
		b.SetOrientation(reporter.OrientationLandscape)
	}

	// Procesar Bloques
	for _, blk := range s.Blocks {
		switch strings.ToLower(blk.Type) {
		case "heading":
			if blk.Heading != nil {
				lvl := reporter.H1
				if blk.Heading.Level == 2 {
					lvl = reporter.H2
				} else if blk.Heading.Level == 3 {
					lvl = reporter.H3
				}
				b.AddHeading(blk.Heading.Text, lvl)
			}

		case "paragraph":
			if blk.Paragraph != nil {
				var opts []reporter.ParagraphOption
				if blk.Paragraph.Bold {
					opts = append(opts, reporter.WithBold(true))
				}
				if blk.Paragraph.Italic {
					opts = append(opts, reporter.WithItalic(true))
				}
				if blk.Paragraph.FontSize > 0 {
					opts = append(opts, reporter.WithFontSize(blk.Paragraph.FontSize))
				}
				if blk.Paragraph.TextColor != "" {
					opts = append(opts, reporter.WithColor(charts.HexColor(blk.Paragraph.TextColor)))
				}
				align := parseAlign(blk.Paragraph.Align)
				opts = append(opts, reporter.WithAlign(align))

				b.AddParagraph(blk.Paragraph.Text, opts...)
			}

		case "table":
			if blk.Table != nil {
				tb := reporter.NewTable().
					SetID(blk.Table.ID).
					SetTitle(blk.Table.Title).
					SetHeaders(blk.Table.Headers...).
					AddRows(blk.Table.Rows)

				if len(blk.Table.ColumnWidths) > 0 {
					tb.SetColumnWidths(blk.Table.ColumnWidths...)
				}
				if len(blk.Table.Alignments) > 0 {
					var aligns []reporter.TextAlign
					for _, a := range blk.Table.Alignments {
						aligns = append(aligns, parseAlign(a))
					}
					tb.SetAlignments(aligns...)
				}
				if blk.Table.ZebraStriping != nil {
					tb.SetZebraStriping(*blk.Table.ZebraStriping)
				}
				b.AddTable(tb.Build())
			}

		case "chart":
			if blk.Chart != nil {
				cfg := charts.Config{
					Type:        charts.ChartType(blk.Chart.Type),
					Title:       blk.Chart.Title,
					Subtitle:    blk.Chart.Subtitle,
					Width:       blk.Chart.Width,
					Height:      blk.Chart.Height,
					Categories:  blk.Chart.Categories,
					Colors:      blk.Chart.Colors,
					DonutRatio:  blk.Chart.DonutRatio,
					ValuePrefix: blk.Chart.ValuePrefix,
					ValueSuffix: blk.Chart.ValueSuffix,
				}
				if blk.Chart.ShowLegend != nil {
					cfg.ShowLegend = *blk.Chart.ShowLegend
				} else {
					cfg.ShowLegend = true
				}
				if blk.Chart.ShowValues != nil {
					cfg.ShowValues = *blk.Chart.ShowValues
				}
				if blk.Chart.ShowGrid != nil {
					cfg.ShowGrid = *blk.Chart.ShowGrid
				} else {
					cfg.ShowGrid = true
				}

				for _, s := range blk.Chart.Series {
					cfg.Series = append(cfg.Series, charts.Series{
						Name:   s.Name,
						Color:  s.Color,
						Values: s.Values,
					})
				}

				b.AddChart(cfg, blk.Chart.Width, blk.Chart.Height)
			}

		case "metric_cards":
			if blk.MetricCards != nil {
				var cards []reporter.MetricCard
				for _, c := range blk.MetricCards.Cards {
					cards = append(cards, reporter.MetricCard{
						Title:       c.Title,
						Value:       c.Value,
						Change:      c.Change,
						IsPositive:  c.IsPositive,
						Description: c.Description,
					})
				}
				b.AddMetricCards(cards...)
			}

		case "divider":
			b.AddDivider()

		case "page_break":
			b.AddPageBreak()
		}
	}

	return b.Build()
}

func parseAlign(s string) reporter.TextAlign {
	switch strings.ToLower(s) {
	case "center", "centro":
		return reporter.AlignCenter
	case "right", "derecha":
		return reporter.AlignRight
	case "justify", "justificado":
		return reporter.AlignJustify
	default:
		return reporter.AlignLeft
	}
}
