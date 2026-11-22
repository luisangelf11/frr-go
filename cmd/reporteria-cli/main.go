package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	workerMode := flag.Bool("worker", false, "Ejecutar en modo Persistent Worker por Stdio IPC para Node.js")
	format := flag.String("format", "pdf", "Formato de salida: pdf, excel, csv")
	specPath := flag.String("spec", "", "Ruta al archivo JSON de especificación del reporte (o '-' para stdin)")
	outPath := flag.String("out", "", "Ruta de archivo de salida (o vacío para stdout)")
	delimiter := flag.String("delimiter", ",", "Delimitador para exportaciones CSV (, o ;)")

	flag.Parse()

	if *workerMode {
		runWorkerLoop()
		return
	}

	if *specPath == "" {
		fmt.Fprintf(os.Stderr, "Uso: reporteria-cli [opciones]\n\n")
		fmt.Fprintf(os.Stderr, "Modos de uso:\n")
		fmt.Fprintf(os.Stderr, "  1. Modo Worker (para Node.js / servicios):\n")
		fmt.Fprintf(os.Stderr, "     reporteria-cli --worker\n\n")
		fmt.Fprintf(os.Stderr, "  2. Modo CLI One-Shot:\n")
		fmt.Fprintf(os.Stderr, "     reporteria-cli -format pdf -spec reporte.json -out salida.pdf\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := generateSingleReport(*format, *specPath, *outPath, *delimiter); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
