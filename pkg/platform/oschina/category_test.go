package oschina

import (
	"fmt"
	"testing"
)

func TestListCategories(t *testing.T) {
	setupClient(t)
	categories, err := client.ListCategories()
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range categories {
		fmt.Printf("name: %s, id: %s\n", c.Name, c.ID)
	}
	fmt.Println(len(categories))
}
