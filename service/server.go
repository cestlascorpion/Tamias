package service

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"

	"github.com/cestlascorpion/Tamias/core"
	"github.com/cestlascorpion/Tamias/proto"
	"github.com/cestlascorpion/Tamias/storage"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	proto.UnimplementedTamiasServer
	cache storage.Cache
}

func NewServer(ctx context.Context, config *core.Config) (*Server, error) {
	cache, err := storage.NewRedis(ctx, config)
	if err != nil {
		log.Errorf("new redis err %+v", err)
		return nil, err
	}

	return &Server{
		cache: cache,
	}, nil
}

func (s *Server) Upload(ctx context.Context, in *proto.UploadReq) (*proto.UploadResp, error) {
	out := &proto.UploadResp{}
	if len(in.FileUrl) == 0 {
		log.Errorf("file url is nil")
		return out, errInvalidParameters
	}

	file, err := core.Download(in.FileUrl)
	if err != nil {
		log.Errorf("download %s err %+v", in.FileUrl, err)
		return out, err
	}

	switch in.Manufacturer {
	case proto.Manufacturer_XM:
		uri, ttl, err := s.xmUpload(ctx, file, in.FileType)
		if err != nil {
			log.Errorf("xm upload err %+v", err)
			return out, err
		}
		out.Uri = uri
		out.Ttl = ttl
		return out, nil
	}

	log.Errorf("unknown manufacturer %s", in.Manufacturer.String())
	return out, errInvalidParameters
}

func (s *Server) Close(ctx context.Context) error {
	return s.cache.Close(ctx)
}

func (s *Server) xmUpload(ctx context.Context, file []byte, fileType proto.FileType) (string, int64, error) {
	var (
		name string
		err  error
	)

	switch fileType {
	case proto.FileType_LARGE_ICON:
		name, err = checkXMLargeIcon(ctx, file)
		if err != nil {
			return "", 0, err
		}
	case proto.FileType_BIG_PICTURE:
		name, err = checkXMBigPicture(ctx, file)
		if err != nil {
			return "", 0, err
		}
	default:
		log.Errorf("unknown xm file type %s", fileType.String())
		return "", 0, errInvalidParameters
	}

	key := fmt.Sprintf("xm-%d-%s", fileType, name)
	log.Debugf("xm file key %s", key)

	uri, ttl, err := s.cache.GetUri(ctx, key)
	if err == nil {
		return uri, ttl, nil
	}

	log.Warnf("get %s from cache failed", key)
	uri, ttl, err = core.XMUpload(ctx, file, name, fileType)
	if err != nil {
		log.Errorf("xm upload err %+v", err)
		return "", 0, err
	}
	_ = s.cache.SetUri(ctx, key, uri, ttl)
	return uri, ttl, nil
}

func checkXMLargeIcon(ctx context.Context, file []byte) (string, error) {
	// < 200KB
	if len(file) >= 200*1024 {
		log.Errorf("xm large icon over size %d", len(file))
		return "", errInvalidParameters
	}

	x, y, name, err := core.Parse(file)
	if err != nil {
		log.Errorf("xm large icon decode image err %+v", err)
		return "", errInvalidParameters
	}

	// PNG/JPEG/JPG
	if name != "png" && name != "jpeg" && name != "jpg" {
		log.Errorf("unknown format %s of xm large icon", name)
		return "", errInvalidParameters
	}

	// 120 X 120px
	if x != 120 || y != 120 {
		log.Errorf("invalid xm large icon size %d %d", x, y)
		return "", errInvalidParameters
	}

	return fmt.Sprintf("%x.%s", md5.Sum(file), name), nil
}

func checkXMBigPicture(ctx context.Context, file []byte) (string, error) {
	// < 1MB
	if len(file) >= 1*1024*1024 {
		log.Errorf("xm big picture over size %d", len(file))
		return "", errInvalidParameters
	}

	x, y, name, err := core.Parse(file)
	if err != nil {
		log.Errorf("decode xm big picture err %+v", err)
		return "", errInvalidParameters
	}

	// PNG/JPEG/JPG
	if name != "png" && name != "jpeg" && name != "jpg" {
		log.Errorf("unknown format %s of xm big picture", name)
		return "", errInvalidParameters
	}

	// 876 X 324px
	if x != 876 || y != 324 {
		log.Errorf("invalid xm big picture size %d %d", x, y)
		return "", errInvalidParameters
	}

	return fmt.Sprintf("%x.%s", md5.Sum(file), name), nil
}

var (
	errInvalidParameters = errors.New("invalid parameters")
)
