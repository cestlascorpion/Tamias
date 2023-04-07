package service

import (
	"context"

	"github.com/cestlascorpion/Tamias/proto"
)

type Server struct {
	proto.UnimplementedTamiasServer
}

func (s *Server) Upload(ctx context.Context, in *proto.UploadReq) (*proto.UploadResp, error) {
	// TODO:
	return nil, nil
}
