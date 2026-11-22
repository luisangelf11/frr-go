using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Linq;
using QuestPDF.Fluent;
using QuestPDF.Helpers;
using QuestPDF.Infrastructure;

namespace QuestPdfBenchmark
{
    class Program
    {
        const int Iterations = 100;
        const int Warmup = 5;

        record RowData(string Codigo, string Cliente, string Region, string Monto, string Estado);

        static void Main(string[] args)
        {
            QuestPDF.Settings.License = LicenseType.Community;

            var sampleRows = new List<RowData>();
            for (int i = 1; i <= 25; i++)
            {
                sampleRows.Add(new RowData(
                    $"TRX-{i:D4}",
                    $"Corporación Cliente {i} S.A.",
                    i % 4 == 0 ? "Norte" : i % 3 == 0 ? "Sur" : i % 2 == 0 ? "Este" : "Oeste",
                    $"${(i * 1234.56):F2}",
                    i % 5 == 0 ? "Pendiente" : "Pagado"
                ));
            }

            Document CreateDoc()
            {
                return Document.Create(container =>
                {
                    container.Page(page =>
                    {
                        page.Size(PageSizes.A4);
                        page.Margin(40);
                        page.DefaultTextStyle(x => x.FontSize(9).FontColor("#1F2937"));

                        page.Header().Column(col =>
                        {
                            col.Item().Text("Reporte Financiero de Rendimiento").FontSize(20).Bold().FontColor("#1E3A8A");
                            col.Item().Text("Consolidado de Auditoría de Operaciones").FontSize(10).FontColor("#6B7280");
                        });

                        page.Content().PaddingVertical(15).Column(col =>
                        {
                            // Tarjetas KPI
                            col.Item().Row(row =>
                            {
                                row.RelativeItem().Border(1).BorderColor("#E2E8F0").Background("#F8FAFC").Padding(10).Column(c =>
                                {
                                    c.Item().Text("$1,450,200.00").FontSize(14).Bold().FontColor("#1E3A8A");
                                    c.Item().Text("Ingresos Totales (+18.4%)").FontSize(8).FontColor("#6B7280");
                                });

                                row.Spacing(15);

                                row.RelativeItem().Border(1).BorderColor("#E2E8F0").Background("#F8FAFC").Padding(10).Column(c =>
                                {
                                    c.Item().Text("1,250").FontSize(14).Bold().FontColor("#1E3A8A");
                                    c.Item().Text("Transacciones (Meta superada)").FontSize(8).FontColor("#6B7280");
                                });
                            });

                            col.Item().PaddingTop(20).Table(table =>
                            {
                                table.ColumnsDefinition(columns =>
                                {
                                    columns.ConstantColumn(80);
                                    columns.RelativeColumn();
                                    columns.ConstantColumn(70);
                                    columns.ConstantColumn(90);
                                    columns.ConstantColumn(70);
                                });

                                table.Header(header =>
                                {
                                    header.Cell().Background("#1E3A8A").Padding(5).Text("Código").Bold().FontColor("#FFFFFF");
                                    header.Cell().Background("#1E3A8A").Padding(5).Text("Cliente").Bold().FontColor("#FFFFFF");
                                    header.Cell().Background("#1E3A8A").Padding(5).Text("Región").Bold().FontColor("#FFFFFF");
                                    header.Cell().Background("#1E3A8A").Padding(5).Text("Monto").Bold().FontColor("#FFFFFF");
                                    header.Cell().Background("#1E3A8A").Padding(5).Text("Estado").Bold().FontColor("#FFFFFF");
                                });

                                for (int r = 0; r < sampleRows.Count; r++)
                                {
                                    var rowData = sampleRows[r];
                                    var bg = r % 2 == 1 ? "#F8FAFC" : "#FFFFFF";

                                    table.Cell().Background(bg).BorderBottom(1).BorderColor("#E2E8F0").Padding(4).Text(rowData.Codigo);
                                    table.Cell().Background(bg).BorderBottom(1).BorderColor("#E2E8F0").Padding(4).Text(rowData.Cliente);
                                    table.Cell().Background(bg).BorderBottom(1).BorderColor("#E2E8F0").Padding(4).Text(rowData.Region);
                                    table.Cell().Background(bg).BorderBottom(1).BorderColor("#E2E8F0").Padding(4).Text(rowData.Monto);
                                    table.Cell().Background(bg).BorderBottom(1).BorderColor("#E2E8F0").Padding(4).Text(rowData.Estado);
                                }
                            });
                        });

                        page.Footer().AlignCenter().Text(x =>
                        {
                            x.Span("Página ");
                            x.CurrentPageNumber();
                            x.Span(" de ");
                            x.TotalPages();
                        });
                    });
                });
            }

            Console.WriteLine($"\n======================================================");
            Console.WriteLine($"🔬 BENCHMARK QUESTPDF .NET ({Iterations} iteraciones)");
            Console.WriteLine($"======================================================\n");

            // Warmup (JIT compilation)
            for (int i = 0; i < Warmup; i++)
            {
                CreateDoc().GeneratePdf();
            }

            var times = new List<double>();
            byte[] lastBytes = null;
            var swTotal = Stopwatch.StartNew();

            for (int i = 0; i < Iterations; i++)
            {
                var sw = Stopwatch.StartNew();
                lastBytes = CreateDoc().GeneratePdf();
                sw.Stop();
                times.Add(sw.Elapsed.TotalMilliseconds);
            }
            swTotal.Stop();

            times.Sort();
            double avg = times.Average();
            double min = times.Min();
            double max = times.Max();
            double p50 = times[(int)(times.Count * 0.5)];
            double p95 = times[(int)(times.Count * 0.95)];

            Console.WriteLine("--- RESULTADOS QUESTPDF (.NET 10 / SkiaSharp) ---");
            Console.WriteLine($"Total ({Iterations} docs): {swTotal.ElapsedMilliseconds} ms");
            Console.WriteLine($"Promedio por doc:   {avg:F2} ms");
            Console.WriteLine($"Mediana P50:        {p50:F2} ms");
            Console.WriteLine($"Percentil P95:      {p95:F2} ms");
            Console.WriteLine($"Mínimo:             {min:F2} ms");
            Console.WriteLine($"Máximo:             {max:F2} ms");
            Console.WriteLine($"Throughput:         {(1000.0 / avg):F0} doc/s");
            Console.WriteLine($"Tamaño PDF:         {(lastBytes?.Length ?? 0) / 1024.0:F1} KB\n");
        }
    }
}
