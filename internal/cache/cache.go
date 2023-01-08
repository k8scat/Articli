package cache

import (
	"encoding/json"
	"os"
	"path"

	"github.com/juju/errors"
	"github.com/k8scat/articli/internal/config"
)

var (
	localCacheFile = path.Join(config.GetConfigDir(), "cache.json")

	GlobalLocalCache *LocalCache
)

func init() {
	var err error
	GlobalLocalCache, err = NewLocalCache(localCacheFile)
	if err != nil {
		panic(err)
	}
}

const (
	KeyOschinaTechnicalFields = "oschina_techinal_fields"
	KeyOschinaCategories      = "oschina_categories"

	KeyJuejinCategories = "juejin_categories"
	KeyJuejinTags       = "juejin_tags"
	KeyJuejinColumns    = "juejin_columns"
)

type LocalCache struct {
	file string
	data map[string]any
}

func NewLocalCache(file string) (*LocalCache, error) {
	c := &LocalCache{
		file: file,
	}
	if err := c.load(); err != nil {
		return nil, errors.Trace(err)
	}
	return c, nil
}

func (c *LocalCache) load() error {
	c.data = make(map[string]any)
	b, err := os.ReadFile(c.file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	err = json.Unmarshal(b, &c.data)
	if err != nil {
		return err
	}
	return nil
}

func (c *LocalCache) Set(key string, data any) error {
	c.data[key] = data
	b, err := json.Marshal(c.data)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(c.file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(b)
	return err
}

func (c *LocalCache) Get(key string, payload any) error {
	v := c.data[key]
	if v != nil {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		err = json.Unmarshal(b, payload)
		return err
	}
	return nil
}
