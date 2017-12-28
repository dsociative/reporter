package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var authTable = []struct {
	key       string
	code      int
	client    string
	network   string
	startDate string
	stopDate  string
}{
	{"successCode", 200, "16", "google", "2017-12-01", "2017-12-03"},
	{"else", 400, "17", "yandex", "2017-12-01", "SQLINJECTIONFU((("},
	{"failCode", 401, "", "", "", ""},
}

func TestHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	requesterMock := NewMockRequester(ctrl)
	h := NewHandler(
		map[string][]string{
			"successCode": {"16", "google"},
			"else":        {"17", "yandex"},
		},
		requesterMock,
	)

	for _, tt := range authTable {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(
			"GET",
			"/?key="+tt.key+"&start_date="+tt.startDate+"&stop_date="+tt.stopDate,
			nil,
		)
		assert.NoError(t, err)
		if tt.code == 200 {
			requesterMock.EXPECT().Do(
				tt.client,
				tt.network,
				tt.startDate,
				tt.stopDate,
			).Return(tt.client, nil)
		}

		h.ServeHTTP(w, r)

		assert.Equal(t, tt.code, w.Code, tt.key)
		if tt.code == 200 {
			assert.Equal(t, tt.client, w.Body.String())
		}
	}
}
