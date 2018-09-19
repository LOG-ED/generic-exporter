package generic

import (
	"fmt"
	"net/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	address string
	path string
)

func Run(_address, _path string) error {
	address = _address
	path = _path

	exporter, err := NewExporter()
	if err != nil {
		return err
	}
	prometheus.MustRegister(exporter)

	http.Handle(path, promhttp.Handler())
	http.HandleFunc("/", rootHandler)
	fmt.Printf("run the exporter on %s%s\n", address, path)
	return http.ListenAndServe(address, nil)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`<html>
		<head><title>Generic Exporter</title></head>
		<body>
		<h1>Generic Exporter</h1>
		<p><a href="` + path + `">Metrics</a></p>
		</body>
		</html>
	`))

}