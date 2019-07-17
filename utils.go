package upvest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"path"

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

func joinURLs(basePath string, paths ...string) (*url.URL, error) {
	u, err := url.Parse(basePath)

	if err != nil {
		return nil, fmt.Errorf("invalid url")
	}

	p2 := append([]string{u.Path}, paths...)

	result := path.Join(p2...)

	u.Path = result

	return u, nil
}
