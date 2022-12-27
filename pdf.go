package main

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

func CreatePDF() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, world")
	err := pdf.OutputFileAndClose("hello.pdf")

	fmt.Println(err)
}
