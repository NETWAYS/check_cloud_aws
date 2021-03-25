# check_cloud_aws

Icinga check plugin to check Amazon AWS resources. At the moment the check only supports EC2 instances.

## Usage

### EC2 - Instances

When one of the states is non-ok, or a machine is stopped, the check will alert.

#### ec2 instances

Check multiple EC2 instances from a region and filtered by name or auto-scaling group.

```
Usage:
  check_cloud_aws ec2 instances [flags]

Flags:
  -a, --autoscale string   Search for ec2 instances by autoscaling group
  -h, --help               help for instances
  -n, --name string        Search for ec2 instances by name (e.g. instance*)

Global Flags:
  -P, --profile string   AWS credential profile (~/.aws/credentials)
  -R, --region string    AWS region name (e.g. eu-central-1)
  -t, --timeout int      Timeout for the check (default 30)

```

```
$ check_cloud_aws ec2 instances
CRITICAL - 4 Instances found - 2 running - 2 stopped

[CRITICAL] i-0f38c870feae9c593 "(none)" stopped (no status)
[OK] i-06378c1b63a42384f "IcingaMaster" running instance=ok system=ok
[CRITICAL] i-0d5c0ad646be4610e "IcingaAgent" stopped (no status)
[OK] i-0a1b3143ae11565fc "IcingaSatellit" running instance=ok system=ok
```

#### ec2 instance

Check the states of a single EC2 instance

```
Usage:
  check_cloud_aws ec2 instance [flags]

Flags:
  -h, --help          help for instance
  -i, --id string     Look for ec2 instance by id
  -n, --name string   Look for ec2 instance by name

Global Flags:
  -P, --profile string   AWS credential profile (~/.aws/credentials)
  -R, --region string    AWS region name (e.g. eu-central-1)
  -t, --timeout int      Timeout for the check (default 30)
```

```
$ check_cloud_aws ec2 instance --name IcingaMaster
OK - "IcingaMaster" running instance=ok system=ok

ID: i-06378c1b63a42384f
Type: t2.micro
AutoScaling: (none)
```

## Authentication

Create a new IAM account for API access, add the required **read-only permissions, and check can monitor resources on
that account.

Permissions:
* `arn:aws:iam::aws:policy/AmazonEC2ReadOnlyAccess`

You can use the common AWS environment variables for the check, but we recommend storing the credentials in
the standard AWS locations on disk.

## Storing Credentials

By default, the check uses the config and credential storage of the AWS CLI.

**~/.aws/credentials**

The INI section header represents the credential profile you can pass to check_cloud_aws.

```
[default]
aws_access_key_id=eXaMPleID
aws_secret_access_key=SeCrEtExAmPleKey
```

**~/.aws/config**

The region setting here influences the default region of check_cloud_aws, which you can override.

```
[default]
region = eu-central-1
```


## Further Documentation

* [AWS SDK for Go](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html)
* [AWS Docs](https://docs.aws.amazon.com/ec2/index.html)

## License

Copyright (c) 2021 [NETWAYS GmbH](mailto:info@netways.de)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see [gnu.org/licenses](https://www.gnu.org/licenses/).
