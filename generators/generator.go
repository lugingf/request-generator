package generators

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type RequestGenerator struct {
	Db *sql.DB
}

func (g *RequestGenerator) MakeGen()  {
	envRPM := os.Getenv("GEN_RPM")

	rpmInt, _ := strconv.Atoi(envRPM)
	rpm := time.Duration(rpmInt)

	if rpm == 0 {
		rpm = 30
	}

	//urlData := getRowsDataDb(g.Db)
	urlData := getRowsData()
	i := 0
	client := &http.Client{Timeout: 60 * time.Second}
	for {
		isEnabled := os.Getenv("GEN_IS_ENABLED")
		if isEnabled != "1" {
			continue
		}
		for _, v := range urlData {
			go func(v string) {
			request, err := http.NewRequest("GET", v, nil)

			query := request.URL.Query()
			from, to, departureDate := getRequestData()
			query.Add("from", strconv.Itoa(from))
			query.Add("to", strconv.Itoa(to))
			query.Add("date", departureDate)

			request.URL.RawQuery = query.Encode()
			resp, err := client.Do(request)
			if err != nil {
				error.Error(err)
			}

			bodyBytes, err := ioutil.ReadAll(resp.Body)
			fmt.Print("Response ", string(bodyBytes), "\n")
			}(v)
		}
		// todo выпилить
		i = i + 1
		fmt.Print("Request count: ", i, "\n")

		sleep := time.Minute / rpm
		time.Sleep(sleep * time.Nanosecond)
	}
}

func getRequestData() (int, int, string) {
	var cities = []string {
		"Rostov-on-Don",
		"Taganrog",
		"Armavir",
		"Moscow",
		"Temriuk",
		"Tashkent",
		"Azov",
		"Eysk",
		"Mukhosransk",
		"Bidlogorsk",
	}

	departureTimestamp := rand.Int63n( 60 * 60 * 24 * (25 - 14)) + time.Now().Unix() + 60 * 60 * 24 * 14
	departureDate := time.Unix(departureTimestamp, 0).Format("2006-01-02")

	departureId := rand.Intn(len(cities) - 1) + 1
	arrivalId := rand.Intn(len(cities) - 1) + 1

	return departureId, arrivalId, departureDate
}

func getRowsDataDb(db *sql.DB) []string {
	rows, _ := db.Query("SELECT * FROM urls")
	var id int
	var url string
	var urlData []string

	for rows.Next() {
		rows.Scan(&id, &url)
		urlData = append(urlData, url)
	}

	return urlData
}

func getRowsData() []string {
	urlData := strings.Split(os.Getenv("GEN_TARGETS"), ",")
	return urlData
}
