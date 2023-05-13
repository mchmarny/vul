package vul

import "time"

type ImageVersion struct {
	Image     string    `json:"image"`
	Digest    string    `json:"digest"`
	Processed time.Time `json:"processed"`

	*Item `json:"-"`
}
