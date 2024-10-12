package cli

import (
	"fabricng/api"
)

type WorkflowBuilderCLI struct {
	Args []string
}

func (o *WorkflowBuilderCLI) Build() (ret *api.Workflow, err error) {
	ret = &api.Workflow{}
	return
}
