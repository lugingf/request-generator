package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

const SearchResponsesDurationMetricName string = "ocws_request_generator_search_responses_duration"
const SearchResponsesDurationMetricLabelTarget string = "target"
const SearchResponsesDurationMetricLabelResponseCode string = "response_code"

var SearchResponsesDurationSummaryVec *prometheus.SummaryVec

func Init() {
	SearchResponsesDurationSummaryVec = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: SearchResponsesDurationMetricName,
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
		},
		[]string{SearchResponsesDurationMetricLabelTarget, SearchResponsesDurationMetricLabelResponseCode},
		)

	prometheus.MustRegister(SearchResponsesDurationSummaryVec)
}

func SendDurationMetric(url string, statusCode int, startTime time.Time)  {
	SearchResponsesDurationSummaryVec.With(prometheus.Labels{SearchResponsesDurationMetricLabelTarget: url, SearchResponsesDurationMetricLabelResponseCode: strconv.Itoa(statusCode)}).Observe(time.Since(startTime).Seconds())
}