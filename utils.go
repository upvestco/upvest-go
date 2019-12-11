package upvest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"path"
	"reflect"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
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

	result := joinPreservingTrailingSlash(p2...)

	u.Path = result

	return u, nil
}

// JoinPreservingTrailingSlash does a path.Join of the specified elements,
// preserving any trailing slash on the last non-empty segment
// credit: https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/util/net/http.go#L40
func joinPreservingTrailingSlash(elem ...string) string {
	// do the basic path join
	result := path.Join(elem...)

	// find the last non-empty segment
	for i := len(elem) - 1; i >= 0; i-- {
		if len(elem[i]) > 0 {
			// if the last segment ended in a slash, ensure our result does as well
			if strings.HasSuffix(elem[i], "/") && !strings.HasSuffix(result, "/") {
				result += "/"
			}
			break
		}
	}

	return result
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Second)
}

func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) //A=65 and Z = 65+25
	}
	return string(bytes)
}

// addOptions adds the parameters in opt as URL query parameters to s.
// opt  must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
