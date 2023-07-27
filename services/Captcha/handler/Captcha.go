package handler

import (
	"bytes"
	"context"
	"image/png"

	pb "github.com/Wuhao-9/IHome/services/Captcha/proto"
	"github.com/afocus/captcha"
)

type CaptchaSrvImpl struct{}

func (srv *CaptchaSrvImpl) GetCaptcha(ctx context.Context, req *pb.EmptyRequest, resp *pb.Img) (err error) {
	// generate a Captcha
	cap := captcha.New()
	err = cap.SetFont("comic.ttf")
	if err != nil {
		return
	}
	cap.SetDisturbance(captcha.NORMAL)
	img, _ := cap.Create(4, captcha.NUM)

	// encode image for transmit
	buf := &bytes.Buffer{}
	if err = png.Encode(buf, img); err != nil {
		return
	}
	resp.CaptchaImage = buf.Bytes()
	return nil
}
