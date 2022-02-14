package gitlab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/juju/errors"
)

type Project struct {
	ID                int    `json:"id"`
	Description       string `json:"description"`
	DefaultBranch     string `json:"default_branch"`
	Visibility        string `json:"visibility"`
	Path              string `json:"path"`
	PathWithNamespace string `json:"path_with_namespace"`
	Name              string `json:"name"`
	NameWithNamespace string `json:"name_with_namespace"`
}

func (p *Project) IsPrivate() bool {
	return p.Visibility != "public"
}

func (c *Client) GetProject(id string) (*Project, error) {
	path := fmt.Sprintf("/projects/%s", URLEncoded(id))
	resp, err := c.Request(http.MethodGet, path, nil, nil, nil)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected response: %s", b)
	}
	var project *Project
	err = json.Unmarshal(b, &project)
	return project, errors.Trace(err)
}
