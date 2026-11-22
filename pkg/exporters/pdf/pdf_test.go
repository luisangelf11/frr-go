package pdf_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/luisangelf11/frr-go/pkg/charts"
	"github.com/luisangelf11/frr-go/pkg/exporters/pdf"
	"github.com/luisangelf11/frr-go/pkg/reporter"
)

func buildSampleReport() *reporter.Document {
	// 1. Gráficos de muestra
	barChart := charts.NewBarChart("Ingresos por Región").
		SetSubtitle("Comparativo 2025 vs 2026 en miles de USD").
		SetCategories("Norte", "Sur", "Este", "Oeste", "Centro").
		AddSeries("2025", []float64{120, 150, 95, 210, 160}).
		AddSeries("2026", []float64{145, 180, 110, 245, 195}).
		SetShowLegend(true).
		SetShowValues(true).
		SetValueFormat("$", "K").
		Config()

	donutChart := charts.NewDonutChart("Distribución de Clientes por Segmento").
		SetData(
			[]string{"Corporativo", "Pymes", "Gobierno", "Educación"},
			[]float64{450, 320, 180, 90},
		).
		SetValueFormat("", " clientes").
		Config()

	// 2. Tabla con 35 filas para forzar salto de página y probar repetición de cabeceras
	tbl := reporter.NewTable().
		SetTitle("Detalle de Transacciones Operativas").
		SetHeaders("ID", "Fecha", "Cliente", "Región", "Monto ($)", "Estado").
		SetAlignments(
			reporter.AlignCenter,
			reporter.AlignCenter,
			reporter.AlignLeft,
			reporter.AlignCenter,
			reporter.AlignRight,
			reporter.AlignCenter,
		).
		SetZebraStriping(true)

	regiones := []string{"Norte", "Sur", "Este", "Oeste", "Centro"}
	estados := []string{"Completado", "Pendiente", "Procesando"}

	for i := 1; i <= 35; i++ {
		tbl.AddRow(
			fmt.Sprintf("TRX-%04d", i),
			fmt.Sprintf("2026-03-%02d", (i%28)+1),
			fmt.Sprintf("Empresa Alpha %d S.A.", i),
			regiones[i%len(regiones)],
			fmt.Sprintf("$%d,%03d.00", (i*3)%10+1, (i*137)%900+100),
			estados[i%len(estados)],
		)
	}

	// 3. Documento del Reporte
	doc := reporter.New("Reporte Ejecutivo de Rendimiento Q1 2026").
		SetSubtitle("Consolidado Integral de Ventas, Clientes y Operaciones").
		SetAuthor("Equipo de Inteligencia de Negocios").
		SetDate("2026-03-31").
		SetTheme(reporter.ThemeCorporateBlue).
		// Tarjetas de Métricas (KPIs)
		AddMetricCards(
			reporter.MetricCard{
				Title:       "INGRESOS TOTALES",
				Value:       "$875,000",
				Change:      "+18.4%",
				IsPositive:  true,
				Description: "vs. Q1 2025",
			},
			reporter.MetricCard{
				Title:       "NUEVOS CLIENTES",
				Value:       "1,040",
				Change:      "+24.1%",
				IsPositive:  true,
				Description: "Objetivo cumplido",
			},
			reporter.MetricCard{
				Title:       "MARGEN OPERATIVO",
				Value:       "32.8%",
				Change:      "-1.2%",
				IsPositive:  false,
				Description: "Dentro de tolerancia",
			},
			reporter.MetricCard{
				Title:       "NPS PROMEDIO",
				Value:       "74/100",
				Change:      "+5 pts",
				IsPositive:  true,
				Description: "Excelente satisfacción",
			},
		).
		AddHeading("1. Resumen Estratégico de Ventas", reporter.H1).
		AddParagraph("Durante el primer trimestre de 2026, la organización experimentó un crecimiento sostenido en todas las regiones comerciales, impulsado principalmente por la adopción de nuevas líneas de servicio y la optimización de procesos de conversión.").
		AddChart(barChart, 180, 80).
		AddDivider().
		AddHeading("2. Segmentación de Clientes", reporter.H1).
		AddParagraph("A continuación se presenta la distribución proporcional de clientes activos según el segmento empresarial correspondiente:").
		AddChart(donutChart, 180, 75).
		AddDivider().
		AddHeading("3. Auditoría de Transacciones", reporter.H1).
		AddParagraph("Registro pormenorizado de las principales operaciones del periodo con trazabilidad de estado:").
		AddTable(tbl.Build()).
		Build()

	return doc
}

func TestPDFGeneration(t *testing.T) {
	doc := buildSampleReport()
	exp := pdf.NewExporter()

	bytesData, err := exp.RenderToBytes(doc)
	if err != nil {
		t.Fatalf("Error renderizando PDF a bytes: %v", err)
	}

	if len(bytesData) == 0 {
		t.Fatalf("El buffer PDF está vacío")
	}

	// Comprobar cabecera mágica de PDF (%PDF-)
	if !strings.HasPrefix(string(bytesData[:10]), "%PDF-") {
		t.Errorf("El archivo generado no tiene cabecera PDF válida")
	}

	// Guardar archivo de muestra para inspección visual
	outDir := filepath.Join("..", "..", "..", "output")
	_ = os.MkdirAll(outDir, 0755)
	outPath := filepath.Join(outDir, "executive_report.pdf")

	if err := exp.Render(doc, outPath); err != nil {
		t.Fatalf("Error guardando PDF en disco: %v", err)
	}

	t.Logf("PDF generado exitosamente en: %s (tamaño: %d bytes)", outPath, len(bytesData))
}

func BenchmarkPDFGeneration(b *testing.B) {
	doc := buildSampleReport()
	exp := pdf.NewExporter()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := exp.RenderToBytes(doc)
		if err != nil {
			b.Fatalf("Error en benchmark: %v", err)
		}
	}
}
