package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/jung-kurt/gofpdf"
)

func TestCreatePDF(t *testing.T) {
	CreatePDF()

}

// ExampleFpdf_CellFormat_tables demonstrates various table styles.
func TestExampleFpdf_CellFormat_tables(t *testing.T) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	type compareType struct {
		orgStr, entStr, policyStr, commentStr,
		statusBool string
	}
	compareList := make([]compareType, 0, 8)
	header := []string{"Organization", "Enterprise", "Policy", "Comment", "Status"}
	loadData := func(fileStr string) {
		fl, err := os.Open(fileStr)
		if err == nil {
			scanner := bufio.NewScanner(fl)
			var c compareType
			for scanner.Scan() {
				// Austria;Vienna;83859;8075
				lineStr := scanner.Text()
				list := strings.Split(lineStr, ";")
				if len(list) == 5 {
					c.orgStr = list[0]
					c.entStr = list[1]
					c.policyStr = list[2]
					c.commentStr = list[3]
					c.statusBool = list[4]
					compareList = append(compareList, c)
				} else {
					err = fmt.Errorf("error tokenizing %s", lineStr)
				}
			}
			fl.Close()
			if len(compareList) == 0 {
				err = fmt.Errorf("error loading data from %s", fileStr)
			}
		}
		if err != nil {
			pdf.SetError(err)
		}
	}
	// Simple table
	basicTable := func() {
		left := (210.0 - 4*40) / 2
		pdf.SetX(left)
		for _, str := range header {
			pdf.CellFormat(40, 7, str, "1", 0, "", false, 0, "")
		}
		pdf.Ln(-1)
		fill := true
		for _, c := range compareList {
			if c.statusBool == "true" {
				pdf.SetFillColor(255, 255, 0)
			} else {
				pdf.SetFillColor(255, 0, 255)
			}

			fmt.Println("statusBool: ", c.statusBool)

			pdf.SetX(left)
			pdf.CellFormat(40, 6, c.orgStr, "1", 0, "", fill, 0, "")
			pdf.CellFormat(40, 6, c.entStr, "1", 0, "", fill, 0, "")
			pdf.CellFormat(40, 6, c.policyStr, "1", 0, "", fill, 0, "")
			pdf.CellFormat(40, 6, c.commentStr, "1", 0, "", fill, 0, "")
			// pdf.CellFormat(40, 6, c.statusBool, "1", 0, "", false, 0, "")
			pdf.Ln(-1)
		}
	}
	// Colored table
	// fancyTable := func() {
	// 	// Colors, line width and bold font
	// 	pdf.SetFillColor(255, 0, 0)
	// 	pdf.SetTextColor(255, 255, 255)
	// 	pdf.SetDrawColor(128, 0, 0)
	// 	pdf.SetLineWidth(.3)
	// 	pdf.SetFont("", "B", 0)
	// 	// 	Header
	// 	w := []float64{40, 35, 40, 45}
	// 	wSum := 0.0
	// 	for _, v := range w {
	// 		wSum += v
	// 	}
	// 	left := (210 - wSum) / 2
	// 	pdf.SetX(left)
	// 	for j, str := range header {
	// 		pdf.CellFormat(w[j], 7, str, "1", 0, "C", true, 0, "")
	// 	}
	// 	pdf.Ln(-1)
	// 	// Color and font restoration
	// 	pdf.SetFillColor(224, 235, 255)
	// 	pdf.SetTextColor(0, 0, 0)
	// 	pdf.SetFont("", "", 0)
	// 	// 	Data
	// 	fill := true
	// 	for _, c := range compareList {
	// 		if c.statusBool == "true" {
	// 			pdf.SetFillColor(0, 255, 0)
	// 		} else {
	// 			pdf.SetFillColor(255, 0, 0)
	// 		}
	// 		pdf.SetX(left)
	// 		pdf.CellFormat(w[0], 6, c.orgStr, "LR", 0, "", fill, 0, "")
	// 		pdf.CellFormat(w[1], 6, c.entStr, "LR", 0, "", fill, 0, "")
	// 		// pdf.CellFormat(w[2], 6, strDelimit(c.policyStr, ",", 3),
	// 		// 	"LR", 0, "R", fill, 0, "")
	// 		// pdf.CellFormat(w[3], 6, strDelimit(c.commentStr, ",", 3),
	// 		// 	"LR", 0, "R", fill, 0, "")
	// 		pdf.Ln(-1)
	// 		// fill = !fill
	// 	}
	// 	pdf.SetX(left)
	// 	pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")
	// }
	loadData("compares.txt")
	pdf.SetFont("Arial", "", 14)
	pdf.AddPage()
	basicTable()
	// pdf.AddPage()
	// improvedTable()
	// pdf.AddPage()
	// fancyTable()
	// fileStr := example.Filename("Fpdf_CellFormat_tables")
	err := pdf.OutputFileAndClose("test.pdf")
	fmt.Println(err)
	// example.Summary(err, fileStr)
	// Output:
	// Successfully generated pdf/Fpdf_CellFormat_tables.pdf
}

// strDelimit converts 'ABCDEFG' to, for example, 'A,BCD,EFG'
func strDelimit(str string, sepstr string, sepcount int) string {
	pos := len(str) - sepcount
	for pos > 0 {
		str = str[:pos] + sepstr + str[pos:]
		pos = pos - sepcount
	}
	return str
}
