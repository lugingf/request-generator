package search

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"stash.tutu.ru/golang/log"
	"stash.tutu.ru/opscore-workshop-admin/request-generator/metrics"
	"strconv"
	"strings"
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

func GetSearchResult(client *http.Client, params Params, searchUrl string) (int, string, string) {
	request, err := http.NewRequest("GET", searchUrl, nil)

	query := request.URL.Query()
	query.Add("from", params.From)
	query.Add("to", params.To)
	query.Add("date", params.DepartureDate)
	request.URL.RawQuery = query.Encode()

	resp, err := client.Do(request)

	respStatus := http.StatusOK
	var respText string
	var errorText string

	if err, ok := err.(net.Error); ok && err.Timeout() {
		errorText = "Got Timeout ERROR: " + error.Error(err)
		respStatus = http.StatusGatewayTimeout
	} else if err != nil {
		errorText = "Got ERROR: " + error.Error(err)
		respStatus = http.StatusInternalServerError
	} else {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		respText = "Response " + strings.Replace(string(bodyBytes), "\n", "", -1) + "\n"
		respStatus = resp.StatusCode

		searchResult := Result{}
		err := json.Unmarshal(bodyBytes, &searchResult)

		if err != nil {
			fmt.Println(err.Error())
		}

		offersCount := len(searchResult.Offers)
		metrics.SendOffersCountMetric(searchUrl, float64(offersCount))
		log.Logger.Info().Msg("Offers count: " + strconv.Itoa(offersCount))

		defer resp.Body.Close()
	}

	return respStatus, respText, errorText
}
