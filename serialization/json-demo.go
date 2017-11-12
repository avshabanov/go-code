package main

import (
	"encoding/json"
	"fmt"
)

type role struct {
	ID       string `json:"id"`
	RoleName string `json:"role-name"`
	Priority int    `json:"priority"`
}

type person struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Roles []*role `json:"roles"`
}

func personToJSON(p *person) (string, error) {
	j, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

func main() {
	var userRole = &role{ID: "R1", RoleName: "USER", Priority: 14}
	var adminRole = &role{ID: "R9", RoleName: "ADMIN", Priority: 0}
	var persons = []*person{
		&person{ID: "1", Name: "Cavin", Roles: []*role{}},
		&person{ID: "2", Name: "Alice", Roles: []*role{userRole}},
		&person{ID: "3", Name: "David", Roles: []*role{userRole, adminRole}},
	}

	for index, person := range persons {
		personJSON, err := personToJSON(person)
		if err != nil {
			fmt.Printf("Unable to convert item to json: %s", err.Error())
			return
		}

		fmt.Printf("person #%d - json=\n\t%s\n", index, personJSON)
	}
}
