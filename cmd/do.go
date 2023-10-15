package cmd

import "encoding/json"

type Arg struct {
	Op    string   `json:"op"`
	Files []string `json:"files"`
	Base  string   `json:"base"`
}

func ArgJSONBuilder(op string, files []string, base string) string {
	j, _ := json.Marshal(Arg{op, files, base})
	return string(j)
}
