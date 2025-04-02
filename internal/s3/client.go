package s3

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var ErrBucketNotFound = errors.New("no such Bucket")

type S3Client struct {
	S3Client *s3.S3
}

func NewS3Client(session *session.Session) *S3Client {
	return &S3Client{s3.New(session)}
}

func (c *S3Client) LoadAllBuckets() (buckets *s3.ListBucketsOutput, err error) {
	buckets, err = c.S3Client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("could not load all buckets: %w", err)
	}

	return buckets, nil
}

func (c *S3Client) LoadBucketByName(name string) (bucket *s3.Bucket, err error) {
	buckets, err := c.LoadAllBuckets()

	if err != nil {
		return nil, err
	}

	for _, bucket = range buckets.Buckets {
		if *bucket.Name == name {
			return bucket, nil
		}
	}

	return nil, ErrBucketNotFound
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
