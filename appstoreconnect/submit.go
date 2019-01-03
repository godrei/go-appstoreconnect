package appstoreconnect

import (
	"net/http"
)

// BetaAppReviewSubmission ...
func (s TestFlightService) BetaAppReviewSubmission(buildID string) (*http.Response, error) {
	u := "betaAppReviewSubmissions"

	type RequestDataRelationshipsBuildData struct {
		ID   string `json:"id,omitempty"`
		Type string `json:"type,omitempty"`
	}

	type RequestDataRelationshipsBuild struct {
		Data RequestDataRelationshipsBuildData `json:"data,omitempty"`
	}

	type RequestDataRelationships struct {
		Build RequestDataRelationshipsBuild `json:"build,omitempty"`
	}

	type RequestData struct {
		Relationships RequestDataRelationships `json:"relationships,omitempty"`
		Type          string                   `json:"type,omitempty"`
	}
	type Request struct {
		Data RequestData `json:"data,omitempty"`
	}

	r := Request{
		Data: RequestData{
			Type: "betaAppReviewSubmissions",
			Relationships: RequestDataRelationships{
				Build: RequestDataRelationshipsBuild{
					Data: RequestDataRelationshipsBuildData{
						ID:   buildID,
						Type: "builds",
					},
				},
			},
		},
	}
	req, err := s.client.NewRequest(http.MethodPost, u, r)
	if err != nil {
		return nil, err
	}

	return s.client.do(req, nil)
}
