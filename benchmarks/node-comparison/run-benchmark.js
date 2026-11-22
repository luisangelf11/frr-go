const { GoWorkerManager, ReportBuilder, Theme } = require('fast-report-render');
const PDFDocument = require('pdfkit');
const { performance, monitorEventLoopDelay } = require('perf_hooks');

const ITERATIONS = 100;
const WARMUP = 5;

// Datos de prueba: 25 transacciones
const sampleRows = [];
for (let i = 1; i <= 25; i++) {
  sampleRows.push([
    `TRX-${String(i).padStart(4, '0')}`,
    `Corporación Cliente ${i} S.A.`,
    i % 4 === 0 ? 'Norte' : i % 3 === 0 ? 'Sur' : i % 2 === 0 ? 'Este' : 'Oeste',
    `$${(i * 1234.56).toFixed(2)}`,
    i % 5 === 0 ? 'Pendiente' : 'Pagado'
  ]);
}

// -------------------------------------------------------------
// 1. Benchmark: fast-report-render (Go Worker IPC)
// -------------------------------------------------------------
async function runFastReportRender(worker) {
  const times = [];
  let lastBuffer = null;

  // Report definition
  const buildReport = () => {
    return new ReportBuilder('Reporte Financiero de Rendimiento')
      .setSubtitle('Consolidado de Auditoría de Operaciones')
      .setTheme(Theme.CorporateBlue)
      .addMetricCards(
        { title: 'Ingresos Totales', value: '$1,450,200.00', change: '+18.4%', isPositive: true },
        { title: 'Transacciones', value: '1,250', description: 'Meta superada' }
      )
      .addTable({
        headers: ['Código', 'Cliente', 'Región', 'Monto', 'Estado'],
        rows: sampleRows,
      });
  };

  // Warmup
  for (let i = 0; i < WARMUP; i++) {
    await worker.execute(buildReport().toJSON(), 'pdf');
  }

  const loopMonitor = monitorEventLoopDelay({ resolution: 10 });
  loopMonitor.enable();

  const startTotal = performance.now();
  for (let i = 0; i < ITERATIONS; i++) {
    const report = buildReport();
    const t0 = performance.now();
    lastBuffer = await worker.execute(report.toJSON(), 'pdf');
    const t1 = performance.now();
    times.push(t1 - t0);
  }
  const endTotal = performance.now();

  loopMonitor.disable();

  return {
    engine: 'fast-report-render (Go Worker IPC)',
    totalMs: endTotal - startTotal,
    avgMs: times.reduce((a, b) => a + b, 0) / times.length,
    minMs: Math.min(...times),
    maxMs: Math.max(...times),
    p50Ms: times.sort((a, b) => a - b)[Math.floor(times.length * 0.5)],
    p95Ms: times.sort((a, b) => a - b)[Math.floor(times.length * 0.95)],
    pdfBytes: lastBuffer ? lastBuffer.length : 0,
    eventLoopLagMeanMs: loopMonitor.mean / 1e6,
    eventLoopLagMaxMs: loopMonitor.max / 1e6,
  };
}

// -------------------------------------------------------------
// 2. Benchmark: PDFKit (Node.js nativo en hilo de eventos)
// -------------------------------------------------------------
function generatePdfKit() {
  return new Promise((resolve, reject) => {
    const doc = new PDFDocument({ margin: 40, size: 'A4' });
    const chunks = [];

    doc.on('data', chunk => chunks.push(chunk));
    doc.on('end', () => resolve(Buffer.concat(chunks)));
    doc.on('error', reject);

    // Título y subtítulo
    doc.fontSize(20).fillColor('#1E3A8A').text('Reporte Financiero de Rendimiento', { align: 'left' });
    doc.fontSize(10).fillColor('#6B7280').text('Consolidado de Auditoría de Operaciones', { align: 'left' });
    doc.moveDown(1);

    // Tarjetas KPI simuladas
    doc.rect(40, doc.y, 240, 45).fillAndStroke('#F8FAFC', '#E2E8F0');
    doc.fillColor('#1E3A8A').fontSize(14).text('$1,450,200.00', 50, doc.y - 35);
    doc.fillColor('#6B7280').fontSize(8).text('Ingresos Totales (+18.4%)', 50, doc.y);

    doc.rect(300, doc.y - 35, 240, 45).fillAndStroke('#F8FAFC', '#E2E8F0');
    doc.fillColor('#1E3A8A').fontSize(14).text('1,250', 310, doc.y - 35);
    doc.fillColor('#6B7280').fontSize(8).text('Transacciones (Meta superada)', 310, doc.y);

    doc.moveDown(3);

    // Tabla con encabezado y filas con zebra striping
    const tableTop = doc.y;
    const colWidths = [80, 200, 70, 90, 70];
    const headers = ['Código', 'Cliente', 'Región', 'Monto', 'Estado'];

    // Cabecera
    doc.rect(40, tableTop, 515, 20).fill('#1E3A8A');
    let x = 45;
    doc.fillColor('#FFFFFF').fontSize(9).font('Helvetica-Bold');
    headers.forEach((h, idx) => {
      doc.text(h, x, tableTop + 5, { width: colWidths[idx] });
      x += colWidths[idx];
    });

    let currentY = tableTop + 20;
    doc.font('Helvetica').fontSize(9);

    sampleRows.forEach((row, rowIdx) => {
      if (rowIdx % 2 === 1) {
        doc.rect(40, currentY, 515, 18).fill('#F8FAFC');
      }
      doc.rect(40, currentY, 515, 18).stroke('#E2E8F0');

      let cellX = 45;
      doc.fillColor('#1F2937');
      row.forEach((cell, colIdx) => {
        doc.text(cell, cellX, currentY + 4, { width: colWidths[colIdx] });
        cellX += colWidths[colIdx];
      });
      currentY += 18;
    });

    doc.end();
  });
}

async function runPdfKit() {
  const times = [];
  let lastBuffer = null;

  // Warmup
  for (let i = 0; i < WARMUP; i++) {
    await generatePdfKit();
  }

  const loopMonitor = monitorEventLoopDelay({ resolution: 10 });
  loopMonitor.enable();

  const startTotal = performance.now();
  for (let i = 0; i < ITERATIONS; i++) {
    const t0 = performance.now();
    lastBuffer = await generatePdfKit();
    const t1 = performance.now();
    times.push(t1 - t0);
  }
  const endTotal = performance.now();

  loopMonitor.disable();

  return {
    engine: 'PDFKit (Pure Node.js V8)',
    totalMs: endTotal - startTotal,
    avgMs: times.reduce((a, b) => a + b, 0) / times.length,
    minMs: Math.min(...times),
    maxMs: Math.max(...times),
    p50Ms: times.sort((a, b) => a - b)[Math.floor(times.length * 0.5)],
    p95Ms: times.sort((a, b) => a - b)[Math.floor(times.length * 0.95)],
    pdfBytes: lastBuffer ? lastBuffer.length : 0,
    eventLoopLagMeanMs: loopMonitor.mean / 1e6,
    eventLoopLagMaxMs: loopMonitor.max / 1e6,
  };
}

async function main() {
  console.log(`\n======================================================`);
  console.log(`🔬 BENCHMARK COMPARATIVO NODE.JS (${ITERATIONS} iteraciones)`);
  console.log(`======================================================\n`);

  console.log('Iniciando Go Worker...');
  const worker = GoWorkerManager.getInstance();

  console.log('1. Ejecutando pruebas con fast-report-render (Go Worker)...');
  const resFrr = await runFastReportRender(worker);
  worker.destroy();

  console.log('2. Ejecutando pruebas con PDFKit (Node.js nativo)...');
  const resPdfKit = await runPdfKit();

  console.log('\n--- RESULTADOS COMPARATIVOS ---');
  console.table([
    {
      'Motor': resFrr.engine,
      'Avg (ms)': resFrr.avgMs.toFixed(2),
      'P50 (ms)': resFrr.p50Ms.toFixed(2),
      'P95 (ms)': resFrr.p95Ms.toFixed(2),
      'Min (ms)': resFrr.minMs.toFixed(2),
      'Max (ms)': resFrr.maxMs.toFixed(2),
      'Total (ms)': resFrr.totalMs.toFixed(0),
      'Throughput (doc/s)': (1000 / resFrr.avgMs).toFixed(0),
      'EventLoop Lag (ms)': resFrr.eventLoopLagMeanMs.toFixed(3),
      'Tamaño PDF': `${(resFrr.pdfBytes / 1024).toFixed(1)} KB`
    },
    {
      'Motor': resPdfKit.engine,
      'Avg (ms)': resPdfKit.avgMs.toFixed(2),
      'P50 (ms)': resPdfKit.p50Ms.toFixed(2),
      'P95 (ms)': resPdfKit.p95Ms.toFixed(2),
      'Min (ms)': resPdfKit.minMs.toFixed(2),
      'Max (ms)': resPdfKit.maxMs.toFixed(2),
      'Total (ms)': resPdfKit.totalMs.toFixed(0),
      'Throughput (doc/s)': (1000 / resPdfKit.avgMs).toFixed(0),
      'EventLoop Lag (ms)': resPdfKit.eventLoopLagMeanMs.toFixed(3),
      'Tamaño PDF': `${(resPdfKit.pdfBytes / 1024).toFixed(1)} KB`
    }
  ]);
}

main().catch(err => {
  console.error('Error en benchmark:', err);
  process.exit(1);
});
