package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/luisangelf11/reporteria-go/pkg/exporters/csv"
	"github.com/luisangelf11/reporteria-go/pkg/exporters/excel"
	"github.com/luisangelf11/reporteria-go/pkg/exporters/pdf"
	"github.com/luisangelf11/reporteria-go/pkg/reporter"
	"github.com/luisangelf11/reporteria-go/pkg/spec"
)

// runWorkerLoop ejecuta el bucle de servicio persistente por Stdio IPC.
func runWorkerLoop() {
	pdfExp := pdf.NewExporter()
	excelExp := excel.NewExporter()

	lenBuf := make([]byte, 4)
	outHeader := make([]byte, 5) // 4 bytes longitud + 1 byte estado

	for {
		// 1. Leer longitud del mensaje entrante (4 bytes Big-Endian)
		if _, err := io.ReadFull(os.Stdin, lenBuf); err != nil {
			if err == io.EOF {
				return // El proceso padre (Node.js) cerró stdin, terminar limpiamente
			}
			return
		}

		msgLen := binary.BigEndian.Uint32(lenBuf)
		if msgLen == 0 || msgLen > 50*1024*1024 { // Límite de seguridad: 50MB
			sendWorkerResponse(1, []byte("longitud de payload inválida"), outHeader)
			continue
		}

		// 2. Leer payload JSON
		payload := make([]byte, msgLen)
		if _, err := io.ReadFull(os.Stdin, payload); err != nil {
			sendWorkerResponse(1, []byte(fmt.Sprintf("error leyendo payload: %v", err)), outHeader)
			continue
		}

		// 3. Procesar petición
		var req spec.RequestSpec
		if err := json.Unmarshal(payload, &req); err != nil {
			sendWorkerResponse(1, []byte(fmt.Sprintf("error deserializando JSON: %v", err)), outHeader)
			continue
		}

		doc := spec.ConvertSpecToDoc(req.Spec)
		var result []byte
		var err error

		switch req.Format {
		case "pdf":
			result, err = pdfExp.RenderToBytes(doc)
		case "excel", "xlsx":
			result, err = excelExp.RenderToBytes(doc)
		case "csv":
			delim := ','
			if len(req.Delimiter) > 0 {
				delim = rune(req.Delimiter[0])
			}
			csvExp := csvexp.NewExporter(csvexp.Config{
				Delimiter:  delim,
				UseCRLF:    true,
				IncludeBOM: true,
			})
			result, err = csvExp.RenderToBytes(doc)
		default:
			err = fmt.Errorf("formato no soportado: %s", req.Format)
		}

		if err != nil {
			sendWorkerResponse(1, []byte(err.Error()), outHeader)
		} else {
			sendWorkerResponse(0, result, outHeader)
		}
	}
}

// sendWorkerResponse escribe [4 bytes longitud] + [1 byte status] + [datos] en os.Stdout.
func sendWorkerResponse(status byte, data []byte, headerBuf []byte) {
	binary.BigEndian.PutUint32(headerBuf[:4], uint32(len(data)))
	headerBuf[4] = status

	_, _ = os.Stdout.Write(headerBuf)
	_, _ = os.Stdout.Write(data)
	_ = os.Stdout.Sync()
}

// generateSingleReport genera un reporte en modo one-shot (CLI clásico).
func generateSingleReport(format, specPath, outPath, delimStr string) error {
	var specData []byte
	var err error

	if specPath == "-" {
		specData, err = io.ReadAll(os.Stdin)
	} else {
		specData, err = os.ReadFile(specPath)
	}
	if err != nil {
		return fmt.Errorf("error leyendo archivo spec: %w", err)
	}

	doc, err := spec.ParseSpecJSON(specData)
	if err != nil {
		return err
	}

	var result []byte
	switch format {
	case "pdf":
		exp := pdf.NewExporter()
		result, err = exp.RenderToBytes(doc)
	case "excel", "xlsx":
		exp := excel.NewExporter()
		result, err = exp.RenderToBytes(doc)
	case "csv":
		delim := ','
		if len(delimStr) > 0 {
			delim = rune(delimStr[0])
		}
		exp := csvexp.NewExporter(csvexp.Config{
			Delimiter:  delim,
			UseCRLF:    true,
			IncludeBOM: true,
		})
		result, err = exp.RenderToBytes(doc)
	default:
		return fmt.Errorf("formato desconocido: %s", format)
	}

	if err != nil {
		return err
	}

	if outPath == "" || outPath == "-" {
		_, err = os.Stdout.Write(result)
		return err
	}

	return os.WriteFile(outPath, result, 0644)
}

// Helper para silenciar imports no usados
var _ = reporter.H1
