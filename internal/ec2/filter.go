package ec2

import (
	"github.com/aws/aws-sdk-go/service/ec2"
)

// DescribeInstancesInput builds an input with one or multiple filters
//
// Useful in conjunction with Filter
func DescribeInstancesInput(filter ...*ec2.Filter) *ec2.DescribeInstancesInput {
	return &ec2.DescribeInstancesInput{Filters: filter}
}

// Filter builds a simple key value filter for EC2
func Filter(name, value string) *ec2.Filter {
	return &ec2.Filter{
		Name:   &name,
		Values: []*string{&value},
	}
}
