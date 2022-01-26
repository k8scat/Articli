package oschina

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSaveDraft(t *testing.T) {
	client, err := NewClient(os.Getenv("ARTICLI_OSCHINA_COOKIE"))
	if err != nil {
		t.Fail()
		return
	}

	bad := &ContentParams{
		Content: "test",
	}
	_, err = client.SaveDraft(bad)
	assert.NotNil(t, err)

	good := &ContentParams{
		DraftID:        "2686735",
		Title:          "test1112222++111",
		Content:        "test",
		Category:       "11723151",
		TechnicalField: "21",
	}
	id, err := client.SaveDraft(good)
	assert.Nil(t, err)
	assert.NotEqual(t, "", id)
}

func TestDeleteDraft(t *testing.T) {
	client, err := NewClient(os.Getenv("ARTICLI_OSCHINA_COOKIE"))
	if err != nil {
		t.Fail()
		return
	}

	err = client.DeleteDraft("2686724")
	assert.Nil(t, err)
}

func TestListDrafts(t *testing.T) {
	client, err := NewClient(os.Getenv("ARTICLI_OSCHINA_COOKIE"))
	if err != nil {
		t.Fail()
		return
	}

	_, _, err = client.ListDrafts(2)
	assert.Nil(t, err)
}

func TestGetDraftDetail(t *testing.T) {
	client, err := NewClient(os.Getenv("ARTICLI_OSCHINA_COOKIE"))
	if err != nil {
		t.Fail()
		return
	}

	_, err = client.GetDraftDetail("2686735")
	assert.Nil(t, err)
}

func TestPublishDraft(t *testing.T) {
	client, err := NewClient(os.Getenv("ARTICLI_OSCHINA_COOKIE"))
	if err != nil {
		t.Fail()
		return
	}

	_, err = client.PublishDraft("2686735")
	assert.Nil(t, err)
}
