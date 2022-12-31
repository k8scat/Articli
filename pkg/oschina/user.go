package oschina

import (
	"fmt"
	"regexp"
)

// parseUser parse user data from html
func (c *Client) parseUser() error {
	raw, err := c.get("https://www.oschina.net/", nil, nil)
	if err != nil {
		return err
	}

	// Parse base url
	r := regexp.MustCompile(`g_user_url" data-value="([^"]+)`).FindStringSubmatch(raw)
	if len(r) != 2 {
		return fmt.Errorf("user url not found: %v", r)
	}
	c.baseURL = r[1]

	// Parse username
	r = regexp.MustCompile(`g_user_name" data-value="([^"]+)`).FindStringSubmatch(raw)
	if len(r) != 2 {
		return fmt.Errorf("user name not found: %v", r)
	}
	c.userName = r[1]

	// Parse space id
	ch := make(chan error, 1)
	go func(c *Client) {
		raw, err := c.get(c.baseURL, nil, nil)
		if err != nil {
			ch <- err
			return
		}
		r := regexp.MustCompile(`space_user_id" data-value="(\d+)`).FindStringSubmatch(raw)
		if len(r) != 2 {
			ch <- fmt.Errorf("space id not found: %v", r)
		}
		c.spaceID = r[1]
		ch <- nil
	}(c)
	if err = <-ch; err != nil {
		return err
	}

	// Parse user code
	r = regexp.MustCompile(`g_user_code" data-value="([0-9a-zA-Z]+)`).FindStringSubmatch(raw)
	if len(r) != 2 {
		return fmt.Errorf("user code not found: %v", r)
	}
	c.userCode = r[1]
	// Parse user id
	r = regexp.MustCompile(`g_user_id" data-value="(\d+)`).FindStringSubmatch(raw)
	if len(r) != 2 {
		return fmt.Errorf("user id not found: %v", r)
	}
	c.userID = r[1]
	return nil
}
