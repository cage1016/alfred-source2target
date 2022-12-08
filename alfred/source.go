package alfred

import (
	"fmt"

	aw "github.com/deanishe/awgo"
)

const ongoingSources = "sources.json"

type Sources map[string]string

func LoadOngoingSources(wf *aw.Workflow) (Sources, error) {
	// fallback load function doing nothing
	nop := func() (interface{}, error) {
		return 0.0, nil
	}

	var s Sources
	if err := wf.Data.LoadOrStoreJSON(ongoingSources, 0, nop, &s); err != nil {
		return Sources{}, fmt.Errorf("error loading the ongoing sources: %w", err)
	}

	return s, nil
}

func StoreOngoingSources(wf *aw.Workflow, prop Sources) error {
	if err := wf.Data.StoreJSON(ongoingSources, prop); err != nil {
		return fmt.Errorf("error storing the ongoing sources: %w", err)
	}

	return nil
}
