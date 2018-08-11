package main

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
)

const sampleTemplate string = `{
	"foo": "{{secret "foo"}}",
	"val1": 5,
	"val2": {{.ValueTwo}}
}`

func runTemplatingDemo1() {
	fmt.Println("## Templating Demo 1")

	var t *template.Template
	var err error

	t = template.New("sampleTemplate").Funcs(template.FuncMap{
		"secret": func(key string) (string, error) {
			if key == "foo" {
				return "BarSecretValue", nil
			}
			return "", fmt.Errorf("unknown secret key=%s", key)
		},
	})

	if t, err = t.Parse(sampleTemplate); err != nil {
		log.Fatalf("unable to parse a template, error=%v", err)
	}

	var buffer bytes.Buffer
	err = t.ExecuteTemplate(&buffer, "sampleTemplate", map[string]interface{}{
		"ValueTwo": 70,
	})
	if err != nil {
		fmt.Printf("Failed to execute template, err=%s", err)
		return
	}

	fmt.Printf("Executed Template = %s\n", buffer.String())
}

func main() {
	fmt.Println("# Plain Templates Demo")

	runTemplatingDemo1()
}
