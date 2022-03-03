package generated

// Code generated by github.com/GannettDigital/jstransform; DO NOT EDIT.

type Nested struct {
	FactCheckClaims []struct {
		AppearanceURLs []struct {
			Original bool   `json:"original,omitempty" description:"The original flag"`
			Url      string `json:"url" description:"The appearance url"`
		} `json:"appearanceURLs,omitempty" description:"The factCheck appearanceURLs"`
		Author string `json:"author,omitempty" description:"The factCheck author"`
		Claim  string `json:"claim,omitempty" description:"The factCheck claim"`
		Date   string `json:"date,omitempty" description:"The date of the factCheck"`
		Rating string `json:"rating,omitempty" description:"The factCheck rating"`
	} `json:"factCheckClaims" description:"Information related to the International Fact-Checking Network (IFCN) program"`
}
