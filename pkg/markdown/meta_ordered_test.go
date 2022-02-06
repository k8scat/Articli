package markdown

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestGet(t *testing.T) {
	s := `
a: 1
b:
  c: 2
  d: 3
e:
  f:
    g: 4
    h:
    - i: "5"
    - j: "6"
`
	var m Meta
	err := yaml.Unmarshal([]byte(s), &m)
	assert.Nil(t, err)
	if m != nil {
		v := m.Get("b.c")
		assert.NotNil(t, v)
		assert.Equal(t, 2, v)
	}
}

func TestGetStringArray(t *testing.T) {
	s := `
a: 1
b:
  c: 2
  d: 3
e:
  f:
    g: 4
    h:
    - "5"
    - "6"
`
	var m Meta
	err := yaml.Unmarshal([]byte(s), &m)
	assert.Nil(t, err)
	if m != nil {
		v := m.GetStringSlice("e.f.h")
		assert.NotNil(t, v)
		assert.Equal(t, []string{"5", "6"}, v)
	}
}

func TestSet(t *testing.T) {
	s := `
a: 1
b:
  c: 2
  d: 3
e:
  f:
    g: 4
    h:
    - "5"
    - "6"
`
	var m Meta
	err := yaml.Unmarshal([]byte(s), &m)
	assert.Nil(t, err)

	if m != nil {
		assert.Equal(t, 2, m.Get("b.c"))
		m.Set("b.c", 4)
		assert.Equal(t, 4, m.Get("b.c"))
		m = m.Set("i", "7")
		assert.Equal(t, "7", m.Get("i"))
		m = m.Set("b.k", 8)
		assert.Equal(t, 8, m.Get("b.k"))
	}
}
