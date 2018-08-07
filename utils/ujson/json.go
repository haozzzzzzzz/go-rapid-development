package ujson

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func UnmarshalJsonFromReader(reader io.Reader, v interface{}) error {

	bytes, err := ioutil.ReadAll(reader)

	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, v)

	if err != nil {
		return err
	}

	return nil
}
