package producers

import "time"

type TwitterDataField struct {
	Text      string    `json:"text"`
	ID        string    `json:"id"`
	AuthorID  string    `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
}

type TwitterRecentData struct {
	Data []TwitterDataField `json:"data"`
	Meta struct {
		NewestID    string `json:"newest_id"`
		NextToken   string `json:"next_token"`
		OldestID    string `json:"oldest_id"`
		ResultCount int    `json:"result_count"`
	} `json:"meta"`
}
