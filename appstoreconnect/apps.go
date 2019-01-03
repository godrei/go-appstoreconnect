package appstoreconnect

import (
	"net/http"
	"net/url"

	"github.com/bitrise-io/go-utils/log"
)

// AppsResponse ...
type AppsResponse struct {
	Apps []App `json:"data,omitempty"`

	Links PagedDocumentLinks `json:"links,omitempty"`
	Meta  PagingInformation  `json:"meta,omitempty"`
}

// AppsOptions ...
type AppsOptions struct {
	BundleIDFilter string `url:"filter[bundleId],omitempty"`
	IDFilter       string `url:"filter[id],omitempty"`
	NameFilter     string `url:"filter[name],omitempty"`
	SKUFilter      string `url:"filter[sku],omitempty"`

	Limit  int    `url:"limit,omitempty"`
	Cursor string `url:"cursor,omitempty"`
	Next   string `url:"-"`
}

// Apps ...
func (s TestFlightService) Apps(opt *AppsOptions) (*AppsResponse, *http.Response, error) {
	if opt != nil && opt.Next != "" {
		u, err := url.Parse(opt.Next)
		if err != nil {
			return nil, nil, err
		}
		cursor := u.Query().Get("cursor")
		log.Debugf("cursor: %s", cursor)
		opt.Cursor = cursor
	}

	u := "apps"
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	r := &AppsResponse{}
	resp, err := s.client.do(req, r)
	if err != nil {
		return nil, resp, err
	}

	return r, resp, err
}
