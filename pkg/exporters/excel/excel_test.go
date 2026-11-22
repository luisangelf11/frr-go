package excel_test

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/luisangelf11/reporteria-go/pkg/exporters/excel"
	"github.com/luisangelf11/reporteria-go/pkg/reporter"
)

func buildSampleExcelReport() *reporter.Document {
	tbl1 := reporter.NewTable().
		SetTitle("Ventas Consolidadas").
		SetHeaders("ID", "Fecha", "Cliente", "Región", "Unidades", "Monto", "Estado").
		AddRow("TRX-001", "2026-03-01", "Empresa Alpha S.A.", "Norte", "150", "12500.50", "Completado").
		AddRow("TRX-002", "2026-03-02", "Logística Beta Ltda.", "Sur", "45", "3800.00", "Completado").
		AddRow("TRX-003", "2026-03-03", "Servicios Gamma Corp.", "Este", "210", "18900.25", "Pendiente").
		AddRow("TRX-004", "2026-03-04", "Soluciones Delta S.A.", "Oeste", "85", "7200.00", "Completado").
		AddRow("TRX-005", "2026-03-05", "Industrias Epsilon", "Centro", "300", "28500.00", "Completado")

	tbl2 := reporter.NewTable().
		SetTitle("Metas y Presupuesto").
		SetHeaders("Región", "Meta Q1", "Logrado Q1", "Cumplimiento").
		AddRow("Norte", "100000", "125000", "125%").
		AddRow("Sur", "80000", "92000", "115%").
		AddRow("Este", "120000", "110000", "91.6%").
		AddRow("Oeste", "90000", "105000", "116.7%").
		AddRow("Centro", "150000", "165000", "110%")

	return reporter.New("Reporte Financiero Excel").
		SetAuthor("Equipo de Finanzas").
		SetDate("2026-03-31").
		SetTheme(reporter.ThemeCorporateBlue).
		AddTable(tbl1.Build()).
		AddTable(tbl2.Build()).
		Build()
}

func TestExcelGeneration(t *testing.T) {
	doc := buildSampleExcelReport()
	exp := excel.NewExporter(excel.DefaultConfig())

	data, err := exp.RenderToBytes(doc)
	if err != nil {
		t.Fatalf("Error exportando a Excel: %v", err)
	}

	if len(data) == 0 {
		t.Fatalf("El buffer Excel generado está vacío")
	}

	// 1. Validar que es un archivo ZIP válido
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("El archivo generado no es un ZIP válido (especificación OpenXML): %v", err)
	}

	// 2. Verificar archivos obligatorios del estándar OpenXML
	expectedFiles := map[string]bool{
		"[Content_Types].xml":          false,
		"_rels/.rels":                  false,
		"xl/_rels/workbook.xml.rels":   false,
		"xl/workbook.xml":              false,
		"xl/styles.xml":                false,
		"xl/sharedStrings.xml":         false,
		"xl/worksheets/sheet1.xml":     false,
		"xl/worksheets/sheet2.xml":     false,
	}

	for _, f := range zipReader.File {
		if _, exists := expectedFiles[f.Name]; exists {
			expectedFiles[f.Name] = true
		}
	}

	for fileName, found := range expectedFiles {
		if !found {
			t.Errorf("Falta archivo obligatorio de OpenXML en el .xlsx: %s", fileName)
		}
	}

	// 3. Guardar archivo de muestra para inspección en Excel / LibreOffice
	outDir := filepath.Join("..", "..", "..", "output")
	_ = os.MkdirAll(outDir, 0755)
	outPath := filepath.Join(outDir, "executive_report.xlsx")
	if err := exp.Render(doc, outPath); err != nil {
		t.Fatalf("Error guardando .xlsx en disco: %v", err)
	}

	t.Logf("Excel generado exitosamente en: %s (tamaño: %d bytes)", outPath, len(data))
}

func BenchmarkExcelGeneration1000Rows(b *testing.B) {
	tbl := reporter.NewTable().
		SetTitle("Grandes Volúmenes").
		SetHeaders("ID", "Código", "Descripción", "Cantidad", "Precio", "Total")

	for i := 1; i <= 1000; i++ {
		tbl.AddRow(
			fmt.Sprintf("%d", i),
			fmt.Sprintf("SKU-%05d", i),
			fmt.Sprintf("Producto Comercial %d de Alta Demanda", i),
			fmt.Sprintf("%d", (i*7)%50+1),
			fmt.Sprintf("%.2f", float64((i*13)%100)+0.99),
			fmt.Sprintf("%.2f", float64(((i*7)%50+1))*float64((i*13)%100)+0.99),
		)
	}

	doc := reporter.New("Benchmark Excel").
		AddTable(tbl.Build()).
		Build()

	exp := excel.NewExporter(excel.DefaultConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := exp.RenderToBytes(doc)
		if err != nil {
			b.Fatalf("Error en benchmark Excel: %v", err)
		}
	}
}
