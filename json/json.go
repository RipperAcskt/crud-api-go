package json

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

type Person struct {
	Id      int
	Name    string `json:"name"`
	Surname string `json:"lastName"`
	Age     int    `json:"age"`
}

func JsonUnmarshal(body io.ReadCloser, s *[]Person) error {
	readedBody, err := ioutil.ReadAll(body)
	if err != nil {
		return fmt.Errorf("readAll faild: %v", err)
	}

	err = json.Unmarshal(readedBody, &s)
	if err != nil {
		return fmt.Errorf("unmarshal faild: %v", err)
	}

	return nil
}

func JsonMarshalResponse(p []Person) ([]byte, error) {
	jsonResp, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return jsonResp, nil
}

func JsonMarshalError(errText string) ([]byte, error) {
	message := make(map[string]string)
	message["Error"] = errText
	jsonResp, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	return jsonResp, nil
}
