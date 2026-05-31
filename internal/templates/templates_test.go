package templates

import "testing"

func TestGetTemplates(t *testing.T) {
	tmpl, err := GetTemplates("index.html")
	if err != nil {
		t.Fatalf("GetTemplates failed: %v", err)
	}
	t.Log(tmpl)
}
