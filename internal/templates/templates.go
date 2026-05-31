package templates

import (
	"fmt"
	"html/template"
	"path/filepath"
	"runtime"
)

func getTemplatesPath(templateName string) (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}

	helpersDir := filepath.Dir(filename)
	templatesPath := filepath.Join(helpersDir, "..", "..", "templates")
	templatesPath = filepath.Clean(templatesPath)
	templatesPath = filepath.Join(templatesPath, templateName)

	return templatesPath, nil
}

func GetTemplates(tmplName string) (*template.Template, error) {
	baseTmplPath, err := getTemplatesPath("base.html")
	if err != nil {
		return nil, fmt.Errorf("failed to get template path: %w", err)
	}

	tmplPath, err := getTemplatesPath(tmplName)
	if err != nil {
		return nil, fmt.Errorf("failed to get template path: %w", err)
	}

	base := template.Must(template.ParseFiles(baseTmplPath))
	tmpl := template.Must(template.Must(base.Clone()).ParseFiles(tmplPath))

	return tmpl, nil
}
