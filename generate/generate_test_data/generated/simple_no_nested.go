package generated

// Code generated by github.com/GannettDigital/jstransform; DO NOT EDIT.

import "time"

type Simple_no_nested struct {
	Contributors []Simple_no_nestedContributors `json:"contributors,omitempty"`
	Height       int64                          `json:"height,omitempty"`
	SomeDateObj  Simple_no_nestedSomeDateObj    `json:"someDateObj,omitempty"`
	Visible      bool                           `json:"visible,omitempty"`
	Width        float64                        `json:"width,omitempty"`
}

type Simple_no_nestedContributors struct {
	ContributorId string `json:"contributorId,omitempty"`
	Id            string `json:"id"`
	Name          string `json:"name"`
}

type Simple_no_nestedSomeDateObj struct {
	Dates []time.Time `json:"dates,omitempty"`
}
