package main


import (
    "time"

    server "github.com/henning70/simple_go_modules/prometheus_server"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
   simple_exporter_metrics = promauto.NewGaugeVec(prometheus.GaugeOpts{
       Name: "simple_exporter_metrics",
       Help: "Simple exporter metrics",
   },
       []string{"name", "country", "company"},
   )

   simple_exporter_error = promauto.NewGaugeVec(prometheus.GaugeOpts{
       Name: "simple_exporter_error",
       Help: "Simple exporter error information",
   },
       []string{"name", "country", "company", "error"},
   )
)

func main() {
    server.Init()

    for {
        simple_exporter_error.WithLabelValues("reports last 15m", "NL", "ing", "429 Too Many Requests").Set(float64(1))

        time.Sleep(5 * time.Second)
    }
}
