package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tpl := template.New(patterns[0])
	tpl = tpl.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return `<input type="hidden" />`
		},
	})

	tpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %v", err)
	}

	return Template{tpl}, nil
}

// func Parse(filePath string) (Template, error) {
// 	tpl, err := template.ParseFiles(filePath)
// 	if err != nil {
// 		return Template{}, fmt.Errorf("parsing template: %v", err)
// 	}

// 	return Template{tpl}, nil
// }

type Template struct {
	htmlTemplate *template.Template
}

func (t Template) Execute(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := t.htmlTemplate.Execute(w, data)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing template", http.StatusInternalServerError)
		return
	}
}
