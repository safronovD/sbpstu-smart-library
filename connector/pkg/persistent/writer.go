package persistent

type Writer interface {
	Write(jsonData []byte, recordId string)
}
