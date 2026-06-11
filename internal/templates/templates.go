package templates

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

func getTemplatesPath(templateName string) (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get current file path")
	}

	exeDir := filepath.Dir(exe)
	templatesPath := filepath.Join(exeDir, "templates")
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
