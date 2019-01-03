package appstoreconnect

import "time"

// TestFlightService ...
type TestFlightService service

// Paging ...
type Paging struct {
	Total int `json:"total,omitempty"`
	Limit int `json:"limit,omitempty"`
}

// PagingInformation ...
type PagingInformation struct {
	Paging Paging `json:"paging,omitempty"`
}

// DocumentLinks ...
type DocumentLinks struct {
	Self string `json:"self,omitempty"`
}

// PagedDocumentLinks ...
type PagedDocumentLinks struct {
	First string `json:"first,omitempty"`
	Next  string `json:"next,omitempty"`
	Self  string `json:"self,omitempty"`
}

// RelationshipLinks ...
type RelationshipLinks struct {
	Self    string `json:"self,omitempty"`
	Related string `json:"related,omitempty"`
}

// ResourceLinks ...
type ResourceLinks struct {
	Self string `json:"self,omitempty"`
}

// RelationshipData ...
type RelationshipData struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

// Relationship ...
type Relationship struct {
	Data  RelationshipData  `json:"data,omitempty"`
	Links RelationshipLinks `json:"links,omitempty"`
	Meta  PagingInformation `json:"meta,omitempty"`
}

// AppRelationships ...
type AppRelationships struct {
	BetaLicenseAgreement Relationship `json:"betaLicenseAgreement,omitempty"`
	PreReleaseVersions   Relationship `json:"preReleaseVersions,omitempty"`
	BetaAppLocalizations Relationship `json:"betaAppLocalizations,omitempty"`
	BetaGroups           Relationship `json:"betaGroups,omitempty"`
	BetaTesters          Relationship `json:"betaTesters,omitempty"`
	Builds               Relationship `json:"builds,omitempty"`
	BetaAppReviewDetail  Relationship `json:"betaAppReviewDetail,omitempty"`
}

// AppAttributes ...
type AppAttributes struct {
	BundleID      string `json:"bundleId,omitempty"`
	Name          string `json:"name,omitempty"`
	PrimaryLocale string `json:"primaryLocale,omitempty"`
	Sku           string `json:"sku,omitempty"`
}

// App ...
type App struct {
	Attributes    AppAttributes    `json:"attributes,omitempty"`
	ID            string           `json:"id,omitempty"`
	Type          string           `json:"type,omitempty"`
	Relationships AppRelationships `json:"relationships,omitempty"`
	Links         ResourceLinks    `json:"links,omitempty"`
}

// ImageAsset ...
type ImageAsset struct {
	TemplateURL string `json:"templateUrl,omitempty"`
	Height      int    `json:"height,omitempty"`
	Width       int    `json:"width,omitempty"`
}

// BuildAttributes ...
type BuildAttributes struct {
	Expired                 bool       `json:"expired,omitempty"`
	IconAssetToken          ImageAsset `json:"iconAssetToken,omitempty"`
	MinOsVersion            string     `json:"minOsVersion,omitempty"`
	ProcessingState         string     `json:"processingState,omitempty"`
	Version                 string     `json:"version,omitempty"`
	UsesNonExemptEncryption bool       `json:"usesNonExemptEncryption,omitempty"`
	UploadedDate            time.Time  `json:"uploadedDate,omitempty"`
	ExpirationDate          time.Time  `json:"expirationDate,omitempty"`
}

// BuildRelationship ...
type BuildRelationship struct {
	App                      Relationship `json:"app,omitempty"`
	AppEncryptionDeclaration Relationship `json:"appEncryptionDeclaration,omitempty"`
	IndividualTesters        Relationship `json:"individualTesters ,omitempty"`
	PreReleaseVersion        Relationship `json:"preReleaseVersion ,omitempty"`
	BetaBuildLocalizations   Relationship `json:"betaBuildLocalizations ,omitempty"`
	BetaGroups               Relationship `json:"betaGroups ,omitempty"`
	BuildBetaDetail          Relationship `json:"buildBetaDetail ,omitempty"`
	BetaAppReviewSubmission  Relationship `json:"betaAppReviewSubmission ,omitempty"`
}

// Build ...
type Build struct {
	Attributes    BuildAttributes   `json:"attributes,omitempty"`
	ID            string            `json:"id,omitempty"`
	Relationships BuildRelationship `json:"relationships,omitempty"`
	Type          string            `json:"type,omitempty"`
	Links         ResourceLinks     `json:"links,omitempty"`
}
