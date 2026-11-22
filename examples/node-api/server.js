const http = require('http');
const path = require('path');
const { ReportBuilder, ChartType, HeadingLevel, Theme } = require('../../nodejs-wrapper/dist/index.js');

const PORT = process.env.PORT || 3000;

function createSampleSalesReport() {
  const tbl = {
    title: 'Detalle de Facturación Q1',
    headers: ['Código', 'Fecha', 'Cliente', 'Región', 'Monto (€)', 'Estado'],
    rows: [
      ['FAC-001', '2026-03-01', 'Corporación Andina S.A.', 'América Latina', '14,500.00 €', 'Pagado'],
      ['FAC-002', '2026-03-02', 'Iberia Logística S.L.', 'Europa', '8,250.50 €', 'Pagado'],
      ['FAC-003', '2026-03-03', 'Nordic Systems AB', 'Europa', '22,100.00 €', 'Pendiente'],
      ['FAC-004', '2026-03-04', 'Pacífico Soluciones', 'América Latina', '6,700.00 €', 'Pagado'],
      ['FAC-005', '2026-03-05', 'Nippon Commerce Ltd.', 'Asia-Pacífico', '31,400.00 €', 'Pagado'],
    ],
  };

  return new ReportBuilder('Reporte de Ventas en Vivo (API Node.js)')
    .setSubtitle('Generación dinámica y streaming de alta velocidad')
    .setAuthor('API Gateway Microservice')
    .setTheme(Theme.CorporateBlue)
    .addMetricCards(
      { title: 'FACTURACIÓN TOTAL', value: '82,950.50 €', change: '+24.5%', isPositive: true },
      { title: 'CLIENTES ATENDIDOS', value: '48', description: 'Meta mensual alcanzada' },
      { title: 'TIEMPO GENERACIÓN', value: '< 3 ms', change: '-90% latencia', isPositive: true }
    )
    .addHeading('1. Comparativa de Ingresos por Región', HeadingLevel.H1)
    .addChart({
      type: ChartType.Bar,
      title: 'Ventas por Región Geográfica',
      categories: ['América Latina', 'Europa', 'Asia-Pacífico'],
      series: [
        { name: 'Meta', values: [18000, 25000, 28000] },
        { name: 'Alcanzado', values: [21200, 30350, 31400] },
      ],
      valuePrefix: '€',
    })
    .addDivider()
    .addHeading('2. Desglose de Transacciones', HeadingLevel.H1)
    .addTable(tbl);
}

const server = http.createServer(async (req, res) => {
  const url = new URL(req.url, `http://${req.headers.host}`);
  const pathname = url.pathname;

  // Endpoint de salud
  if (pathname === '/health' || pathname === '/') {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    return res.end(JSON.stringify({ status: 'ok', engine: 'frr-go Stdio Worker', version: '1.0.0' }));
  }

  // 1. Endpoint PDF
  if (pathname === '/api/reports/sales/pdf') {
    const t0 = performance.now();
    try {
      const report = createSampleSalesReport();
      const pdfBuf = await report.exportToBuffer('pdf');
      const latency = (performance.now() - t0).toFixed(2);

      res.writeHead(200, {
        'Content-Type': 'application/pdf',
        'Content-Disposition': 'inline; filename="reporte_ventas.pdf"',
        'Content-Length': pdfBuf.length,
        'X-Engine': 'frr-go',
        'X-Generation-Time-Ms': latency,
      });
      return res.end(pdfBuf);
    } catch (err) {
      res.writeHead(500, { 'Content-Type': 'application/json' });
      return res.end(JSON.stringify({ error: err.message }));
    }
  }

  // 2. Endpoint Excel
  if (pathname === '/api/reports/sales/excel') {
    const t0 = performance.now();
    try {
      const report = createSampleSalesReport();
      const xlsxBuf = await report.exportToBuffer('excel');
      const latency = (performance.now() - t0).toFixed(2);

      res.writeHead(200, {
        'Content-Type': 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
        'Content-Disposition': 'attachment; filename="reporte_ventas.xlsx"',
        'Content-Length': xlsxBuf.length,
        'X-Generation-Time-Ms': latency,
      });
      return res.end(xlsxBuf);
    } catch (err) {
      res.writeHead(500, { 'Content-Type': 'application/json' });
      return res.end(JSON.stringify({ error: err.message }));
    }
  }

  // 3. Endpoint CSV
  if (pathname === '/api/reports/sales/csv') {
    const t0 = performance.now();
    try {
      const report = createSampleSalesReport();
      const csvBuf = await report.exportToBuffer('csv', ';');
      const latency = (performance.now() - t0).toFixed(2);

      res.writeHead(200, {
        'Content-Type': 'text/csv; charset=utf-8',
        'Content-Disposition': 'attachment; filename="reporte_ventas.csv"',
        'Content-Length': csvBuf.length,
        'X-Generation-Time-Ms': latency,
      });
      return res.end(csvBuf);
    } catch (err) {
      res.writeHead(500, { 'Content-Type': 'application/json' });
      return res.end(JSON.stringify({ error: err.message }));
    }
  }

  // 4. Benchmark en vivo
  if (pathname === '/api/benchmark') {
    const count = parseInt(url.searchParams.get('n') || '50', 10);
    const times = [];

    for (let i = 0; i < count; i++) {
      const report = createSampleSalesReport();
      const t0 = performance.now();
      await report.exportToBuffer('excel');
      times.push(performance.now() - t0);
    }

    const avg = times.reduce((a, b) => a + b, 0) / times.length;
    const min = Math.min(...times);
    const max = Math.max(...times);

    res.writeHead(200, { 'Content-Type': 'application/json' });
    return res.end(JSON.stringify({
      totalRequests: count,
      avgLatencyMs: Number(avg.toFixed(3)),
      minLatencyMs: Number(min.toFixed(3)),
      maxLatencyMs: Number(max.toFixed(3)),
      throughputRps: Number((1000 / avg).toFixed(1)),
    }));
  }

  res.writeHead(404, { 'Content-Type': 'application/json' });
  res.end(JSON.stringify({ error: 'Ruta no encontrada' }));
});

if (require.main === module) {
  server.listen(PORT, () => {
    console.log(`Servidor API de reportería escuchando en http://localhost:${PORT}`);
  });
}

module.exports = { server, createSampleSalesReport };
