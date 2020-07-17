package json

import (
	"encoding/json"
	"os"

	"github.com/avegner/utils/files"
)

func UnmarshalFile(path string, v interface{}) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer files.Close(f, &err)

	d := json.NewDecoder(f)
	return d.Decode(v)
}

func MarshalFile(v interface{}, path string) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer files.Close(f, &err)

	e := json.NewEncoder(f)
	return e.Encode(v)
}
