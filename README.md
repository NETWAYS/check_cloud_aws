# check_cloud_aws

An Icinga check plugin to check Amazon AWS resources.

## Usage

### Health Status

A general status based on the RSS feed on the AWS Health Page.

```
Usage:
  check_cloud_aws status [flags]

Flags:
  -u, --url string       The AWS Status Page URL (default "https://status.aws.amazon.com")
  -s, --service string   The AWS Service to check (default "ec2")
  -h, --help             help for status
```

```
check_cloud_aws status --service cloudfront
OK - Service cloudfront is operating normally

check_cloud_aws --region us-west-1 status --service cloudwatch
WARNING - Information available for cloudwatch in us-west-1

check_cloud_aws --region eu-west-1 status
CRITICAL - Service disruption for ec2 in eu-west-1

check_cloud_aws --region "" status iam
CRITICAL - WARNING - Information available for iam (Global)
```

### EC2

When one of the states is non-OK, or a machine is stopped, the check will alert.

#### EC2 Instances

Check multiple EC2 instances from a region and filtered by name or auto-scaling group.

```
Usage:
  check_cloud_aws ec2 instances [flags]

Flags:
  -a, --autoscale string   Search for ec2 instances by autoscaling group
  -h, --help               help for instances
  -n, --name string        Search for ec2 instances by name (e.g. instance*)
```

```
$ check_cloud_aws ec2 instances
CRITICAL - 4 Instances: 2 running - 2 stopped
 \_[CRITICAL] i-0f38c870feae9c593 "(none)" stopped (no status)
 \_[OK] i-06378c1b63a42384f "IcingaMaster" running instance=ok system=ok
 \_[CRITICAL] i-0d5c0ad646be4610e "IcingaAgent" stopped (no status)
 \_[OK] i-0a1b3143ae11565fc "IcingaSatellit" running instance=ok system=ok
```

#### EC2 Instance

Check the states of a single EC2 instance.

```
Usage:
  check_cloud_aws ec2 instance [flags]

Flags:
  -h, --help          help for instance
  -i, --id string     Search for ec2 instance by id
  -n, --name string   Search for ec2 instance by name
```

```
$ check_cloud_aws ec2 instance --name IcingaMaster
OK - "IcingaMaster" running instance=ok system=ok
 \_ID: i-06378c1b63a42384f
 \_Type: t2.micro
 \_AutoScaling: (none)
```

### S3

In the **bucket** context, calculates the size of the bucket and alerts if its size reaches the threshold.
In the **object** context, the size of each object inside a bucket will be checked against the threshold.

#### S3 Bucket

```
Usage:
  check_cloud_aws s3 bucket [flags]

Flags:
  -b, --buckets strings           Name of the S3 bucket. If '--buckets' is empty, all buckets will be evaluated.
  -c, --crit-bucket-size string   Critical threshold for the size of the specified bucket. Alerts if the size is greater than the critical threshold.
                                  Possible units are MB, GB or TB (defaults to MB if none is specified). (default "20Gb")
  -w, --warn-bucket-size string   Warning threshold for the size of the specified bucket. Alerts if the size is greater than the warning threshold.
                                  Possible units are MB, GB or TB (defaults to MB if none is specified). (default "10Gb")
      --empty-ok                  Return OK if there are no buckets
  -h, --help                      help for bucket
```

```
check_cloud_aws s3 bucket
OK - 1 Buckets: 0 Critical - 0 Warning - 1 OK
 \_[OK] my-bucket - value: 100MiB | my-bucket=100MB;10240;20480

check_cloud_aws s3 bucket --crit-bucket-size 10
CRITICAL - 1 Buckets: 1 Critical - 0 Warning - 0 OK
 \_[CRITICAL] my-bucket - value: 100MiB | my-bucket=100MB;10240;10

check_cloud_aws s3 bucket --crit-bucket-size 5GB
CRITICAL - 1 Buckets: 1 Critical - 0 Warning - 0 OK
 \_[CRITICAL] my-bucket - value: 100MiB | my-bucket=100MB;10240;10
```

### S3 Object

```
Usage:
  check_cloud_aws s3 object [flags]

Flags:
  -b, --buckets strings           Name of one or multiple S3 buckets. If '--buckets' is empty, all buckets will be evaluated.
      --prefix string             Limits the response to keys that begin with the specified prefix, e.G. '--prefix test' filters all objects which starts with 'test'.
                                  NOTE: Keep in mind, that objects beneath a directory will be ignored.
  -c, --crit-object-size string   Critical threshold for the size of the object. Alerts if the size is greater than the critical threshold.
                                  Possible units are MB, GB or TB (defaults to MB if none is specified). (default "1gb")
  -w, --warn-object-size string   Critical threshold for the size of the object. Alerts if the size is greater than the warning threshold.
                                  Possible units are MB, GB or TB (defaults to MB if none is specified). (default "800mb")
  -p, --perfdata                  Displays perfdata and lists all objects in the specified bucket.
      --empty-ok                  Return OK if there are no buckets
  -h, --help                      help for object
```

```
check_cloud_aws s3 object
OK - 2 Objects: 0 Critical - 0 Warning - 2 OK
 \_[my-bucket]:
     \_[OK] foo.fs: 100MiB
     \_[OK] bar.fs: 100MiB

check_cloud_aws s3 object --prefix file
OK - 3 Objects: 0 Critical - 0 Warning - 3 OK
 \_[my-bucket]:
     \_[OK] file_1.fs: 10MiB
     \_[OK] file_2.fs: 20MiB
     \_[OK] file_3.fs: 30MiB

check_cloud_aws s3 object --crit-object-size 10KB --buckets foo --buckets bar
CRITICAL - 2 Objects: 2 Critical - 0 Warning - 0 OK
 \_[foo]:
     \_[CRITICAL] file.fs: 100MiB
 \_[bar]:
     \_[CRITICAL] file.fs: 100MiB
```

### Cloudfront

Checks a specific or multiple cloudfront distributions from a region. When the state is `disabled` or `InProgress`,
the check will alert.

```
Usage:
  check_cloud_aws cloudfront [flags]

Flags:
  -e, --etag strings   Etag name of one or multiple distributions. If '--etag' is empty, all distributions will be evaluated.
  -h, --help           help for cloudfront

Global Flags:
  -C, --credentials-file string   Path to the credentials file (default "$HOME/.aws/credentials")
  -P, --profile string            The AWS profile name, which represents a separate credential profile in the credential file (default "default")
  -R, --region string             The AWS region to send requests to (default "eu-central-1")
  -t, --timeout int               Timeout for the check (default 30)
```

```
check_cloud_aws cloudfront
WARNING - 1 Distributions: 0 Disabled - 1 InProgress - 0 Enabled
 \_[WARNING] E2BD5GDFJZXKWC status=InProgress enabled=true
 | E2BD5GDFJZXKWC=inprogress
```

## Authentication

Create a new IAM account for API access, add the required **read-only permissions, and check can monitor resources on
that account.

Permission EC2:

* `arn:aws:iam::aws:policy/AmazonEC2ReadOnlyAccess`

Permission Cloudfront:

* `arn:aws:iam::aws:policy/CloudFrontReadOnlyAccess`

Permission S3:

* `arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess`

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
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see [gnu.org/licenses](https://www.gnu.org/licenses/).
