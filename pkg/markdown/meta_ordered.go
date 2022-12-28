package markdown

import (
	"gopkg.in/yaml.v2"
	"strings"
)

type Meta yaml.MapSlice

// Get returns the value by path
//
//	{
//	  "a": {
//	    "b": {
//	      "c": 1
//	    }
//	  }
//	}
//
// e.g. Get("a.b.c") returns 1
func (m Meta) Get(path string) interface{} {
	var subpath string
	if i := strings.Index(path, "."); i != -1 {
		subpath = path[i+1:]
		path = path[:i]
	}
	for _, item := range m {
		k, ok := item.Key.(string)
		if !ok {
			continue
		}
		if k == path {
			if subpath == "" {
				return item.Value
			}
			if v, ok := item.Value.(Meta); ok {
				return v.Get(subpath)
			}
			return nil
		}
	}
	return nil
}

func (m Meta) GetBool(path string) bool {
	v := m.Get(path)
	if v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}

func (m Meta) GetString(path string) string {
	v := m.Get(path)
	if v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}

func (m Meta) GetStringSlice(path string) []string {
	v := m.Get(path)
	if v == nil {
		return nil
	}
	s, _ := v.([]interface{})
	if s == nil {
		return nil
	}
	ss := make([]string, 0)
	for _, v := range s {
		i, ok := v.(string)
		if !ok {
			continue
		}
		ss = append(ss, i)
	}
	return ss
}

func (m Meta) Set(path string, value interface{}) Meta {
	var subpath string
	if i := strings.Index(path, "."); i != -1 {
		subpath = path[i+1:]
		path = path[:i]
	}
	for i, item := range m {
		k, ok := item.Key.(string)
		if !ok {
			continue
		}
		if k == path {
			if subpath == "" {
				m[i] = yaml.MapItem{Key: path, Value: value}
				return m
			}
			if v, ok := item.Value.(Meta); ok {
				t := v.Set(subpath, value)
				m[i] = yaml.MapItem{Key: path, Value: t}
				return m
			}

			root := m
			t := m
			parts := strings.Split(subpath, ".")
			for j, s := range parts {
				n := Meta{}
				t = append(t, yaml.MapItem{Key: s, Value: n})
				if j < len(parts)-1 {
					t = n
				}
			}
			return root
		}
	}
	return append(m, yaml.MapItem{Key: path, Value: value})
}
