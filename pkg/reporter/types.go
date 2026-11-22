package reporter

import "github.com/luisangelf11/frr-go/pkg/charts"

// PageSize representa el tamaño de papel físico.
type PageSize string

const (
	PageSizeA4     PageSize = "A4"
	PageSizeLetter PageSize = "Letter"
	PageSizeLegal  PageSize = "Legal"
)

// Orientation representa la orientación de página.
type Orientation string

const (
	OrientationPortrait  Orientation = "Portrait"
	OrientationLandscape Orientation = "Landscape"
)

// TextAlign representa la alineación del texto.
type TextAlign string

const (
	AlignLeft    TextAlign = "left"
	AlignCenter  TextAlign = "center"
	AlignRight   TextAlign = "right"
	AlignJustify TextAlign = "justify"
)

// HeadingLevel representa el nivel de jerarquía de un encabezado.
type HeadingLevel int

const (
	H1 HeadingLevel = 1
	H2 HeadingLevel = 2
	H3 HeadingLevel = 3
)

// Theme define la paleta corporativa y estilos visuales del reporte.
type Theme struct {
	Name              string
	PrimaryColor      charts.Color
	SecondaryColor    charts.Color
	HeaderBgColor     charts.Color
	HeaderTextColor   charts.Color
	TextColor         charts.Color
	MutedTextColor    charts.Color
	AlternateRowColor charts.Color
	BorderColor       charts.Color
}

// Temas predefinidos modernos y elegantes.
var (
	ThemeCorporateBlue = Theme{
		Name:              "Corporate Blue",
		PrimaryColor:      charts.HexColor("#1E40AF"),
		SecondaryColor:    charts.HexColor("#3B82F6"),
		HeaderBgColor:     charts.HexColor("#1E3A8A"),
		HeaderTextColor:   charts.HexColor("#FFFFFF"),
		TextColor:         charts.HexColor("#1F2937"),
		MutedTextColor:    charts.HexColor("#6B7280"),
		AlternateRowColor: charts.HexColor("#F8FAFC"),
		BorderColor:       charts.HexColor("#E2E8F0"),
	}

	ThemeEmerald = Theme{
		Name:              "Emerald",
		PrimaryColor:      charts.HexColor("#065F46"),
		SecondaryColor:    charts.HexColor("#10B981"),
		HeaderBgColor:     charts.HexColor("#047857"),
		HeaderTextColor:   charts.HexColor("#FFFFFF"),
		TextColor:         charts.HexColor("#1F2937"),
		MutedTextColor:    charts.HexColor("#6B7280"),
		AlternateRowColor: charts.HexColor("#F0FDF4"),
		BorderColor:       charts.HexColor("#E5E7EB"),
	}

	ThemeDarkSlate = Theme{
		Name:              "Dark Slate",
		PrimaryColor:      charts.HexColor("#0F172A"),
		SecondaryColor:    charts.HexColor("#334155"),
		HeaderBgColor:     charts.HexColor("#1E293B"),
		HeaderTextColor:   charts.HexColor("#FFFFFF"),
		TextColor:         charts.HexColor("#0F172A"),
		MutedTextColor:    charts.HexColor("#64748B"),
		AlternateRowColor: charts.HexColor("#F8FAFC"),
		BorderColor:       charts.HexColor("#CBD5E1"),
	}
)
