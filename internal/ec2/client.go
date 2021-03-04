package ec2

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Client implementation to offer various load functions for getting data from the API
type Client struct {
	Client *ec2.EC2
}

// NewClient sets up a Client with a AWS session.Session
func NewClient(session *session.Session) *Client {
	return &Client{ec2.New(session)}
}

// LoadInstance returns a single Instance looking for its id
func (c *Client) LoadInstance(id string) (instance *Instance, err error) {
	return c.LoadInstanceByFilter(DescribeInstancesInput(Filter("instance-id", id)))
}

// LoadInstanceByName returns a single Instance looking for a name
//
// Name is not required to be unique, but our interface expects it is.
func (c *Client) LoadInstanceByName(name string) (instance *Instance, err error) {
	return c.LoadInstanceByFilter(DescribeInstancesInput(Filter("tag:Name", name)))
}

// LoadInstanceByFilter returns a single Instance using a ec2.DescribeInstancesInput with a ec2.Filter in it
//
// The function expects the result to have exactly one match.
//
// Also see our Filter and DescribeInstancesInput
func (c *Client) LoadInstanceByFilter(filter *ec2.DescribeInstancesInput) (instance *Instance, err error) {
	instances, err := c.Client.DescribeInstances(filter)

	if err != nil {
		err = fmt.Errorf("could not load instances: %w", err)
		return nil, err
	}

	// Check count of VMs returned
	l := 0

	for _, res := range instances.Reservations {
		for range res.Instances {
			l++
		}
	}

	if l == 0 {
		return nil, fmt.Errorf("no instance found matching filter: %s", filter)
	} else if l > 1 {
		return nil, fmt.Errorf("more than one instance found matching filter: %s", filter)
	}

	instance = &Instance{}
	instance.Instance = instances.Reservations[0].Instances[0]

	// Load status for instance
	instance.Status, err = c.LoadInstanceStatus(*instance.Instance.InstanceId)
	if err != nil {
		return
	}

	return
}

// LoadInstanceStatus returns the ec2.Instance for an id
func (c *Client) LoadInstanceStatus(id string) (status *ec2.InstanceStatus, err error) {
	d, err := c.Client.DescribeInstanceStatus(&ec2.DescribeInstanceStatusInput{
		InstanceIds: []*string{&id},
	})
	if err != nil {
		err = fmt.Errorf("could not load instance status for '%s': %w", id, err)
		return
	}

	// Check result length
	l := len(d.InstanceStatuses)

	// We would have a nil status when the VM is stopped
	if l == 1 {
		status = d.InstanceStatuses[0]
	} else if l > 1 {
		return nil, fmt.Errorf("more than one instance status found for %s", id)
	}

	return
}

// LoadAllInstancesByFilter returns Instances with a list of Instance to work with
//
// Also see our Filter and DescribeInstancesInput
func (c *Client) LoadAllInstancesByFilter(filter *ec2.DescribeInstancesInput) (instances *Instances, err error) {
	d, err := c.Client.DescribeInstances(filter)
	if err != nil {
		err = fmt.Errorf("could not load instances: %w", err)
		return nil, err
	}

	instances = &Instances{}

	for index := range d.Reservations {
		for _, i := range d.Reservations[index].Instances {
			instance := &Instance{Instance: i}

			instance.Status, err = c.LoadInstanceStatus(*i.InstanceId)
			if err != nil {
				return
			}

			instances.Instances = append(instances.Instances, instance)
		}
	}

	return
}
