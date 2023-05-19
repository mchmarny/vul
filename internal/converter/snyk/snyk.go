package snyk

import (
	"github.com/Jeffail/gabs/v2"
	"github.com/mchmarny/vul/internal/data"
	"github.com/mchmarny/vul/internal/parser"
	"github.com/pkg/errors"
)

// Convert converts JSON to a list of common vulnerabilities.
func Convert(c *gabs.Container) ([]*data.Vulnerability, error) {
	if c == nil {
		return nil, errors.New("source required")
	}

	v := c.Search("vulnerabilities")
	if !v.Exists() {
		return nil, errors.New("unable to find vulnerabilities in source data")
	}

	list := make([]*data.Vulnerability, 0)

	for _, r := range v.Children() {
		vul := mapVulnerability(r)
		if vul == nil {
			continue
		}

		list = append(list, vul)
	}

	return list, nil
}

func mapVulnerability(v *gabs.Container) *data.Vulnerability {
	c := v.Search("cvssDetails")
	if !c.Exists() {
		return nil
	}

	item := &data.Vulnerability{
		Exposure: parser.ToString(v.Search("identifiers", "CVE").Index(0).Data()),
		Package:  parser.String(v, "name"),
		Version:  parser.String(v, "version"),
		Severity: parser.GetFirstString(v, "nvdSeverity", "severity"),
		Score:    parser.ToFloat32(v.Search("cvssScore").Data()),
		IsFixed:  parser.ToBool(v.Search("isUpgradable").Data()),
	}

	if item.Score == 0 {
		item.Score = parser.ToFloat32(c.Search("baseScore").Data())
	}

	return item
}