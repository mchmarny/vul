package scanner

import (
	"fmt"
	"os/exec"
	"path"
)

const (
	trivyReportName = "trivy.json"
	snykReportName  = "snyk.json"
	grypeReportName = "grype.json"
)

type scannerCmd struct {
	name string
	path string
	cmd  *exec.Cmd
}

func (c *scannerCmd) String() string {
	return fmt.Sprintf("name:%s, target:%s, cmd:%s", c.name, c.path, c.cmd.String())
}

func makeScannerCommands(digest, targetPath string) []*scannerCmd {
	return []*scannerCmd{
		makeTrivyCmd(digest, targetPath),
		makeSnykCmd(digest, targetPath),
		makeGrypeCmd(digest, targetPath),
	}
}

func makeTrivyCmd(digest, targetPath string) *scannerCmd {
	t := path.Join(targetPath, trivyReportName)
	c := exec.Command("trivy", "image", "--quiet", "--security-checks", "vuln", "--format", "json", "--no-progress", "--output", t, digest)

	return &scannerCmd{
		name: "trivy",
		path: t,
		cmd:  c,
	}
}

func makeSnykCmd(digest, targetPath string) *scannerCmd {
	t := path.Join(targetPath, snykReportName)
	jfo := fmt.Sprintf("--json-file-output=%s", t)
	c := exec.Command("snyk", "container", "test", "--app-vulns", jfo, digest)

	return &scannerCmd{
		name: "snyk",
		path: t,
		cmd:  c,
	}
}

func makeGrypeCmd(digest, targetPath string) *scannerCmd {
	t := path.Join(targetPath, grypeReportName)
	c := exec.Command("grype", "-q", "--add-cpes-if-none", "-s", "AllLayers", "-o", "json", "--file", t, digest)

	return &scannerCmd{
		name: "trivy",
		path: t,
		cmd:  c,
	}
}
