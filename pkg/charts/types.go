package charts

import (
	"fmt"
	"strconv"
	"strings"
)

// ChartType define los tipos de gráficos disponibles.
type ChartType string

const (
	ChartTypeBar   ChartType = "bar"
	ChartTypeLine  ChartType = "line"
	ChartTypePie   ChartType = "pie"
	ChartTypeDonut ChartType = "donut"
)

// Color representa un color RGB.
type Color struct {
	R, G, B int
	A       float64 // 0.0 transparente a 1.0 opaco (por defecto 1.0)
}

// RGB crea un Color con alpha 1.0.
func RGB(r, g, b int) Color {
	return Color{R: r, G: g, B: b, A: 1.0}
}

// RGBA crea un Color con transparencia especificada (0.0 a 1.0).
func RGBA(r, g, b int, a float64) Color {
	return Color{R: r, G: g, B: b, A: a}
}

// HexColor convierte un string hexadecimal tipo "#1E3A8A" o "1E3A8A" a Color.
func HexColor(hex string) Color {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) == 3 {
		r, _ := strconv.ParseInt(string([]byte{hex[0], hex[0]}), 16, 64)
		g, _ := strconv.ParseInt(string([]byte{hex[1], hex[1]}), 16, 64)
		b, _ := strconv.ParseInt(string([]byte{hex[2], hex[2]}), 16, 64)
		return RGB(int(r), int(g), int(b))
	}
	if len(hex) == 6 {
		r, _ := strconv.ParseInt(hex[0:2], 16, 64)
		g, _ := strconv.ParseInt(hex[2:4], 16, 64)
		b, _ := strconv.ParseInt(hex[4:6], 16, 64)
		return RGB(int(r), int(g), int(b))
	}
	// Fallback por defecto: gris oscuro
	return RGB(50, 50, 50)
}

// Hex retorna la representación hexadecimal en formato "#RRGGBB".
func (c Color) Hex() string {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}

// Paleta corporativa moderna por defecto (accesible y de alto contraste).
var DefaultPalette = []Color{
	HexColor("#2563EB"), // Azul Real
	HexColor("#10B981"), // Esmeralda
	HexColor("#F59E0B"), // Ámbar / Dorado
	HexColor("#8B5CF6"), // Púrpura / Violeta
	HexColor("#EC4899"), // Rosa / Magenta
	HexColor("#06B6D4"), // Cian
	HexColor("#F97316"), // Naranja
	HexColor("#6366F1"), // Índigo
}

// Series representa un conjunto de datos asociados a una etiqueta o dimensión.
type Series struct {
	Name   string    `json:"name"`
	Color  string    `json:"color,omitempty"`
	Values []float64 `json:"values"`
}

// Config representa la configuración completa para renderizar un gráfico.
type Config struct {
	Type        ChartType `json:"type"`
	Title       string    `json:"title,omitempty"`
	Subtitle    string    `json:"subtitle,omitempty"`
	Width       float64   `json:"width"`      // Ancho en puntos o mm (ej: 500)
	Height      float64   `json:"height"`     // Alto en puntos o mm (ej: 300)
	Categories  []string  `json:"categories"` // Nombres de categorías para eje X (ej: ["Ene", "Feb", "Mar"])
	Series      []Series  `json:"series"`     // Datos numéricos por serie
	ShowLegend  bool      `json:"show_legend"`
	ShowValues  bool      `json:"show_values"`
	ShowGrid    bool      `json:"show_grid"`
	Colors      []string  `json:"colors,omitempty"` // Paleta opcional personalizada en hex
	DonutRatio  float64   `json:"donut_ratio"`      // Para Donut: relación radio interno (0.4 - 0.7)
	ValuePrefix string    `json:"value_prefix"`     // Ej: "$"
	ValueSuffix string    `json:"value_suffix"`     // Ej: "%" o " USD"
}

// DrawOpType define las operaciones vectoriales elementales.
type DrawOpType int

const (
	OpRect DrawOpType = iota
	OpLine
	OpCircle
	OpPolygon
	OpText
)

// TextAlign define la alineación del texto.
type TextAlign string

const (
	AlignLeft   TextAlign = "left"
	AlignCenter TextAlign = "center"
	AlignRight  TextAlign = "right"
)

// Point representa una coordenada 2D.
type Point struct {
	X, Y float64
}

// DrawOp es la unidad atómica de dibujo vectorial generada por el motor.
type DrawOp struct {
	Type        DrawOpType
	X, Y        float64     // Posición base
	Width       float64     // Para Rect
	Height      float64     // Para Rect
	X2, Y2      float64     // Para Line
	Radius      float64     // Para Circle
	Points      []Point     // Para Polygon
	FillColor   *Color      // Color de relleno (nil si no se rellena)
	StrokeColor *Color      // Color de contorno (nil si no se contornea)
	StrokeWidth float64     // Grosor de contorno
	Text        string      // Contenido para OpText
	FontSize    float64     // Tamaño tipográfico en puntos
	FontBold    bool        // Indica si la fuente es negrita
	TextAlign   TextAlign   // Alineación horizontal
}
