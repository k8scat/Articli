package csdn

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestUploadImage(t *testing.T) {
	client, err := NewClient(os.Getenv("ARTICLI_CSDN_COOKIE"))
	if err != nil {
		t.Error(err)
		return
	}

	downloadURL, err := client.UploadImage("/Users/hsowan/workspace/articli/images/go.png")
	if err != nil {
		t.Error(err)
		return
	}
	assert.NotEqual(t, "", downloadURL)
}
