package search

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"stash.tutu.ru/golang/log"
	"stash.tutu.ru/opscore-workshop-admin/request-generator/metrics"
	"strconv"
	"strings"
	"time"
)

type Params struct {
	From          string
	To            string
	DepartureDate string
}

type Offer struct {
	From          string      `json:"from"`
	To            string      `json:"to"`
	CarrierName   string      `json:"carrierName"`
	CarrierRate   int         `json:"carrierRate"`
	DepartureTime string      `json:"departureTime"`
	ArrivalTime   string      `json:"arrivalTime"`
	Id            int         `json:"id"`
	OfferDetail   OfferDetail `json:"offerDetail"`
}

type OfferDetail struct {
	Price int `json:"price"`
	Seats int `json:"seats"`
}

type Result struct {
	Offers []Offer `json:"offers"`
}

func GetSearchResult(client *http.Client, params Params, target string) (string, string) {
	baseUrlParts := strings.Split(os.Getenv("BASE_SEARCH_URL"), ",")
	targetUrl := baseUrlParts[0] + target + baseUrlParts[1]
	request, err := http.NewRequest("GET", targetUrl, nil)

	query := request.URL.Query()
	query.Add("from", params.From)
	query.Add("to", params.To)
	query.Add("date", params.DepartureDate)
	request.URL.RawQuery = query.Encode()
	requestId := " RequestId: "+ strconv.Itoa(rand.Intn(1000000))
	log.Logger.Info().Msg("Send request to " + target + " Params: " + request.URL.RawQuery + requestId)

	startTime := time.Now()
	resp, err := client.Do(request)

	log.Logger.Info().Msg("Finish request to " + target + requestId)

	respStatus := http.StatusOK
	var respText string
	var errorText string

	if err, ok := err.(net.Error); ok && err.Timeout() {
		errorText = "Got Timeout ERROR: " + error.Error(err) +  requestId + "Timeout was: " + strconv.Itoa(int(client.Timeout))
		respStatus = http.StatusGatewayTimeout
	} else if err != nil {
		errorText = "Got ERROR: " + error.Error(err) + requestId
		respStatus = http.StatusInternalServerError
	} else {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		respText = "Response " + strings.Replace(string(bodyBytes), "\n", "", -1) + requestId
		respStatus = resp.StatusCode

		searchResult := Result{}
		err := json.Unmarshal(bodyBytes, &searchResult)

		if err != nil {
			log.Logger.Info().Msg("Got resp parsing error: " + err.Error() + " from " + target + requestId)
		}

		offersCount := len(searchResult.Offers)
		metrics.SendOffersCountMetric(target, float64(offersCount))
		log.Logger.Info().Msg("Offers count from "+ target + " Count: " + strconv.Itoa(offersCount) + requestId)

		defer resp.Body.Close()
	}
	metrics.SendDurationMetric(target, respStatus, startTime)

	return respText, errorText
}
