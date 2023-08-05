package alfred

import (
	aw "github.com/deanishe/awgo"
)

const (
	Debug           = "DEBUG"
	Exclude         = "EXCLUDE"
	Type            = "TYPE"
	MaxQueryResults = "MAX_QUERY_RESULTS"
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

func GetMaxQueryResults(wf *aw.Workflow) int {
	return wf.Config.GetInt(MaxQueryResults)
}
