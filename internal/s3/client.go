package s3

import (
	"fmt"
	"github.com/NETWAYS/go-check"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Client struct {
	S3Client *s3.S3
}

func NewS3Client(session *session.Session) *S3Client {
	return &S3Client{s3.New(session)}
}

func (c *S3Client) LoadAllBuckets() (buckets *s3.ListBucketsOutput, err error) {
	buckets, err = c.S3Client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		err = fmt.Errorf("could not load all buckets: %w", err)
		return nil, err
	}

	return
}

func (c *S3Client) LoadBucketByName(name string) (bucket *s3.Bucket) {
	buckets, err := c.LoadAllBuckets()
	if err != nil {
		check.ExitError(err)
	}

	for _, bucket = range buckets.Buckets {
		if *bucket.Name == name {
			return bucket
		}
	}

	return
}

func (c *S3Client) LoadAllObjectsFromBucket(bucket, prefix string) (objects *s3.ListObjectsV2Output, err error) {
	if prefix != "" {
		objects, err = c.S3Client.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
			Prefix: aws.String(prefix),
		})
	} else {
		objects, err = c.S3Client.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
		})
	}

	if err != nil {
		err = fmt.Errorf("could not load all bucket objects: %w", err)
		return nil, err
	}

	return
}
