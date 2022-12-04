package alfred

import (
	"fmt"

	aw "github.com/deanishe/awgo"
)

const ongoingTargets = "targets.json"

// type Targets struct {
// 	Name string `json:"name"`
// 	Path string `json:"path"`
// }

type Targets map[string]string

func LoadOngoingTargets(wf *aw.Workflow) (Targets, error) {
	// fallback load function doing nothing
	nop := func() (interface{}, error) {
		return 0.0, nil
	}

	var ta Targets
	if err := wf.Data.LoadOrStoreJSON(ongoingTargets, 0, nop, &ta); err != nil {
		return Targets{}, fmt.Errorf("error loading the ongoing targets: %w", err)
	}

	return ta, nil
}

func StoreOngoingTargets(wf *aw.Workflow, prop Targets) error {
	if err := wf.Data.StoreJSON(ongoingTargets, prop); err != nil {
		return fmt.Errorf("error storing the ongoing targets: %w", err)
	}

	return nil
}
