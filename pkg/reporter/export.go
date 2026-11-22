package reporter

// PDFExporterFunc define la firma del exportador PDF.
type PDFExporterFunc func(doc *Document, outPath string) error

// GlobalPDFExporter almacena el exportador registrado para desacoplar paquetes.
var GlobalPDFExporter func(doc *Document) ([]byte, error)

// GlobalCSVExporter almacena el exportador CSV registrado.
var GlobalCSVExporter func(doc *Document, delimiter rune, includeBOM bool) ([]byte, error)

// GlobalExcelExporter almacena el exportador Excel registrado.
var GlobalExcelExporter func(doc *Document) ([]byte, error)

// RegisterPDFExporter registra la implementación del exportador PDF.
func RegisterPDFExporter(fn func(doc *Document) ([]byte, error)) {
	GlobalPDFExporter = fn
}

// RegisterCSVExporter registra la implementación del exportador CSV.
func RegisterCSVExporter(fn func(doc *Document, delimiter rune, includeBOM bool) ([]byte, error)) {
	GlobalCSVExporter = fn
}

// RegisterExcelExporter registra la implementación del exportador Excel.
func RegisterExcelExporter(fn func(doc *Document) ([]byte, error)) {
	GlobalExcelExporter = fn
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

// ExportCSV exporta las tablas del documento a un archivo CSV.
func (d *Document) ExportCSV(outPath string, delimiter ...rune) error {
	b, err := d.ExportCSVBytes(delimiter...)
	if err != nil {
		return err
	}
	return writeToFile(outPath, b)
}

// ExportCSVBytes exporta las tablas del documento a un slice de bytes en formato CSV.
func (d *Document) ExportCSVBytes(delimiter ...rune) ([]byte, error) {
	delim := ','
	if len(delimiter) > 0 && delimiter[0] != 0 {
		delim = delimiter[0]
	}
	if GlobalCSVExporter != nil {
		return GlobalCSVExporter(d, delim, true)
	}
	return nil, nil
}

// ExportExcel exporta las tablas del documento a un archivo .xlsx de Excel.
func (d *Document) ExportExcel(outPath string) error {
	b, err := d.ExportExcelBytes()
	if err != nil {
		return err
	}
	return writeToFile(outPath, b)
}

// ExportExcelBytes exporta las tablas del documento a un slice de bytes en formato .xlsx.
func (d *Document) ExportExcelBytes() ([]byte, error) {
	if GlobalExcelExporter != nil {
		return GlobalExcelExporter(d)
	}
	return nil, nil
}

func writeToFile(path string, data []byte) error {
	return writeBytes(path, data)
}
