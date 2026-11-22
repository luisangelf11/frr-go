package excel

import (
	"fmt"
	"html"
	"math"
	"strconv"
	"strings"

	"github.com/luisangelf11/reporteria-go/pkg/reporter"
)

// buildSheetXML construye el XML de una hoja de cálculo individual para una tabla de datos.
func buildSheetXML(tbl *reporter.TableBlock, freezeHeader, autoWidth, zebraStriping bool) string {
	numCols := len(tbl.Headers)
	if numCols == 0 && len(tbl.Rows) > 0 {
		numCols = len(tbl.Rows[0])
	}
	if numCols == 0 {
		return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">
  <sheetData/>
</worksheet>`
	}

	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n")
	sb.WriteString(`<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">` + "\n")

	// 1. Inmovilización de paneles (Fila 1 fija)
	if freezeHeader {
		sb.WriteString(`  <sheetViews>` + "\n")
		sb.WriteString(`    <sheetView tabSelected="1" workbookViewId="0">` + "\n")
		sb.WriteString(`      <pane ySplit="1" topLeftCell="A2" activePane="bottomLeft" state="frozen"/>` + "\n")
		sb.WriteString(`    </sheetView>` + "\n")
		sb.WriteString(`  </sheetViews>` + "\n")
	}

	// 2. Cálculo de ancho de columnas (Auto-width)
	colWidths := make([]float64, numCols)
	for i := 0; i < numCols; i++ {
		w := 12.0
		if i < len(tbl.Headers) {
			w = math.Max(w, float64(len(tbl.Headers[i]))+4.0)
		}
		colWidths[i] = w
	}

	// Muestrear filas para ancho óptimo
	sampleRows := len(tbl.Rows)
	if sampleRows > 100 {
		sampleRows = 100
	}
	for r := 0; r < sampleRows; r++ {
		row := tbl.Rows[r]
		for c := 0; c < numCols; c++ {
			if c < len(row) {
				l := float64(len(row[c])) + 3.0
				if l > colWidths[c] {
					colWidths[c] = math.Min(50.0, l) // máximo 50 caracteres de ancho
				}
			}
		}
	}

	if autoWidth {
		sb.WriteString(`  <cols>` + "\n")
		for i, w := range colWidths {
			sb.WriteString(fmt.Sprintf(`    <col min="%d" max="%d" width="%.2f" customWidth="1"/>`+"\n", i+1, i+1, w))
		}
		sb.WriteString(`  </cols>` + "\n")
	}

	// 3. Datos de la hoja (sheetData)
	sb.WriteString(`  <sheetData>` + "\n")

	// Fila 1: Cabeceras
	sb.WriteString(`    <row r="1" ht="26" customHeight="1">` + "\n")
	for colIdx := 0; colIdx < numCols; colIdx++ {
		headerText := ""
		if colIdx < len(tbl.Headers) {
			headerText = tbl.Headers[colIdx]
		}
		cellRef := fmt.Sprintf("%s1", colIndexToLetters(colIdx))
		escaped := html.EscapeString(headerText)
		sb.WriteString(fmt.Sprintf(`      <c r="%s" s="1" t="inlineStr"><is><t>%s</t></is></c>`+"\n", cellRef, escaped))
	}
	sb.WriteString(`    </row>` + "\n")

	// Filas 2..N: Datos
	for rowIdx, row := range tbl.Rows {
		rowNum := rowIdx + 2
		isZebra := (rowIdx%2 == 1) && zebraStriping

		sb.WriteString(fmt.Sprintf(`    <row r="%d" ht="20" customHeight="1">`+"\n", rowNum))

		for colIdx := 0; colIdx < numCols; colIdx++ {
			cellVal := ""
			if colIdx < len(row) {
				cellVal = row[colIdx]
			}

			cellRef := fmt.Sprintf("%s%d", colIndexToLetters(colIdx), rowNum)

			// Determinar si es número puro
			numVal, isNum := tryParseNumber(cellVal)

			// Determinar alineación
			align := reporter.AlignLeft
			if colIdx < len(tbl.ColumnAlignments) {
				align = tbl.ColumnAlignments[colIdx]
			}

			// Seleccionar estilo de celda (0..7 según styles.xml)
			styleID := 2 // Normal texto
			if isZebra {
				styleID = 3 // Zebra texto
			}

			if isNum || align == reporter.AlignRight {
				styleID = 4 // Normal número
				if isZebra {
					styleID = 5 // Zebra número
				}
			} else if align == reporter.AlignCenter {
				styleID = 6 // Normal centro
				if isZebra {
					styleID = 7 // Zebra centro
				}
			}

			if isNum {
				sb.WriteString(fmt.Sprintf(`      <c r="%s" s="%d"><v>%s</v></c>`+"\n", cellRef, styleID, numVal))
			} else {
				escaped := html.EscapeString(cellVal)
				sb.WriteString(fmt.Sprintf(`      <c r="%s" s="%d" t="inlineStr"><is><t>%s</t></is></c>`+"\n", cellRef, styleID, escaped))
			}
		}

		sb.WriteString(`    </row>` + "\n")
	}

	sb.WriteString(`  </sheetData>` + "\n")

	// 4. AutoFiltro en la tabla para Excel
	totalRows := len(tbl.Rows) + 1
	if totalRows > 1 && numCols > 0 {
		lastCol := colIndexToLetters(numCols - 1)
		sb.WriteString(fmt.Sprintf(`  <autoFilter ref="A1:%s%d"/>`+"\n", lastCol, totalRows))
	}

	sb.WriteString(`</worksheet>`)
	return sb.String()
}

func tryParseNumber(s string) (string, bool) {
	clean := strings.TrimSpace(s)
	if clean == "" {
		return "", false
	}
	// Si contiene letras o símbolos como $, €, %, no es número puro de celda
	if _, err := strconv.ParseInt(clean, 10, 64); err == nil {
		return clean, true
	}
	if _, err := strconv.ParseFloat(clean, 64); err == nil {
		return clean, true
	}
	return "", false
}
