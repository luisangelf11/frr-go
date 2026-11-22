package charts

import (
	"fmt"
	"html"
	"strings"
)

// RenderSVG convierte una lista de DrawOp en un documento SVG estándar en formato string.
func RenderSVG(ops []DrawOp, width, height float64) string {
	if width <= 0 {
		width = 500
	}
	if height <= 0 {
		height = 300
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %.1f %.1f" width="%.1f" height="%.1f" style="font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;">`+"\n",
		width, height, width, height,
	))

	// Fondo blanco limpio por defecto
	sb.WriteString(fmt.Sprintf(
		`  <rect width="100%%" height="100%%" fill="#FFFFFF"/>`+"\n",
	))

	for _, op := range ops {
		switch op.Type {
		case OpRect:
			fill := formatSVGColor(op.FillColor, "none")
			stroke := formatSVGColor(op.StrokeColor, "none")
			strokeAttr := ""
			if stroke != "none" {
				strokeAttr = fmt.Sprintf(` stroke="%s" stroke-width="%.2f"`, stroke, op.StrokeWidth)
			}
			sb.WriteString(fmt.Sprintf(
				`  <rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" fill="%s"%s rx="2" ry="2"/>`+"\n",
				op.X, op.Y, op.Width, op.Height, fill, strokeAttr,
			))

		case OpLine:
			stroke := formatSVGColor(op.StrokeColor, "#000000")
			widthVal := op.StrokeWidth
			if widthVal <= 0 {
				widthVal = 1.0
			}
			sb.WriteString(fmt.Sprintf(
				`  <line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="%.2f" stroke-linecap="round"/>`+"\n",
				op.X, op.Y, op.X2, op.Y2, stroke, widthVal,
			))

		case OpCircle:
			fill := formatSVGColor(op.FillColor, "none")
			stroke := formatSVGColor(op.StrokeColor, "none")
			strokeAttr := ""
			if stroke != "none" {
				strokeAttr = fmt.Sprintf(` stroke="%s" stroke-width="%.2f"`, stroke, op.StrokeWidth)
			}
			sb.WriteString(fmt.Sprintf(
				`  <circle cx="%.2f" cy="%.2f" r="%.2f" fill="%s"%s/>`+"\n",
				op.X, op.Y, op.Radius, fill, strokeAttr,
			))

		case OpPolygon:
			if len(op.Points) == 0 {
				continue
			}
			var pts []string
			for _, p := range op.Points {
				pts = append(pts, fmt.Sprintf("%.2f,%.2f", p.X, p.Y))
			}
			fill := formatSVGColor(op.FillColor, "none")
			stroke := formatSVGColor(op.StrokeColor, "none")
			strokeAttr := ""
			if stroke != "none" {
				strokeAttr = fmt.Sprintf(` stroke="%s" stroke-width="%.2f"`, stroke, op.StrokeWidth)
			}
			sb.WriteString(fmt.Sprintf(
				`  <polygon points="%s" fill="%s"%s/>`+"\n",
				strings.Join(pts, " "), fill, strokeAttr,
			))

		case OpText:
			fill := formatSVGColor(op.FillColor, "#1F2937")
			anchor := "start"
			if op.TextAlign == AlignCenter {
				anchor = "middle"
			} else if op.TextAlign == AlignRight {
				anchor = "end"
			}
			weight := "normal"
			if op.FontBold {
				weight = "bold"
			}
			fontSize := op.FontSize
			if fontSize <= 0 {
				fontSize = 9.0
			}
			escapedText := html.EscapeString(op.Text)
			sb.WriteString(fmt.Sprintf(
				`  <text x="%.2f" y="%.2f" font-size="%.2f" font-weight="%s" text-anchor="%s" fill="%s">%s</text>`+"\n",
				op.X, op.Y, fontSize, weight, anchor, fill, escapedText,
			))
		}
	}

	sb.WriteString("</svg>\n")
	return sb.String()
}

func formatSVGColor(c *Color, fallback string) string {
	if c == nil {
		return fallback
	}
	if c.A < 0.999 {
		return fmt.Sprintf("rgba(%d,%d,%d,%.2f)", c.R, c.G, c.B, c.A)
	}
	return c.Hex()
}
