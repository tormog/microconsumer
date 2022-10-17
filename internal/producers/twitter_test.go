package producers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	ch "microconsumer/internal/queue"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func TestProduce(t *testing.T) {
	content, err := ioutil.ReadFile("../../test/twitter/twitter.json")
	if err != nil {
		t.Error(err)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log("Got request:" + r.RequestURI)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
		w.Header().Set("Content-Type", "application/json")
		w.Write(content)
	}))
	defer ts.Close()

	mr, err := miniredis.Run()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mr.Close()

	rdClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	rdCache := ch.RedisCache{rdClient}
	queueStoreName := "data"

	tws := twitterSource{
		producerData: ProducerData{
			queue:          rdCache,
			queueStoreName: queueStoreName,
			resourceName:   "twitter",
			producerID:     "someid",
		},
		searchURL:   ts.URL,
		searchQuery: "somequery",
		bearerToken: "token",
	}

	assert.Equal(t, int64(0), rdCache.Length(queueStoreName))
	tws.Produce("")
	assert.Equal(t, int64(10), rdCache.Length(queueStoreName))
}
