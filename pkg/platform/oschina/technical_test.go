package oschina

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestListTechnicalFields(t *testing.T) {
	client, err := NewClient(os.Getenv("ARTICLI_OSCHINA_COOKIE"))
	if err != nil {
		t.Fail()
		return
	}

	fields, err := client.ListTechnicalFields()
	assert.Nil(t, err)
	assert.Greater(t, 0, len(fields))
}
