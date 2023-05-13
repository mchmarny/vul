package vul

import (
	"fmt"
)

type ImageDigestExposures struct {
	Image    string                       `json:"image"`
	Digest   string                       `json:"digest"`
	Packages map[string]*PackageExposures `json:"packages"`

	*Item `json:"-"`
}

type PackageExposures struct {
	Versions map[string]*PackageVersionExposures `json:"versions"`

	*Item `json:"-"`
}

type PackageVersionExposures struct {
	Sources map[string]*SourceExposures `json:"sources"`

	*Item `json:"-"`
}

func (p *PackageVersionExposures) UniqueExposures() bool {
	if p.Sources == nil || len(p.Sources) == 0 {
		return false
	}

	m := make(map[string]bool, 0)

	for _, sev := range p.Sources {
		for e, exp := range sev.Exposures {
			m[fmt.Sprintf("%s-%s", e, exp.severityScoreHash())] = true
		}
	}

	return len(m) > 1
}

type SourceExposures struct {
	Exposures map[string]*Exposures `json:"exposures"`

	*Item `json:"-"`
}

type Exposures struct {
	Severity string  `json:"severity"`
	Score    float64 `json:"score"`
	Fixed    bool    `json:"fixed"`

	*Item `json:"-"`
}

func (e *Exposures) severityScoreHash() string {
	return fmt.Sprintf("%s-%.2f", e.Severity, e.Score)
}
