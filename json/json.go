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

type DeleteInfo struct {
	Id             int
	DeleteAllTable bool `json:"all"`
}

func JsonUnmarshal(body io.ReadCloser, s interface{}, delete bool) string {
	readedBody, err := ioutil.ReadAll(body)
	if err != nil {
		return fmt.Sprintf("Error while reading body: %v", err)
	}
	if delete {
		err = json.Unmarshal(readedBody, s.(*DeleteInfo))
	} else {
		err = json.Unmarshal(readedBody, s.(*Person))
	}

	if err != nil {
		return fmt.Sprintf("Error while unmarshal: %v", err)
	}

	return ""
}

func JsonMarshalResponse(p []Person) ([]byte, error) {
	jsonResp, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return jsonResp, nil
}
