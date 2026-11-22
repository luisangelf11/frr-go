package charts

import (
	"fmt"
	"math"
)

// renderBarChart genera las operaciones de dibujo para un gráfico de barras.
func renderBarChart(cfg Config) ([]DrawOp, error) {
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

	// Márgenes del gráfico
	marginTop := 30.0
	marginBottom := 40.0
	marginLeft := 55.0
	marginRight := 25.0

	// Espacio para Título y Subtítulo
	if cfg.Title != "" {
		marginTop += 20
		titleColor := RGB(31, 41, 55) // #1F2937
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
			subColor := RGB(107, 114, 128) // #6B7280
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

	// Espacio para Leyenda (arriba a la derecha o debajo del título)
	if cfg.ShowLegend && len(cfg.Series) > 0 {
		legendY := 22.0
		if cfg.Title != "" {
			legendY = 22.0
		}
		legendOps := renderLegend(cfg.Series, palette, width-marginRight, legendY, AlignRight)
		ops = append(ops, legendOps...)
	}

	plotLeft := marginLeft
	plotRight := width - marginRight
	plotTop := marginTop
	plotBottom := height - marginBottom
	plotWidth := plotRight - plotLeft
	plotHeight := plotBottom - plotTop

	if plotWidth <= 0 || plotHeight <= 0 {
		return ops, fmt.Errorf("área de trazado inválida (dimensiones muy pequeñas)")
	}

	// Determinar valores min/max
	minVal, maxVal := getSeriesMinMax(cfg.Series)
	scale := CalculateNiceScale(minVal, maxVal, 5, true)

	// 1. Dibujar líneas de cuadrícula y etiquetas de eje Y
	gridColor := RGB(229, 231, 235) // #E5E7EB
	axisLabelColor := RGB(107, 114, 128)
	gridStroke := 0.6

	for _, tick := range scale.Ticks {
		valRatio := (tick - scale.Min) / (scale.Max - scale.Min)
		if valRatio < 0 || valRatio > 1.0001 {
			continue
		}
		yPos := plotBottom - (valRatio * plotHeight)

		// Línea de cuadrícula horizontal
		if cfg.ShowGrid {
			ops = append(ops, DrawOp{
				Type:        OpLine,
				X:           plotLeft,
				Y:           yPos,
				X2:          plotRight,
				Y2:          yPos,
				StrokeColor: &gridColor,
				StrokeWidth: gridStroke,
			})
		}

		// Etiqueta del eje Y
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

	// Eje X (línea base)
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

	// 2. Dibujar Barras por categoría
	numCategories := len(cfg.Categories)
	if numCategories == 0 && len(cfg.Series) > 0 {
		// Inferir de la primera serie
		numCategories = len(cfg.Series[0].Values)
	}

	if numCategories == 0 || len(cfg.Series) == 0 {
		return ops, nil
	}

	numSeries := len(cfg.Series)
	catSlotWidth := plotWidth / float64(numCategories)
	groupPadding := catSlotWidth * 0.25 // 25% de espacio entre grupos
	groupWidth := catSlotWidth - groupPadding
	barWidth := groupWidth / float64(numSeries)

	for catIdx := 0; catIdx < numCategories; catIdx++ {
		groupStart := plotLeft + (float64(catIdx) * catSlotWidth) + (groupPadding / 2.0)

		// Etiqueta de la categoría en eje X
		catName := ""
		if catIdx < len(cfg.Categories) {
			catName = cfg.Categories[catIdx]
		} else {
			catName = fmt.Sprintf("C%d", catIdx+1)
		}

		catCenterX := groupStart + (groupWidth / 2.0)
		ops = append(ops, DrawOp{
			Type:      OpText,
			X:         catCenterX,
			Y:         plotBottom + 16,
			Text:      catName,
			FontSize:  8.5,
			TextAlign: AlignCenter,
			FillColor: &axisLabelColor,
		})

		// Dibujar cada barra de la serie en esta categoría
		for serIdx, ser := range cfg.Series {
			val := 0.0
			if catIdx < len(ser.Values) {
				val = ser.Values[catIdx]
			}

			barX := groupStart + (float64(serIdx) * barWidth)
			barH := 0.0
			if scale.Max > scale.Min {
				barH = ((val - scale.Min) / (scale.Max - scale.Min)) * plotHeight
			}
			if barH < 0 {
				barH = 0
			}
			barY := plotBottom - barH

			color := palette[serIdx%len(palette)]
			if ser.Color != "" {
				color = HexColor(ser.Color)
			}

			// Barra
			ops = append(ops, DrawOp{
				Type:      OpRect,
				X:         barX + (barWidth * 0.05), // margen sutil interno
				Y:         barY,
				Width:     barWidth * 0.9,
				Height:    barH,
				FillColor: &color,
			})

			// Valor numérico sobre la barra
			if cfg.ShowValues && barH > 0 {
				valText := formatValue(val, cfg.ValuePrefix, cfg.ValueSuffix)
				ops = append(ops, DrawOp{
					Type:      OpText,
					X:         barX + (barWidth / 2.0),
					Y:         barY - 4,
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

// Helpers generales
func resolvePalette(cfg Config) []Color {
	if len(cfg.Colors) > 0 {
		var custom []Color
		for _, c := range cfg.Colors {
			custom = append(custom, HexColor(c))
		}
		return custom
	}
	return DefaultPalette
}

func getSeriesMinMax(series []Series) (float64, float64) {
	if len(series) == 0 {
		return 0, 100
	}
	minVal := math.MaxFloat64
	maxVal := -math.MaxFloat64

	hasData := false
	for _, s := range series {
		for _, v := range s.Values {
			hasData = true
			if v < minVal {
				minVal = v
			}
			if v > maxVal {
				maxVal = v
			}
		}
	}
	if !hasData {
		return 0, 100
	}
	return minVal, maxVal
}

func formatValue(val float64, prefix, suffix string) string {
	var s string
	if math.Abs(val-math.Round(val)) < 0.0001 {
		s = fmt.Sprintf("%.0f", val)
	} else {
		s = fmt.Sprintf("%.1f", val)
	}
	return prefix + s + suffix
}

func renderLegend(series []Series, palette []Color, startX, startY float64, align TextAlign) []DrawOp {
	var ops []DrawOp

	curX := startX
	fontSize := 8.5
	textColor := RGB(55, 65, 81)

	// Si alineamos a la derecha, calculamos de derecha a izquierda
	for i := len(series) - 1; i >= 0; i-- {
		s := series[i]
		color := palette[i%len(palette)]
		if s.Color != "" {
			color = HexColor(s.Color)
		}

		name := s.Name
		if name == "" {
			name = fmt.Sprintf("Serie %d", i+1)
		}

		textWidthApprox := float64(len(name)) * 5.0
		itemWidth := 14.0 + textWidthApprox + 10.0 // cuadro + texto + espacio

		rectX := curX - itemWidth
		textX := rectX + 14.0

		// Cuadradito de color
		ops = append(ops, DrawOp{
			Type:      OpRect,
			X:         rectX,
			Y:         startY - 7,
			Width:     9,
			Height:    9,
			FillColor: &color,
		})

		// Texto del nombre de la serie
		ops = append(ops, DrawOp{
			Type:      OpText,
			X:         textX,
			Y:         startY,
			Text:      name,
			FontSize:  fontSize,
			TextAlign: AlignLeft,
			FillColor: &textColor,
		})

		curX -= itemWidth
	}

	return ops
}
