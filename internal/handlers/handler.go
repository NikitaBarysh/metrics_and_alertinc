package handlers

import (
	"net/http"
	"strconv"
	"strings"
)

var (
	gauges   = make(map[string]float64)
	counters = make(map[string]int64)
)

func Router(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		post(rw, r)
	}
}

func post(rw http.ResponseWriter, r *http.Request) {
	res := strings.Split(r.URL.Path, "/")

	if len(res) < 5 {
		http.Error(rw, ",", http.StatusNotFound)
		return
	}

	update := res[1]
	if update != "update" {
		http.Error(rw, "not update", http.StatusNotFound)
		return
	}

	metricType := res[2]

	metricName := res[3]

	metricValue := res[4]

	switch metricType {
	case "counter":
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(rw, "wrong counter type", http.StatusBadRequest)
			return
		}
		counters[metricName] += value
	case "gauge":
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(rw, "wrong gauge type", http.StatusBadRequest)
			return
		}
		gauges[metricName] = value
	default:
		http.Error(rw, "unknown metric type", http.StatusNotImplemented)
		return
	}

	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusOK)
}
