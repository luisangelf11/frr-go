export enum HeadingLevel {
  H1 = 1,
  H2 = 2,
  H3 = 3,
}

export enum ChartType {
  Bar = 'bar',
  Line = 'line',
  Pie = 'pie',
  Donut = 'donut',
}

export enum TextAlign {
  Left = 'left',
  Center = 'center',
  Right = 'right',
  Justify = 'justify',
}

export enum PageSize {
  A4 = 'A4',
  Letter = 'Letter',
  Legal = 'Legal',
}

export enum Orientation {
  Portrait = 'portrait',
  Landscape = 'landscape',
}

export enum Theme {
  CorporateBlue = 'corporate_blue',
  Emerald = 'emerald',
  DarkSlate = 'dark_slate',
}

export interface MetricCard {
  title: string;
  value: string;
  change?: string;
  isPositive?: boolean;
  description?: string;
}

export interface ChartSeries {
  name: string;
  color?: string;
  values: number[];
}

export interface ChartConfig {
  type: ChartType;
  title?: string;
  subtitle?: string;
  width?: number;
  height?: number;
  categories?: string[];
  series: ChartSeries[];
  showLegend?: boolean;
  showValues?: boolean;
  showGrid?: boolean;
  colors?: string[];
  donutRatio?: number;
  valuePrefix?: string;
  valueSuffix?: string;
}

export interface TableConfig {
  id?: string;
  title?: string;
  headers: string[];
  rows: string[][];
  columnWidths?: number[];
  alignments?: TextAlign[];
  zebraStriping?: boolean;
}

export interface ParagraphOptions {
  bold?: boolean;
  italic?: boolean;
  align?: TextAlign;
  fontSize?: number;
  textColor?: string;
}

export interface BlockSpec {
  type: string;
  heading?: { text: string; level: number };
  paragraph?: {
    text: string;
    bold?: boolean;
    italic?: boolean;
    align?: string;
    font_size?: number;
    text_color?: string;
  };
  table?: {
    id?: string;
    title?: string;
    headers: string[];
    rows: string[][];
    column_widths?: number[];
    alignments?: string[];
    zebra_striping?: boolean;
  };
  chart?: {
    type: string;
    title?: string;
    subtitle?: string;
    width?: number;
    height?: number;
    categories?: string[];
    series: { name: string; color?: string; values: number[] }[];
    show_legend?: boolean;
    show_values?: boolean;
    show_grid?: boolean;
    colors?: string[];
    donut_ratio?: number;
    value_prefix?: string;
    value_suffix?: string;
  };
  metric_cards?: {
    cards: {
      title: string;
      value: string;
      change?: string;
      is_positive: boolean;
      description?: string;
    }[];
  };
}

export interface ReportSpec {
  title: string;
  subtitle?: string;
  author?: string;
  date?: string;
  theme?: string;
  page_size?: string;
  orientation?: string;
  blocks: BlockSpec[];
}

export type ExportFormat = 'pdf' | 'excel' | 'xlsx' | 'csv';
