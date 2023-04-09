package log

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
    write = []byte("hello world")
    width = uint64(len(write)) + lenWidth
)

func TestStoreAppendRead(t *testing.T){
    f, err := ioutil.TempFile("", "store_append_read_test")
    require.NoError(t, err)

    defer os.Remove(f.Name())
}
