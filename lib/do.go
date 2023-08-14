package lib

import (
	"bytes"
	"os/exec"
	"regexp"
	"strconv"
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
	// prepare fd command flags
	flags := lo.Compact[string](lo.ReduceRight([][]string{
		strings.Split(cfg.Arg, " "),
		strings.Split(cfg.Type, "\n"),
		lo.FlatMap(lo.Compact[string](strings.Split(cfg.Exclude, "\n")), func(arg string, index int) []string {
			return []string{"-E", arg}
		}),
	}, func(agg []string, item []string, _ int) []string {
		return append(agg, item...)
	}, []string{}))

	flags = append(flags, "-X", "ls", "-lt")                                    // exec batch
	flags = append(flags, "|", "head", "-n", strconv.Itoa(cfg.MaxQueryResults)) // limit results

	// prepare fd command
	cmd := exec.Command("sh", "-c", "fd "+strings.Join(flags, " "))
	cmd.Dir = cfg.Source

	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()
	if err != nil {
		logrus.Fatal(err)
	}

	documents := []string{}
	for _, line := range strings.Split(out.String(), "\n") {
		if line == "" {
			continue
		}

		logrus.Debugf("line: %s", line)
		if strings.HasPrefix(line, "-") {
			matches := regex.FindAllString(line, -1)
			logrus.Debugf("got: %s", matches[0])
			documents = append(documents, matches[0])
		} else {
			documents = append(documents, line)
		}
	}

	for _, line := range strings.Split(errOut.String(), "\n") {
		if line == "" {
			continue
		}

		documents = append(documents, line)
	}

	return lo.Uniq(lo.WithoutEmpty(documents))
}
