package core

import (
	"context"
	"fmt"
	"testing"
	"time"

	_ "image/jpeg"
	_ "image/png"

	"github.com/cestlascorpion/Tamias/proto"
)

var (
	testUrl          = "https://i.328888.xyz/2023/04/11/ipFoxE.png"
	testXMLargeIcon  = "https://i.328888.xyz/2023/04/11/ipFSlQ.jpeg"
	testXMBigPicture = "https://i.328888.xyz/2023/04/11/ipFoxE.png"
)

func TestDownloadAndParse(t *testing.T) {
	bs, err := Download(testUrl)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(len(bs))

	width, height, name, err := Parse(bs)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(width, height, name)
}

func TestXMUpload(t *testing.T) {
	bs, err := Download(testXMLargeIcon)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(len(bs))

	width, height, name, err := Parse(bs)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(width, height, name)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	uri, ttl, err := XMUpload(ctx, bs, fmt.Sprintf("icon.%s", name), proto.FileType_LARGE_ICON)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(uri, ttl)
}

func TestXMUpload2(t *testing.T) {
	bs, err := Download(testXMBigPicture)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(len(bs))

	width, height, name, err := Parse(bs)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(width, height, name)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	uri, ttl, err := XMUpload(ctx, bs, fmt.Sprintf("pic.%s", name), proto.FileType_BIG_PICTURE)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(uri, ttl)
}
