package main

import (
	"Hybrid_Cloud/metricHandler/server/handler"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterHandler() http.Handler {
	mux := mux.NewRouter()
	// metric
	mux.HandleFunc("/metrics/clusters/{clustername}/nodes/{nodename}", handler.GetNodeMetric).Methods("GET")
	//mux.HandleFunc("/metrics/clusters/{clustername}/nodes/{podname}", handler.GetPodMetric).Methods("GET")
	return mux
}
func main() {

	http.ListenAndServe(":8090", RegisterHandler())
}
