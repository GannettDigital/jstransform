package nonest

// Code generated by github.com/GannettDigital/jstransform; DO NOT EDIT.

import "time"

type Complex struct {
	Simple

	Caption string `json:"caption"`
	Credit  string `json:"credit"`
	// The available cropped images
	Crops          []ComplexCrops      `json:"crops"`
	Cutline        string              `json:"cutline,omitempty"`
	DatePhotoTaken time.Time           `json:"datePhotoTaken"`
	Orientation    string              `json:"orientation"`
	OriginalSize   ComplexOriginalSize `json:"originalSize"`
	// a type
	Type string `json:"type"`
	// Universal Resource Locator
	URL ComplexURL `json:"URL"`
}

type ComplexCrops struct {
	Height float64 `json:"height"`
	Name   string  `json:"name"`
	// full path to the cropped image file
	Path string `json:"path"`
	// a long
	// multi-line description
	RelativePath string  `json:"relativePath"`
	Width        float64 `json:"width"`
}

type ComplexOriginalSize struct {
	Height float64 `json:"height"`
	Width  float64 `json:"width"`
}

type ComplexURL struct {
	// The full Canonical URL
	Absolute string          `json:"absolute"`
	Meta     *ComplexURLMeta `json:"meta,omitempty"`
	Publish  string          `json:"publish"`
}

type ComplexURLMeta struct {
	Description string `json:"description"`
	SiteName    string `json:"siteName"`
}
