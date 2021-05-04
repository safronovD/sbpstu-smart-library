package loader

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/spbstu-smart-library/connector/pkg/persistent"
)

func DownloadSamples(outputDir string) {
	connectURL := "https://ruslan.library.spbstu.ru/rrs-web/db/"
	httpClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	commonPath := path.Join(".", outputDir, "samples")
	if err := os.Mkdir(commonPath, os.ModePerm); err != nil {
		log.Panic(err)
	}

	req, err := http.NewRequest("PROPFIND", connectURL, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Accept", "application/json")

	data, err := downloadJSON(&httpClient, req)
	if err != nil {
		log.Panic(err)
	}

	persistent.SaveJSON(data, path.Join(commonPath, "db_list"+".json"))

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		log.Panic(err)
	}

	dbList := result["response"].([]interface{})

	for _, db := range dbList {
		href := db.(map[string]interface{})["href"].([]interface{})[0].(string)
		splitedHref := strings.Split(href, "/")
		dbName := splitedHref[len(splitedHref)-1]

		dbPath := path.Join(commonPath, dbName)
		if err := os.Mkdir(dbPath, os.ModePerm); err != nil {
			log.Panic(err)
		}

		if dbName == "db" {
			continue
		}

		req, err := http.NewRequest(http.MethodGet, href, nil)
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("Accept", "application/json")

		data, err := downloadJSON(&httpClient, req)
		if err != nil {
			log.Panic(err)
		}

		persistent.SaveJSON(data, path.Join(dbPath, "db_index"+".json"))

		req, err = http.NewRequest(http.MethodGet, href, nil)
		if err != nil {
			log.Panic(err)
		}

		req.Header.Set("Accept", "application/json")

		q := req.URL.Query()
		q.Add("query", "(bib.volume=*)")
		q.Add("maximumRecords", "2")
		q.Add("startRecord", "1")
		req.URL.RawQuery = q.Encode()

		data, err = downloadJSON(&httpClient, req)
		if err != nil {
			log.Panic(err)
		}

		persistent.SaveJSON(data, path.Join(dbPath, "db_records_list"+".json"))

		rl := RecordsList{}

		err = json.Unmarshal(data, &rl)
		if err != nil {
			log.Panic(err)
		}

		log.Printf("DB: %s Number of records: %d\n", dbName, rl.NumberOfRecords)

		for i, val := range rl.InnerRecordsList.Records {
			downloadURL := href + "/" + url.PathEscape(val.RecordIdentifier)

			req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
			if err != nil {
				log.Panic(err)
			}

			req.Header.Set("Accept", "application/json")

			q := req.URL.Query()
			q.Add("recordSchema", "gost-7.0.100")
			req.URL.RawQuery = q.Encode()

			jsonData, err := downloadJSON(&httpClient, req)
			if err != nil {
				log.Println(err)
				continue
			}

			persistent.SaveJSON(jsonData, path.Join(dbPath, "record"+strconv.Itoa(i)+".json"))
		}
	}
}
