package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

var (
	configPath = flag.String("config", ".", "config path")

	version     string
	startMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "start",
			Help: "start service notify",
		},
		[]string{"service", "version"},
	)
)

func init() {
	prometheus.MustRegister(startMetric)
}

func readInConfig() error {
	flag.Parse()
	viper.SetConfigName("config")
	viper.AddConfigPath(*configPath)

	viper.SetDefault("prometheus", ":9092")
	viper.SetDefault("https_bind", "localhost:9443")

	viper.SetDefault("cert", "bridge.sysmtem.cer")
	viper.SetDefault("key", "bridge.sysmtem.key")

	viper.SetDefault("api_host", "ortb_api.bridge.systems:10500")
	viper.SetDefault("api_cert", "api.cer")
	viper.SetDefault("site", "ortb_api.bridge.systems")
	viper.SetDefault("admin_key", "SOMEKEYFORAPI")

	viper.SetDefault("clickhouse", "http://localhost:8123/")
	viper.SetDefault("access", "./access.json")

	viper.AutomaticEnv()
	return viper.ReadInConfig()
}

func main() {
	startMetric.WithLabelValues("reporter", version).SetToCurrentTime()
	go http.ListenAndServe(viper.GetString("prometheus"), promhttp.Handler())

	if err := readInConfig(); err != nil {
		log.Fatal("config read: ", err)
	}

	access := map[string][]string{}
	rawAccess, err := ioutil.ReadFile(viper.GetString("access"))
	if err != nil {
		log.Fatal("read acess file", err)
	}

	err = json.Unmarshal(rawAccess, &access)
	if err != nil {
		log.Fatal("unmarshall access file", err)
	}

	handler := NewHandler(
		access,
		NewClickhouseRequester(&http.Client{}, viper.GetString("clickhouse")),
	)
	log.Fatal(http.ListenAndServeTLS(
		viper.GetString("https_bind"),
		viper.GetString("cert"),
		viper.GetString("key"),
		handler,
	))
}
