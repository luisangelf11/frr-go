package reporter

// PDFExporterFunc define la firma del exportador PDF.
type PDFExporterFunc func(doc *Document, outPath string) error

// GlobalPDFExporter almacena el exportador registrado para desacoplar paquetes.
var GlobalPDFExporter func(doc *Document) ([]byte, error)

// RegisterPDFExporter registra la implementación del exportador PDF.
func RegisterPDFExporter(fn func(doc *Document) ([]byte, error)) {
	GlobalPDFExporter = fn
}

// ExportPDF exporta el documento a un archivo PDF.
func (d *Document) ExportPDF(outPath string) error {
	if GlobalPDFExporter != nil {
		b, err := GlobalPDFExporter(d)
		if err != nil {
			return err
		}
		return writeToFile(outPath, b)
	}
	return nil
}

// ExportPDFBytes exporta el documento a un slice de bytes en formato PDF.
func (d *Document) ExportPDFBytes() ([]byte, error) {
	if GlobalPDFExporter != nil {
		return GlobalPDFExporter(d)
	}
	return nil, nil
}

func writeToFile(path string, data []byte) error {
	return writeBytes(path, data)
}
