package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

const timeLayout = "02/01/2006 15:04"

var (
	urlAddress string
	dbPath     string
	records    int
)

func init() {
	flag.StringVar(&urlAddress, "url", "", "Its the web url for data to be analysed, If not set program stops.")
	flag.StringVar(&dbPath, "path", "", "Database path where sqlite db is store. Defaults to current path.")
	flag.IntVar(&records, "limit", 10, "Number of records to be returned starting with the most recent.")
}

type matchInfo struct {
	ID       string    `json:"id"`
	Division string    `json:"div"`
	Date     time.Time `json:"date_time"`
	HomeTeam string    `json:"home_team"`
	AwayTeam string    `json:"away_team"`
	FTHG     int8      `json:"fthg"` // Full Time Home Goals
	FTAG     int8      `json:"ftag"` // Full Time Away Goals
}

// queryMatchesData fetches the data from the provided url and reads it into an array.
func queryMatchesData(urlAddress string) ([]matchInfo, error) {
	resp, err := http.Get(urlAddress)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid status code found: " + resp.Status + " expected status 200 OK")
	}

	data := make([]matchInfo, 0)
	r := csv.NewReader(resp.Body)

	rowsCount := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if rowsCount == 0 {
			rowsCount++
			continue
		}

		if err != nil {
			return nil, errors.New("error reading file: " + err.Error())
		}

		t, err := time.Parse(timeLayout, record[1]+" "+record[2])
		if err != nil {
			return nil, errors.New("invalid time format: " + err.Error())
		}

		fthg, err := strconv.Atoi(record[5])
		if err != nil {
			return nil, errors.New("invalid FTHG value: " + err.Error())
		}

		ftag, err := strconv.Atoi(record[6])
		if err != nil {
			return nil, errors.New("invalid FTAG value: " + err.Error())
		}

		data = append(data, matchInfo{
			Division: record[0],
			Date:     t,
			HomeTeam: record[3],
			AwayTeam: record[4],
			FTHG:     int8(fthg),
			FTAG:     int8(ftag),
		})
	}

	return data, nil
}

func main() {
	flag.Parse()

	fmt.Println("Teams Matches sort running, please wait...")
	if urlAddress == "" {
		fmt.Println("ERROR (init): missing data URL Address expected")
		os.Exit(0)
	}

	// pull data from the online location
	fmt.Printf("Scrapping data from (%s)...\n", urlAddress)
	data, err := queryMatchesData(urlAddress)
	if err != nil {
		log.Fatal("ERROR(queryMatchesData): ", err)
	}

	fmt.Printf("Sort %d records returned from the data source\n", len(data))
	// Sort data from the most recent to the earliest
	sort.Slice(data, func(i, j int) bool { return data[i].Date.After(data[j].Date) })

	// Set up database
	fmt.Println("Setting up the database connection...")
	db, err := setUpDatabase(dbPath)
	if err != nil {
		log.Fatal("ERROR (setUpDatabase): ", err)
	}

	fmt.Println("Inserting the data collected into the database...")
	if err = db.insertData(data); err != nil {
		log.Fatal("ERROR (insertData): ", err)
	}

	fmt.Printf("Fetching %d records starting with the most recent...\n", records)
	info, err := db.readData(records)
	if err != nil {
		log.Fatal("ERROR (readData): ", err)
	}

	strData, err := json.Marshal(info)
	if err != nil {
		log.Fatal("ERROR (info Marshal): ", err)
	}

	var out bytes.Buffer
	if err := json.Indent(&out, strData, "", "  "); err != nil {
		log.Fatal("ERROR (info Indent): ", err)
	}

	fmt.Printf(">>> Data: \n %s", out.Bytes())
}
