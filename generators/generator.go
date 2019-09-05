package generators

import (
	"math"
	"math/rand"
	"net/http"
	"stash.tutu.ru/golang/log"
	"stash.tutu.ru/opscore-workshop-admin/request-generator/resources"
	"stash.tutu.ru/opscore-workshop-admin/request-generator/search"
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
			case "IsBadRequestsEnabled":
				g.Config.IsBadRequestsEnabled = changedNode.Value
				log.Logger.Info().Msg(configChangedMessage + "IsBadRequestsEnabled set to " + g.Config.IsEnabled)
			case "rpm":
				g.Config.Rpm = changedNode.Value
				rpm = getDurationTypeFromString(g.Config.Rpm)
				log.Logger.Info().Msg(configChangedMessage + "Parameter requests per minute was changed to " + g.Config.Rpm)
			case "GenRequestTimeout":
				g.Config.GenRequestTimeout = changedNode.Value
				timeout := getDurationTypeFromString(g.Config.GenRequestTimeout)
				client.Timeout = timeout * time.Millisecond
				log.Logger.Info().Msg(configChangedMessage + "Parameter request timeout was changed to " + g.Config.GenRequestTimeout)
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
	for _, target := range urlData {
		go func(target string) {
			from, to, departureDate := getRequestData()
			if g.Config.IsBadRequestsEnabled == "1" {
				from, to, departureDate = getFailedRequestData()
			}

			searchParams := search.Params{From: from, To: to, DepartureDate: departureDate}

			respText, err := search.GetSearchResult(client, searchParams, target)

			if g.Config.LogResponsesEnabled == "1" {
				log.Logger.Info().Msg(respText)
			}

			if err != "" {
				log.Logger.Info().Msg(err)
			}

		}(target)
	}
}

func getFailedRequestData() (string, string, string)  {
	return getRandomStringForRequest(), getRandomStringForRequest(), getRandomStringForRequest()
}

func getRandomStringForRequest() string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 5)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func getRequestData() (string, string, string) {
	var citiesCount = 10

	departureTimestamp := rand.Int63n(60*60*24*(45)) + time.Now().Unix() + 60*60*24*14
	departureDate := time.Unix(departureTimestamp, 0).Format("2006-01-02")

	departureId := rand.Intn(citiesCount - 1) + 1
	arrivalId := rand.Intn(citiesCount -1) + 1

	return strconv.Itoa(departureId), strconv.Itoa(arrivalId), departureDate
}

func getDurationTypeFromString(strValue string) time.Duration {
	valueInt, _ := strconv.Atoi(strValue)
	return time.Duration(valueInt)
}
