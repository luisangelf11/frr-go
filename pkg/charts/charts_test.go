package charts_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/luisangelf11/frr-go/pkg/charts"
)

func TestBarChartGeneration(t *testing.T) {
	b := charts.NewBarChart("Ventas Trimestrales").
		SetSubtitle("Comparativa 2025 vs 2026").
		SetCategories("Q1", "Q2", "Q3", "Q4").
		AddSeries("2025", []float64{120, 150, 180, 220}).
		AddSeries("2026", []float64{140, 190, 210, 260}).
		SetShowLegend(true).
		SetShowValues(true).
		SetValueFormat("$", "K")

	ops, err := charts.Generate(b.Config())
	if err != nil {
		t.Fatalf("Error generando gráfico de barras: %v", err)
	}

	if len(ops) == 0 {
		t.Fatalf("Se esperaban operaciones de dibujo, se obtuvieron 0")
	}

	rectCount := 0
	textCount := 0
	for _, op := range ops {
		if op.Type == charts.OpRect {
			rectCount++
		}
		if op.Type == charts.OpText {
			textCount++
		}
	}

	// 8 barras (4 categorias * 2 series) + 2 cuadros de leyenda = al menos 10 rectángulos
	if rectCount < 10 {
		t.Errorf("Se esperaban al menos 10 rectángulos, se obtuvieron %d", rectCount)
	}

	// Generar SVG y validar
	svgStr, err := charts.GenerateSVG(b.Config())
	if err != nil {
		t.Fatalf("Error generando SVG: %v", err)
	}
	if (!strings.HasPrefix(svgStr, "<svg") && !strings.HasPrefix(svgStr, "<?xml")) || !strings.HasSuffix(strings.TrimSpace(svgStr), "</svg>") {
		t.Errorf("Salida SVG no válida")
	}

	saveSampleSVG(t, "bar_chart.svg", svgStr)
}

func TestLineChartGeneration(t *testing.T) {
	b := charts.NewLineChart("Crecimiento de Usuarios").
		SetSubtitle("Métricas de Adquisición").
		SetCategories("Ene", "Feb", "Mar", "Abr", "May", "Jun").
		AddSeries("Orgánico", []float64{1200, 1500, 1800, 2200, 2700, 3100}).
		AddSeries("Referidos", []float64{400, 600, 850, 1100, 1400, 1750}).
		SetShowLegend(true).
		SetShowValues(true)

	ops, err := charts.Generate(b.Config())
	if err != nil {
		t.Fatalf("Error generando gráfico de líneas: %v", err)
	}

	hasCircles := false
	hasLines := false
	for _, op := range ops {
		if op.Type == charts.OpCircle {
			hasCircles = true
		}
		if op.Type == charts.OpLine {
			hasLines = true
		}
	}

	if !hasCircles || !hasLines {
		t.Errorf("El gráfico de líneas debe contener círculos y líneas")
	}

	svgStr, err := charts.GenerateSVG(b.Config())
	if err != nil {
		t.Fatalf("Error generando SVG de líneas: %v", err)
	}
	saveSampleSVG(t, "line_chart.svg", svgStr)
}

func TestPieAndDonutChartGeneration(t *testing.T) {
	// 1. Pastel
	pie := charts.NewPieChart("Distribución de Presupuesto").
		SetSubtitle("Gastos Operativos 2026").
		SetData(
			[]string{"Ingeniería", "Marketing", "Ventas", "Operaciones", "Legal"},
			[]float64{45000, 25000, 30000, 15000, 5000},
		).
		SetValueFormat("$", "")

	pieSVG, err := charts.GenerateSVG(pie.Config())
	if err != nil {
		t.Fatalf("Error generando SVG de pastel: %v", err)
	}
	saveSampleSVG(t, "pie_chart.svg", pieSVG)

	// 2. Dona
	donut := charts.NewDonutChart("Distribución de Clientes por Región").
		SetData(
			[]string{"América del Norte", "Europa", "América Latina", "Asia-Pacífico"},
			[]float64{420, 280, 190, 110},
		).
		SetValueFormat("", " clientes")

	donutSVG, err := charts.GenerateSVG(donut.Config())
	if err != nil {
		t.Fatalf("Error generando SVG de dona: %v", err)
	}
	saveSampleSVG(t, "donut_chart.svg", donutSVG)
}

func saveSampleSVG(t *testing.T, filename, content string) {
	t.Helper()
	outDir := filepath.Join("..", "..", "output")
	_ = os.MkdirAll(outDir, 0755)
	filePath := filepath.Join(outDir, filename)
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Logf("Aviso: no se pudo guardar archivo de muestra %s: %v", filename, err)
	}
}

// --- Benchmarks de Velocidad ---

func BenchmarkBarChart(b *testing.B) {
	cfg := charts.NewBarChart("Benchmark Bar").
		SetCategories("A", "B", "C", "D", "E").
		AddSeries("S1", []float64{10, 20, 30, 40, 50}).
		AddSeries("S2", []float64{15, 25, 35, 45, 55}).
		Config()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = charts.Generate(cfg)
	}
}

func BenchmarkLineChart(b *testing.B) {
	cfg := charts.NewLineChart("Benchmark Line").
		SetCategories("A", "B", "C", "D", "E").
		AddSeries("S1", []float64{10, 20, 30, 40, 50}).
		Config()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = charts.Generate(cfg)
	}
}

func BenchmarkDonutChart(b *testing.B) {
	cfg := charts.NewDonutChart("Benchmark Donut").
		SetData([]string{"A", "B", "C", "D"}, []float64{40, 30, 20, 10}).
		Config()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = charts.Generate(cfg)
	}
}

func BenchmarkGenerateSVG(b *testing.B) {
	cfg := charts.NewBarChart("Benchmark SVG").
		SetCategories("A", "B", "C", "D", "E").
		AddSeries("S1", []float64{10, 20, 30, 40, 50}).
		Config()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = charts.GenerateSVG(cfg)
	}
}
