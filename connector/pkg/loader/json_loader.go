package loader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/spbstu-smart-library/connector/pkg/config"
	"github.com/spbstu-smart-library/connector/pkg/persistent"
)

type Record struct {
	RecordIdentifier string `json:"recordIdentifier"`
}

type RecordsList struct {
	InnerRecordsList struct {
		Records []Record `json:"record"`
	} `json:"records"`

	NumberOfRecords    int `json:"numberOfRecords"`
	NextRecordPosition int `json:"nextRecordPosition"`
}

func downloadJSON(client *http.Client, req *http.Request) ([]byte, error) {
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if (res.StatusCode != http.StatusOK) && (res.StatusCode != http.StatusMultiStatus) {
		return nil, fmt.Errorf("Response failed: %s Status code: %d", res.Request.URL, res.StatusCode)
	}

	if res.Body != nil {
		defer res.Body.Close()
	} else {
		return nil, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func DownloadRecords(config *config.JSONConfig, outputDir string) {
	var writers []persistent.Writer

	if config.Output.Elasticsearch.Enable {
		esConf := config.Output.Elasticsearch
		esWriter, err := persistent.NewESWriter(esConf.Host, esConf.Login, esConf.Login, esConf.Index)
		if err != nil {
			log.Printf("ES connection failed: %s", err)
		} else {
			writers = append(writers, esWriter)
		}
	}

	if config.Output.CsvFile.Enable {
		csvWriter, err := persistent.NewCSVWriter(path.Join(outputDir, config.Output.CsvFile.FileName))
		if err != nil {
			log.Printf("CSV Writer failed: %s", err)
		} else {
			writers = append(writers, csvWriter)
		}
	}

	if config.Output.CsvFile.Enable {
		csvWriter, err := persistent.NewCSVWriter(path.Join(outputDir, config.Output.CsvFile.FileName))
		if err != nil {
			log.Printf("CSV Writer failed: %s", err)
		} else {
			writers = append(writers, csvWriter)
		}
	}

	if config.Output.FileSystem.Enable {
		fsWriter := persistent.NewFileSystemWriter(path.Join(outputDir, config.Output.FileSystem.JSONDir))
		writers = append(writers, fsWriter)
	}

	if config.Output.ConvertEnable {

	}

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	retryClient.RetryWaitMin = 10 * time.Second
	retryClient.RetryWaitMax = 30 * time.Second
	httpClient := retryClient.StandardClient()

	connectURL := config.Connection.URL + config.Connection.DB

	n := 0
	maxDownloads := config.Connection.DownloadListMaxsize

	for n < maxDownloads {
		req, err := http.NewRequest(http.MethodGet, connectURL, nil)
		if err != nil {
			log.Panic(err)
		}

		req.Header.Set("Accept", "application/json")

		q := req.URL.Query()
		q.Add("query", config.Connection.Query)

		if len(config.Connection.Fcq) != 0 {
			q.Add("fcq", config.Connection.Fcq)
		}

		q.Add("maximumRecords", strconv.Itoa(config.Connection.DownloadBatchSize))
		q.Add("startRecord", strconv.Itoa(n+1))
		req.URL.RawQuery = q.Encode()

		data, err := downloadJSON(httpClient, req)
		if err != nil {
			log.Panic(err)
		}

		rl := RecordsList{}

		err = json.Unmarshal(data, &rl)
		if err != nil {
			log.Panic(err)
		}

		if rl.NumberOfRecords < maxDownloads {
			maxDownloads = rl.NumberOfRecords
		}

		log.Println(fmt.Sprintf("Start to download [%d-%d]/%d", n+1, n+len(rl.InnerRecordsList.Records), maxDownloads))

		for _, val := range rl.InnerRecordsList.Records {
			n++

			downloadURL := connectURL + "/" + url.PathEscape(val.RecordIdentifier)

			req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
			if err != nil {
				log.Panic(err)
			}

			req.Header.Set("Accept", "application/json")

			q := req.URL.Query()
			q.Add("recordSchema", "gost-7.0.100")
			req.URL.RawQuery = q.Encode()

			jsonData, err := downloadJSON(httpClient, req)
			if err != nil {
				log.Println(err)
				continue
			}

			formattedID := formatID(val.RecordIdentifier)

			if config.Output.ConvertEnable {
				convertedJSON, err := convertJSONData(jsonData)
				if err != nil {
					log.Printf("JSON converted failed: %s", err)
				} else {
					jsonData = convertedJSON
				}
			}

			for _, w := range writers {
				w.Write(jsonData, formattedID)
			}
		}

		log.Printf(fmt.Sprintf("Downloaded %d/%d. Next record number is %d", n, maxDownloads, rl.NextRecordPosition))
	}
}

func formatID(id string) string {
	return strings.ReplaceAll(id, "\\", "_")
}

// TODO refactor with Golang
func convertJSONData(jsonData []byte) ([]byte, error) {
	err := persistent.SaveJSON(jsonData, "./tmp1.json")
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("python3", "./utils/json_converter3.py", "./tmp1.json", "./res.json")
	if err = cmd.Run(); err != nil {
		return nil, err
	}

	convertedJSON, err := ioutil.ReadFile("./res.json")
	if err != nil {
		return nil, err
	}

	return convertedJSON, nil
}
