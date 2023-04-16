package log

import (
	"io/ioutil"
	"os"
	"testing"

	api "github.com/mpbxyz/proglog/api/v1"
	"github.com/stretchr/testify/require"
	//"google.golang.org/protobuf/proto"
)

func TestLog(t *testing.T)  {
    for scenario, fn := range map[string]func(
        t *testing.T, log *Log,
    ) {
        "append and read a record succeeds": testAppendRead,
        "validates out of range error happens": testOutOfRange,
        "Init with existing segments": testInitExisting,
    }{
        t.Run(scenario, func(t *testing.T) {
            dir, err := ioutil.TempDir("", "store-test")
            require.NoError(t, err)
            defer os.RemoveAll(dir)

            c := Config{}
            c.Segment.MaxStoreBytes = 32
            log, err := newLog(dir, c)
            require.NoError(t, err)

            fn(t,log)

        })
    }
}

func testAppendRead(t *testing.T, log *Log){
    write := &api.Record {
       Value: []byte("Hello World"), 
    }

    off, err := log.Append(write)
    require.NoError(t, err)
    require.Equal(t, off, uint64(0))

    got, err := log.Read(off)
    require.NoError(t, err)
    require.Equal(t, write.Value, got.Value)
    
}

func testOutOfRange(t *testing.T, log *Log){
    got, err := log.Read(uint64(55))
    require.Nil(t, got)
    require.Error(t, err)
}

func testInitExisting(t *testing.T, log *Log)  {
}
