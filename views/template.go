package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"

	"github.com/AkifhanIlgaz/lenslocked/context"
	"github.com/AkifhanIlgaz/lenslocked/models"
	"github.com/gorilla/csrf"
)

type public interface {
	Public() string
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tpl := template.New(filepath.Base(patterns[0]))
	tpl = tpl.Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return `<!-- Please implement this function -->`, fmt.Errorf("csrfField isn't implemented")
		},
		"currentUser": func() (template.HTML, error) {
			return "", fmt.Errorf("currentUser not implemented")
		},
		"errors": func() []string {
			return nil
		},
	})

	tpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %v", err)
	}

	return Template{tpl}, nil
}

type Template struct {
	htmlTemplate *template.Template
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data any, errs ...error) {
	tpl, err := t.htmlTemplate.Clone()
	if err != nil {
		log.Printf("cloning template: %w", err)
		http.Error(w, "There was an error rendering the page", http.StatusInternalServerError)
		return
	}

	errMsgs := errMessages(errs...)
	tpl = tpl.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrf.TemplateField(r)
		},
		"currentUser": func() *models.User {
			return context.User(r.Context())
		},
		"errors": func() []string {
			return errMsgs
		},
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// You can't set HTTP status code twice in Go. So when tpl.Execute() is called status code is set to 200 OK
	// But if there is an error while executing the template we are trying to set the status code to 500 InternalServerError
	// So, when it happens we get 'superfluous response.WriteHEader call'
	// In order to prevent this error:
	// First we are copying the template to our buffer.
	// If there was an error while writing to buffer we return an error with HTTP 500 status code
	// Otherwise, return 200 OK with successfully executed template

	// Downside of this pattern is that, we are writing the template to memory
	// So, if you are executing a large html template you might run into memory issues
	var buff bytes.Buffer
	err = tpl.Execute(&buff, data)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing template", http.StatusInternalServerError)
		return
	}

	io.Copy(w, &buff)
	// w.Write(buff.Bytes())
}

func errMessages(errs ...error) []string {
	var messages []string
	for _, err := range errs {
		if pErr, ok := err.(public); ok {
			messages = append(messages, pErr.Public())
		} else {
			fmt.Println(err)
			messages = append(messages, "Something went wrong!")
		}
	}
	return messages
}
