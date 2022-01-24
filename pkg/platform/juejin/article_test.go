package juejin

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateArticle(t *testing.T) {
	setupClient(t)
	articleID, draftID, err := client.SaveArticle("", "", "Docker Intro", "Docker 是一個開放原始碼軟體，是一個開放平台，用於開發應用、交付應用、執行應用。 Docker允許使用者將基礎設施中的應用單獨分割出來，形成更小的顆粒，從而提高交付軟體的速度。 Docker容器與虛擬機器類似，但二者在原理上不同。",
		"https://www.cloudsigma.com/wp-content/uploads/cgroups-docker.jpg",
		"6809637769959178254", "", []string{"6809637776909139982"},
		false)
	assert.Nil(t, err)
	fmt.Println(draftID)
	fmt.Println(fmt.Sprintf(ArticleURLFormat, articleID))
}

func TestDeleteArticle(t *testing.T) {
	setupClient(t)
	err := client.DeleteArticle("6947311736529633294")
	assert.Nil(t, err)
}
