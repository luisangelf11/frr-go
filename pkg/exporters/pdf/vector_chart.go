package pdf

import (
	"github.com/go-pdf/fpdf"
	"github.com/luisangelf11/reporteria-go/pkg/charts"
)

// renderVectorChart dibuja las operaciones vectoriales de un gráfico directamente sobre el lienzo de fpdf.
func renderVectorChart(pdf *fpdf.Fpdf, tr func(string) string, ops []charts.DrawOp, cfg charts.Config, startX, startY, targetW, targetH float64) {
	origW := cfg.Width
	if origW <= 0 {
		origW = 500
	}
	origH := cfg.Height
	if origH <= 0 {
		origH = 300
	}

	scaleX := targetW / origW
	scaleY := targetH / origH

	for _, op := range ops {
		x := startX + (op.X * scaleX)
		y := startY + (op.Y * scaleY)

		switch op.Type {
		case charts.OpRect:
			w := op.Width * scaleX
			h := op.Height * scaleY
			style := ""

			if op.FillColor != nil {
				pdf.SetFillColor(op.FillColor.R, op.FillColor.G, op.FillColor.B)
				style += "F"
			}
			if op.StrokeColor != nil && op.StrokeWidth > 0 {
				pdf.SetDrawColor(op.StrokeColor.R, op.StrokeColor.G, op.StrokeColor.B)
				pdf.SetLineWidth(op.StrokeWidth * 0.25)
				style += "D"
			}
			if style != "" {
				pdf.Rect(x, y, w, h, style)
			}

		case charts.OpLine:
			x2 := startX + (op.X2 * scaleX)
			y2 := startY + (op.Y2 * scaleY)
			if op.StrokeColor != nil {
				pdf.SetDrawColor(op.StrokeColor.R, op.StrokeColor.G, op.StrokeColor.B)
			} else {
				pdf.SetDrawColor(150, 150, 150)
			}
			lineWidth := 0.2
			if op.StrokeWidth > 0 {
				lineWidth = op.StrokeWidth * 0.25
			}
			pdf.SetLineWidth(lineWidth)
			pdf.Line(x, y, x2, y2)

		case charts.OpCircle:
			r := op.Radius * scaleX
			style := ""
			if op.FillColor != nil {
				pdf.SetFillColor(op.FillColor.R, op.FillColor.G, op.FillColor.B)
				style += "F"
			}
			if op.StrokeColor != nil && op.StrokeWidth > 0 {
				pdf.SetDrawColor(op.StrokeColor.R, op.StrokeColor.G, op.StrokeColor.B)
				pdf.SetLineWidth(op.StrokeWidth * 0.25)
				style += "D"
			}
			if style != "" {
				pdf.Circle(x, y, r, style)
			}

		case charts.OpPolygon:
			if len(op.Points) == 0 {
				continue
			}
			fpdfPts := make([]fpdf.PointType, len(op.Points))
			for i, p := range op.Points {
				fpdfPts[i] = fpdf.PointType{
					X: startX + (p.X * scaleX),
					Y: startY + (p.Y * scaleY),
				}
			}
			style := ""
			if op.FillColor != nil {
				pdf.SetFillColor(op.FillColor.R, op.FillColor.G, op.FillColor.B)
				style += "F"
			}
			if op.StrokeColor != nil && op.StrokeWidth > 0 {
				pdf.SetDrawColor(op.StrokeColor.R, op.StrokeColor.G, op.StrokeColor.B)
				pdf.SetLineWidth(op.StrokeWidth * 0.25)
				style += "D"
			}
			if style != "" {
				pdf.Polygon(fpdfPts, style)
			}

		case charts.OpText:
			if op.FillColor != nil {
				pdf.SetTextColor(op.FillColor.R, op.FillColor.G, op.FillColor.B)
			} else {
				pdf.SetTextColor(31, 41, 55)
			}
			fontStyle := ""
			if op.FontBold {
				fontStyle = "B"
			}
			fontSize := 8.0
			if op.FontSize > 0 {
				fontSize = op.FontSize * 0.75
			}
			pdf.SetFont("Arial", fontStyle, fontSize)

			translatedText := tr(op.Text)
			textX := x
			textW := pdf.GetStringWidth(translatedText)
			if op.TextAlign == charts.AlignCenter {
				textX = x - (textW / 2.0)
			} else if op.TextAlign == charts.AlignRight {
				textX = x - textW
			}

			pdf.Text(textX, y, translatedText)
		}
	}
}
