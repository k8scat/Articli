package oschina

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSaveDraft(t *testing.T) {
	setupClient(t)

	bad := &ContentParams{
		Content: "test",
	}
	_, err := client.SaveDraft(bad)
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
	fmt.Println(id)
}

func TestDeleteDraft(t *testing.T) {
	setupClient(t)

	err := client.DeleteDraft("2686724")
	assert.Nil(t, err)
}

func TestListDrafts(t *testing.T) {
	setupClient(t)

	drafts, err := client.ListDrafts(2)
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range drafts {
		fmt.Printf("name: %s, id: %s\n", d.Title, d.ID)
	}
}

func TestGetDraftDetail(t *testing.T) {
	setupClient(t)

	draft, err := client.GetDraftDetail("2686735")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", draft)
}

func TestPublishDraft(t *testing.T) {
	setupClient(t)

	articleID, err := client.PublishDraft("2686735")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(client.BuildArticleURL(articleID))
}
