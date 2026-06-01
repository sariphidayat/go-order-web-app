package web

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"
)

//go:embed templates/*.html templates/partials/*.html
var templateFiles embed.FS

type Renderer struct {
	pages map[string]*template.Template
}

func NewRenderer() (*Renderer, error) {
	pages := map[string]*template.Template{}
	for _, page := range []string{"login", "dashboard", "customers", "orders", "tracks"} {
		tmpl, err := template.New("").Funcs(funcMap()).ParseFS(
			templateFiles,
			"templates/base.html",
			"templates/"+page+".html",
			"templates/partials/*.html",
		)
		if err != nil {
			return nil, err
		}
		pages[page] = tmpl
	}
	return &Renderer{pages: pages}, nil
}

func (r *Renderer) Render(w http.ResponseWriter, status int, page string, data any) {
	tmpl, ok := r.pages[page]
	if !ok {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (r *Renderer) RenderPartial(w http.ResponseWriter, status int, page string, partial string, data any) {
	tmpl, ok := r.pages[page]
	if !ok {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := tmpl.ExecuteTemplate(w, partial, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func formatRupiah(value float64) string {
	s := fmt.Sprintf("%.2f", value)

	parts := strings.Split(s, ".")

	intPart := parts[0]
	decPart := parts[1]

	var result []string

	for len(intPart) > 3 {
		result = append([]string{intPart[len(intPart)-3:]}, result...)
		intPart = intPart[:len(intPart)-3]
	}

	if intPart != "" {
		result = append([]string{intPart}, result...)
	}

	return "Rp " + strings.Join(result, ".") + "," + decPart
}

func funcMap() template.FuncMap {
	return template.FuncMap{
		"currency": formatRupiah,
		// "currency": func(value float64) string {
		// 	return fmt.Sprintf("Rp %.2f", value)
		// },
		"dateTime": func(value time.Time) string {
			if value.IsZero() {
				return "-"
			}
			return value.Format("02 Jan 2006 15:04")
		},
		"dateTimeInput": func(value time.Time) string {
			if value.IsZero() {
				return time.Now().Format("2006-01-02T15:04")
			}
			return value.Format("2006-01-02T15:04")
		},
		"selected": func(actual any, expected any) template.HTMLAttr {
			if fmt.Sprint(actual) == fmt.Sprint(expected) {
				return "selected"
			}
			return ""
		},
	}
}
