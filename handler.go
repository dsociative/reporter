package main

import (
	"net/http"
	"time"
)

const timeFormat = "2006-01-02"

type Requester interface {
	Do(client, network, start_date, stop_date string) (string, error)
}

type handler struct {
	access    map[string][]string
	requester Requester
}

func NewHandler(
	access map[string][]string,
	requester Requester,
) handler {
	return handler{access: access, requester: requester}
}

func validate(d string) bool {
	_, err := time.Parse(timeFormat, d)
	return err == nil
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	key := r.Form.Get("key")
	var client string
	var network string

	if clientNetworkPair, ok := h.access[key]; ok {
		client = clientNetworkPair[0]
		network = clientNetworkPair[1]
	}

	if client != "" {
		startDate := r.Form.Get("start_date")
		stopDate := r.Form.Get("stop_date")
		if validate(startDate) && validate(stopDate) {
			if data, err := h.requester.Do(client, network, startDate, stopDate); err == nil {
				w.Write([]byte(data))
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
