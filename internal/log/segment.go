package log

import (
	"fmt"
	"os"
	"path"
    api "github.com/mpbxyz/proglog/api/v1"
    "google.golang.org/protobuf/proto"
)

type segment struct{
    store *store
    index *index
    baseOffset, nextOffset uint64
    config Config
}

func newSegment(dir string, baseOffset uint64, c Config)(*segment, error){
    s := &segment{
        baseOffset: baseOffset,
        config: c,
    }

    var err error

    storeFile, err := os.OpenFile(
        path.Join(dir, fmt.Sprintf("%d%s", baseOffset, ".store")),
        os.O_RDWR|os.O_CREATE|os.O_APPEND,
        0644,
    )
    if err != nil{
        return nil, err
    }

    if s.store, err = newStore(storeFile); err != nil {
        return nil, err
    }


    indexFile, err := os.OpenFile(
        path.Join(dir, fmt.Sprintf("%d%s", baseOffset, ".index")),
        os.O_RDWR|os.O_CREATE|os.O_APPEND,
        0644,
    )
    if err != nil{
        return nil, err
    }
    
    if s.index, err = newIndex(indexFile,c); err != nil {
        return nil,err
    }

    if out, _, err := s.index.Read(-1);err != nil {
        s.nextOffset = baseOffset        
    }else {
        s.nextOffset = baseOffset + 1 + uint64(out)
    }

    return s, nil
}

func (s *segment) Append(record *api.Record) (offset uint64, err error){
    cur := s.nextOffset
    record.Offset = cur

    p, err := proto.Marshal(record)
    if err != nil {
        return 0, err
    }
    
    _, pos, err := s.store.Append(p)
    if err != nil{
        return 0, err
    }

    if err = s.index.Write(
        uint32(s.nextOffset-uint64(s.baseOffset)),
        pos,
    ); err != nil {
        return 0,err
    }

    s.nextOffset++
    return cur, nil
}