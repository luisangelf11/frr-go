# fast-report-render

Ultra-fast reporting engine for Node.js powered by Go (`frr-go`).

Generates **PDF**, **Excel (.xlsx)**, and **CSV** with vector charts in **< 1-3 ms** using a persistent Go worker over stdio binary IPC.

## Installation

```bash
npm install fast-report-render
```

The package bundles pre-compiled Go worker binaries for Windows, Linux, and macOS (x64 / arm64).

## Quick Example

```javascript
const { GoWorkerManager, ReportBuilder } = require('fast-report-render');

async function main() {
  const worker = new GoWorkerManager();
  await worker.start();

  const report = new ReportBuilder()
    .setTitle('Monthly Sales')
    .setSubtitle('Operations summary')
    .addKpi('Revenue', '$45,000')
    .addTable(
      ['Product', 'Qty', 'Total'],
      [
        ['Widget A', '100', '$1,000'],
        ['Widget B', '250', '$5,000']
      ]
    );

  // Generate in under 2ms
  const pdfBuffer = await worker.generatePdf(report);
  const xlsxBuffer = await worker.generateExcel(report);
  const csvBuffer = await worker.generateCsv(report);

  await worker.stop();
}

main();
```

## Features

- **Ultra-low latency**: 0.5 - 2 ms per report via IPC binary protocol.
- **Zero dependencies**: No headless browsers, no Puppeteer, no JVM.
- **Accents & UTF-8**: Native support for Spanish accents (`á, é, í, ó, ú, ñ`) and symbols.
- **Multi-format**: PDF with vector charts, Excel (.xlsx) with styles, CSV with UTF-8 BOM.

## License

MIT © luisangelf11
