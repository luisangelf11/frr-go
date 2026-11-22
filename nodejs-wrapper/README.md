# fast-report-render

[![npm version](https://img.shields.io/npm/v/fast-report-render.svg)](https://www.npmjs.com/package/fast-report-render)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

**fast-report-render** (`frr`) es el motor de generación de reportes ultra-rápido para Node.js y TypeScript impulsado por Go (`frr-go`).

Genera documentos ejecutivos en **PDF** (con gráficos vectoriales nativos), hojas de cálculo **Excel (.xlsx)** con estilos y zebra striping, y archivos **CSV** en **menos de 1 a 3 milisegundos**.

---

## ⚡ ¿Por qué fast-report-render?

- **Rendimiento Extremo**: Tiempos de respuesta de **~0.5 a 2 ms** por reporte gracias a un worker persistente en Go comunicado por IPC binario en memoria.
- **Cero Dependencias Pesadas**: Sin navegadores headless (adiós a Puppeteer/Chromium de 400MB), sin Java ni librerías de C pesadas.
- **Gráficos Vectoriales Puros**: Gráficos de Barras, Líneas, Pastel y Donas dibujados con primitivas vectoriales en el PDF (nitidez infinita y peso en KBs).
- **Soporte Perfecto de Acentos y UTF-8**: Manejo nativo de caracteres en español (`á, é, í, ó, ú, ñ, ¿, ¡, €`).
- **Binarios Pre-compilados Incluidos**: Funciona de inmediato en Windows (x64/arm64), Linux (x64/arm64) y macOS (Intel/Apple Silicon).

---

## 📦 Instalación

```bash
npm install fast-report-render
```

---

## 🚀 Inicio Rápido

```javascript
const { GoWorkerManager, ReportBuilder, Theme, ChartType } = require('fast-report-render');
const fs = require('fs');

async function main() {
  // 1. Iniciar el worker persistente de Go
  const worker = new GoWorkerManager();
  await worker.start();

  // 2. Construir el reporte con la API fluida
  const report = new ReportBuilder()
    .setTitle('Reporte Ejecutivo de Ventas')
    .setSubtitle('Consolidado de Operaciones Q4')
    .setTheme(Theme.CorporateBlue)
    
    // Tarjetas KPI
    .addKpiCard({
      title: 'Ingresos Totales',
      value: '$1,250,400',
      change: '+14.8%',
      isPositive: true,
      description: 'vs. trimestre anterior'
    })
    .addKpiCard({
      title: 'Margen Operativo',
      value: '28.4%',
      change: '+2.1%',
      isPositive: true
    })

    // Gráfico de Barras Vectorial
    .addChart({
      type: ChartType.Bar,
      title: 'Ventas por Región',
      categories: ['Norte', 'Centro', 'Sur', 'Este', 'Oeste'],
      series: [
        { name: 'Meta', values: [100, 150, 120, 180, 200] },
        { name: 'Real', values: [120, 140, 135, 195, 210] }
      ]
    })

    // Tabla con estilos automáticos
    .addTable(
      ['Región', 'Vendedor', 'Operaciones', 'Total Facturado'],
      [
        ['Norte', 'Ana Pérez', '45', '$145,200'],
        ['Centro', 'Carlos Gómez', '62', '$210,800'],
        ['Sur', 'Lucía Fernández', '38', '$98,500'],
        ['Este', 'Marcos Ruiz', '51', '$175,000']
      ]
    );

  // 3. Generar en el formato deseado (< 2ms)
  const pdfBuffer = await worker.generatePdf(report);
  const xlsxBuffer = await worker.generateExcel(report);
  const csvBuffer = await worker.generateCsv(report);

  fs.writeFileSync('reporte.pdf', pdfBuffer);
  fs.writeFileSync('reporte.xlsx', xlsxBuffer);
  fs.writeFileSync('reporte.csv', csvBuffer);

  // 4. Detener el worker cuando la app finalice
  await worker.stop();
}

main();
```

---

## 🎨 Temas Visuales Disponibles

`fast-report-render` incluye temas diseñados profesionalmente para ajustar paletas de colores en PDF y Excel:

| Tema | Identificador | Descripción | Paleta Principal |
| :--- | :--- | :--- | :--- |
| **Corporate Blue** | `Theme.CorporateBlue` | Azul ejecutivo corporativo *(por defecto)* | Azul Marino (`#1E3A8A`), Azul Real (`#3B82F6`) |
| **Emerald** | `Theme.Emerald` | Verde esmeralda moderno para finanzas, salud y crecimiento | Verde Bosque (`#065F46`), Esmeralda (`#10B981`) |
| **Dark Slate** | `Theme.DarkSlate` | Pizarra oscura sobria para tecnología e ingeniería | Pizarra Oscura (`#0F172A`), Carbón (`#334155`) |

### Aplicar un Tema:
```typescript
import { ReportBuilder, Theme } from 'fast-report-render';

// Aplicar tema Emerald
const report = new ReportBuilder()
  .setTitle('Métricas de Crecimiento')
  .setTheme(Theme.Emerald);
```

### Paleta Personalizada (Custom Hex):
Puedes definir tus propios colores corporativos:
```typescript
const report = new ReportBuilder()
  .setTitle('Reporte Personalizado')
  .setCustomPalette({
    primary: '#7C3AED',    // Púrpura corporativo
    secondary: '#A78BFA',
    headerBg: '#6D28D9',
    headerText: '#FFFFFF',
    text: '#1F2937',
    alternateRow: '#F5F3FF' // Fila zebra suave
  });
```

---

## 📊 Tipos de Gráficos Vectoriales

Soporta 4 tipos de gráficos vectoriales que se integran en PDF con nitidez infinita o se pueden exportar a SVG:

```typescript
import { ChartType } from 'fast-report-render';

// 1. Gráfico de Barras
report.addBarChart('Ventas Mensuales', ['Ene', 'Feb', 'Mar'], [
  { name: '2025', values: [30, 45, 60] },
  { name: '2026', values: [40, 55, 75] }
]);

// 2. Gráfico de Líneas
report.addLineChart('Tendencia de Tráfico', ['Sem 1', 'Sem 2', 'Sem 3'], [
  { name: 'Visitas', values: [1200, 1900, 2400] }
]);

// 3. Gráfico de Pastel
report.addPieChart('Distribución por Canal', [
  { name: 'Directo', value: 45 },
  { name: 'Orgánico', value: 35 },
  { name: 'Referidos', value: 20 }
]);

// 4. Gráfico de Dona
report.addDonutChart('Participación de Mercado', [
  { name: 'Producto A', value: 50 },
  { name: 'Producto B', value: 30 },
  { name: 'Otros', value: 20 }
]);
```

---

## 🌐 Ejemplo en API Web (Express.js)

En un servidor web, mantén una instancia global de `GoWorkerManager` para atender miles de solicitudes por segundo:

```javascript
const express = require('express');
const { GoWorkerManager, ReportBuilder, Theme } = require('fast-report-render');

const app = express();
const worker = new GoWorkerManager();

// Descargar reporte dinámico en PDF
app.get('/api/reportes/ventas/pdf', async (req, res) => {
  try {
    const report = new ReportBuilder()
      .setTitle('Ventas en Vivo')
      .setTheme(Theme.CorporateBlue)
      .addKpi('Ingresos', '$85,400')
      .addTable(
        ['ID', 'Cliente', 'Monto'],
        [
          ['#101', 'Corporación Alpha', '$12,400'],
          ['#102', 'Beta Systems', '$8,950']
        ]
      );

    const pdfBuffer = await worker.generatePdf(report);

    res.setHeader('Content-Type', 'application/pdf');
    res.setHeader('Content-Disposition', 'attachment; filename="ventas.pdf"');
    res.end(pdfBuffer);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

// Descargar reporte dinámico en Excel
app.get('/api/reportes/ventas/excel', async (req, res) => {
  try {
    const report = new ReportBuilder()
      .setTitle('Ventas en Vivo')
      .addTable(
        ['ID', 'Cliente', 'Monto'],
        [
          ['#101', 'Corporación Alpha', '$12,400'],
          ['#102', 'Beta Systems', '$8,950']
        ]
      );

    const xlsxBuffer = await worker.generateExcel(report);

    res.setHeader('Content-Type', 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet');
    res.setHeader('Content-Disposition', 'attachment; filename="ventas.xlsx"');
    res.end(xlsxBuffer);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

worker.start().then(() => {
  app.listen(3000, () => console.log('Servidor listo en http://localhost:3000'));
});
```

---

## 📄 Licencia

MIT © [Luis Angel](https://github.com/luisangelf11)
