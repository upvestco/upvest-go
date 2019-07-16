package upvest

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/mitchellh/mapstructure"
)

func mapstruct(data interface{}, v interface{}) error {
	config := &mapstructure.DecoderConfig{
		Result:           v,
		TagName:          "json",
		WeaklyTypedInput: true,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	err = decoder.Decode(data)
	return err
}

func jsonEncode(data interface{}) (io.ReadWriter, error) {
	var buf io.ReadWriter
	buf = new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(data)
	return buf, err
}
