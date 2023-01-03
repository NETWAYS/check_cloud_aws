package ec2

import (
	"fmt"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
	"github.com/aws/aws-sdk-go/service/ec2"
	"net/url"
)

type Instance struct {
	Instance *ec2.Instance
	Status   *ec2.InstanceStatus
}

func (i *Instance) GetOutput() (out string) {
	instance := i.Instance

	name := "(none)"

	for _, keys := range instance.Tags {
		if *keys.Key == "Name" {
			name = url.QueryEscape(*keys.Value)
			break
		}
	}

	out = fmt.Sprintf(`"%s" %s`, name, *instance.State.Name)

	if i.Status != nil {
		if i.Status.InstanceStatus != nil {
			out += " instance=" + *i.Status.InstanceStatus.Status
		}

		if i.Status.SystemStatus != nil {
			out += " system=" + *i.Status.SystemStatus.Status
		}
	} else {
		out += " (no status)"
	}

	return
}

func (i *Instance) GetLongOutput() (out string) {
	autoscaling := "(none)"

	for _, keys := range i.Instance.Tags {
		if *keys.Key == "aws:autoscaling:groupName" {
			autoscaling = url.QueryEscape(*keys.Value)
			break
		}
	}

	out += " \\_ID: " + *i.Instance.InstanceId + "\n"
	out += " \\_Type: " + *i.Instance.InstanceType + "\n"
	out += " \\_AutoScaling: " + autoscaling + "\n"

	return
}

// * instance-state-name - The state of the instance (pending | running |
// shutting-down | terminated | stopping | stopped).
//
// * instance-status.status - The status of the instance (ok | impaired |
// initializing | insufficient-data | not-applicable).
//
// * system-status.status - The system status of the instance (ok | impaired
// | initializing | insufficient-data | not-applicable).
func (i *Instance) GetStatus() int {
	states := []int{3, 3, 3}

	switch *i.Instance.State.Name {
	case "running":
		states[0] = check.OK
	case "pending", "shutting-down", "terminated", "stopping", "stopped":
		states[0] = check.Critical
	default:
		states[0] = check.Critical
	}

	if i.Status != nil {
		switch *i.Status.InstanceStatus.Status {
		case "ok":
			states[1] = check.OK
		case "impaired", "initializing", "insufficient-data", "not-applicable":
			states[1] = check.Critical
		default:
			states[1] = check.Critical
		}

		switch *i.Status.SystemStatus.Status {
		case "ok":
			states[2] = check.OK
		case "impaired", "initializing", "insufficient-data", "not-applicable":
			states[2] = check.Critical
		default:
			states[2] = check.Critical
		}
	} else {
		states[1] = check.Critical
		states[2] = check.Critical
	}

	return result.WorstState(states...)
}
