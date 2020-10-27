package avro_test_data

// Code generated by github.com/GannettDigital/jstransform; DO NOT EDIT.

import "time"

type TimesPointers struct {
	// Information related to the International Fact-Checking Network (IFCN) program
	FactCheckClaims []*TimesPointersFactCheckClaims `json:"factCheckClaims,omitempty"`
	NonRequiredDate *time.Time                      `json:"nonRequiredDate,omitempty"`
	RequiredDate    time.Time                       `json:"requiredDate"`
}

type TimesPointersFactCheckClaims struct {
	// The factCheck appearanceURLs
	AppearanceURLs []*TimesPointersFactCheckClaimsAppearanceURLs `json:"appearanceURLs,omitempty"`
	// The factCheck author
	Author string `json:"author,omitempty"`
	// The factCheck claim
	Claim string `json:"claim,omitempty"`
	// The date of the factCheck
	Date string `json:"date,omitempty"`
	// The factCheck rating
	Rating string `json:"rating,omitempty"`
}

type TimesPointersFactCheckClaimsAppearanceURLs struct {
	// The original flag
	Original bool `json:"original,omitempty"`
	// The appearance url
	Url string `json:"url"`
}
