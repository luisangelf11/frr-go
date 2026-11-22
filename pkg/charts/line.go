package charts

import (
	"fmt"
)

// renderLineChart genera las operaciones de dibujo para un gráfico de líneas y áreas.
func renderLineChart(cfg Config) ([]DrawOp, error) {
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

	marginTop := 30.0
	marginBottom := 40.0
	marginLeft := 55.0
	marginRight := 25.0

	// Título y Subtítulo
	if cfg.Title != "" {
		marginTop += 20
		titleColor := RGB(31, 41, 55)
		ops = append(ops, DrawOp{
			Type:      OpText,
			X:         marginLeft,
			Y:         25,
			Text:      cfg.Title,
			FontSize:  14,
			FontBold:  true,
			TextAlign: AlignLeft,
			FillColor: &titleColor,
		})

		if cfg.Subtitle != "" {
			marginTop += 14
			subColor := RGB(107, 114, 128)
			ops = append(ops, DrawOp{
				Type:      OpText,
				X:         marginLeft,
				Y:         40,
				Text:      cfg.Subtitle,
				FontSize:  9,
				TextAlign: AlignLeft,
				FillColor: &subColor,
			})
		}
	}

	// Leyenda
	if cfg.ShowLegend && len(cfg.Series) > 0 {
		legendOps := renderLegend(cfg.Series, palette, width-marginRight, 22.0, AlignRight)
		ops = append(ops, legendOps...)
	}

	plotLeft := marginLeft
	plotRight := width - marginRight
	plotTop := marginTop
	plotBottom := height - marginBottom
	plotWidth := plotRight - plotLeft
	plotHeight := plotBottom - plotTop

	if plotWidth <= 0 || plotHeight <= 0 {
		return ops, fmt.Errorf("área de trazado inválida")
	}

	minVal, maxVal := getSeriesMinMax(cfg.Series)
	// Para gráficos de línea podemos forzar o no cero según los datos
	scale := CalculateNiceScale(minVal, maxVal, 5, minVal >= 0)

	// Cuadrícula y eje Y
	gridColor := RGB(229, 231, 235)
	axisLabelColor := RGB(107, 114, 128)

	for _, tick := range scale.Ticks {
		valRatio := (tick - scale.Min) / (scale.Max - scale.Min)
		if valRatio < 0 || valRatio > 1.0001 {
			continue
		}
		yPos := plotBottom - (valRatio * plotHeight)

		if cfg.ShowGrid {
			ops = append(ops, DrawOp{
				Type:        OpLine,
				X:           plotLeft,
				Y:           yPos,
				X2:          plotRight,
				Y2:          yPos,
				StrokeColor: &gridColor,
				StrokeWidth: 0.6,
			})
		}

		labelText := formatValue(tick, cfg.ValuePrefix, cfg.ValueSuffix)
		ops = append(ops, DrawOp{
			Type:      OpText,
			X:         plotLeft - 8,
			Y:         yPos + 3,
			Text:      labelText,
			FontSize:  8,
			TextAlign: AlignRight,
			FillColor: &axisLabelColor,
		})
	}

	// Línea base eje X
	axisBaselineColor := RGB(156, 163, 175)
	ops = append(ops, DrawOp{
		Type:        OpLine,
		X:           plotLeft,
		Y:           plotBottom,
		X2:          plotRight,
		Y2:          plotBottom,
		StrokeColor: &axisBaselineColor,
		StrokeWidth: 1.0,
	})

	numCategories := len(cfg.Categories)
	if numCategories == 0 && len(cfg.Series) > 0 {
		numCategories = len(cfg.Series[0].Values)
	}

	if numCategories == 0 || len(cfg.Series) == 0 {
		return ops, nil
	}

	slotWidth := plotWidth / float64(numCategories)

	// Etiquetas del eje X
	for catIdx := 0; catIdx < numCategories; catIdx++ {
		catCenterX := plotLeft + (float64(catIdx) + 0.5)*slotWidth
		catName := ""
		if catIdx < len(cfg.Categories) {
			catName = cfg.Categories[catIdx]
		} else {
			catName = fmt.Sprintf("C%d", catIdx+1)
		}

		ops = append(ops, DrawOp{
			Type:      OpText,
			X:         catCenterX,
			Y:         plotBottom + 16,
			Text:      catName,
			FontSize:  8.5,
			TextAlign: AlignCenter,
			FillColor: &axisLabelColor,
		})
	}

	// Trazar cada serie de líneas
	for serIdx, ser := range cfg.Series {
		color := palette[serIdx%len(palette)]
		if ser.Color != "" {
			color = HexColor(ser.Color)
		}

		var points []Point
		for catIdx := 0; catIdx < numCategories; catIdx++ {
			val := 0.0
			if catIdx < len(ser.Values) {
				val = ser.Values[catIdx]
			}

			ptX := plotLeft + (float64(catIdx)+0.5)*slotWidth
			ptY := plotBottom
			if scale.Max > scale.Min {
				ratio := (val - scale.Min) / (scale.Max - scale.Min)
				ptY = plotBottom - (ratio * plotHeight)
			}
			points = append(points, Point{X: ptX, Y: ptY})
		}

		if len(points) == 0 {
			continue
		}

		// 1. Opcional: Relleno de área suave bajo la curva
		areaColor := RGBA(color.R, color.G, color.B, 0.12)
		var polyPoints []Point
		polyPoints = append(polyPoints, Point{X: points[0].X, Y: plotBottom})
		polyPoints = append(polyPoints, points...)
		polyPoints = append(polyPoints, Point{X: points[len(points)-1].X, Y: plotBottom})

		ops = append(ops, DrawOp{
			Type:      OpPolygon,
			Points:    polyPoints,
			FillColor: &areaColor,
		})

		// 2. Líneas conectando puntos
		for i := 0; i < len(points)-1; i++ {
			ops = append(ops, DrawOp{
				Type:        OpLine,
				X:           points[i].X,
				Y:           points[i].Y,
				X2:          points[i+1].X,
				Y2:          points[i+1].Y,
				StrokeColor: &color,
				StrokeWidth: 2.2,
			})
		}

		// 3. Puntos (círculos) y valores
		whiteColor := RGB(255, 255, 255)
		for i, pt := range points {
			// Círculo exterior (borde de color)
			ops = append(ops, DrawOp{
				Type:        OpCircle,
				X:           pt.X,
				Y:           pt.Y,
				Radius:      3.5,
				FillColor:   &whiteColor,
				StrokeColor: &color,
				StrokeWidth: 2.0,
			})

			if cfg.ShowValues && i < len(ser.Values) {
				valText := formatValue(ser.Values[i], cfg.ValuePrefix, cfg.ValueSuffix)
				ops = append(ops, DrawOp{
					Type:      OpText,
					X:         pt.X,
					Y:         pt.Y - 8,
					Text:      valText,
					FontSize:  7.5,
					TextAlign: AlignCenter,
					FillColor: &axisLabelColor,
				})
			}
		}
	}

	return ops, nil
}
