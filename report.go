package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"text/template"
	"strings"
)

func render_report(video_list []Video, date string, yearly bool) {

	var template_filename string
	if yearly {
		template_filename = "templates/tmpl_report_yearly.md"
	} else {
		template_filename = "templates/tmpl_report.md"
	}

	tmpl, err := template.ParseFiles(template_filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := tmpl.Execute(os.Stdout, video_list); err != nil {
		panic(err)
	}

	var report_filename string
	date_clean := strings.ReplaceAll(date, "-", "")

	if yearly {
		report_filename = "reports/showint_report_yearly_" + date_clean + ".md"
	} else {
		report_filename = "reports/showint_report_" + date_clean + ".md"
	}

	f, err := os.Create(report_filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = tmpl.Execute(f, video_list)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Generate JSON for PDF report generation
	if !yearly {
		generate_json_for_pdf(video_list, date_clean)
		generate_pdf_report(video_list, date_clean)
	}
}

// generate_json_for_pdf exports video list as JSON for Python PDF generation script
func generate_json_for_pdf(video_list []Video, date string) {
	json_filename := "reports/report_data_" + date + ".json"

	data, err := json.MarshalIndent(video_list, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	err = os.WriteFile(json_filename, data, 0644)
	if err != nil {
		fmt.Printf("Error writing JSON file: %v\n", err)
		return
	}

	fmt.Printf("JSON report data saved: %s\n", json_filename)
}

// generate_pdf_report calls Python script to generate 2-column PDF
func generate_pdf_report(video_list []Video, date string) {
	// Ensure reports_pdf directory exists
	pdfDir := "reports_pdf"
	if _, err := os.Stat(pdfDir); os.IsNotExist(err) {
		os.MkdirAll(pdfDir, 0755)
	}

	json_filename := "reports/report_data_" + date + ".json"
	pdf_filename := pdfDir + "/showint_report_" + date + ".pdf"

	// Prefer Homebrew Python over virtualenv
	pythonPath := "python3"

	// Try Homebrew location first (macOS)
	homebrewPath := "/opt/homebrew/bin/python3"
	if _, err := os.Stat(homebrewPath); err == nil {
		pythonPath = homebrewPath
	} else if path, err := exec.LookPath("python3"); err == nil {
		pythonPath = path
	}

	// Call Python script: python3 generate_pdf.py <json_input> <pdf_output>
	cmd := exec.Command(pythonPath, "generate_pdf.py", json_filename, pdf_filename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Generating PDF report using: %s\n", pythonPath)
	fmt.Printf("Generating PDF report: %s\n", pdf_filename)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Warning: PDF generation failed: %v\n", err)
		fmt.Printf("Python path used: %s\n", pythonPath)
		fmt.Println("Make sure generate_pdf.py is in the current directory and reportlab is installed")
		return
	}

	fmt.Printf("PDF report generated successfully: %s\n", pdf_filename)
}
