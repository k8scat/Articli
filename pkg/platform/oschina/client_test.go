package oschina

import (
	"log"
	"os"
	"testing"
)

var (
	client *Client

	cookie = "_user_behavior_=31ba6082-c0fe-4c33-ba06-8b7a92400343; _ga=GA1.2.704792458.1640619283; oscid=InA%2F3ilXagEF1N0E%2FiHqqeAxKuxCLSKysKwCm%2B67wMW%2B4jiVEBUoO%2FY7j7TCKEwkHhueKRv1TyjjRKFxI8uG9G1D86e9zOk58rttEVuTkpUtmr0BHUsDw7cNM%2F258nD2twixuuj1iydvCZ%2F69Fj%2B8jsthuf0PPYN; Hm_lvt_a411c4d1664dd70048ee98afe7b28f0b=1640755313,1641404618,1641404709,1642991616; Hm_lpvt_a411c4d1664dd70048ee98afe7b28f0b=1643113319"
)

func setupClient(t *testing.T) {
	if cookie == "" {
		cookie = os.Getenv("ARTICLI_OSCHINA_COOKIE")
	}
	var err error
	client, err = NewClient(cookie)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewClient(t *testing.T) {
	setupClient(t)
	log.Printf("client: %+v", client)
}
