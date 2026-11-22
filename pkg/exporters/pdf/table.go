package pdf

import (
	"math"
	"unicode/utf8"

	"github.com/go-pdf/fpdf"
	"github.com/luisangelf11/reporteria-go/pkg/reporter"
)

// renderTable renderiza una tabla con soporte de auto-wrap, anchos calculados, zebra striping y repetición de cabeceras.
func renderTable(pdf *fpdf.Fpdf, tr func(string) string, t *reporter.TableBlock, theme reporter.Theme, startX, startY, availWidth, pageH, bottomMargin float64) float64 {
	numCols := len(t.Headers)
	if numCols == 0 && len(t.Rows) > 0 {
		numCols = len(t.Rows[0])
	}
	if numCols == 0 {
		return startY
	}

	// 1. Calcular anchos de columna
	colWidths := calculateColumnWidths(pdf, t, availWidth, numCols)

	// Alturas base
	lineHeight := 4.5
	cellPadding := 2.0
	curY := startY

	// Función auxiliar para dibujar la cabecera
	drawHeader := func(y float64) float64 {
		headerH := 8.0
		pdf.SetFillColor(theme.HeaderBgColor.R, theme.HeaderBgColor.G, theme.HeaderBgColor.B)
		pdf.SetTextColor(theme.HeaderTextColor.R, theme.HeaderTextColor.G, theme.HeaderTextColor.B)
		pdf.SetDrawColor(theme.BorderColor.R, theme.BorderColor.G, theme.BorderColor.B)
		pdf.SetLineWidth(0.2)
		pdf.SetFont("Arial", "B", 8.5)

		curX := startX
		for colIdx := 0; colIdx < numCols; colIdx++ {
			rawText := ""
			if colIdx < len(t.Headers) {
				rawText = t.Headers[colIdx]
			}
			headerText := tr(rawText)
			w := colWidths[colIdx]

			pdf.Rect(curX, y, w, headerH, "FD")

			// Texto alineado
			align := "L"
			textX := curX + cellPadding
			textW := pdf.GetStringWidth(headerText)

			if colIdx < len(t.ColumnAlignments) {
				if t.ColumnAlignments[colIdx] == reporter.AlignCenter {
					align = "C"
					textX = curX + (w / 2.0) - (textW / 2.0)
				} else if t.ColumnAlignments[colIdx] == reporter.AlignRight {
					align = "R"
					textX = curX + w - cellPadding - textW
				}
			}
			_ = align

			pdf.Text(textX, y+headerH-2.5, headerText)
			curX += w
		}
		return y + headerH
	}

	// Título opcional de la tabla
	if t.Title != "" {
		if curY+12 > pageH-bottomMargin {
			pdf.AddPage()
			curY = 25.0
		}
		pdf.SetFont("Arial", "B", 10.5)
		pdf.SetTextColor(theme.PrimaryColor.R, theme.PrimaryColor.G, theme.PrimaryColor.B)
		pdf.Text(startX, curY+4, tr(t.Title))
		curY += 7.0
	}

	// Comprobar si cabe la cabecera
	if curY+10 > pageH-bottomMargin {
		pdf.AddPage()
		curY = 25.0
	}

	curY = drawHeader(curY)

	// 2. Dibujar Filas de Datos
	pdf.SetFont("Arial", "", 8.0)
	pdf.SetLineWidth(0.15)
	pdf.SetDrawColor(theme.BorderColor.R, theme.BorderColor.G, theme.BorderColor.B)

	for rowIdx, row := range t.Rows {
		// Calcular la altura requerida para esta fila (según el número de líneas por celda)
		maxLines := 1
		cellLines := make([][][]byte, numCols)

		for colIdx := 0; colIdx < numCols; colIdx++ {
			rawVal := ""
			if colIdx < len(row) {
				rawVal = row[colIdx]
			}
			cellVal := tr(rawVal)
			w := colWidths[colIdx] - (cellPadding * 2)
			lines := pdf.SplitLines([]byte(cellVal), w)
			if len(lines) > maxLines {
				maxLines = len(lines)
			}
			cellLines[colIdx] = lines
		}

		rowH := math.Max(6.5, float64(maxLines)*lineHeight+3.0)

		// Paginación inteligente: si no cabe, salto de página y repetir cabecera
		if curY+rowH > pageH-bottomMargin {
			pdf.AddPage()
			curY = 25.0
			curY = drawHeader(curY)
			pdf.SetFont("Arial", "", 8.0)
			pdf.SetLineWidth(0.15)
			pdf.SetDrawColor(theme.BorderColor.R, theme.BorderColor.G, theme.BorderColor.B)
		}

		// Color de fondo (Zebra Striping)
		isAlt := (rowIdx%2 == 1) && t.ZebraStriping
		if isAlt {
			pdf.SetFillColor(theme.AlternateRowColor.R, theme.AlternateRowColor.G, theme.AlternateRowColor.B)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}
		pdf.SetTextColor(theme.TextColor.R, theme.TextColor.G, theme.TextColor.B)

		curX := startX
		for colIdx := 0; colIdx < numCols; colIdx++ {
			w := colWidths[colIdx]
			pdf.Rect(curX, curY, w, rowH, "FD")

			lines := cellLines[colIdx]
			textStartY := curY + cellPadding + 2.5

			for lineIdx, line := range lines {
				lineStr := string(line)
				lineW := pdf.GetStringWidth(lineStr)
				textX := curX + cellPadding

				if colIdx < len(t.ColumnAlignments) {
					if t.ColumnAlignments[colIdx] == reporter.AlignCenter {
						textX = curX + (w / 2.0) - (lineW / 2.0)
					} else if t.ColumnAlignments[colIdx] == reporter.AlignRight {
						textX = curX + w - cellPadding - lineW
					}
				}

				pdf.Text(textX, textStartY+float64(lineIdx)*lineHeight, lineStr)
			}

			curX += w
		}

		curY += rowH
	}

	return curY
}

// calculateColumnWidths deduce anchos de columna basados en contenido o configuración.
func calculateColumnWidths(pdf *fpdf.Fpdf, t *reporter.TableBlock, availWidth float64, numCols int) []float64 {
	// Si el usuario especificó anchos explícitos
	if len(t.ColumnWidths) == numCols {
		totalSet := 0.0
		for _, w := range t.ColumnWidths {
			totalSet += w
		}
		if totalSet > 0 {
			res := make([]float64, numCols)
			for i, w := range t.ColumnWidths {
				res[i] = (w / totalSet) * availWidth
			}
			return res
		}
	}

	// Auto-cálculo basado en cantidad real de caracteres (runas UTF-8)
	maxLengths := make([]float64, numCols)
	for i := 0; i < numCols; i++ {
		headerLen := 5.0
		if i < len(t.Headers) {
			headerLen = math.Max(headerLen, float64(utf8.RuneCountInString(t.Headers[i])))
		}
		maxLengths[i] = headerLen
	}

	evalRows := len(t.Rows)
	if evalRows > 25 {
		evalRows = 25
	}

	for r := 0; r < evalRows; r++ {
		row := t.Rows[r]
		for c := 0; c < numCols; c++ {
			if c < len(row) {
				l := float64(utf8.RuneCountInString(row[c]))
				if l > maxLengths[c] {
					maxLengths[c] = math.Min(35.0, l)
				}
			}
		}
	}

	totalLength := 0.0
	for _, l := range maxLengths {
		totalLength += l
	}

	res := make([]float64, numCols)
	minColWidth := 16.0
	for i := 0; i < numCols; i++ {
		calculated := (maxLengths[i] / totalLength) * availWidth
		if calculated < minColWidth {
			calculated = minColWidth
		}
		res[i] = calculated
	}

	sum := 0.0
	for _, w := range res {
		sum += w
	}
	for i := range res {
		res[i] = (res[i] / sum) * availWidth
	}

	return res
}
