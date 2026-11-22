package excel

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/luisangelf11/reporteria-go/pkg/reporter"
)

// Config define las opciones de formato para la exportación a Excel.
type Config struct {
	FreezeHeader  bool // Si es true, inmoviliza la fila de cabecera
	AutoWidth     bool // Si es true, calcula el ancho óptimo de columnas
	ZebraStriping bool // Si es true, aplica alternancia de colores a las filas
}

// DefaultConfig retorna la configuración recomendada para Excel.
func DefaultConfig() Config {
	return Config{
		FreezeHeader:  true,
		AutoWidth:     true,
		ZebraStriping: true,
	}
}

// Exporter se encarga de convertir un Document en un archivo o buffer .xlsx estándar.
type Exporter struct {
	cfg Config
}

// NewExporter crea una nueva instancia del exportador de Excel.
func NewExporter(cfg ...Config) *Exporter {
	c := DefaultConfig()
	if len(cfg) > 0 {
		c = cfg[0]
	}
	return &Exporter{cfg: c}
}

// Render exporta el reporte a un archivo .xlsx en disco.
func (e *Exporter) Render(doc *reporter.Document, outPath string) error {
	data, err := e.RenderToBytes(doc)
	if err != nil {
		return err
	}
	return os.WriteFile(outPath, data, 0644)
}

// RenderToWriter escribe los bytes del .xlsx directamente en un io.Writer.
func (e *Exporter) RenderToWriter(doc *reporter.Document, w io.Writer) error {
	data, err := e.RenderToBytes(doc)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

// RenderToBytes genera el archivo .xlsx completo en memoria con cero dependencias externas y tabla de strings compartidos.
func (e *Exporter) RenderToBytes(doc *reporter.Document) ([]byte, error) {
	// 1. Extraer o sintetizar tablas a exportar
	type sheetInfo struct {
		name  string
		table *reporter.TableBlock
	}

	var sheets []sheetInfo

	for _, blk := range doc.Blocks {
		if blk.Type == reporter.BlockTypeTable && blk.Table != nil {
			name := blk.Table.Title
			if name == "" {
				name = fmt.Sprintf("Tabla %d", len(sheets)+1)
			}
			sheets = append(sheets, sheetInfo{name: name, table: blk.Table})
		}
	}

	// Si no hay tablas, crear una hoja de resumen con los metadatos y KPIs del reporte
	if len(sheets) == 0 {
		summaryTbl := reporter.NewTable().
			SetTitle(doc.Title).
			SetHeaders("Campo", "Valor").
			AddRow("Título", doc.Title).
			AddRow("Subtítulo", doc.Subtitle).
			AddRow("Autor", doc.Author).
			AddRow("Fecha", doc.Date).
			Build()
		sheets = append(sheets, sheetInfo{name: "Resumen", table: &summaryTbl})
	}

	numSheets := len(sheets)
	sheetNames := make([]string, numSheets)
	for i, s := range sheets {
		sheetNames[i] = s.name
	}

	// 2. Construir el contenido de las hojas y recolectar cadenas en SharedStringsTable (SST)
	sst := newSST()
	sheetContents := make([]string, numSheets)
	for i, s := range sheets {
		sheetContents[i] = buildSheetXML(s.table, sst, e.cfg.FreezeHeader, e.cfg.AutoWidth, e.cfg.ZebraStriping)
	}

	hasSST := len(sst.list) > 0

	// 3. Empaquetar el archivo ZIP OpenXML (.xlsx)
	var zipBuf bytes.Buffer
	zw := zip.NewWriter(&zipBuf)

	// Helper para agregar archivos al ZIP
	addFile := func(name, content string) error {
		w, err := zw.Create(name)
		if err != nil {
			return err
		}
		_, err = io.WriteString(w, content)
		return err
	}

	// [Content_Types].xml
	if err := addFile("[Content_Types].xml", buildContentTypesXML(numSheets, hasSST)); err != nil {
		return nil, fmt.Errorf("error creando [Content_Types].xml: %w", err)
	}

	// _rels/.rels
	if err := addFile("_rels/.rels", buildRelsXML()); err != nil {
		return nil, fmt.Errorf("error creando _rels/.rels: %w", err)
	}

	// xl/_rels/workbook.xml.rels
	if err := addFile("xl/_rels/workbook.xml.rels", buildWorkbookRelsXML(numSheets, hasSST)); err != nil {
		return nil, fmt.Errorf("error creando workbook.xml.rels: %w", err)
	}

	// xl/workbook.xml
	if err := addFile("xl/workbook.xml", buildWorkbookXML(sheetNames)); err != nil {
		return nil, fmt.Errorf("error creando workbook.xml: %w", err)
	}

	// xl/styles.xml
	theme := doc.Theme
	stylesXML := buildStylesXML(theme.HeaderBgColor.Hex(), theme.AlternateRowColor.Hex(), theme.BorderColor.Hex())
	if err := addFile("xl/styles.xml", stylesXML); err != nil {
		return nil, fmt.Errorf("error creando styles.xml: %w", err)
	}

	// xl/sharedStrings.xml (si hay cadenas)
	if hasSST {
		sstXML := buildSharedStringsXML(sst.list, sst.total)
		if err := addFile("xl/sharedStrings.xml", sstXML); err != nil {
			return nil, fmt.Errorf("error creando sharedStrings.xml: %w", err)
		}
	}

	// xl/worksheets/sheet{1..N}.xml
	for i := range sheets {
		sheetPath := fmt.Sprintf("xl/worksheets/sheet%d.xml", i+1)
		if err := addFile(sheetPath, sheetContents[i]); err != nil {
			return nil, fmt.Errorf("error creando %s: %w", sheetPath, err)
		}
	}

	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("error cerrando contenedor .xlsx: %w", err)
	}

	return zipBuf.Bytes(), nil
}
