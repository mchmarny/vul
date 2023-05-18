package scanner

import (
	"fmt"
	"os/exec"
)

func makeTrivyCmd(digest, targetPath string) *exec.Cmd {
	return exec.Command("trivy", "image", "--quiet", "--security-checks", "vuln", "--format", "json", "--no-progress", "--output", targetPath, digest)
}

func makeSnykCmd(digest, targetPath string) *exec.Cmd {
	jfo := fmt.Sprintf("--json-file-output=%s", targetPath)
	return exec.Command("snyk", "container", "test", "--app-vulns", jfo, digest)
}

func makeGrypeCmd(digest, targetPath string) *exec.Cmd {
	return exec.Command("grype", "-q", "--add-cpes-if-none", "-s", "AllLayers", "-o", "json", "--file", targetPath, digest)
}
