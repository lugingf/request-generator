package search

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"stash.tutu.ru/golang/log"
	"strconv"
	"strings"
)

type SearchParams struct {
	From          string
	To            string
	DepartureDate string
}

type Offer struct {
	From          string
	To            string
	CarrierName   string
	CarrierRate   int
	DepartureTime string
	ArrivalTime   string
	Id            int
	OfferDetail   OfferDetail
}

type OfferDetail struct {
	Price int
	Seats int
}

type SearchResult struct {
	Offers []Offer
}

func GetSearchResult(client *http.Client, params SearchParams, searchUrl string) (int, string, string) {
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

		searchResult := SearchResult{}
		getOffers(resp, &searchResult)

		offersCount := len(searchResult.Offers)
		log.Logger.Info().Msg("Offers count: " + strconv.Itoa(offersCount))

		defer resp.Body.Close()
	}

	return respStatus, respText, errorText
}

func getOffers(resp *http.Response, target *SearchResult) error {
	return json.NewDecoder(resp.Body).Decode(target)
}
