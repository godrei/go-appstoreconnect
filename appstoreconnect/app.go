package appstoreconnect

import "net/http"

// AppResponse ...
type AppResponse struct {
	App   App           `json:"data,omitempty"`
	Links DocumentLinks `json:"links,omitempty"`
}

// App ...
func (s TestFlightService) App(id string) (*AppResponse, *http.Response, error) {
	u := "apps/" + id
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	r := &AppResponse{}
	resp, err := s.client.do(req, r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, err
}
