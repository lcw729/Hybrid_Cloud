package handler

import (
	"Hybrid_Cloud/metricHandler/util"
	"net/http"

	"github.com/gorilla/mux"
)

func GetNodeMetric(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cluster_name := vars["clustername"]
	node_name := vars["nodename"]

	util.GetResource
}
