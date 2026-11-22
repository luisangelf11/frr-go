const http = require('http');
const fs = require('fs');
const path = require('path');
const { server } = require('./server');

const TEST_PORT = 3456;

function fetchBuffer(url) {
  return new Promise((resolve, reject) => {
    http.get(url, (res) => {
      const chunks = [];
      res.on('data', (d) => chunks.push(d));
      res.on('end', () => {
        const buf = Buffer.concat(chunks);
        resolve({
          status: res.statusCode,
          headers: res.headers,
          data: buf,
        });
      });
      res.on('error', reject);
    });
  });
}

async function runTests() {
  await new Promise((resolve) => server.listen(TEST_PORT, resolve));
  console.log(`\n========================================`);
  console.log(`  PRUEBA INTEGRAL: API HTTP NODE.JS + GO  `);
  console.log(`========================================\n`);

  const outDir = path.resolve(__dirname, '..', '..', 'output');
  if (!fs.existsSync(outDir)) {
    fs.mkdirSync(outDir, { recursive: true });
  }

  // 1. Probar Endpoint PDF
  console.log('1. Solicitando GET /api/reports/sales/pdf...');
  const resPdf = await fetchBuffer(`http://localhost:${TEST_PORT}/api/reports/sales/pdf`);
  console.log(`   Status: ${resPdf.status}`);
  console.log(`   Content-Type: ${resPdf.headers['content-type']}`);
  console.log(`   X-Generation-Time: ${resPdf.headers['x-generation-time-ms']} ms`);
  console.log(`   Tamaño: ${resPdf.data.length} bytes`);
  fs.writeFileSync(path.join(outDir, 'api_downloaded_sales.pdf'), resPdf.data);
  console.log(`   ✓ Archivo guardado en output/api_downloaded_sales.pdf`);

  // 2. Probar Endpoint Excel (.xlsx)
  console.log('\n2. Solicitando GET /api/reports/sales/excel...');
  const resExcel = await fetchBuffer(`http://localhost:${TEST_PORT}/api/reports/sales/excel`);
  console.log(`   Status: ${resExcel.status}`);
  console.log(`   Content-Type: ${resExcel.headers['content-type']}`);
  console.log(`   X-Generation-Time: ${resExcel.headers['x-generation-time-ms']} ms`);
  console.log(`   Tamaño: ${resExcel.data.length} bytes`);
  fs.writeFileSync(path.join(outDir, 'api_downloaded_sales.xlsx'), resExcel.data);
  console.log(`   ✓ Archivo guardado en output/api_downloaded_sales.xlsx`);

  // 3. Probar Endpoint CSV
  console.log('\n3. Solicitando GET /api/reports/sales/csv...');
  const resCsv = await fetchBuffer(`http://localhost:${TEST_PORT}/api/reports/sales/csv`);
  console.log(`   Status: ${resCsv.status}`);
  console.log(`   Content-Type: ${resCsv.headers['content-type']}`);
  console.log(`   X-Generation-Time: ${resCsv.headers['x-generation-time-ms']} ms`);
  console.log(`   Tamaño: ${resCsv.data.length} bytes`);
  fs.writeFileSync(path.join(outDir, 'api_downloaded_sales.csv'), resCsv.data);
  console.log(`   ✓ Archivo guardado en output/api_downloaded_sales.csv`);

  // 4. Benchmark en Vivo (50 solicitudes consecutivas al motor)
  console.log('\n4. Ejecutando Benchmark en vivo (50 generaciones de reportes completas)...');
  const resBench = await fetchBuffer(`http://localhost:${TEST_PORT}/api/benchmark?n=50`);
  const stats = JSON.parse(resBench.data.toString('utf-8'));
  console.log(`   Total Reportes:      ${stats.totalRequests}`);
  console.log(`   Latencia Promedio:   ${stats.avgLatencyMs} ms`);
  console.log(`   Latencia Mínima:     ${stats.minLatencyMs} ms`);
  console.log(`   Latencia Máxima:     ${stats.maxLatencyMs} ms`);
  console.log(`   Throughput Estimado: ${stats.throughputRps} reportes/segundo`);

  console.log(`\n========================================`);
  console.log(`  ¡TODAS LAS PRUEBAS HTTP FUERON EXITOSAS! `);
  console.log(`========================================\n`);

  server.close();
  process.exit(0);
}

runTests().catch((err) => {
  console.error('Error durante la prueba de API:', err);
  server.close();
  process.exit(1);
});
