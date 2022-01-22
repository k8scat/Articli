package juejin

import (
	"fmt"
	"testing"

	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
)

func TestUploadImage(t *testing.T) {
	setupClient(t)
	imageURL, err := client.UploadImage(RegionCNNorth, "/home/hsowan/workspace/articli/images/go.png")
	assert.Nil(t, errors.Trace(err))
	fmt.Printf("imageURL: %s\n", imageURL)
}
