package persistent

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"path"
)

type FileSystemWriter struct {
	jsonPath string
}

func NewFileSystemWriter(jsonPath string) *FileSystemWriter {
	return &FileSystemWriter{
		jsonPath: jsonPath,
	}
}

func (w *FileSystemWriter) Write(jsonData []byte, recordId string) {
	if err := SaveJSON(jsonData, path.Join(w.jsonPath, recordId+".json")); err != nil {
		log.Printf("Failed to save json with err: %s", err)
	}
}

func SaveJSON(data []byte, path string) error {
	jsonFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, data, "", "    ")
	if err != nil {
		return err
	}

	_, err = prettyJSON.WriteTo(jsonFile)
	if err != nil {
		os.Remove(jsonFile.Name())
		return err
	}

	log.Printf("Json file \"%s\" saved", jsonFile.Name())
	return nil
}
