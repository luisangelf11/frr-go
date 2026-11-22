package csvexp_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	csvexp "github.com/luisangelf11/frr-go/pkg/exporters/csv"
	"github.com/luisangelf11/frr-go/pkg/reporter"
)

func buildSampleReportWithTable() *reporter.Document {
	tbl := reporter.NewTable().
		SetTitle("Detalle de Ventas e Inventario").
		SetHeaders("Código", "Descripción del Artículo", "Categoría", "Stock", "Precio Unitario (€)", "Estado").
		AddRow("ART-001", "Café Tostado Especial 500g", "Bebidas", "150", "8.50 €", "Disponible").
		AddRow("ART-002", "Té Verde Matcha Orgánico", "Bebidas", "80", "14.20 €", "Disponible").
		AddRow("ART-003", "Chocolate Amargo 85% Cacao", "Dulces", "200", "3.95 €", "Disponible").
		AddRow("ART-004", "Galletas de Avena y Miel", "Panadería", "45", "2.80 €", "Poco Stock").
		AddRow("ART-005", "Aceite de Oliva Virgen Extra", "Aceites", "120", "11.50 €", "Disponible")

	return reporter.New("Reporte de Inventario").
		AddTable(tbl.Build()).
		Build()
}

func TestCSVExportCommaAndBOM(t *testing.T) {
	doc := buildSampleReportWithTable()
	exp := csvexp.NewExporter(csvexp.Config{
		Delimiter:  ',',
		UseCRLF:    true,
		IncludeBOM: true,
	})

	data, err := exp.RenderToBytes(doc)
	if err != nil {
		t.Fatalf("Error exportando a CSV: %v", err)
	}

	// Verificar BOM UTF-8
	if !bytes.HasPrefix(data, []byte{0xEF, 0xBB, 0xBF}) {
		t.Errorf("El archivo CSV no contiene el prefijo BOM UTF-8")
	}

	content := string(data[3:]) // contenido sin BOM
	lines := strings.Split(content, "\r\n")

	if len(lines) < 6 {
		t.Errorf("Se esperaban al menos 6 líneas (cabecera + 5 filas), se obtuvieron %d", len(lines))
	}

	// Verificar que contiene tildes y caracteres especiales
	if !strings.Contains(content, "Café Tostado Especial 500g") {
		t.Errorf("No se preservaron los acentos en el contenido CSV")
	}
	if !strings.Contains(content, "8.50 €") {
		t.Errorf("No se preservó el símbolo de moneda €")
	}

	// Guardar en output para inspección
	outDir := filepath.Join("..", "..", "..", "output")
	_ = os.MkdirAll(outDir, 0755)
	outPath := filepath.Join(outDir, "executive_report.csv")
	if err := exp.Render(doc, outPath); err != nil {
		t.Fatalf("Error guardando CSV en disco: %v", err)
	}
}

func TestCSVExportSemicolon(t *testing.T) {
	doc := buildSampleReportWithTable()
	exp := csvexp.NewExporter(csvexp.Config{
		Delimiter:  ';',
		UseCRLF:    true,
		IncludeBOM: false,
	})

	data, err := exp.RenderToBytes(doc)
	if err != nil {
		t.Fatalf("Error exportando CSV con punto y coma: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "Código;Descripción del Artículo;Categoría") {
		t.Errorf("El delimitador punto y coma no se aplicó correctamente: %s", content)
	}
}

func TestMultiTableCSVExport(t *testing.T) {
	tbl1 := reporter.NewTable().
		SetID("tabla-1").
		SetTitle("Resumen Q1").
		SetHeaders("Trimestre", "Meta").
		AddRow("Q1", "100K").
		Build()

	tbl2 := reporter.NewTable().
		SetID("tabla-2").
		SetTitle("Resumen Q2").
		SetHeaders("Trimestre", "Meta").
		AddRow("Q2", "120K").
		Build()

	doc := reporter.New("Doble Tabla").
		AddTable(tbl1).
		AddTable(tbl2).
		Build()

	// 1. Exportar ambas tablas
	expAll := csvexp.NewExporter(csvexp.DefaultConfig())
	dataAll, err := expAll.RenderToBytes(doc)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if !strings.Contains(string(dataAll), "# Resumen Q1") || !strings.Contains(string(dataAll), "# Resumen Q2") {
		t.Errorf("El CSV multitabla debe incluir títulos de sección")
	}

	// 2. Exportar solo tabla-2
	expSingle := csvexp.NewExporter(csvexp.Config{
		Delimiter: ',',
		TableID:   "tabla-2",
	})
	dataSingle, err := expSingle.RenderToBytes(doc)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if strings.Contains(string(dataSingle), "Resumen Q1") {
		t.Errorf("No debería contener la tabla 1 cuando se filtra por tabla-2")
	}
	if !strings.Contains(string(dataSingle), "Q2") {
		t.Errorf("Debe contener datos de la tabla-2")
	}
}

func BenchmarkCSVExport1000Rows(b *testing.B) {
	tbl := reporter.NewTable().
		SetHeaders("ID", "Nombre", "Email", "Puntaje", "Activo")

	for i := 1; i <= 1000; i++ {
		tbl.AddRow(
			fmt.Sprintf("USR-%05d", i),
			fmt.Sprintf("Usuario Nombre %d", i),
			fmt.Sprintf("usuario%d@empresa.com", i),
			fmt.Sprintf("%d", (i*37)%100),
			"true",
		)
	}

	doc := reporter.New("Benchmark CSV").
		AddTable(tbl.Build()).
		Build()

	exp := csvexp.NewExporter(csvexp.DefaultConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := exp.RenderToBytes(doc)
		if err != nil {
			b.Fatalf("Error en benchmark CSV: %v", err)
		}
	}
}
