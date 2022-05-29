package dao

import (
	"encoding/json"
	"github.com/rocboss/paopao-ce/global"

	//"errors"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111" // 引入sms
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/rocboss/paopao-ce/internal/model"
)

type JuhePhoneCaptchaRsp struct {
	ErrorCode int    `json:"error_code"`
	Reason    string `json:"reason"`
}

// 根据用户ID获取用户
func (d *Dao) GetUserByID(id int64) (*model.User, error) {
	user := &model.User{
		Model: &model.Model{
			ID: id,
		},
	}

	return user.Get(d.engine)
}

// 根据用户名获取用户
func (d *Dao) GetUserByUsername(username string) (*model.User, error) {
	user := &model.User{
		Username: username,
	}

	return user.Get(d.engine)
}

// 根据手机号获取用户
func (d *Dao) GetUserByPhone(phone string) (*model.User, error) {
	user := &model.User{
		Phone: phone,
	}

	return user.Get(d.engine)
}

// 根据IDs获取用户列表
func (d *Dao) GetUsersByIDs(ids []int64) ([]*model.User, error) {
	user := &model.User{}

	return user.List(d.engine, &model.ConditionsT{
		"id IN ?": ids,
	}, 0, 0)
}

// 根据关键词模糊获取用户列表
func (d *Dao) GetUsersByKeyword(keyword string) ([]*model.User, error) {
	user := &model.User{}

	if strings.Trim(keyword, "") == "" {
		return user.List(d.engine, &model.ConditionsT{
			"ORDER": "id ASC",
		}, 0, 6)
	} else {

		return user.List(d.engine, &model.ConditionsT{
			"username LIKE ?": strings.Trim(keyword, "") + "%",
		}, 0, 6)
	}
}

// 根据关键词模糊获取用户列表
func (d *Dao) GetTagsByKeyword(keyword string) ([]*model.Tag, error) {
	tag := &model.Tag{}

	if strings.Trim(keyword, "") == "" {
		return tag.List(d.engine, &model.ConditionsT{
			"ORDER": "quote_num DESC",
		}, 0, 6)
	} else {

		return tag.List(d.engine, &model.ConditionsT{
			"tag LIKE ?": "%" + strings.Trim(keyword, "") + "%",
			"ORDER":      "quote_num DESC",
		}, 0, 6)
	}
}

// 创建用户
func (d *Dao) CreateUser(user *model.User) (*model.User, error) {
	return user.Create(d.engine)
}

// 更新用户
func (d *Dao) UpdateUser(user *model.User) error {
	return user.Update(d.engine)
}

// 获取最新短信验证码
func (d *Dao) GetLatestPhoneCaptcha(phone string) (*model.Captcha, error) {
	return (&model.Captcha{
		Phone: phone,
	}).Get(d.engine)
}

// 更新短信验证码
func (d *Dao) UsePhoneCaptcha(captcha *model.Captcha) error {
	captcha.UseTimes++
	return captcha.Update(d.engine)
}

// 发送短信验证码
func (d *Dao) SendPhoneCaptcha(phone string) error {
	rand.Seed(time.Now().UnixNano())
	captcha := rand.Intn(900000) + 100000
	m := 5

	//gateway := "https://v.juhe.cn/sms/send"
	//
	//client := resty.New()
	//client.DisableWarn = true
	//resp, err := client.R().
	//	SetFormData(map[string]string{
	//		"mobile":    phone,
	//		"tpl_id":    global.AppSetting.SmsJuheTplID,
	//		"tpl_value": fmt.Sprintf(global.AppSetting.SmsJuheTplVal, captcha, m),
	//		"key":       global.AppSetting.SmsJuheKey,
	//	}).Post(gateway)
	//if err != nil {
	//	return err
	//}
	//
	//if resp.StatusCode() != http.StatusOK {
	//	return errors.New(resp.Status())
	//}
	//
	//result := &JuhePhoneCaptchaRsp{}
	//err = json.Unmarshal(resp.Body(), result)
	//if err != nil {
	//	return err
	//}
	//
	//if result.ErrorCode != 0 {
	//	return errors.New(result.Reason)
	//}

	result := sendTencentSms(phone, strconv.Itoa(captcha))
	if result != nil {
		return result
	}

	// 写入表
	captchaModel := &model.Captcha{
		Phone:     phone,
		Captcha:   strconv.Itoa(captcha),
		ExpiredOn: time.Now().Add(time.Minute * time.Duration(m)).Unix(),
	}
	captchaModel.Create(d.engine)
	return nil
}

func sendTencentSms(phone string, code string) error {
	credential := common.NewCredential(
		global.AppSetting.TencentSecretId,
		global.AppSetting.TencentSecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	cpf.SignMethod = "HmacSHA1"
	client, _ := sms.NewClient(credential, "ap-guangzhou", cpf)
	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr(global.AppSetting.TencentSmsSdkAppId)
	request.SignName = common.StringPtr(global.AppSetting.TencentSignName)
	request.TemplateId = common.StringPtr(global.AppSetting.TencentTemplateId)
	request.TemplateParamSet = common.StringPtrs([]string{code})
	request.PhoneNumberSet = common.StringPtrs([]string{"+86" + phone})
	response, err := client.SendSms(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		global.Logger.Errorf("sms send err: %v", err)
		return err
	}
	if err != nil {
		panic(err)
		global.Logger.Errorf("sms send err: %v", err)
		return err
	}
	b, _ := json.Marshal(response.Response)
	fmt.Printf("%s", b)
	return nil
}
