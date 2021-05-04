package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path"

	"github.com/spbstu-smart-library/connector/pkg/config"
	"github.com/spbstu-smart-library/connector/pkg/loader"
)

func main() {
	launchMod := flag.String("launch-mod", "download-json", "a string")
	configFileName := flag.String("config-file", "config.yaml", "a string")
	logFileName := flag.String("log-file", "connector.log", "a string")
	outputDir := flag.String("output-dir", "output", "a string")
	flag.Parse()

	if _, err := os.Stat(*outputDir); os.IsNotExist(err) {
		if err := os.Mkdir(*outputDir, os.ModePerm); err != nil {
			log.Panic(err)
		}
	}

	logFilePath := path.Join(".", *outputDir, *logFileName)

	logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	config, err := config.NewConfig(*configFileName)
	if err != nil {
		log.Panic(err)
	}

	switch *launchMod {
	case "download-json":
		loader.DownloadRecords(&config.JSONConfig, *outputDir)
	case "download-pdf":
		loader.DownloadPDFFiles(&config.PDFConfig, *outputDir)
	case "samples":
		loader.DownloadSamples(*outputDir)
	default:
		log.Panic("Launch mod is not correct")
	}
}
