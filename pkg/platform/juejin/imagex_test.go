package juejin

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadImage(t *testing.T) {
	client, err := NewClient(os.Getenv("ARTICLI_JUEJIN_COOKIE"))
	if err != nil {
		t.Fail()
		return
	}

	_, err = client.UploadImage(RegionCNNorth, "./images/go.png")
	assert.Nil(t, err)
}
