package lib

import (
	"bufio"
	"io"
	"os/exec"
	"regexp"
	"strings"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

var regex = regexp.MustCompile(`\./.+`)

type DoConfig struct {
	Source          string
	Arg             string
	Type            string
	Exclude         string
	MaxQueryResults int
}

func FdExecute(cfg DoConfig) []string {
	// prepare flags
	flags := lo.Compact[string](lo.ReduceRight([][]string{
		strings.Split(cfg.Arg, " "),
		strings.Split(cfg.Type, "\n"),
		lo.FlatMap(lo.Compact[string](strings.Split(cfg.Exclude, "\n")), func(arg string, index int) []string {
			return []string{"-E", arg}
		}),
	}, func(agg []string, item []string, _ int) []string {
		return append(agg, item...)
	}, []string{}))
	flags = append(flags, "-X", "ls", "-lt") // exec batch

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
			matches := regex.FindAllString(line, -1)
			logrus.Debugf("got: %s", matches[0])
			documents = append(documents, matches[0])
		} else {
			documents = append(documents, line)
		}
	})

	cmd.Start()
	<-done
	cmd.Wait()

	// know issue
	// how to break scanner.Scan loop with lock by cmd.Wait()?
	if len(documents) > cfg.MaxQueryResults {
		return documents[:cfg.MaxQueryResults]
	} else {
		return documents
	}
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
