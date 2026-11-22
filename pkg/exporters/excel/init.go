package excel

import "github.com/luisangelf11/frr-go/pkg/reporter"

func init() {
	reporter.RegisterExcelExporter(func(doc *reporter.Document) ([]byte, error) {
		exp := NewExporter()
		return exp.RenderToBytes(doc)
	})
}
