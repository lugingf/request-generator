package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

const SearchResponsesDurationMetricName string = "ocws_request_generator_search_responses_duration"
const SearchResponsesDurationMetricLabelTarget string = "target"
const SearchResponsesDurationMetricLabelResponseCode string = "response_code"

const SearchOffersCountMetricName string = "ocws_request_generator_search_offers_count"
const SearchOffersCountMetricLabelTarget string = "target"

var SearchResponsesDurationSummaryVec *prometheus.SummaryVec
var SearchOffersCountCounterVec *prometheus.CounterVec

func Init() {
	SearchResponsesDurationSummaryVec = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       SearchResponsesDurationMetricName,
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
		},
		[]string{SearchResponsesDurationMetricLabelTarget, SearchResponsesDurationMetricLabelResponseCode},
	)
	SearchOffersCountCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: SearchOffersCountMetricName,
		},
		[]string{SearchOffersCountMetricLabelTarget},
	)

	prometheus.MustRegister(SearchResponsesDurationSummaryVec)
	prometheus.MustRegister(SearchOffersCountCounterVec)
}

func SendDurationMetric(url string, statusCode int, startTime time.Time) {
	SearchResponsesDurationSummaryVec.With(prometheus.Labels{SearchResponsesDurationMetricLabelTarget: url, SearchResponsesDurationMetricLabelResponseCode: strconv.Itoa(statusCode)}).Observe(time.Since(startTime).Seconds())
}

func SendOffersCountMetric(url string, countValue float64) {
	SearchOffersCountCounterVec.With(prometheus.Labels{SearchOffersCountMetricLabelTarget: url}).Add(countValue)
}
