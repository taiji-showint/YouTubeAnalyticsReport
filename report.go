package main

import (
	"fmt"
	"os"
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
	date = strings.ReplaceAll(date, "-", "")

	if yearly {
		report_filename = "reports/showint_report_yearly_" + date + ".md"
	} else {
		report_filename = "reports/showint_report_" + date + ".md"
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
}
