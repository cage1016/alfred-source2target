package alfred

import (
	aw "github.com/deanishe/awgo"
)

const (
	Debug     = "DEBUG"
	Exclude   = "EXCLUDE"
	Type      = "TYPE"
	ExecBatch = "EXEC_BATCH"
)

func GetDebug(wf *aw.Workflow) bool {
	return wf.Config.GetBool(Debug)
}

func GetExclude(wf *aw.Workflow) string {
	return wf.Config.Get(Exclude)
}

func GetType(wf *aw.Workflow) string {
	return wf.Config.Get(Type)
}

func GetExecBatch(wf *aw.Workflow) string {
	return wf.Config.Get(ExecBatch)
}
