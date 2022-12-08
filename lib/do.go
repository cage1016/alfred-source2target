package lib

import (
	"bufio"
	"io"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

type DoConfig struct {
	Source    string
	Arg       string
	Type      string
	Exclude   string
	ExecBatch string
}

func FdExecute(cfg DoConfig) []string {
	flags := []string{}
	for _, arg := range strings.Split(cfg.Arg, " ") {
		if arg != "" {
			flags = append(flags, arg)
		}
	}

	// type
	for _, ts := range strings.Split(cfg.Type, "\n") {
		if ts != "" {
			flags = append(flags, ts)
		}
	}

	// exclude
	for _, es := range strings.Split(cfg.Exclude, "\n") {
		if es != "" {
			flags = append(flags, "-E", es)
		}
	}

	// exec batch
	for _, bs := range strings.Split(cfg.ExecBatch, " ") {
		if bs != "" {
			flags = append(flags, bs)
		}
	}

	cmd := exec.Command("fd", flags...)
	cmd.Dir = cfg.Source
	logrus.Debugf("fd: %s", cmd)

	r, _ := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	documents := []string{}
	done := make(chan struct{})
	worker(done, r, func(line string) {
		logrus.Debugf("line: %s", line)
		if strings.HasPrefix(line, "-") {
			x := strings.Split(line, "./")
			documents = append(documents, x[1])
		} else {
			documents = append(documents, line)
		}
	})

	cmd.Start()
	<-done
	cmd.Wait()

	return documents
}

func worker(done chan struct{}, r io.ReadCloser, fn func(string)) {
	scanner := bufio.NewScanner(r)
	go func() {
		for scanner.Scan() {
			fn(scanner.Text())
		}
		done <- struct{}{}
	}()
}
