package appstoreconnect

import (
	"net/http"
	"net/url"

	"github.com/bitrise-io/go-utils/log"
)

// BuildsResponse ...
type BuildsResponse struct {
	Builds []Build `json:"data,omitempty"`

	Links PagedDocumentLinks `json:"links,omitempty"`
	Meta  PagingInformation  `json:"meta,omitempty"`
}

// BuildsOptions ...
type BuildsOptions struct {
	AppFilter             string `url:"filter[app],omitempty"`
	ExpiredFilter         string `url:"filter[expired],omitempty"`
	IDFilter              string `url:"filter[id],omitempty"`
	ProcessingStateFilter string `url:"filter[processingState],omitempty"`
	VersionFilter         string `url:"filter[version],omitempty"`

	Limit  int    `url:"limit,omitempty"`
	Cursor string `url:"cursor,omitempty"`
	Next   string `url:"-"`
}

// Builds ...
func (s TestFlightService) Builds(opt *BuildsOptions) (*BuildsResponse, *http.Response, error) {
	if opt != nil && opt.Next != "" {
		u, err := url.Parse(opt.Next)
		if err != nil {
			return nil, nil, err
		}
		cursor := u.Query().Get("cursor")
		log.Debugf("cursor: %s", cursor)
		opt.Cursor = cursor
	}

	u := "builds"
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	r := &BuildsResponse{}
	resp, err := s.client.do(req, r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, err
}
