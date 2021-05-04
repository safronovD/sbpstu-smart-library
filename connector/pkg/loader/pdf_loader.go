package loader

import (
	"encoding/csv"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/spbstu-smart-library/connector/pkg/config"
)

func downloadFile(client *http.Client, cookies [2]http.Cookie, url string, path string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Failed to download " + url)
		}
	}()

	url += "/download"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Panic(err)
	}

	for i := range cookies {
		req.AddCookie(&cookies[i])
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, bt")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")

	rsp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		log.Panic(errors.New("Failed to GET " + url + " Code: " + rsp.Status))
	}

	file, err := os.Create(path)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		os.Remove(path)
		log.Panic(err)
	}

	if _, err := file.Write(data); err != nil {
		log.Println("Write failed")
	}

	log.Println(url + " downloaded")
}

func DownloadPDFFiles(config *config.PDFConfig, outputDir string) {
	commonPath := path.Join(".", outputDir, config.DB)
	pdfPath := path.Join(commonPath, config.Dir)

	if _, err := os.Stat(commonPath); err != nil {
		log.Panic(err)
	}

	if err := os.Mkdir(pdfPath, os.ModePerm); err != nil {
		log.Panic(err)
	}

	//TODO: Set in config full path to csv
	csvFile, err := os.Open(path.Join(outputDir, config.CsvFile))
	if err != nil {
		log.Panic(err)
	}
	defer csvFile.Close()
	csvReader := csv.NewReader(csvFile)

	httpClient := http.Client{
		Timeout: time.Second * 60, // Timeout after 2 seconds
	}

	aspxauth := config.Auth.ASPXAUTH
	sessionID := config.Auth.ASPNETSessionID

	var cookies [2]http.Cookie
	cookies[0] = http.Cookie{Name: ".ASPXAUTH", Value: aspxauth}
	cookies[1] = http.Cookie{Name: "ASP.NET_SessionId", Value: sessionID}

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Panic(err)
		}

		downloadFile(&httpClient, cookies, record[1], path.Join(pdfPath, record[0]+".pdf"))
	}
}
