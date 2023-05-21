package scanner

import (
	"bytes"
	"os"
	"strings"
	"sync"

	"github.com/mchmarny/vul/internal/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Scan runs vulnerability scan on the provided image.
func Scan(cnf config.Scanner, imageURI, targetDirPath string) ([]string, error) {
	log.Info().Msgf("scanning image %s to %s", imageURI, targetDirPath)

	var wg sync.WaitGroup

	commands := makeScannerCommands(imageURI, targetDirPath)
	list := make([]string, 0, len(commands))

	for _, c := range commands {
		wg.Add(1)
		list = append(list, c.path)
		go runCmd(&wg, c, cnf.EnvVars)
	}

	wg.Wait()

	files, err := os.ReadDir(targetDirPath)
	if err != nil {
		return nil, errors.Wrap(err, "error reading scan dir")
	}

	if len(files) != len(commands) {
		return nil, errors.Errorf("expected %d files, got %d, see logs for details", len(commands), len(files))
	}

	return list, nil
}

func runCmd(wg *sync.WaitGroup, c *scannerCmd, envVars []string) {
	defer wg.Done()

	log.Debug().
		Str("name", c.name).
		Str("report", c.path).
		Str("cmd", c.cmd.String()).
		Str("env", strings.Join(envVars, ",")).
		Msg("running scanner")

	var outb, errb bytes.Buffer
	c.cmd.Stdout = &outb
	c.cmd.Stderr = &errb
	if len(envVars) > 0 {
		c.cmd.Env = append(c.cmd.Env, envVars...)
	}

	err := c.cmd.Run()
	if _, e := os.Stat(c.path); errors.Is(e, os.ErrNotExist) {
		// only err if the file doesn't exist
		// some scanners (snyk) will return 1 when they find vulnerabilities
		log.Error().
			Err(err).
			Str("name", c.name).
			Str("report", c.path).
			Str("cmd", c.cmd.String()).
			Str("out", outb.String()).
			Str("err", errb.String()).
			Msgf("error executing %s scanner command", c.name)
	}
}
