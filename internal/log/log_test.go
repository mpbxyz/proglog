package log

import (
	"io/ioutil"
	"os"
	"testing"

	api "github.com/mpbxyz/proglog/api/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestLog(t *testing.T)  {
    for scenario, fn := range map[string]func(
        t *testing.T, log *Log,
    ) {
        "append and read a record succeeds": testAppendRead,
        "validates out of range error happens": testOutOfRange,
        "Init with existing segments": testInitExisting,
        "Test reader": testReader,
        "Test that truncate equals right size": testTruncate,
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
    write := &api.Record {
       Value: []byte("Hello World"), 
    }

    for i := 0; i < 3; i++ {
        _, err := log.Append(write)
        require.NoError(t, err)
    }

    require.NoError(t, log.Close())
    
    off, err := log.LowestOffset()
    require.NoError(t, err)
    require.Equal(t, uint64(0), off)

    off, err = log.HighestOffset()
    require.NoError(t, err)
    require.Equal(t, uint64(2), off)
    
    n, err := newLog(log.Dir, log.Config)

    require.NoError(t, n.Close())
    
    off, err = n.LowestOffset()
    require.NoError(t, err)
    require.Equal(t, uint64(0), off)

    off, err = n.HighestOffset()
    require.NoError(t, err)
    require.Equal(t, uint64(2), off)
} 

func testReader(t *testing.T, log *Log)  {
    write := &api.Record {
       Value: []byte("Hello World"), 
    }
    off, err := log.Append(write)
    require.NoError(t, err)
    require.Equal(t, uint64(0), off)

    reader := log.Reader()

    b, err := ioutil.ReadAll(reader)
    require.NoError(t, err)
    read := &api.Record{}
    err = proto.Unmarshal(b[lenWidth:], read)

    require.NoError(t, err)
    require.Equal(t, write.Value, read.Value)
}

func testTruncate(t *testing.T, log *Log){
        
    write := &api.Record {
       Value: []byte("Hello World"), 
    }

    for i := 0; i < 3; i++ {
        _, err := log.Append(write)
        require.NoError(t, err)
    }

    err := log.Truncate(1)
    require.NoError(t, err)

    Base,err := log.LowestOffset()
    require.Equal(t, uint64(2), Base)
}
