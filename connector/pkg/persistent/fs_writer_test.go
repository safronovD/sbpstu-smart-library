package persistent

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStructJSON struct {
	Field1 string   `json:"field1"`
	Field2 int      `json:"field2"`
	Field3 bool     `json:"field3"`
	Field4 []string `json:"field4"`
}

func TestSaveJSON(t *testing.T) {
	const (
		jsonPath = "/tmp/test.json"
	)

	// Create Test struct and set some fields
	testStruct := createTestStruct()

	// Write data to file
	jsonData, err := json.Marshal(testStruct)
	assert.Nil(t, err)

	err = SaveJSON(jsonData, jsonPath)
	assert.Nil(t, err)

	// Read data from file and check that structures are the same
	readTestStruct, err := readTestStruct(jsonPath)
	assert.Nil(t, err)
	assert.Equal(t, testStruct, readTestStruct)
}

func TestNewFileSystemWriter(t *testing.T) {
	const (
		jsonDir = "/tmp/json"
	)

	fsw1 := NewFileSystemWriter(jsonDir)
	fsw2 := NewFileSystemWriter(jsonDir)

	assert.Equal(t, fsw1, fsw2)
}

func TestFileSystemWriter_Write(t *testing.T) {
	const (
		jsonDir  = "/tmp/json"
		jsonId   = "id"
		jsonFile = jsonId + ".json"
	)

	// Create Test struct and set some fields
	testStruct := createTestStruct()

	// Write data to file
	jsonData, err := json.Marshal(testStruct)
	assert.Nil(t, err)

	fsw := NewFileSystemWriter(jsonDir)
	fsw.Write(jsonData, jsonId)

	// Read data from file and check that structures are the same
	fswJsonPath := path.Join(jsonDir, jsonFile)
	readTestStruct, err := readTestStruct(fswJsonPath)
	assert.Nil(t, err)
	assert.Equal(t, testStruct, readTestStruct)
}

func createTestStruct() *TestStructJSON {
	return &TestStructJSON{
		Field1: "string",
		Field2: 23,
		Field3: true,
		Field4: []string{"str1", "str2"},
	}
}

func readTestStruct(path string) (*TestStructJSON, error) {
	readData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	readTestStruct := &TestStructJSON{}
	if err = json.Unmarshal(readData, readTestStruct); err != nil {
		return nil, err
	}

	return readTestStruct, nil
}
