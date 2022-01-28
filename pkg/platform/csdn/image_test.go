package csdn

import (
	"fmt"
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

	downloadURL, err := client.UploadImage("/Users/hsowan/Downloads/mmap.pNg")
	if err != nil {
		t.Error(err)
		return
	}
	assert.NotEqual(t, "", downloadURL)
	fmt.Println(downloadURL)
}
