package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// CIOWebhook is a wrapper for the customer.io webhook format
type CIOWebhook struct {
	EventType string                 `json:"event_type"`
	EventID   string                 `json:"event_id"`
	Timestamp int                    `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

var routes = map[string]string{
	// "email_drafted": "http://www.example.com",
	"email_sent": "http://www.example.com",
	// "email_delivered":       "http://www.example.com",
	// "email_opened":          "http://www.example.com",
	// "email_clicked":         "http://www.example.com",
	// "email_bounced":         "http://www.example.com",
	// "email_spammed":         "http://www.example.com",
	// "email_dropped":         "http://www.example.com",
	// "email_failed":          "http://www.example.com",
	// "customer_unsubscribed": "http://www.example.com",
	// "customer_subscribed":   "http://www.example.com",
}

type request struct {
	url   string
	hook  *CIOWebhook
	tries int
}

func main() {

	var urlInvalid bool
	for eventType, outbound := range routes {
		u, err := url.Parse(outbound)
		if err != nil {
			fmt.Println(eventType, "--", err)
			urlInvalid = true
			continue
		} else if !u.IsAbs() {
			fmt.Println(eventType, "malformed -- URLs must be fully qualified domains (ex. http://customer.io)")
			urlInvalid = true
		}
	}

	if urlInvalid {
		return
	}

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {

		buf := make([]byte, r.ContentLength)
		r.Body.Read(buf)

		var webhook *CIOWebhook
		err := json.Unmarshal(buf, &webhook)

		if err != nil {
			log.Println(err, r)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("bad request"))
			return
		}

		if routes[webhook.EventType] != "" {
			response, err := http.Post(routes[webhook.EventType], "application/json", bytes.NewReader(buf))
			if err != nil {
				log.Println(webhook.EventType, routes[webhook.EventType], err, r)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			log.Println(webhook.EventType, routes[webhook.EventType], response.StatusCode, r)
			w.WriteHeader(response.StatusCode)
			return
		}

		log.Println(webhook.EventType, nil, 200)
		w.WriteHeader(http.StatusOK)
	})

	log.Println("Listening on :8080 for incoming webhooks to reroute")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
