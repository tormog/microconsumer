package producers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	md "microconsumer/internal/models"
	"net/http"
)

var (
	standardHeaders = map[string]string{"User-Agent": "v2RecentSearchGo", "Accept": "application/json", "Connection": "keep-alive"}
)

type twitterSource struct {
	producerData ProducerData
	searchURL    string
	searchQuery  string
	bearerToken  string
}

func getTwitterRecentData(bearerToken, nextToken, newestID, url, query string, hits int) (TwitterRecentData, error) {
	var data TwitterRecentData
	req, _ := http.NewRequest("GET", url, nil)

	for k, v := range standardHeaders {
		req.Header.Add(k, v)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

	q := req.URL.Query()
	q.Add("query", query)
	if nextToken != "" {
		q.Add("next_token", nextToken)
	}
	if newestID != "" {
		q.Add("since_id", newestID)
	}
	q.Add("max_results", fmt.Sprintf("%d", hits))
	req.URL.RawQuery = q.Encode()
	http.DefaultClient.Do(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return data, fmt.Errorf("returned with status code %d", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &data); err != nil {
		return data, err
	}
	return data, nil
}

func (tw *twitterSource) Produce(fromID string) (string, error) {
	nextToken := ""
	newestID := ""
	for {
		data, err := getTwitterRecentData(tw.bearerToken, nextToken, fromID, tw.searchURL, tw.searchQuery, 10)
		if err != nil {
			return "", err
		}
		if data.Meta.ResultCount == 0 {
			return fromID, nil
		}
		if newestID == "" {
			newestID = data.Meta.NewestID
		}
		log.Printf("Producer ID %s => NewestID:%s OldestID:%s NextToken:%s ResultCount:%d\n",
			tw.producerData.producerID, data.Meta.NewestID, data.Meta.OldestID, data.Meta.NextToken, data.Meta.ResultCount)
		for _, field := range data.Data {
			dataStore := md.DataStore{
				Type: tw.producerData.resourceName,
				ID:   field.ID,
				Data: field,
			}
			bytes, err := json.Marshal(dataStore)
			if err != nil {
				return "", err
			}
			tw.producerData.queue.Push(tw.producerData.queueStoreName, bytes)
			log.Printf("Producer ID %s pushed twitter id:%s\n", tw.producerData.producerID, field.ID)
		}
		if data.Meta.NextToken != "" {
			nextToken = data.Meta.NextToken
			continue
		}
		break
	}
	return newestID, nil
}
