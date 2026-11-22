package charts

import "math"

// ScaleInfo contiene los valores calculados para la graduación de un eje.
type ScaleInfo struct {
	Min   float64
	Max   float64
	Step  float64
	Ticks []float64
}

// CalculateNiceScale calcula un rango y pasos con números redondos (1, 2, 5, 10, etc.).
func CalculateNiceScale(minVal, maxVal float64, targetTicks int, forceZeroMin bool) ScaleInfo {
	if forceZeroMin && minVal > 0 {
		minVal = 0
	}

	// Manejo de caso donde todos los valores son iguales o cero
	if minVal == maxVal {
		if minVal == 0 {
			maxVal = 10
		} else if minVal > 0 {
			minVal = 0
			maxVal = maxVal * 1.2
		} else {
			maxVal = 0
			minVal = minVal * 1.2
		}
	}

	if targetTicks < 2 {
		targetTicks = 5
	}

	rawRange := maxVal - minVal
	rawStep := rawRange / float64(targetTicks)

	exponent := math.Floor(math.Log10(rawStep))
	fraction := rawStep / math.Pow(10, exponent)

	var niceFraction float64
	if fraction <= 1.0 {
		niceFraction = 1.0
	} else if fraction <= 2.0 {
		niceFraction = 2.0
	} else if fraction <= 5.0 {
		niceFraction = 5.0
	} else {
		niceFraction = 10.0
	}

	niceStep := niceFraction * math.Pow(10, exponent)
	niceMin := math.Floor(minVal/niceStep) * niceStep
	niceMax := math.Ceil(maxVal/niceStep) * niceStep

	if forceZeroMin && minVal >= 0 {
		niceMin = 0
	}

	// Generar marcas (ticks)
	var ticks []float64
	// Protección contra bucles infinitos con epsilon
	epsilon := niceStep * 0.0001
	for t := niceMin; t <= niceMax+epsilon; t += niceStep {
		// Redondeo de precisión para evitar 0.30000000000000004
		rounded := math.Round(t/niceStep) * niceStep
		ticks = append(ticks, rounded)
	}

	return ScaleInfo{
		Min:   niceMin,
		Max:   niceMax,
		Step:  niceStep,
		Ticks: ticks,
	}
}
