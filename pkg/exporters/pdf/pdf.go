package pdf

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/go-pdf/fpdf"
	"github.com/luisangelf11/reporteria-go/pkg/charts"
	"github.com/luisangelf11/reporteria-go/pkg/reporter"
)

// Exporter maneja la conversión de un Document en un archivo o buffer PDF.
type Exporter struct{}

// NewExporter crea una nueva instancia del exportador PDF.
func NewExporter() *Exporter {
	return &Exporter{}
}

// Render exporta el documento a una ruta de archivo en disco.
func (e *Exporter) Render(doc *reporter.Document, outPath string) error {
	buf, err := e.RenderToBytes(doc)
	if err != nil {
		return err
	}
	return os.WriteFile(outPath, buf, 0644)
}

// RenderToWriter escribe los bytes del PDF en un io.Writer.
func (e *Exporter) RenderToWriter(doc *reporter.Document, w io.Writer) error {
	buf, err := e.RenderToBytes(doc)
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	return err
}

// RenderToBytes genera el PDF completo en memoria y lo retorna como []byte.
func (e *Exporter) RenderToBytes(doc *reporter.Document) ([]byte, error) {
	orientation := "P"
	if doc.Orientation == reporter.OrientationLandscape {
		orientation = "L"
	}

	pageSize := "A4"
	if doc.PageSize == reporter.PageSizeLetter {
		pageSize = "Letter"
	} else if doc.PageSize == reporter.PageSizeLegal {
		pageSize = "Legal"
	}

	pdf := fpdf.New(orientation, "mm", pageSize, "")
	pdf.SetTitle(doc.Title, true)
	if doc.Author != "" {
		pdf.SetAuthor(doc.Author, true)
	}

	pageW, pageH := pdf.GetPageSize()
	marginLeft := doc.MarginLeft
	marginRight := doc.MarginRight
	marginTop := doc.MarginTop
	marginBottom := doc.MarginBottom
	availWidth := pageW - marginLeft - marginRight

	theme := doc.Theme

	// Traductor Unicode para codificación de fuentes PDF estándar (ISO-8859-1 / CP1252)
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	// Configuración de Header & Footer automáticos
	pdf.AliasNbPages("{nb}")

	if doc.ShowFooter {
		pdf.SetFooterFunc(func() {
			pdf.SetY(-marginBottom + 3)
			// Línea sutil sobre el pie
			pdf.SetDrawColor(theme.BorderColor.R, theme.BorderColor.G, theme.BorderColor.B)
			pdf.SetLineWidth(0.2)
			pdf.Line(marginLeft, pageH-marginBottom+1, pageW-marginRight, pageH-marginBottom+1)

			pdf.SetFont("Arial", "", 7.5)
			pdf.SetTextColor(theme.MutedTextColor.R, theme.MutedTextColor.G, theme.MutedTextColor.B)

			// Texto izquierdo (Título del reporte o confidencialidad)
			pdf.SetX(marginLeft)
			footerLeft := doc.Title
			if doc.Author != "" {
				footerLeft = fmt.Sprintf("%s • %s", doc.Title, doc.Author)
			}
			pdf.CellFormat(availWidth/2.0, 6, tr(footerLeft), "", 0, "L", false, 0, "")

			// Texto derecho (Paginación)
			pageStr := fmt.Sprintf("Página %d de {nb}", pdf.PageNo())
			pdf.CellFormat(availWidth/2.0, 6, tr(pageStr), "", 0, "R", false, 0, "")
		})
	}

	if doc.ShowHeader {
		pdf.SetHeaderFunc(func() {
			// En la página 1 dejamos que el título principal brille; en páginas posteriores mostramos encabezado compacto
			if pdf.PageNo() > 1 {
				pdf.SetY(marginTop - 6)
				pdf.SetFont("Arial", "I", 7.5)
				pdf.SetTextColor(theme.MutedTextColor.R, theme.MutedTextColor.G, theme.MutedTextColor.B)

				pdf.SetX(marginLeft)
				pdf.CellFormat(availWidth/2.0, 5, tr(doc.Title), "", 0, "L", false, 0, "")

				if doc.Date != "" {
					pdf.CellFormat(availWidth/2.0, 5, tr(doc.Date), "", 0, "R", false, 0, "")
				}

				pdf.SetDrawColor(theme.BorderColor.R, theme.BorderColor.G, theme.BorderColor.B)
				pdf.SetLineWidth(0.2)
				pdf.Line(marginLeft, marginTop-1, pageW-marginRight, marginTop-1)
			}
		})
	}

	pdf.SetMargins(marginLeft, marginTop, marginRight)
	pdf.SetAutoPageBreak(false, marginBottom)
	pdf.AddPage()

	curY := marginTop

	// 1. Portada / Banner de Título Principal en la Primera Página
	curY = renderTitleBanner(pdf, tr, doc, theme, marginLeft, curY, availWidth)

	// 2. Procesar Bloques de Contenido
	for _, blk := range doc.Blocks {
		switch blk.Type {
		case reporter.BlockTypeHeading:
			curY = renderHeading(pdf, tr, blk.Heading, theme, marginLeft, curY, availWidth, pageH, marginBottom)

		case reporter.BlockTypeParagraph:
			curY = renderParagraph(pdf, tr, blk.Paragraph, theme, marginLeft, curY, availWidth, pageH, marginBottom)

		case reporter.BlockTypeMetricCards:
			curY = renderMetricCards(pdf, tr, blk.MetricCards, theme, marginLeft, curY, availWidth, pageH, marginBottom)

		case reporter.BlockTypeTable:
			curY = renderTable(pdf, tr, blk.Table, theme, marginLeft, curY, availWidth, pageH, marginBottom)
			curY += 6.0 // Margen inferior tras tabla

		case reporter.BlockTypeChart:
			curY = renderChartBlock(pdf, tr, blk.Chart, marginLeft, curY, availWidth, pageH, marginBottom)

		case reporter.BlockTypeDivider:
			curY = renderDivider(pdf, theme, marginLeft, curY, availWidth, pageH, marginBottom)

		case reporter.BlockTypePageBreak:
			pdf.AddPage()
			curY = marginTop + 4.0

		case reporter.BlockTypeImage:
			curY = renderImage(pdf, blk.Image, marginLeft, curY, availWidth, pageH, marginBottom)
		}
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("error generando salida PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// renderTitleBanner dibuja el encabezado principal en la primera página con una barra de acento.
func renderTitleBanner(pdf *fpdf.Fpdf, tr func(string) string, doc *reporter.Document, theme reporter.Theme, startX, startY, availWidth float64) float64 {
	curY := startY

	// Barra vertical de acento corporativo
	accentBarW := 3.5
	titleAreaH := 20.0
	if doc.Subtitle != "" {
		titleAreaH = 26.0
	}

	pdf.SetFillColor(theme.PrimaryColor.R, theme.PrimaryColor.G, theme.PrimaryColor.B)
	pdf.Rect(startX, curY, accentBarW, titleAreaH, "F")

	textStartX := startX + accentBarW + 4.0

	// Título principal
	pdf.SetFont("Arial", "B", 18)
	pdf.SetTextColor(theme.PrimaryColor.R, theme.PrimaryColor.G, theme.PrimaryColor.B)
	pdf.Text(textStartX, curY+7, tr(doc.Title))

	// Subtítulo
	if doc.Subtitle != "" {
		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(theme.MutedTextColor.R, theme.MutedTextColor.G, theme.MutedTextColor.B)
		pdf.Text(textStartX, curY+13.5, tr(doc.Subtitle))
	}

	// Metadatos (Autor y Fecha) a la derecha
	if doc.Date != "" || doc.Author != "" {
		pdf.SetFont("Arial", "", 8.5)
		pdf.SetTextColor(theme.MutedTextColor.R, theme.MutedTextColor.G, theme.MutedTextColor.B)
		metaY := curY + 6.0
		if doc.Date != "" {
			metaStr := "Fecha: " + doc.Date
			trMeta := tr(metaStr)
			w := pdf.GetStringWidth(trMeta)
			pdf.Text(startX+availWidth-w, metaY, trMeta)
			metaY += 5.5
		}
		if doc.Author != "" {
			metaStr := "Autor: " + doc.Author
			trMeta := tr(metaStr)
			w := pdf.GetStringWidth(trMeta)
			pdf.Text(startX+availWidth-w, metaY, trMeta)
		}
	}

	curY += titleAreaH + 6.0

	// Línea sutil divisoria bajo el banner
	pdf.SetDrawColor(theme.BorderColor.R, theme.BorderColor.G, theme.BorderColor.B)
	pdf.SetLineWidth(0.25)
	pdf.Line(startX, curY, startX+availWidth, curY)

	return curY + 6.0
}

func renderHeading(pdf *fpdf.Fpdf, tr func(string) string, h *reporter.HeadingBlock, theme reporter.Theme, startX, curY, availW, pageH, bottomM float64) float64 {
	fontSize := 14.0
	spacingTop := 6.0
	spacingBottom := 4.0

	switch h.Level {
	case reporter.H1:
		fontSize = 14.0
		spacingTop = 7.0
		spacingBottom = 4.5
		pdf.SetTextColor(theme.PrimaryColor.R, theme.PrimaryColor.G, theme.PrimaryColor.B)
	case reporter.H2:
		fontSize = 11.5
		spacingTop = 5.0
		spacingBottom = 3.5
		pdf.SetTextColor(theme.SecondaryColor.R, theme.SecondaryColor.G, theme.SecondaryColor.B)
	case reporter.H3:
		fontSize = 9.5
		spacingTop = 4.0
		spacingBottom = 2.5
		pdf.SetTextColor(theme.TextColor.R, theme.TextColor.G, theme.TextColor.B)
	}

	if curY+spacingTop+fontSize+spacingBottom > pageH-bottomM {
		pdf.AddPage()
		curY = 25.0
	}

	curY += spacingTop
	pdf.SetFont("Arial", "B", fontSize)
	pdf.Text(startX, curY+4, tr(h.Text))

	return curY + spacingBottom + 4.0
}

func renderParagraph(pdf *fpdf.Fpdf, tr func(string) string, p *reporter.ParagraphBlock, theme reporter.Theme, startX, curY, availW, pageH, bottomM float64) float64 {
	fontStyle := ""
	if p.Bold {
		fontStyle += "B"
	}
	if p.Italic {
		fontStyle += "I"
	}
	fontSize := p.FontSize
	if fontSize <= 0 {
		fontSize = 9.0
	}

	if p.TextColor != nil {
		pdf.SetTextColor(p.TextColor.R, p.TextColor.G, p.TextColor.B)
	} else {
		pdf.SetTextColor(theme.TextColor.R, theme.TextColor.G, theme.TextColor.B)
	}

	pdf.SetFont("Arial", fontStyle, fontSize)
	lineH := fontSize * 0.45

	translatedText := tr(p.Text)
	lines := pdf.SplitLines([]byte(translatedText), availW)
	for _, l := range lines {
		if curY+lineH > pageH-bottomM {
			pdf.AddPage()
			curY = 25.0
		}
		align := "L"
		if p.Align == reporter.AlignCenter {
			align = "C"
		} else if p.Align == reporter.AlignRight {
			align = "R"
		} else if p.Align == reporter.AlignJustify {
			align = "J"
		}
		pdf.SetXY(startX, curY)
		pdf.CellFormat(availW, lineH, string(l), "", 1, align, false, 0, "")
		curY += lineH
	}

	return curY + 3.0
}

func renderMetricCards(pdf *fpdf.Fpdf, tr func(string) string, m *reporter.MetricCardsBlock, theme reporter.Theme, startX, curY, availW, pageH, bottomM float64) float64 {
	numCards := len(m.Cards)
	if numCards == 0 {
		return curY
	}

	cardH := 22.0
	if curY+cardH > pageH-bottomM {
		pdf.AddPage()
		curY = 25.0
	}

	gap := 3.5
	cardW := (availW - (float64(numCards-1) * gap)) / float64(numCards)

	for i, card := range m.Cards {
		cX := startX + float64(i)*(cardW+gap)

		// Fondo de la tarjeta
		pdf.SetFillColor(theme.AlternateRowColor.R, theme.AlternateRowColor.G, theme.AlternateRowColor.B)
		pdf.SetDrawColor(theme.BorderColor.R, theme.BorderColor.G, theme.BorderColor.B)
		pdf.SetLineWidth(0.2)
		pdf.Rect(cX, curY, cardW, cardH, "FD")

		// Título de la métrica (arriba, pequeño)
		pdf.SetFont("Arial", "B", 7.0)
		pdf.SetTextColor(theme.MutedTextColor.R, theme.MutedTextColor.G, theme.MutedTextColor.B)
		pdf.Text(cX+3.0, curY+5.5, tr(card.Title))

		// Valor numérico destacado
		pdf.SetFont("Arial", "B", 13.0)
		pdf.SetTextColor(theme.PrimaryColor.R, theme.PrimaryColor.G, theme.PrimaryColor.B)
		pdf.Text(cX+3.0, curY+13.5, tr(card.Value))

		// Cambio / Delta (+15%)
		if card.Change != "" {
			badgeColor := charts.HexColor("#10B981") // Verde
			if !card.IsPositive {
				badgeColor = charts.HexColor("#EF4444") // Rojo
			}
			pdf.SetFont("Arial", "B", 7.5)
			pdf.SetTextColor(badgeColor.R, badgeColor.G, badgeColor.B)
			trChange := tr(card.Change)
			w := pdf.GetStringWidth(trChange)
			pdf.Text(cX+cardW-w-3.0, curY+13.5, trChange)
		}

		// Descripción corta
		if card.Description != "" {
			pdf.SetFont("Arial", "", 6.5)
			pdf.SetTextColor(theme.MutedTextColor.R, theme.MutedTextColor.G, theme.MutedTextColor.B)
			pdf.Text(cX+3.0, curY+19.0, tr(card.Description))
		}
	}

	return curY + cardH + 6.0
}

func renderChartBlock(pdf *fpdf.Fpdf, tr func(string) string, ch *reporter.ChartBlock, startX, curY, availW, pageH, bottomM float64) float64 {
	targetW := ch.Width
	if targetW <= 0 || targetW > availW {
		targetW = availW
	}
	targetH := ch.Height
	if targetH <= 0 {
		targetH = 80.0
	}

	// Comprobar si cabe en la página
	if curY+targetH > pageH-bottomM {
		pdf.AddPage()
		curY = 25.0
	}

	// Centrar gráfico si su ancho es menor que el ancho disponible
	chartX := startX + ((availW - targetW) / 2.0)

	ops, err := charts.Generate(ch.Config)
	if err == nil {
		renderVectorChart(pdf, tr, ops, ch.Config, chartX, curY, targetW, targetH)
	}

	return curY + targetH + 8.0
}

func renderDivider(pdf *fpdf.Fpdf, theme reporter.Theme, startX, curY, availW, pageH, bottomM float64) float64 {
	if curY+5 > pageH-bottomM {
		pdf.AddPage()
		curY = 25.0
	}
	pdf.SetDrawColor(theme.BorderColor.R, theme.BorderColor.G, theme.BorderColor.B)
	pdf.SetLineWidth(0.2)
	pdf.Line(startX, curY+2, startX+availW, curY+2)
	return curY + 6.0
}

func renderImage(pdf *fpdf.Fpdf, img *reporter.ImageBlock, startX, curY, availW, pageH, bottomM float64) float64 {
	if img.Path == "" && len(img.Data) == 0 {
		return curY
	}

	w := img.Width
	if w <= 0 || w > availW {
		w = math.Min(100.0, availW)
	}
	h := img.Height
	if h <= 0 {
		h = 40.0
	}

	if curY+h > pageH-bottomM {
		pdf.AddPage()
		curY = 25.0
	}

	imgX := startX
	if img.Align == reporter.AlignCenter {
		imgX = startX + ((availW - w) / 2.0)
	} else if img.Align == reporter.AlignRight {
		imgX = startX + availW - w
	}

	if img.Path != "" {
		pdf.ImageOptions(img.Path, imgX, curY, w, h, false, fpdf.ImageOptions{ImageType: ""}, 0, "")
	} else if len(img.Data) > 0 {
		reader := bytes.NewReader(img.Data)
		imgName := fmt.Sprintf("mem_img_%d", len(img.Data))
		pdf.RegisterImageOptionsReader(imgName, fpdf.ImageOptions{ImageType: "PNG"}, reader)
		pdf.ImageOptions(imgName, imgX, curY, w, h, false, fpdf.ImageOptions{ImageType: "PNG"}, 0, "")
	}

	return curY + h + 5.0
}
