package github

import (
	"fmt"
	"os"
	"testing"
)

var (
	client *Client
	token  = "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"
)

func setupClient(t *testing.T) {
	if token == "" {
		token = os.Getenv("ARTICLI_GITHUB_TOKEN")
	}

	var err error
	client, err = NewClient(token)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(client.User.Name)
}

func TestNewClient(t *testing.T) {
	setupClient(t)
	fmt.Println(client.User.Login)
}
