package alfred

import (
	aw "github.com/deanishe/awgo"
)

const (
	Debug = "DEBUG"
)

func GetDebug(wf *aw.Workflow) bool {
	return wf.Config.GetBool(Debug)
}
