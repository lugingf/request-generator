package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const NginxErrorWeightSumMetricName string = "errorcontrol_nginx_500_error_weight_sum"
const NginxErrorCounterMetricName string = "errorcontrol_nginx_error_count"

var NginxErrorWeightSumCounterVec *prometheus.CounterVec
var NginxErrorCounterVec *prometheus.CounterVec

func Init() {
	NginxErrorWeightSumCounterVec = prometheus.NewCounterVec(prometheus.CounterOpts{Name: NginxErrorWeightSumMetricName},
	[]string{"rule_name"})
	NginxErrorCounterVec = prometheus.NewCounterVec(prometheus.CounterOpts{Name: NginxErrorCounterMetricName},
		[]string{"rule_name"})

	prometheus.MustRegister(NginxErrorWeightSumCounterVec)
	prometheus.MustRegister(NginxErrorCounterVec)
}