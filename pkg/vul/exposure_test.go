package vul

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExposureUniquenessFalse(t *testing.T) {
	p := &PackageVersionExposures{
		Sources: map[string]*SourceExposures{
			"snyk": {
				Exposures: map[string]*Exposures{
					"CVE-2011-3374": {
						Severity: "high",
						Score:    1.0,
					},
				},
			},
			"grype": {
				Exposures: map[string]*Exposures{
					"CVE-2011-3374": {
						Severity: "high",
						Score:    1.0,
					},
				},
			},
			"trivy": {
				Exposures: map[string]*Exposures{
					"CVE-2011-3374": {
						Severity: "high",
						Score:    1.0,
					},
				},
			},
		},
	}

	assert.False(t, p.UniqueExposures())
}

func TestExposureUniquenessTrue(t *testing.T) {
	p := &PackageVersionExposures{
		Sources: map[string]*SourceExposures{
			"snyk": {
				Exposures: map[string]*Exposures{
					"CVE-2011-3374": {
						Severity: "high",
						Score:    1.0,
					},
				},
			},
			"grype": {
				Exposures: map[string]*Exposures{
					"CVE-2011-3374": {
						Severity: "high",
						Score:    1.0,
					},
				},
			},
			"trivy": {
				Exposures: map[string]*Exposures{
					"CVE-2011-3374": {
						Severity: "high",
						Score:    1.1,
					},
				},
			},
		},
	}

	assert.True(t, p.UniqueExposures())
}
