const { ReportBuilder, ChartType, HeadingLevel, Theme } = require('./dist/index.js');
const path = require('path');
const fs = require('fs');

async function run() {
  console.log('--- Iniciando prueba de humo Node.js + Go Worker ---');
  const start = Date.now();

  const report = new ReportBuilder('Reporte de Prueba desde Node.js')
    .setSubtitle('Generado en milisegundos mediante Stdio IPC')
    .setAuthor('Ingeniería de Software')
    .setTheme(Theme.CorporateBlue)
    .addMetricCards(
      { title: 'LATENCIA GO', value: '2.5 ms', change: '-95%', isPositive: true },
      { title: 'DEPENDENCIAS', value: '1', description: 'Zero-deps footprint' },
      { title: 'THROUGHPUT', value: '400 req/s', change: '+300%', isPositive: true }
    )
    .addHeading('1. Gráfico de Rendimiento', HeadingLevel.H1)
    .addParagraph('Este reporte fue solicitado desde Node.js y compilado directamente por el motor nativo de Go en memoria.')
    .addChart({
      type: ChartType.Bar,
      title: 'Métricas de Ejecución',
      categories: ['Ene', 'Feb', 'Mar', 'Abr'],
      series: [
        { name: 'Node Directo', values: [120, 140, 160, 200] },
        { name: 'Go IPC Worker', values: [350, 420, 480, 560] },
      ],
      valuePrefix: '$',
    })
    .addDivider()
    .addHeading('2. Auditoría Tabular', HeadingLevel.H1)
    .addTable({
      title: 'Registros Operativos',
      headers: ['Código', 'Operación', 'Región', 'Monto (€)', 'Estado'],
      rows: [
        ['OP-101', 'Exportación PDF', 'América Latina', '1,250.00 €', 'Éxito'],
        ['OP-102', 'Exportación Excel', 'Europa', '3,400.50 €', 'Éxito'],
        ['OP-103', 'Exportación CSV', 'Norteamérica', '890.00 €', 'Éxito'],
      ],
    });

  const outDir = path.resolve(__dirname, '..', 'output');
  if (!fs.existsSync(outDir)) {
    fs.mkdirSync(outDir, { recursive: true });
  }

  // 1. Exportar PDF
  const t0 = performance.now();
  const pdfBuf = await report.exportToBuffer('pdf');
  const tPdf = performance.now() - t0;
  fs.writeFileSync(path.join(outDir, 'node_sample.pdf'), pdfBuf);
  console.log(`✓ PDF generado en ${tPdf.toFixed(2)} ms (${pdfBuf.length} bytes)`);

  // 2. Exportar Excel
  const t1 = performance.now();
  const xlsxBuf = await report.exportToBuffer('excel');
  const tXlsx = performance.now() - t1;
  fs.writeFileSync(path.join(outDir, 'node_sample.xlsx'), xlsxBuf);
  console.log(`✓ Excel generado en ${tXlsx.toFixed(2)} ms (${xlsxBuf.length} bytes)`);

  // 3. Exportar CSV
  const t2 = performance.now();
  const csvBuf = await report.exportToBuffer('csv', ';');
  const tCsv = performance.now() - t2;
  fs.writeFileSync(path.join(outDir, 'node_sample.csv'), csvBuf);
  console.log(`✓ CSV generado en ${tCsv.toFixed(2)} ms (${csvBuf.length} bytes)`);

  console.log(`--- Prueba finalizada exitosamente en ${(Date.now() - start)} ms totales ---`);
  process.exit(0);
}

run().catch((err) => {
  console.error('Error en prueba de humo:', err);
  process.exit(1);
});
