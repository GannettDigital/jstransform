package nonest

// Code generated by github.com/GannettDigital/jstransform; DO NOT EDIT.

type Nested_to_primitive struct {
	Categories []Nested_to_primitiveCategories `json:"categories"`
}

type Nested_to_primitiveCategories struct {
	Category    string   `json:"category,omitempty"`
	MarketTypes []string `json:"marketTypes,omitempty"`
}
