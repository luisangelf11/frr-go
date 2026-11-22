package main

import (
	"fmt"
	"testing"

	_ "github.com/luisangelf11/frr-go/pkg/exporters/pdf"
	"github.com/luisangelf11/frr-go/pkg/reporter"
)

var sampleRows [][]string

func init() {
	sampleRows = make([][]string, 0, 25)
	for i := 1; i <= 25; i++ {
		region := "Oeste"
		if i%4 == 0 {
			region = "Norte"
		} else if i%3 == 0 {
			region = "Sur"
		} else if i%2 == 0 {
			region = "Este"
		}
		estado := "Pagado"
		if i%5 == 0 {
			estado = "Pendiente"
		}
		sampleRows = append(sampleRows, []string{
			fmt.Sprintf("TRX-%04d", i),
			fmt.Sprintf("Corporación Cliente %d S.A.", i),
			region,
			fmt.Sprintf("$%.2f", float64(i)*1234.56),
			estado,
		})
	}
}

func BenchmarkFrrGoNative(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tabla := reporter.NewTable().
			SetHeaders("Código", "Cliente", "Región", "Monto", "Estado").
			AddRows(sampleRows).
			Build()

		doc := reporter.New("Reporte Financiero de Rendimiento").
			SetSubtitle("Consolidado de Auditoría de Operaciones").
			SetTheme(reporter.ThemeCorporateBlue).
			AddMetricCards(
				reporter.MetricCard{Title: "Ingresos Totales", Value: "$1,450,200.00", Change: "+18.4%", IsPositive: true},
				reporter.MetricCard{Title: "Transacciones", Value: "1,250", Description: "Meta superada"},
			).
			AddTable(tabla).
			Build()

		_, err := doc.ExportPDFBytes()
		if err != nil {
			b.Fatalf("Error exportando PDF: %v", err)
		}
	}
}
