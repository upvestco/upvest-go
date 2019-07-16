package upvest

type service struct {
	client *Client
	auth   AuthProvider
}

// Get makes an HTTP GET request to Upvest API
func (s *service) Get(path string, body, v interface{}, auth AuthProvider) error {
	return s.client.Call("GET", path, body, v, s.auth)
}

// Post makes an HTTP Post request to Upvest API
func (s *service) Post(path string, body, v interface{}, auth AuthProvider) error {
	return s.client.Call("POST", path, body, v, auth)
}

// Patch makes an HTTP PATCH request to Upvest API
func (s *service) Patch(path string, body, v interface{}, auth AuthProvider) error {
	return s.client.Call("PATCH", path, body, v, s.auth)
}

// DELETE makes an HTTP DELETE request to Upvest API
func (s *service) Delete(path string, body, v interface{}, auth AuthProvider) error {
	return s.client.Call("DELETE", path, body, v, s.auth)
}
