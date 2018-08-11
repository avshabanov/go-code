package main

import (
	"bytes"
	"fmt"
	ht "html/template"
)

type user struct {
	ID    string
	Name  string
	Age   int
	Roles []string
}

const userTemplate string = `
<ul>
	<li>User ID:		{{.ID}}</li>
	<li>User Name: 	{{.Name}}</li>
	<li>User Age: 	{{.Age}}</li>
	<li>User Roles: {{.Roles}}</li>
</ul>
`

func htmlTemplatesDemo() {
	fmt.Println("## Html Templates")

	t := ht.New("userTemplate")
	t.Parse(userTemplate)

	var buffer bytes.Buffer
	err := t.ExecuteTemplate(&buffer, "userTemplate", &user{
		ID:    "123",
		Name:  "Alice",
		Age:   25,
		Roles: []string{"generic/user", "site/editor"},
	})
	if err != nil {
		fmt.Printf("Failed to execute template, err=%s", err)
		return
	}

	fmt.Printf("Executed Template = %s\n", buffer.String())
}

func main() {
	fmt.Println("# Templates Demo")

	htmlTemplatesDemo()
}
