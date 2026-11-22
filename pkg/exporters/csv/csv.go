package csvexp

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/luisangelf11/frr-go/pkg/reporter"
)

// Config define las opciones para la exportación a formato CSV / TSV.
type Config struct {
	Delimiter         rune   // Delimitador de columnas: ',' o ';' o '\t' (por defecto ',')
	UseCRLF           bool   // Si es true, usa saltos de línea Windows (\r\n); si es false, (\n)
	IncludeBOM        bool   // Si es true, incluye el prefijo UTF-8 BOM (\xEF\xBB\xBF) para Excel
	TableID           string // Si está definido, exporta únicamente la tabla con este ID
	UseExcelSepHeader bool   // Si es true, incluye la directiva "sep=," para forzar a Excel a reconocer las columnas
}

// DefaultConfig retorna una configuración estándar optimizada para compatibilidad internacional.
func DefaultConfig() Config {
	return Config{
		Delimiter:  ',',
		UseCRLF:    true,
		IncludeBOM: true,
	}
}

// Exporter se encarga de convertir las tablas de un Document en formato CSV.
type Exporter struct {
	cfg Config
}

// NewExporter crea una nueva instancia del exportador con la configuración indicada.
func NewExporter(cfg ...Config) *Exporter {
	c := DefaultConfig()
	if len(cfg) > 0 {
		c = cfg[0]
		if c.Delimiter == 0 {
			c.Delimiter = ','
		}
	}
	return &Exporter{cfg: c}
}

// Render exporta las tablas del documento a un archivo CSV en disco.
func (e *Exporter) Render(doc *reporter.Document, outPath string) error {
	data, err := e.RenderToBytes(doc)
	if err != nil {
		return err
	}
	return os.WriteFile(outPath, data, 0644)
}

// RenderToWriter escribe el contenido CSV directamente sobre un io.Writer.
func (e *Exporter) RenderToWriter(doc *reporter.Document, w io.Writer) error {
	data, err := e.RenderToBytes(doc)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

// RenderToBytes genera el contenido CSV en memoria y lo retorna como slice de bytes.
func (e *Exporter) RenderToBytes(doc *reporter.Document) ([]byte, error) {
	var buf bytes.Buffer

	// 1. Prefijo UTF-8 BOM si está habilitado (permite a Excel abrir caracteres como tildes y ñ sin corrupción)
	if e.cfg.IncludeBOM {
		buf.Write([]byte{0xEF, 0xBB, 0xBF})
	}

	// Directiva sep= para forzar a Excel en Windows a separar columnas correctamente
	if e.cfg.UseExcelSepHeader {
		buf.WriteString(fmt.Sprintf("sep=%c\r\n", e.cfg.Delimiter))
	}

	// 2. Extraer tablas del documento
	var tables []*reporter.TableBlock
	for _, blk := range doc.Blocks {
		if blk.Type == reporter.BlockTypeTable && blk.Table != nil {
			if e.cfg.TableID == "" || blk.Table.ID == e.cfg.TableID {
				tables = append(tables, blk.Table)
			}
		}
	}

	if len(tables) == 0 {
		return buf.Bytes(), nil
	}

	writer := csv.NewWriter(&buf)
	writer.Comma = e.cfg.Delimiter
	writer.UseCRLF = e.cfg.UseCRLF

	for tblIdx, tbl := range tables {
		// Si hay múltiples tablas y no se especificó TableID, agregar un separador informativo
		if len(tables) > 1 && tbl.Title != "" {
			if tblIdx > 0 {
				_ = writer.Write([]string{}) // Línea en blanco entre tablas
			}
			_ = writer.Write([]string{"# " + tbl.Title})
		}

		// Escribir cabeceras
		if len(tbl.Headers) > 0 {
			if err := writer.Write(tbl.Headers); err != nil {
				return nil, fmt.Errorf("error escribiendo cabeceras CSV: %w", err)
			}
		}

		// Escribir filas de datos
		for _, row := range tbl.Rows {
			if err := writer.Write(row); err != nil {
				return nil, fmt.Errorf("error escribiendo fila CSV: %w", err)
			}
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("error finalizando escritura CSV: %w", err)
	}

	return buf.Bytes(), nil
}
