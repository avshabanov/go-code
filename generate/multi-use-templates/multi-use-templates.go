package main

import (
	"fmt"
	"html/template"
	"log"
	"strings"
)

type pageModel struct {
	Title   string
	Type    string
	Content interface{}
}

func (p *pageModel) MatchType(expectedType string) bool {
	return p.Type == expectedType
}

type demoContent struct {
	Description string
}

func main() {
	fmt.Println("Multi-Use Templates")

	t, err := template.ParseFiles(
		// base files (top-level ones)
		"templates/base/base1.html",
		"templates/base/base2.html",

		// page files (leaves)
		"templates/page/foo.html",
		"templates/page/bar.html",
	)
	if err != nil {
		log.Fatalf("unable to parse templates: %v", err)
	}

	b := strings.Builder{}

	model := pageModel{
		Title:   "Demo",
		Type:    "unknown",
		Content: &demoContent{Description: "Test"},
	}
	if err = t.ExecuteTemplate(&b, "base1", &model); err != nil {
		log.Fatalf("unable to execute base1 template: %v", err)
	}

	fmt.Println("base1 template unfolded:")
	fmt.Println(b.String())
	fmt.Println()

	model.Type = "foo"
	b.Reset()
	if err = t.ExecuteTemplate(&b, "base2", &model); err != nil {
		log.Fatalf("unable to execute base2 template: %v", err)
	}

	fmt.Println("================================")
	fmt.Println("base2 template (regular content) unfolded:")
	fmt.Println(b.String())
	fmt.Println()

	model.Type = "bar"
	b.Reset()
	if err = t.ExecuteTemplate(&b, "base2", &model); err != nil {
		log.Fatalf("unable to execute base2 template: %v", err)
	}

	fmt.Println("================================")
	fmt.Println("base2 template (bar content) unfolded:")
	fmt.Println(b.String())
	fmt.Println()
}
