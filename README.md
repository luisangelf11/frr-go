# ⚡ frr-go (Fast Report Render)

> **Motor de reportería ultra-rápido, ergonómico y ligero en Go, con puente nativo de alto rendimiento para Node.js.**  
> Diseñado para generar reportes profesionales con tablas estilizadas, gráficos vectoriales, tarjetas de métricas e imágenes en **PDF**, **Excel (.xlsx)** y **CSV** en **1 a 3 milisegundos**.

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Node.js](https://img.shields.io/badge/Node.js-18+-339933?style=flat&logo=nodedotjs)](https://nodejs.org)
[![TypeScript](https://img.shields.io/badge/TypeScript-Ready-3178C6?style=flat&logo=typescript)](https://www.typescriptlang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Dependencies](https://img.shields.io/badge/Go%20Dependencies-1%20only-brightgreen)]()

---

## 🌟 Características Principales

- 🚀 **Rendimiento Extremo**: Generación de reportes completos en memoria en **1 a 3 ms** (~1,700 reportes/segundo a través de HTTP en Node.js).
- 🪶 **Filosofía de Mínimas Dependencias**:
  - **Gráficos**: Motor vectorial nativo en Go (**0 dependencias externas**).
  - **Excel (.xlsx)**: Generador OpenXML con `sharedStrings.xml` y empaquetador ZIP estándar (**0 dependencias externas**).
  - **CSV**: Biblioteca estándar `encoding/csv` con UTF-8 BOM (**0 dependencias externas**).
  - **PDF**: Motor de maquetación sobre `github.com/go-pdf/fpdf` (**1 única dependencia externa en todo el proyecto**).
- 🎨 **Gráficos Vectoriales Nativos**: Barras agrupadas, líneas con áreas suaves, gráficos de pastel y dona dibujados directamente con primitivas vectoriales en el PDF (nitidez infinita sin pixelado, peso en KBs) y exportables a SVG.
- 📊 **Tablas Inteligentes**: Ancho de columna auto-calculado, cabeceras con color temático, filas alternadas (*Zebra Striping*), ajuste automático de texto multilínea y **repetición automática de cabeceras en saltos de página**.
- 🌐 **Soporte de Caracteres Internacionales**: Tratamiento impecable de tildes (`á, é, í, ó, ú`), `ñ`, diéresis (`ü`), signos (`¿, ¡`) y monedas (`€, $`) en todos los formatos.
- ⚡ **Puente de Doble Modo para Node.js**:
  - **Modo Persistent Worker (Stdio IPC)**: Node.js arranca el binario una sola vez. Las solicitudes se transmiten en microsegundos mediante un protocolo binario de longitud prefijada (`uint32`). Sin riesgo de caídas de CGO y sin necesidad de compiladores C en producción.

---

## 📦 Instalación

### En Go:
```bash
go get github.com/luisangelf11/frr-go
```

### En Node.js (TypeScript / JavaScript):
```bash
npm install fast-report-render
```

---

## 🚀 Guía Rápida en Go

```go
package main

import (
	"fmt"
	"log"

	"github.com/luisangelf11/frr-go/pkg/charts"
	_ "github.com/luisangelf11/frr-go/pkg/exporters/csv"
	_ "github.com/luisangelf11/frr-go/pkg/exporters/excel"
	_ "github.com/luisangelf11/frr-go/pkg/exporters/pdf"
	"github.com/luisangelf11/frr-go/pkg/reporter"
)

func main() {
	// 1. Gráfico de Barras Vectorial
	barChart := charts.NewBarChart("Ingresos por Región").
		SetSubtitle("Comparativo 2025 vs 2026 en miles de USD").
		SetCategories("Norte", "Sur", "Este", "Oeste").
		AddSeries("2025", []float64{120, 150, 95, 210}).
		AddSeries("2026", []float64{145, 180, 110, 245}).
		SetShowLegend(true).
		SetShowValues(true).
		SetValueFormat("$", "K").
		Config()

	// 2. Tabla de Datos con Iteración Dinámica
	type Transaccion struct {
		Codigo  string
		Cliente string
		Region  string
		Monto   float64
		Estado  string
	}

	// Ejemplo de datos recuperados (ej. desde una base de datos o API)
	datos := []Transaccion{
		{"TRX-001", "Empresa Alpha S.A.", "Norte", 14500.00, "Pagado"},
		{"TRX-002", "Iberia Logística S.L.", "Sur", 8250.50, "Pagado"},
		{"TRX-003", "Nordic Systems AB", "Este", 22100.00, "Pendiente"},
	}

	tablaBuilder := reporter.NewTable().
		SetTitle("Detalle de Facturación").
		SetHeaders("Código", "Cliente", "Región", "Monto", "Estado").
		SetAlignments(
			reporter.AlignCenter,
			reporter.AlignLeft,
			reporter.AlignCenter,
			reporter.AlignRight,
			reporter.AlignCenter,
		)

	// Iteración para agregar filas dinámicamente con AddRow
	for _, t := range datos {
		tablaBuilder.AddRow(
			t.Codigo,
			t.Cliente,
			t.Region,
			fmt.Sprintf("$%.2f", t.Monto),
			t.Estado,
		)
	}

	// (Nota: también puedes usar tablaBuilder.AddRows(filas) si ya tienes un [][]string)
	tabla := tablaBuilder.Build()

	// 3. Construir Reporte con API Fluida
	doc := reporter.New("Reporte Financiero Q1").
		SetSubtitle("Consolidado Integral de Operaciones").
		SetAuthor("Departamento de Finanzas").
		SetDate("2026-03-31").
		SetTheme(reporter.ThemeCorporateBlue).
		// Tarjetas de Métricas (KPIs)
		AddMetricCards(
			reporter.MetricCard{Title: "VENTAS TOTALES", Value: "$44,850.50", Change: "+18.4%", IsPositive: true},
			reporter.MetricCard{Title: "NUEVOS CLIENTES", Value: "128", Description: "Meta mensual cumplida"},
		).
		AddHeading("1. Resumen Estratégico", reporter.H1).
		AddParagraph("A continuación se presentan las métricas de rendimiento consolidadas:").
		AddChart(barChart, 180, 80).
		AddDivider().
		AddHeading("2. Auditoría de Transacciones", reporter.H1).
		AddTable(tabla).
		Build()

	// 4. Exportar a cualquier formato en milisegundos:
	if err := doc.ExportPDF("reporte.pdf"); err != nil {
		log.Fatalf("Error en PDF: %v", err)
	}

	if err := doc.ExportExcel("reporte.xlsx"); err != nil {
		log.Fatalf("Error en Excel: %v", err)
	}

	if err := doc.ExportCSV("reporte.csv"); err != nil {
		log.Fatalf("Error en CSV: %v", err)
	}

	log.Println("¡Reportes generados con éxito!")
}
```

---

## ⚡ Guía Rápida en Node.js (TypeScript / JavaScript)

El SDK de Node.js provee una API fluida idéntica y se comunica en segundo plano con el proceso de Go mediante **Stdio IPC**, logrando tiempos de respuesta de **1 a 3 milisegundos**.

```typescript
import { ReportBuilder, ChartType, HeadingLevel, Theme } from 'fast-report-render';
import * as fs from 'fs';

async function generate() {
  const report = new ReportBuilder("Reporte de Ventas Q1")
    .setSubtitle("Generado en tiempo real desde Node.js")
    .setAuthor("API Microservice")
    .setTheme(Theme.CorporateBlue)
    // Métricas destacadas (KPIs)
    .addMetricCards(
      { title: "INGRESOS", value: "82,950 €", change: "+24.5%", isPositive: true },
      { title: "TIEMPO", value: "< 2 ms", isPositive: true }
    )
    .addHeading("1. Rendimiento por Región", HeadingLevel.H1)
    .addChart({
      type: ChartType.Bar,
      title: "Ventas por Zona",
      categories: ["Norte", "Sur", "Este", "Oeste"],
      series: [
        { name: "Meta", values: [18000, 25000, 20000, 30000] },
        { name: "Alcanzado", values: [21200, 30350, 19500, 34000] },
      ],
      valuePrefix: "€",
    })
    .addDivider()
    .addHeading("2. Desglose de Operaciones", HeadingLevel.H1)
    .addTable({
      title: "Transacciones",
      headers: ["Código", "Cliente", "Región", "Monto", "Estado"],
      rows: [
        ["TRX-01", "Corporación Andina", "Norte", "21,200.00 €", "Completado"],
        ["TRX-02", "Iberia Logística", "Sur", "30,350.00 €", "Completado"],
        ["TRX-03", "Soluciones Delta", "Este", "19,500.00 €", "Pendiente"],
      ],
    });

  // Generar directamente a Buffer en memoria (ideal para Express / Fastify / NestJS)
  const pdfBuffer = await report.exportToBuffer('pdf');
  const xlsxBuffer = await report.exportToBuffer('excel');
  const csvBuffer = await report.exportToBuffer('csv', ';');

  // O guardar directamente a archivo:
  await report.exportToFile('pdf', './salida.pdf');
  await report.exportToFile('excel', './salida.xlsx');
  await report.exportToFile('csv', './salida.csv');

  console.log('¡Reportes generados en milisegundos!');
}

generate();
```

### Ejemplo en Express / Fastify: Streaming en Vivo

```typescript
import express from 'express';
import { ReportBuilder } from 'fast-report-render';

const app = express();

app.get('/api/reports/sales/pdf', async (req, res) => {
  const report = buildMyReport(); // Instancia de ReportBuilder
  const buffer = await report.exportToBuffer('pdf');

  res.setHeader('Content-Type', 'application/pdf');
  res.setHeader('Content-Disposition', 'inline; filename="reporte.pdf"');
  res.send(buffer);
});

app.get('/api/reports/sales/excel', async (req, res) => {
  const report = buildMyReport();
  const buffer = await report.exportToBuffer('excel');

  res.setHeader('Content-Type', 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet');
  res.setHeader('Content-Disposition', 'attachment; filename="reporte.xlsx"');
  res.send(buffer);
});
```

---

## 📊 Benchmarks de Rendimiento

Pruebas ejecutadas en hardware idéntico (Intel Core i5-11400H @ 2.70GHz, Windows 11 Pro):

### Comparativa vs Motores Populares (PDFKit y QuestPDF):
*Prueba sobre 100 iteraciones consecutivas de un reporte ejecutivo A4 con 2 KPIs y tabla con estilos zebra de 25 filas.*

| Motor / Biblioteca | Runtime / Ecosistema | Latencia Media | Mediana (P50) | Throughput | Peso PDF | Dependencias Runtime |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 🥇 **`frr-go` (Nativo)** | Go 1.26 (Compilado) | **0.82 ms** | **0.81 ms** | **~1,218 doc/s** | **3.5 KB** | **0** (solo `fpdf`, binario 3.6MB) |
| 🥈 **`fast-report-render`** | Node.js (Go Worker IPC) | **0.98 ms** | **0.97 ms** | **~1,023 doc/s** | **3.5 KB** | **0** (binario pre-compilado en Go) |
| 🥉 **`PDFKit`** | Node.js (V8 JS puro) | **6.15 ms** | **5.79 ms** | **~163 doc/s** | **3.0 KB** | Node.js stdlib |
| 4️⃣ **`QuestPDF`** | .NET 10 (C# / SkiaSharp) | **11.99 ms** | **11.43 ms** | **~83 doc/s** | **48.3 KB** | .NET Runtime + Skia C++ DLLs |

> [!TIP]
> Consulta el informe técnico exhaustivo con gráficas y análisis de cuellos de botella en [.agents/docs/BENCHMARK_COMPARISON.md](file:///.agents/docs/BENCHMARK_COMPARISON.md).

### Métricas Internas por Componente:
| Operación | Tiempo por Reporte | Rendimiento Estimado | Dependencias Externas |
| :--- | :--- | :--- | :--- |
| **Gráfico de Barras Vectorial** | **0.0057 ms** (5.7 µs) | ~175,000 / seg | **0** |
| **Gráfico de Líneas con Área** | **0.0071 ms** (7.1 µs) | ~140,000 / seg | **0** |
| **Exportación CSV (1,000 filas con BOM)** | **0.129 ms** (129 µs) | ~7,700 / seg | **0** |
| **Exportación Excel (1,000 filas estilizadas)** | **1.35 ms** | ~750 / seg | **0** |
| **Exportación PDF (Multi-página con gráficos)** | **3.60 ms** | ~280 / seg | **1** (`fpdf`) |
| **Node.js HTTP API (Round-trip completo)** | **0.56 ms promedio** | **~1,760 solicitudes / seg** | - |

---

## 🛠️ Uso del CLI (`frr-cli`)

El binario ejecutable `frr-cli` puede ser utilizado directamente desde la terminal o scripts de CI/CD:

```bash
# Modo 1: CLI One-Shot
frr-cli -format pdf -spec reporte.json -out salida.pdf
frr-cli -format excel -spec reporte.json -out salida.xlsx
frr-cli -format csv -spec reporte.json -out salida.csv

# Modo 2: Persistent Worker (para Node.js y microservicios)
frr-cli --worker
```

---

## 🏗️ Arquitectura del Sistema

```
frr-go/
├── cmd/
│   └── frr-cli/                # Binario CLI y Persistent Worker Stdio IPC
├── pkg/
│   ├── reporter/               # API Fluida para Go y Modelo de Documento
│   ├── charts/                 # Motor de Gráficos Vectoriales y SVG (0 deps)
│   ├── exporters/              # Renderizadores de salida
│   │   ├── pdf/                # Renderizador PDF sobre fpdf
│   │   ├── excel/              # Generador OpenXML .xlsx nativo (0 deps)
│   │   └── csv/                # Generador CSV con UTF-8 BOM (0 deps)
│   └── spec/                   # Parser y esquema JSON compartido Go <-> Node
├── nodejs-wrapper/             # SDK de Node.js / TypeScript (fast-report-render)
└── examples/
    └── node-api/               # API de demostración y pruebas en vivo
```

---

## 📄 Licencia

Este proyecto está bajo la Licencia **MIT**. Consulta el archivo `LICENSE` para más información.
