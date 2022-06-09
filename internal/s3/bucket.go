package s3

import (
	"fmt"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/convert"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Bucket struct {
	Bucket *s3.Bucket
}

type V2Output struct {
	V2Output *s3.ListObjectsV2Output
}

func (v *V2Output) GetBucketOutput(size int64, status int) (output string) {
	output = fmt.Sprintf("[%s] %s: %s",
		check.StatusText(status),
		*v.V2Output.Name,
		convert.BytesIEC(size))

	return
}

func (v *V2Output) CalculateBucketSize() (size int64) {
	for _, content := range v.V2Output.Contents {
		size += *content.Size
	}

	return size
}
