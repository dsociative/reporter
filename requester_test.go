package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRequest(t *testing.T) {
	addr := "testaddr"
	ctrl := gomock.NewController(t)
	c := NewMockHTTPClient(ctrl)
	r := NewClickhouseRequester(c, addr)

	b, err := render("223", "google", "STARTDATE", "STOPDATE")

	c.EXPECT().Post(addr, "text", b).Return(&http.Response{Body: ioutil.NopCloser(bytes.NewReader([]byte("table")))}, nil)
	data, err := r.Do("223", "google", "STARTDATE", "STOPDATE")
	assert.NoError(t, err)
	assert.Equal(t, "table", data)
}

func TestRender(t *testing.T) {
	b, err := render("2", "yandex", "STARTDATE", "STOPDATE")
	assert.NoError(t, err)
	assert.Equal(
		t,
		`
    SELECT
        toDate(time, 'America/Araguaina') as date,
        network,
        region,
        prefix,
        count() as imps,
        round(sum(media_cost)/1000,2) as traffic_cost
        FROM impression_data
        WHERE toDate('STARTDATE') <= toDate(time, 'America/Araguaina') AND toDate(time, 'America/Araguaina') <= toDate('STOPDATE')
        AND (network = 'yandex')
    AND (client_id = '2')
        GROUP BY
            date,
            region,
            prefix,
            network
        ORDER BY
        date
    FORMAT CSVWithNames
`,
		b.String(),
	)
}
