package templates

import (
	"embed"
	"html/template"
	"log"
	"net/http"
)

//go:embed *.html
var templateFS embed.FS

var (
	templates *template.Template
	baseTmpl  *template.Template
)

func init() {
	var err error
	templates, err = template.ParseFS(templateFS, "*.html")
	baseTmpl = template.Must(template.ParseFS(templateFS, "base.html"))

	for _, t := range templates.Templates() {
		log.Println(t.Name())
	}

	if err != nil {
		log.Printf("template error: %v", err)
	}
}

func RenderTemplate(w http.ResponseWriter, tmplName string, data any) {
	tmpl := template.Must(template.Must(baseTmpl.Clone()).ParseFS(templateFS, tmplName))

	w.Header().Set("Content-Type", "text/html")

	err := tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("template error: %v", err)
		return
	}
}
