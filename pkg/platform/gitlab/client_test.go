package gitlab

import "os"

var (
	token  = ""
	client *Client
)

func setupClient() {
	if token == "" {
		token = os.Getenv("ARTICLI_GITLAB_TOKEN")
	}
	var err error
	client, err = NewClient(BaseURLJihuLab, token)
	if err != nil {
		panic(err)
	}
}
