package excel

import (
	"fmt"
	"strings"
)

// colIndexToLetters convierte un índice numérico de columna (0 = A, 1 = B, 26 = AA, etc.) a letras Excel.
func colIndexToLetters(idx int) string {
	res := ""
	idx++ // Base 1
	for idx > 0 {
		idx--
		res = string(rune('A'+(idx%26))) + res
		idx /= 26
	}
	return res
}

// buildContentTypesXML genera el archivo [Content_Types].xml para N hojas de cálculo.
func buildContentTypesXML(numSheets int) string {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n")
	sb.WriteString(`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">` + "\n")
	sb.WriteString(`  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>` + "\n")
	sb.WriteString(`  <Default Extension="xml" ContentType="application/xml"/>` + "\n")
	sb.WriteString(`  <Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>` + "\n")
	sb.WriteString(`  <Override PartName="/xl/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.styles+xml"/>` + "\n")

	for i := 1; i <= numSheets; i++ {
		sb.WriteString(fmt.Sprintf(`  <Override PartName="/xl/worksheets/sheet%d.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>`+"\n", i))
	}

	sb.WriteString(`</Types>`)
	return sb.String()
}

// buildRelsXML genera el archivo raíz _rels/.rels.
func buildRelsXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/>
</Relationships>`
}

// buildWorkbookRelsXML genera xl/_rels/workbook.xml.rels conectando las hojas y estilos.
func buildWorkbookRelsXML(numSheets int) string {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n")
	sb.WriteString(`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` + "\n")

	for i := 1; i <= numSheets; i++ {
		sb.WriteString(fmt.Sprintf(`  <Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet%d.xml"/>`+"\n", i, i))
	}

	sb.WriteString(fmt.Sprintf(`  <Relationship Id="rIdStyles" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>` + "\n"))
	sb.WriteString(`</Relationships>`)
	return sb.String()
}

// buildWorkbookXML genera xl/workbook.xml listando los nombres de cada hoja.
func buildWorkbookXML(sheetNames []string) string {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n")
	sb.WriteString(`<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">` + "\n")
	sb.WriteString(`  <sheets>` + "\n")

	for i, name := range sheetNames {
		cleanName := strings.ReplaceAll(name, ":", " ")
		cleanName = strings.ReplaceAll(cleanName, "\\", " ")
		cleanName = strings.ReplaceAll(cleanName, "/", " ")
		cleanName = strings.ReplaceAll(cleanName, "?", " ")
		cleanName = strings.ReplaceAll(cleanName, "*", " ")
		cleanName = strings.ReplaceAll(cleanName, "[", " ")
		cleanName = strings.ReplaceAll(cleanName, "]", " ")
		if len(cleanName) > 31 {
			cleanName = cleanName[:31]
		}
		if cleanName == "" {
			cleanName = fmt.Sprintf("Hoja%d", i+1)
		}
		sb.WriteString(fmt.Sprintf(`    <sheet name="%s" sheetId="%d" r:id="rId%d"/>`+"\n", cleanName, i+1, i+1))
	}

	sb.WriteString(`  </sheets>` + "\n")
	sb.WriteString(`</workbook>`)
	return sb.String()
}

// buildStylesXML genera xl/styles.xml con fuentes, bordes, rellenos y formatos de celda.
func buildStylesXML(headerHex, altRowHex, borderHex string) string {
	headerHex = strings.TrimPrefix(headerHex, "#")
	if len(headerHex) == 6 {
		headerHex = "FF" + strings.ToUpper(headerHex)
	} else {
		headerHex = "FF1E3A8A"
	}

	altRowHex = strings.TrimPrefix(altRowHex, "#")
	if len(altRowHex) == 6 {
		altRowHex = "FF" + strings.ToUpper(altRowHex)
	} else {
		altRowHex = "FFF8FAFC"
	}

	borderHex = strings.TrimPrefix(borderHex, "#")
	if len(borderHex) == 6 {
		borderHex = "FF" + strings.ToUpper(borderHex)
	} else {
		borderHex = "FFD1D5DB"
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">
  <fonts count="2">
    <font>
      <sz val="10.5"/>
      <color rgb="FF1F2937"/>
      <name val="Calibri"/>
    </font>
    <font>
      <b/>
      <sz val="11"/>
      <color rgb="FFFFFFFF"/>
      <name val="Calibri"/>
    </font>
  </fonts>
  <fills count="4">
    <fill>
      <patternFill patternType="none"/>
    </fill>
    <fill>
      <patternFill patternType="gray125"/>
    </fill>
    <fill>
      <patternFill patternType="solid">
        <fgColor rgb="%s"/>
      </patternFill>
    </fill>
    <fill>
      <patternFill patternType="solid">
        <fgColor rgb="%s"/>
      </patternFill>
    </fill>
  </fills>
  <borders count="2">
    <border>
      <left/><right/><top/><bottom/><diagonal/>
    </border>
    <border>
      <left style="thin"><color rgb="%s"/></left>
      <right style="thin"><color rgb="%s"/></right>
      <top style="thin"><color rgb="%s"/></top>
      <bottom style="thin"><color rgb="%s"/></bottom>
      <diagonal/>
    </border>
  </borders>
  <cellStyleXfs count="1">
    <xf numFmtId="0" fontId="0" fillId="0" borderId="0"/>
  </cellStyleXfs>
  <cellXfs count="8">
    <!-- 0: Default general -->
    <xf numFmtId="0" fontId="0" fillId="0" borderId="0" xfId="0"/>
    <!-- 1: Header (Negrita, Fondo primario, Borde, Centrado) -->
    <xf numFmtId="0" fontId="1" fillId="2" borderId="1" xfId="0" applyFont="1" applyFill="1" applyBorder="1" applyAlignment="1">
      <alignment horizontal="center" vertical="center" wrapText="1"/>
    </xf>
    <!-- 2: Fila normal (Texto izquierda, Borde fino) -->
    <xf numFmtId="0" fontId="0" fillId="0" borderId="1" xfId="0" applyFont="1" applyBorder="1" applyAlignment="1">
      <alignment horizontal="left" vertical="center"/>
    </xf>
    <!-- 3: Fila Zebra (Texto izquierda, Fondo suave, Borde fino) -->
    <xf numFmtId="0" fontId="0" fillId="3" borderId="1" xfId="0" applyFont="1" applyFill="1" applyBorder="1" applyAlignment="1">
      <alignment horizontal="left" vertical="center"/>
    </xf>
    <!-- 4: Fila normal número (Derecha, Borde fino) -->
    <xf numFmtId="0" fontId="0" fillId="0" borderId="1" xfId="0" applyFont="1" applyBorder="1" applyAlignment="1">
      <alignment horizontal="right" vertical="center"/>
    </xf>
    <!-- 5: Fila Zebra número (Derecha, Fondo suave, Borde fino) -->
    <xf numFmtId="0" fontId="0" fillId="3" borderId="1" xfId="0" applyFont="1" applyFill="1" applyBorder="1" applyAlignment="1">
      <alignment horizontal="right" vertical="center"/>
    </xf>
    <!-- 6: Fila normal centrado (Centro, Borde fino) -->
    <xf numFmtId="0" fontId="0" fillId="0" borderId="1" xfId="0" applyFont="1" applyBorder="1" applyAlignment="1">
      <alignment horizontal="center" vertical="center"/>
    </xf>
    <!-- 7: Fila Zebra centrado (Centro, Fondo suave, Borde fino) -->
    <xf numFmtId="0" fontId="0" fillId="3" borderId="1" xfId="0" applyFont="1" applyFill="1" applyBorder="1" applyAlignment="1">
      <alignment horizontal="center" vertical="center"/>
    </xf>
  </cellXfs>
</styleSheet>`, headerHex, altRowHex, borderHex, borderHex, borderHex, borderHex)
}
