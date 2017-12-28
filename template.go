package main

import (
	"bytes"
	"text/template"
)

var (
	tmpl = template.Must(
		template.New("request").Parse(`
    SELECT
        toDate(time, 'America/Araguaina') as date,
        network,
        region,
        prefix,
        count() as imps,
        round(sum(media_cost)/1000,2) as traffic_cost
        FROM impression_data
        WHERE toDate('{{.StartDate}}') <= toDate(time, 'America/Araguaina') AND toDate(time, 'America/Araguaina') <= toDate('{{.StopDate}}')
        AND (network = '{{.Network}}')
    AND (client_id = '{{.ClientID}}')
        GROUP BY
            date,
            region,
            prefix,
            network
        ORDER BY
        date
    FORMAT CSVWithNames
`),
	)
)

func render(client, network, startDate, stopDate string) (*bytes.Buffer, error) {
	b := &bytes.Buffer{}
	err := tmpl.Execute(
		b,
		Query{
			ClientID:  client,
			Network:   network,
			StartDate: startDate,
			StopDate:  stopDate,
		},
	)
	return b, err
}
