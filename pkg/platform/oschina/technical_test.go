package oschina

import (
	"fmt"
	"testing"
)

func TestListTechnicalFields(t *testing.T) {
	setupClient(t)
	fields, err := client.ListTechnicalFields()
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range fields {
		fmt.Printf("name: %s, id: %s\n", f.Name, f.ID)
	}
	fmt.Println(len(fields))
}
