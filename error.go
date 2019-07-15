package upvest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// APIError represents an error response from the Paystack API server
type APIError map[string]interface{}

// APIError supports the error interface
func (aerr *APIError) Error() string {
	ret, _ := json.Marshal(aerr)
	return string(ret)
}

func newAPIError(resp *http.Response) *APIError {
	p, _ := ioutil.ReadAll(resp.Body)

	var upvestErrorResp APIError
	_ = json.Unmarshal(p, &upvestErrorResp)
	return &upvestErrorResp
}
