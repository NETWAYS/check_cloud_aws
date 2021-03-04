package ec2_test

import (
	"fmt"
	"github.com/NETWAYS/check_cloud_aws/internal/ec2"
)

func ExampleDescribeInstancesInput() {
	input := ec2.DescribeInstancesInput(
		ec2.Filter("tag:aws:autoscaling:groupName", "magic"),
		ec2.Filter("tag:Name", "web*"),
	)
	fmt.Println(input)
	// Output:
	// {
	//   Filters: [{
	//       Name: "tag:aws:autoscaling:groupName",
	//       Values: ["magic"]
	//     },{
	//       Name: "tag:Name",
	//       Values: ["web*"]
	//     }]
	// }
}

func ExampleFilter() {
	filter := ec2.Filter("tag:Name", "vm1")
	fmt.Println(filter)
	// Output:
	// {
	//   Name: "tag:Name",
	//   Values: ["vm1"]
	// }
}
