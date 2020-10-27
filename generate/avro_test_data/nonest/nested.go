package nonest

// Code generated by github.com/GannettDigital/jstransform; DO NOT EDIT.

type Nested struct {
	// Information related to the International Fact-Checking Network (IFCN) program
	FactCheckClaims []NestedFactCheckClaims `json:"factCheckClaims"`
}

type NestedFactCheckClaims struct {
	// The factCheck appearanceURLs
	AppearanceURLs []NestedFactCheckClaimsAppearanceURLs `json:"appearanceURLs,omitempty"`
	// The factCheck author
	Author string `json:"author,omitempty"`
	// The factCheck claim
	Claim string `json:"claim,omitempty"`
	// The date of the factCheck
	Date string `json:"date,omitempty"`
	// The factCheck rating
	Rating string `json:"rating,omitempty"`
}

type NestedFactCheckClaimsAppearanceURLs struct {
	// The original flag
	Original bool `json:"original,omitempty"`
	// The appearance url
	Url string `json:"url"`
}
