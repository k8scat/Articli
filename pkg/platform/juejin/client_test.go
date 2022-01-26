package juejin

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	_, err := NewClient(os.Getenv("ARTICLI_JUEJIN_COOKIE"))
	assert.Nil(t, err)
}
