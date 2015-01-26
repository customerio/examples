package main

import (
	"encoding/json"
	"fmt"
	"github.com/segmentio/analytics-go"
	"log"
	"net/http"
	"time"
)

const SEGMENT_WRITE_KEY = "<YOUR KEY HERE>"

type CIOWebhook struct {
	EventType string                 `json:"event_type"`
	EventID   string                 `json:"event_id"`
	Timestamp int                    `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

func main() {

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {

		buf := make([]byte, r.ContentLength)
		r.Body.Read(buf)

		var webhook *CIOWebhook
		err := json.Unmarshal(buf, &webhook)

		if err != nil {
			log.Println(err, r)
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("bad request"))
			return
		}

		customerID := webhook.Data["customer_id"].(string)

		segment := analytics.New(SEGMENT_WRITE_KEY)
		segment.Track(map[string]interface{}{
			"userId":     customerID,
			"event":      fmt.Sprintf("customerio:%v", webhook.EventType),
			"properties": webhook.Data,
			"context": map[string]interface{}{
				"event_id": webhook.EventID,
			},
			"timestamp": time.Unix(int64(webhook.Timestamp), 0).Format(time.RFC3339),
		})

		log.Println("ok", r)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))

	})

	log.Info("Listening on :8080 for incoming webhooks to forward to segment.com")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
