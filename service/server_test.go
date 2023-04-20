package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/cestlascorpion/Tamias/core"
	"github.com/cestlascorpion/Tamias/proto"
	"github.com/jinzhu/configor"
)

const (
	testXMLargeIcon  = "https://i.328888.xyz/2023/04/20/iGiqTC.jpeg"
	testXMBigPicture = "https://i.328888.xyz/2023/04/20/iGVZWP.png"
)

var (
	testServer *Server
)

func init() {
	conf := &core.Config{}
	err := configor.Load(conf, "/tmp/config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	svr, err := NewServer(context.Background(), conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	testServer = svr
}

func TestServer_Upload_XM_LargeIcon(t *testing.T) {
	if testServer == nil {
		return
	}

	resp, err := testServer.Upload(context.Background(), &proto.UploadReq{
		Manufacturer: proto.Manufacturer_XM,
		FileType:     proto.FileType_LARGE_ICON,
		FileUrl:      testXMLargeIcon,
	})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(resp)
}

func TestServer_Upload_XM_BigPicture(t *testing.T) {
	if testServer == nil {
		return
	}

	resp, err := testServer.Upload(context.Background(), &proto.UploadReq{
		Manufacturer: proto.Manufacturer_XM,
		FileType:     proto.FileType_BIG_PICTURE,
		FileUrl:      testXMBigPicture,
	})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(resp)
}
