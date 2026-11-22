package csvexp

import "github.com/luisangelf11/reporteria-go/pkg/reporter"

func init() {
	reporter.RegisterCSVExporter(func(doc *reporter.Document, delimiter rune, includeBOM bool) ([]byte, error) {
		exp := NewExporter(Config{
			Delimiter:  delimiter,
			UseCRLF:    true,
			IncludeBOM: includeBOM,
		})
		return exp.RenderToBytes(doc)
	})
}
