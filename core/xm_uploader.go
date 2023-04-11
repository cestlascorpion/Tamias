package core

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cestlascorpion/Tamias/proto"
	log "github.com/sirupsen/logrus"
)

func XMUpload(ctx context.Context, file []byte, name string, fileType proto.FileType) (string, int64, error) {
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("key=%s", xmAppSecret)

	fields := make(map[string]string)
	fields["is_global"] = "false"
	if fileType == proto.FileType_LARGE_ICON {
		fields["is_icon"] = "true"
	} else {
		fields["is_icon"] = "false"
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(time.Second * 10)
	}
	timeout := deadline.Sub(time.Now())

	body, err := Upload(xmUploadURL, headers, fields, name, file, timeout)
	if err != nil {
		log.Errorf("upload file err %+v", err)
		return "", 0, err
	}

	resp := &xmResp{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		log.Errorf("json unmarshal err %+v", err)
		return "", 0, err
	}

	if resp.Code != 0 {
		log.Errorf("code %d result %s desc %s", resp.Code, resp.Result, resp.Desc)
		return "", 0, fmt.Errorf("code %d result %s desc %s", resp.Code, resp.Result, resp.Desc)
	}

	if fileType == proto.FileType_LARGE_ICON {
		return resp.Data.IconUrl, int64(xmTtl.Seconds() - delta.Seconds()), nil
	} else {
		return resp.Data.PicUrl, int64(xmTtl.Seconds() - delta.Seconds()), nil
	}
}

const (
	xmUploadURL = "https://api.xmpush.xiaomi.com/media/upload/image"
)

const (
	xmTtl = time.Hour * 24 * 30 * 3 // 3个月
	delta = time.Hour
)

const (
	xmAppSecret = "" // TODO: add it
)

type xmResp struct {
	Result  string `json:"result"`
	TraceId string `json:"trace_id"`
	Code    int    `json:"code"`
	Data    struct {
		IconUrl string `json:"icon_url,omitempty"`
		PicUrl  string `json:"pic_url,omitempty"`
		Sha1    string `json:"sha1,omitempty"`
	} `json:"data,omitempty"`
	Desc string `json:"description"`
}
