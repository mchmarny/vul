package scanner

import (
	"bytes"
	"os"
	"os/exec"
	"path"
	"sync"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	numberOfScanners = 3

	TrivyReportName = "trivy.json"
	SnykReportName  = "snyk.json"
	GrypeReportName = "grype.json"
)

// Scan runs vulnerability scan on the provided image.
func Scan(imageURI, targetDirPath string) ([]string, error) {
	log.Info().Msgf("scanning image %s to %s", imageURI, targetDirPath)

	var wg sync.WaitGroup

	f1 := path.Join(targetDirPath, TrivyReportName)
	wg.Add(1)
	go runCmd(&wg, makeTrivyCmd(imageURI, f1), f1)

	f2 := path.Join(targetDirPath, SnykReportName)
	wg.Add(1)
	go runCmd(&wg, makeSnykCmd(imageURI, f2), f2)

	f3 := path.Join(targetDirPath, GrypeReportName)
	wg.Add(1)
	go runCmd(&wg, makeGrypeCmd(imageURI, f3), f3)

	wg.Wait()

	files, err := os.ReadDir(targetDirPath)
	if err != nil {
		return nil, errors.Wrap(err, "error reading scan dir")
	}

	if len(files) != numberOfScanners {
		return nil, errors.Errorf("expected %d files, got %d, see logs for details", numberOfScanners, len(files))
	}

	list := make([]string, 0, numberOfScanners)
	for _, f := range files {
		list = append(list, path.Join(targetDirPath, f.Name()))
	}

	return list, nil
}

func runCmd(wg *sync.WaitGroup, cmd *exec.Cmd, path string) {
	defer wg.Done()
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()

	if _, e := os.Stat(path); errors.Is(e, os.ErrNotExist) {
		// only err if the file doesn't exist
		// some scanners (snyk) will return 1 when they find vulnerabilities
		log.Error().
			Err(err).
			Str("cmd", cmd.String()).
			Str("out", outb.String()).
			Str("err", errb.String()).
			Msgf("error executing scanner command")
	}
}
