package persistent

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	href     = "some_link"
	csv_path = "/tmp/test.csv"
	recordID = "id"
)

type TestStructCSV struct {
	PdfLink string `json:"pdfLink"`
}

func TestGetHref(t *testing.T) {
	testStruct := &TestStructCSV{PdfLink: href}
	jsonData, err := json.Marshal(testStruct)
	assert.Nil(t, err)

	recognizedHref, err := getHref(jsonData)
	assert.Nil(t, err)
	assert.Equal(t, href, recognizedHref)
}

func TestGetHref_Error(t *testing.T) {
	type CorruptedTestStruct struct {
		NotPdfLink string `json:"notPdfLink"`
	}

	// JSON contains no pdfLink field
	testStruct := &CorruptedTestStruct{NotPdfLink: href}
	jsonData, err := json.Marshal(testStruct)
	assert.Nil(t, err)

	_, err = getHref(jsonData)
	assert.NotNil(t, err)

	// Given data is not JSON
	corruptedData := []byte("some_bytes")
	_, err = getHref(corruptedData)
	assert.NotNil(t, err)
}

func TestNewCSVWriter(t *testing.T) {
	csvw, err := NewCSVWriter(csv_path)
	assert.Nil(t, err)
	assert.NotNil(t, csvw)
}

func TestCSVWriter_Write(t *testing.T) {
	testStruct := &TestStructCSV{PdfLink: href}
	jsonData, err := json.Marshal(testStruct)
	assert.Nil(t, err)

	csvw, err := NewCSVWriter(csv_path)
	assert.Nil(t, err)

	csvw.Write(jsonData, recordID)

	csvFile, err := os.Open(csv_path)
	assert.Nil(t, err)
	r := csv.NewReader(csvFile)

	records, err := r.ReadAll()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(records))
	assert.Equal(t, []string{recordID, href}, records[0])
}
