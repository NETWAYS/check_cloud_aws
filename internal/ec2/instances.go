package ec2

import (
	"fmt"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
)

type Instances struct {
	Instances []*Instance
}

func (i Instances) GetStatus() int {
	var states []int

	for _, instance := range i.Instances {
		states = append(states, instance.GetStatus())
	}

	return result.WorstState(states...)
}

func (i Instances) GetOutput() (output string) {
	for _, instance := range i.Instances {
		output += fmt.Sprintf("[%s] %s %s\n",
			check.StatusText(instance.GetStatus()),
			*instance.Instance.InstanceId,
			instance.GetOutput())
	}

	return
}
