package handler

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"

	stdErrors "errors"
	"image/png"
	"log"
	"time"

	"github.com/Wuhao-9/IHome/services/authentication/model"
	dao "github.com/Wuhao-9/IHome/services/authentication/model/DAO"
	pb "github.com/Wuhao-9/IHome/services/authentication/proto"
	"github.com/afocus/captcha"
	"github.com/alibabacloud-go/tea/tea"
	"go-micro.dev/v4/errors"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
)

type AuthSrvImpl struct {
}

func (srv *AuthSrvImpl) GetImgCaptcha(ctx context.Context, requ *pb.ImgCaptchaRequ, resp *pb.ImgCaptchaResp) error {
	// generate a Captcha
	capt := captcha.New()
	err := capt.SetFont("/home/wuhao/myProj/goProj/IHome/services/authentication/conf/comic.ttf")
	if err != nil {
		return errors.InternalServerError("", "failed to set font for image Captcha, detail: %s", err)
	}

	capt.SetDisturbance(captcha.NORMAL)
	img, code := capt.Create(4, captcha.ALL)

	DbHandle := dao.NewCacheDaoRedis(dao.RedisConnPool.Get())
	defer DbHandle.Close()

	err = DbHandle.CreateWithExpire(model.ImageCaptchaPair{Uuid: requ.GetUuid(), Code: code}, (time.Minute * 5).Seconds())
	if err != nil {
		log.Println(err)
		return errors.InternalServerError("", "failed to create the Captcha record")
	}

	buf := &bytes.Buffer{}
	if err = png.Encode(buf, img); err != nil {
		return errors.InternalServerError("", "failed to encode the image, detail: %s", err)
	}
	resp.Image = buf.Bytes()
	return nil
}

// 调用SMS前提先要成功通过图片验证码
func (srv *AuthSrvImpl) GetSmsCaptcha(ctx context.Context, requ *pb.SmsCaptchaRequ, resp *pb.SmsCaptchaResp) error {
	imageCode := requ.GetImageCode()
	uuid := requ.GetUuid()

	DbHandle := dao.NewCacheDaoRedis(dao.RedisConnPool.Get())
	defer DbHandle.Close()

	pair := model.ImageCaptchaPair{Uuid: uuid}
	err := DbHandle.Query(&pair)
	if err != nil {
		if stdErrors.Is(err, dao.NotExistError) {
			return errors.BadRequest("", `invalied uuid`)
		} else {
			return errors.InternalServerError("", "failed to query, detail: %s", err)
		}
	}

	if !strings.EqualFold(imageCode, pair.Code) {
		return errors.BadRequest("", `invalied image-code`)
	}

	smsCaptcha, err := srv.sms(requ.PhoneNum)
	if err != nil {
		return errors.InternalServerError("", "failed to SMS, detail: %s", err)
	}

	DbHandle.CreateWithExpire(model.SmsCaptchaPair{PhoneNum: requ.PhoneNum, Code: smsCaptcha}, (time.Minute * 5).Seconds())
	if err != nil {
		return errors.InternalServerError("", "failed to create the SMS record")
	}

	return nil
}

// short message service
func (srv *AuthSrvImpl) sms(phoneNum string) (string, error) {
	accessKeyId := os.Getenv("ALIYUN_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ALIYUN_ACCESS_KEY_SECRET")
	signture := os.Getenv("ALIYUN_SMS_SIGNTURE")
	templateCode := os.Getenv("ALIYUN_SMS_TEMPLATE_CODE")
	config := &openapi.Config{
		AccessKeyId:     &accessKeyId,
		AccessKeySecret: &accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")

	// instantiate a client
	client, _ := dysmsapi.NewClient(config)

	// create sms request
	request := &dysmsapi.SendSmsRequest{}
	request.SetPhoneNumbers(phoneNum)
	request.SetSignName(signture)
	request.SetTemplateCode(templateCode)

	r := rand.New(rand.NewSource(time.Now().Unix()))
	random := r.Int31n(1000000)
	randomStr := fmt.Sprintf("%06d", random)
	request.SetTemplateParam(`{"code":"` + randomStr + `"}`)

	response, err := client.SendSms(request)
	if err != nil {
		return "", err
	}
	log.Println(response)
	return randomStr, nil
}
