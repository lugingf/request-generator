package handlers

import (
	"net/http"
)

//type Target struct {
//	Target string
//	Timeout string
//}

func (h *Handler) Test(w http.ResponseWriter, request *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello world!"))
}
	//	tu := Target{}
	//	decoder := json.NewDecoder(request.Body)
	//	decoder.Decode(&tu)
	//
	//	url := tu.Target
	//	request, err := http.NewRequest("GET", url, nil)
	//
	//	tmStr := tu.Timeout
	//	tmInt, _ := strconv.Atoi(tmStr)
	//	tm := time.Duration(tmInt)
	//
	//	departureTimestamp := rand.Int63n(60*60*24*(25-14)) + time.Now().Unix() + 60*60*24*14
	//	departureDate := time.Unix(departureTimestamp, 0).Format("2006-01-02")
	//
	//
	//	query := request.URL.Query()
	//	query.Add("from", "1")
	//	query.Add("to", "2")
	//	query.Add("date", departureDate)
	//	request.URL.RawQuery = query.Encode()
	//
	//	startTime := time.Now()
	//
	//	client := &http.Client{Timeout: tm * time.Millisecond}
	//	fmt.Println("DOING " + request.URL.Path)
	//	resp, err := client.Do(request)
	//
	//	respStatus := http.StatusOK
	//	var respText string
	//
	//	if err, ok := err.(net.Error); ok && err.Timeout() {
	//		respText = "Got Timeout ERROR Error: " + error.Error(err)
	//		respStatus = http.StatusGatewayTimeout
	//	} else if err != nil {
	//		respText = "Got ERROR Error: " + error.Error(err)
	//		respStatus = http.StatusInternalServerError
	//	} else {
	//		bodyBytes, _ := ioutil.ReadAll(resp.Body)
	//		respText = "Response " + strings.Replace(string(bodyBytes), "\n", "", -1) + "\n"
	//		fmt.Print(respText)
	//		defer resp.Body.Close()
	//	}
	//
	//	log.Logger.Info().Msg(respText)
	//	sendDurationMetric(url, respStatus, startTime)
	//	w.WriteHeader(respStatus)
	//	w.Write([]byte(respText))
	//}
	//
	//func sendDurationMetric(url string, statusCode int, startTime time.Time)  {
	//	metrics.SearchResponsesDurationSummaryVec.With(prometheus.Labels{metrics.SearchResponsesDurationMetricLabelTarget: url, metrics.SearchResponsesDurationMetricLabelResponseCode: strconv.Itoa(statusCode)}).Observe(time.Since(startTime).Seconds())
	//}
