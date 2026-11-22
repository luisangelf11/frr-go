package pdf

import "github.com/luisangelf11/reporteria-go/pkg/reporter"

func init() {
	exp := NewExporter()
	reporter.RegisterPDFExporter(func(doc *reporter.Document) ([]byte, error) {
		return exp.RenderToBytes(doc)
	})
}
