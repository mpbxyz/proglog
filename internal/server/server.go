package server

import (
    "context"
    
    api "github.com/mpbxyz/proglog/api/v1"
    "google.golang.org/grpc"
)

type Config struct {
    CommitLog CommitLog
}

var _ api.LogServer = (*grpcServer)(nil)

type grpcServer struct {
    api.UnimplementedLogServer
    *Config
}

func newgrpcServer(config *Config) (svr *grpcServer, err error) {
    svr = &grpcServer{
        Config: Config,
    }
    return svr,nil
    
}

func (s *grpcServer) Produce(ctx context.Context, req *api.ProduceRequest) (*api.ProduceResponse, error) {
    offset, err := s.CommitLog.Append(req.Record) 
    if err != nil {
        return nil, err
    }
    return &api.ProduceResponse{Offset: offset}, nil
}

func (s *grpcServer) Consume(ctx context.Context, req *api.ConsumeRequest) (*api.ConsumeResponse, error) {
    record, err := s.CommitLog.Read(req.Offset)
    if err != nil {
        return nil, err
    }
    return &api.ConsumeResponse{Record: record}, nil
}

func (s *grpcServer) ProduceStream()  {
    
}
