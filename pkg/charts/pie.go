package charts

import (
	"fmt"
	"math"
)

// renderPieChart genera las operaciones de dibujo para gráficos de pastel y dona.
func renderPieChart(cfg Config, isDonut bool) ([]DrawOp, error) {
	var ops []DrawOp

	width := cfg.Width
	if width <= 0 {
		width = 500
	}
	height := cfg.Height
	if height <= 0 {
		height = 300
	}

	palette := resolvePalette(cfg)

	// Extraer valores y nombres
	type sliceData struct {
		name       string
		value      float64
		percentage float64
		color      Color
	}

	var slices []sliceData
	total := 0.0

	// Tomar la primera serie o combinar valores
	if len(cfg.Series) > 0 {
		s := cfg.Series[0]
		for i, v := range s.Values {
			if v < 0 {
				v = 0 // En pastel no se admiten valores negativos
			}
			name := fmt.Sprintf("Item %d", i+1)
			if i < len(cfg.Categories) && cfg.Categories[i] != "" {
				name = cfg.Categories[i]
			}
			c := palette[i%len(palette)]
			slices = append(slices, sliceData{
				name:  name,
				value: v,
				color: c,
			})
			total += v
		}
	}

	if total > 0 {
		for i := range slices {
			slices[i].percentage = (slices[i].value / total) * 100.0
		}
	}

	// Título y Subtítulo
	topOffset := 25.0
	if cfg.Title != "" {
		titleColor := RGB(31, 41, 55)
		ops = append(ops, DrawOp{
			Type:      OpText,
			X:         35,
			Y:         topOffset,
			Text:      cfg.Title,
			FontSize:  14,
			FontBold:  true,
			TextAlign: AlignLeft,
			FillColor: &titleColor,
		})
		topOffset += 18

		if cfg.Subtitle != "" {
			subColor := RGB(107, 114, 128)
			ops = append(ops, DrawOp{
				Type:      OpText,
				X:         35,
				Y:         topOffset,
				Text:      cfg.Subtitle,
				FontSize:  9,
				TextAlign: AlignLeft,
				FillColor: &subColor,
			})
			topOffset += 15
		}
	}

	// Disposición: Centro del pastel a la izquierda, Leyenda a la derecha
	availWidth := width - 70
	availHeight := height - topOffset - 25

	legendWidth := 160.0
	if !cfg.ShowLegend || len(slices) == 0 {
		legendWidth = 0
	}

	pieAreaWidth := availWidth - legendWidth
	cx := 35.0 + (pieAreaWidth / 2.0)
	cy := topOffset + (availHeight / 2.0)

	radius := math.Min(pieAreaWidth, availHeight) / 2.0 * 0.85
	if radius < 15 {
		radius = 15
	}

	innerRadius := 0.0
	if isDonut {
		ratio := cfg.DonutRatio
		if ratio <= 0.1 || ratio >= 0.9 {
			ratio = 0.55 // Ratio por defecto elegante
		}
		innerRadius = radius * ratio
	}

	whiteBorder := RGB(255, 255, 255)

	// Dibujar sectores
	startAngle := -math.Pi / 2.0 // Comenzar a las 12 en punto (-90°)
	for _, sl := range slices {
		if sl.value <= 0 || total <= 0 {
			continue
		}

		angleSpan := (sl.value / total) * 2.0 * math.Pi
		endAngle := startAngle + angleSpan

		// Aproximar arco mediante polígono de n pasos
		steps := int(math.Max(8, angleSpan*24)) // más pasos para arcos más amplios
		var polyPoints []Point

		if !isDonut {
			// Pastel tradicional: centro -> arco exterior -> centro
			polyPoints = append(polyPoints, Point{X: cx, Y: cy})
			for s := 0; s <= steps; s++ {
				curA := startAngle + (float64(s)/float64(steps))*angleSpan
				px := cx + radius*math.Cos(curA)
				py := cy + radius*math.Sin(curA)
				polyPoints = append(polyPoints, Point{X: px, Y: py})
			}
		} else {
			// Dona: arco exterior -> arco interior inverso
			for s := 0; s <= steps; s++ {
				curA := startAngle + (float64(s)/float64(steps))*angleSpan
				px := cx + radius*math.Cos(curA)
				py := cy + radius*math.Sin(curA)
				polyPoints = append(polyPoints, Point{X: px, Y: py})
			}
			for s := steps; s >= 0; s-- {
				curA := startAngle + (float64(s)/float64(steps))*angleSpan
				px := cx + innerRadius*math.Cos(curA)
				py := cy + innerRadius*math.Sin(curA)
				polyPoints = append(polyPoints, Point{X: px, Y: py})
			}
		}

		slColor := sl.color
		ops = append(ops, DrawOp{
			Type:        OpPolygon,
			Points:      polyPoints,
			FillColor:   &slColor,
			StrokeColor: &whiteBorder,
			StrokeWidth: 1.5,
		})

		// Etiquetas dentro del sector si ShowValues está activo y el sector es amplio (> 8%)
		if cfg.ShowValues && sl.percentage >= 8.0 {
			midA := startAngle + (angleSpan / 2.0)
			labelR := radius * 0.7
			if isDonut {
				labelR = (radius + innerRadius) / 2.0
			}
			lx := cx + labelR*math.Cos(midA)
			ly := cy + labelR*math.Sin(midA) + 3

			whiteText := RGB(255, 255, 255)
			ops = append(ops, DrawOp{
				Type:      OpText,
				X:         lx,
				Y:         ly,
				Text:      fmt.Sprintf("%.0f%%", sl.percentage),
				FontSize:  8.0,
				FontBold:  true,
				TextAlign: AlignCenter,
				FillColor: &whiteText,
			})
		}

		startAngle = endAngle
	}

	// Si es Dona, podemos colocar en el centro una etiqueta discreta de Total
	if isDonut && total > 0 {
		totalText := formatValue(total, cfg.ValuePrefix, cfg.ValueSuffix)
		labelColor := RGB(107, 114, 128)
		totalColor := RGB(31, 41, 55)

		ops = append(ops, DrawOp{
			Type:      OpText,
			X:         cx,
			Y:         cy - 2,
			Text:      "Total",
			FontSize:  7.5,
			TextAlign: AlignCenter,
			FillColor: &labelColor,
		})
		ops = append(ops, DrawOp{
			Type:      OpText,
			X:         cx,
			Y:         cy + 10,
			Text:      totalText,
			FontSize:  9.5,
			FontBold:  true,
			TextAlign: AlignCenter,
			FillColor: &totalColor,
		})
	}

	// Leyenda lateral derecha
	if cfg.ShowLegend && len(slices) > 0 {
		legX := width - legendWidth
		legY := cy - (float64(len(slices))*18.0)/2.0
		if legY < topOffset {
			legY = topOffset
		}

		nameColor := RGB(55, 65, 81)
		metaColor := RGB(107, 114, 128)

		for _, sl := range slices {
			slC := sl.color
			// Cuadradito de color
			ops = append(ops, DrawOp{
				Type:      OpRect,
				X:         legX,
				Y:         legY - 7,
				Width:     9,
				Height:    9,
				FillColor: &slC,
			})

			// Nombre de la categoría
			ops = append(ops, DrawOp{
				Type:      OpText,
				X:         legX + 15,
				Y:         legY,
				Text:      sl.name,
				FontSize:  8.5,
				TextAlign: AlignLeft,
				FillColor: &nameColor,
			})

			// Valor y porcentaje
			valStr := fmt.Sprintf("%s (%.1f%%)", formatValue(sl.value, cfg.ValuePrefix, cfg.ValueSuffix), sl.percentage)
			ops = append(ops, DrawOp{
				Type:      OpText,
				X:         legX + 15,
				Y:         legY + 9,
				Text:      valStr,
				FontSize:  7.5,
				TextAlign: AlignLeft,
				FillColor: &metaColor,
			})

			legY += 22
		}
	}

	return ops, nil
}
