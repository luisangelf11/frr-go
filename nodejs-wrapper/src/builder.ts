import * as fs from 'fs';
import {
  ReportSpec,
  HeadingLevel,
  ParagraphOptions,
  TableConfig,
  ChartConfig,
  MetricCard,
  Theme,
  PageSize,
  Orientation,
  ExportFormat,
} from './types';
import { GoWorkerManager } from './worker';

export class ReportBuilder {
  private spec: ReportSpec;

  constructor(title: string = '') {
    this.spec = {
      title,
      date: new Date().toISOString().split('T')[0],
      theme: Theme.CorporateBlue,
      page_size: PageSize.A4,
      orientation: Orientation.Portrait,
      blocks: [],
    };
  }

  public setTitle(title: string): this {
    this.spec.title = title;
    return this;
  }

  public setSubtitle(sub: string): this {
    this.spec.subtitle = sub;
    return this;
  }

  public setAuthor(author: string): this {
    this.spec.author = author;
    return this;
  }

  public setDate(date: string): this {
    this.spec.date = date;
    return this;
  }

  public setTheme(theme: Theme): this {
    this.spec.theme = theme;
    return this;
  }

  public setPageSize(size: PageSize): this {
    this.spec.page_size = size;
    return this;
  }

  public setOrientation(orientation: Orientation): this {
    this.spec.orientation = orientation;
    return this;
  }

  public addHeading(text: string, level: HeadingLevel = HeadingLevel.H1): this {
    this.spec.blocks.push({
      type: 'heading',
      heading: { text, level },
    });
    return this;
  }

  public addParagraph(text: string, options?: ParagraphOptions): this {
    this.spec.blocks.push({
      type: 'paragraph',
      paragraph: {
        text,
        bold: options?.bold,
        italic: options?.italic,
        align: options?.align,
        font_size: options?.fontSize,
        text_color: options?.textColor,
      },
    });
    return this;
  }

  public addMetricCards(...cards: MetricCard[]): this {
    this.spec.blocks.push({
      type: 'metric_cards',
      metric_cards: {
        cards: cards.map((c) => ({
          title: c.title,
          value: c.value,
          change: c.change,
          is_positive: c.isPositive ?? true,
          description: c.description,
        })),
      },
    });
    return this;
  }

  public addKpiCard(card: MetricCard): this {
    return this.addMetricCards(card);
  }

  public addKpi(title: string, value: string, change?: string, isPositive: boolean = true, description?: string): this {
    return this.addMetricCards({ title, value, change, isPositive, description });
  }

  public addTable(tableOrHeaders: TableConfig | string[], maybeRows?: string[][]): this {
    if (Array.isArray(tableOrHeaders)) {
      return this.addTable({
        headers: tableOrHeaders,
        rows: maybeRows || [],
      });
    }
    const table = tableOrHeaders;
    this.spec.blocks.push({
      type: 'table',
      table: {
        id: table.id,
        title: table.title,
        headers: table.headers,
        rows: table.rows,
        column_widths: table.columnWidths,
        alignments: table.alignments,
        zebra_striping: table.zebraStriping ?? true,
      },
    });
    return this;
  }

  public addChart(chart: ChartConfig): this {
    this.spec.blocks.push({
      type: 'chart',
      chart: {
        type: chart.type,
        title: chart.title,
        subtitle: chart.subtitle,
        width: chart.width,
        height: chart.height,
        categories: chart.categories,
        series: chart.series,
        show_legend: chart.showLegend,
        show_values: chart.showValues,
        show_grid: chart.showGrid,
        colors: chart.colors,
        donut_ratio: chart.donutRatio,
        value_prefix: chart.valuePrefix,
        value_suffix: chart.valueSuffix,
      },
    });
    return this;
  }

  public addDivider(): this {
    this.spec.blocks.push({ type: 'divider' });
    return this;
  }

  public addPageBreak(): this {
    this.spec.blocks.push({ type: 'page_break' });
    return this;
  }

  public toJSON(): ReportSpec {
    return JSON.parse(JSON.stringify(this.spec));
  }

  public async exportToBuffer(format: ExportFormat, delimiter?: string): Promise<Buffer> {
    const worker = GoWorkerManager.getInstance();
    return worker.execute(this.spec, format, delimiter);
  }

  public async exportToFile(format: ExportFormat, filePath: string, delimiter?: string): Promise<void> {
    const buf = await this.exportToBuffer(format, delimiter);
    await fs.promises.writeFile(filePath, buf);
  }
}
