package markdown

import (
	"bufio"
	"github.com/juju/errors"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"regexp"
	"strings"
)

var (
	metaSeparatorPattern = regexp.MustCompile(`^--- *$`)
	moreSeparatorPattern = regexp.MustCompile(`(?i)^ *<!-- *more *--> *$`)
)

type Mark struct {
	Meta    Meta
	Raw     []byte
	Content []byte
	Brief   []byte
}

func (m *Mark) WriteFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return errors.Trace(err)
	}
	defer f.Close()

	f.Write([]byte("---\n"))
	b, err := yaml.Marshal(m.Meta)
	if err != nil {
		return errors.Trace(err)
	}
	f.Write(b)
	f.Write([]byte("---\n"))
	f.Write(m.Content)
	return nil
}

func Parse(filepath string) (result *Mark, err error) {
	f, err := os.Open(filepath)
	if err != nil {
		err = errors.Trace(err)
		return
	}
	defer f.Close()

	br := bufio.NewReader(f)
	var metaBytes []byte
	metaSeparatorCount := 0
	moreSeparatorCount := 0
	afterMeta := false
	afterMore := false
	result = new(Mark)
	for {
		if metaSeparatorCount == 2 {
			afterMeta = true
		}
		if moreSeparatorCount == 1 {
			afterMore = true
		}

		line, isPrefix, re := br.ReadLine()
		if re == io.EOF {
			break
		}
		result.Raw = append(result.Raw, line...)
		if !isPrefix {
			result.Raw = append(result.Raw, '\n')
		}

		if metaSeparatorPattern.Match(line) {
			metaSeparatorCount += 1
		}
		if moreSeparatorPattern.Match(line) {
			moreSeparatorCount += 1
		}
		if metaSeparatorCount == 1 {
			metaBytes = append(metaBytes, line...)
			if !isPrefix {
				metaBytes = append(metaBytes, '\n')
			}
		}
		if afterMeta && !afterMore {
			result.Brief = append(result.Brief, line...)
			if !isPrefix {
				result.Brief = append(result.Brief, '\n')
			}
		}
		if afterMeta {
			result.Content = append(result.Content, line...)
			if !isPrefix {
				result.Content = append(result.Content, '\n')
			}
		}
	}
	var meta Meta
	err = yaml.Unmarshal(metaBytes, &meta)
	err = errors.Trace(err)
	result.Meta = meta
	return
}

func UpdateMapSlice(ms yaml.MapSlice, path string, value interface{}) (yaml.MapSlice, error) {
	paths := strings.Split(path, ".")
	slices := []yaml.MapSlice{ms}
	indexes := make([]int, 0)
	found := false
	for i, p := range paths {
		for j, item := range ms {
			if item.Key.(string) == p {
				if i == len(paths)-1 {
					item.Value = value
					ms[j] = item
					found = true
					break
				} else {
					v, ok := item.Value.(yaml.MapSlice)
					if !ok {
						return nil, errors.Errorf("wrong path '%s'", path)
					}
					indexes = append(indexes, j)
					ms = v
					slices = append(slices, ms)
					break
				}
			}
		}
	}
	if !found {
		slices[len(slices)-1] = append(slices[len(slices)-1], yaml.MapItem{
			Key:   paths[len(paths)-1],
			Value: value,
		})
	}

	var result yaml.MapSlice
	for i := len(indexes) - 1; i >= 0; i-- {
		index := indexes[i]
		result = slices[i]
		result[index].Value = slices[i+1]
	}
	return result, nil
}

func GetValue(name string, m1, m2 map[string]interface{}) (interface{}, bool) {
	v, ok := m1[name]
	if ok {
		return v, true
	}
	v, ok = m2[name]
	return v, ok
}

func GetStringValue(name string, m1, m2 map[string]interface{}) (string, error) {
	v, ok := GetValue(name, m1, m2)
	if !ok {
		return "", errors.NotFoundf("key '%s'", name)
	}
	s, ok := v.(string)
	if !ok {
		return "", errors.Errorf("invalid value '%v'", v)
	}
	return s, nil
}

func GetValueFromMapSlice(ms yaml.MapSlice, key string) (interface{}, bool) {
	for _, m := range ms {
		k, ok := m.Key.(string)
		if !ok {
			continue
		}
		if k == key {
			return m.Value, true
		}
	}
	return nil, false
}

func GetStringArray(m map[string]interface{}, key string) ([]string, error) {
	v, ok := m[key]
	if !ok {
		return nil, errors.NotFoundf("key '%v'", key)
	}
	arr, ok := v.([]interface{})
	if !ok {
		return nil, errors.Errorf("invalid value type '%T' for key '%s'", v, key)
	}

	result := make([]string, len(arr))
	for i, v := range arr {
		s, ok := v.(string)
		if !ok {
			return nil, errors.Errorf("invalid value type '%T' for key '%s'", v, key)
		}
		result[i] = s
	}
	return result, nil
}

// ConvertMapSlice convert yaml.MapSlice into map[string]interface{}
func ConvertMapSlice(ms yaml.MapSlice) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for _, m := range ms {
		k, ok := m.Key.(string)
		if !ok {
			return nil, errors.Errorf("invalid type '%T' of key '%s'", k, k)
		}
		result[k] = m.Value
	}
	return result, nil
}
