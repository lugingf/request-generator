package generators

import (
	"io/ioutil"
	"math"
	"math/rand"
	"net"
	"net/http"
	"stash.tutu.ru/golang/log"
	"stash.tutu.ru/opscore-workshop-admin/request-generator/metrics"
	"stash.tutu.ru/opscore-workshop-admin/request-generator/resources"
	"strconv"
	"strings"
	"time"
)

type RequestGenerator struct {
	Config *resources.GeneratorConfig
}

func (g *RequestGenerator) MakeGen(ch chan resources.ChangedNode) {
	rpm := getDurationTypeFromString(g.Config.Rpm)
	if rpm == 0 {
		rpm = 30
		log.Logger.Info().Msg("Generator RPM by default: " + strconv.Itoa(int(rpm)))
	}

	urlData := g.Config.UrlTargetList
	timeout := getDurationTypeFromString(g.Config.GenRequestTimeout)
	client := &http.Client{Timeout: timeout * time.Millisecond}
	configChangedMessage := "Config Changed: "
	for {
		select {
		case changedNode := <-ch:
			switch changedNode.Node {
			case "isEnabled":
				g.Config.IsEnabled = changedNode.Value
				log.Logger.Info().Msg(configChangedMessage + "Generator enabling set to " + g.Config.IsEnabled)
			case "LogResponsesEnabled":
				g.Config.LogResponsesEnabled = changedNode.Value
				log.Logger.Info().Msg(configChangedMessage + "LogResponsesEnabled set to " + g.Config.IsEnabled)
			case "rpm":
				g.Config.Rpm = changedNode.Value
				rpm = getDurationTypeFromString(g.Config.Rpm)
				log.Logger.Info().Msg(configChangedMessage + "Parameter requests per minute was changed to " + g.Config.Rpm)
			case "GenRequestTimeout":
				g.Config.GenRequestTimeout = changedNode.Value
				timeout := getDurationTypeFromString(g.Config.GenRequestTimeout)
				client.Timeout = timeout * time.Millisecond
				log.Logger.Info().Msg(configChangedMessage + "Parameter request timeot was changed to " + g.Config.GenRequestTimeout)
			case "UrlTargetList":
				g.Config.UrlTargetList = strings.Split(changedNode.Value, ",")
				log.Logger.Info().Msg(configChangedMessage + "URL target list was updated to " + changedNode.Value)
				urlData = g.Config.UrlTargetList
			}
		default:
			if g.Config.IsEnabled != "1" {
				continue
			}
			g.doRequest(client, urlData)
			sleep := time.Duration(int64(math.Round(float64(time.Minute / rpm))))
			time.Sleep(sleep * time.Nanosecond)
		}
	}
}

func (g *RequestGenerator) doRequest(client *http.Client, urlData []string) {
	for _, url := range urlData {
		go func(url string) {
			request, err := http.NewRequest("GET", url, nil)

			query := request.URL.Query()
			from, to, departureDate := getRequestData()
			query.Add("from", strconv.Itoa(from))
			query.Add("to", strconv.Itoa(to))
			query.Add("date", departureDate)
			request.URL.RawQuery = query.Encode()

			startTime := time.Now()

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
				defer resp.Body.Close()
			}

			if g.Config.LogResponsesEnabled == "1" {
				log.Logger.Info().Msg(respText)
			}

			if errorText != "" {
				log.Logger.Info().Msg(errorText)
			}

			metrics.SendDurationMetric(url, respStatus, startTime)
		}(url)
	}
}

func getRequestData() (int, int, string) {
	var citiesCount = 10

	departureTimestamp := rand.Int63n(60*60*24*(45)) + time.Now().Unix() + 60*60*24*14
	departureDate := time.Unix(departureTimestamp, 0).Format("2006-01-02")

	departureId := rand.Intn(citiesCount - 1) + 1
	arrivalId := rand.Intn(citiesCount -1) + 1

	return departureId, arrivalId, departureDate
}

func getDurationTypeFromString(strValue string) time.Duration {
	valueInt, _ := strconv.Atoi(strValue)
	return time.Duration(valueInt)
}
