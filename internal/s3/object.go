package s3

import (
	"fmt"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/convert"
)

func (v *V2Output) GetObjectOutput(size int64, status int, path string) (output string) {
	var tmp_size uint64

	if size >= 0 {
		tmp_size = uint64(size)
	}

	output += fmt.Sprintf("   \\_[%s] %s: %s\n",
		check.StatusText(status),
		path,
		convert.BytesIEC(tmp_size))

	return
}
