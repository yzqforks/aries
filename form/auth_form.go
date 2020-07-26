package form

import "aries/model"

// 登录表单
type LoginForm struct {
	Username   string `json:"username" binding:"required,min=3,max=30" label:"用户名"`   // 用户名
	Pwd        string `json:"pwd" binding:"required,min=6,max=20" label:"密码"`         // 密码
	CaptchaId  string `json:"captcha_id" binding:"required" label:"验证码 ID"`           // 验证码 ID
	CaptchaVal string `json:"captcha_val" binding:"required,min=4,max=4" label:"验证码"` // 验证码
}

// 注册表单
type RegisterForm struct {
	Username  string `json:"username" binding:"required,min=3,max=30" label:"用户名"`                // 用户名
	Nickname  string `json:"nickname" binding:"required,min=3,max=30" label:"昵称"`                 // 用户名
	Pwd       string `json:"pwd" binding:"required,min=6,max=20" label:"密码"`                      // 密码
	SecondPwd string `json:"second_pwd" binding:"required,min=6,max=20,eqfield=Pwd" label:"确认密码"` // 确认密码
	Email     string `json:"email" binding:"required,max=30,email" label:"邮箱"`                    // 邮箱
	SiteUrl   string `json:"site_url" binding:"required,max=255,url" label:"网站地址"`                // 网址
}

// 忘记密码表单
type ForgetPwdForm struct {
	Email string `json:"email" binding:"required,max=30,email" label:"邮箱"` // 邮箱
}

// 绑定登录表单到实体结构
func (form LoginForm) BindToModel() model.User {
	return model.User{
		Username: form.Username,
		Pwd:      form.Pwd,
	}
}

// 绑定注册表单到实体结构
func (form RegisterForm) BindToModel() model.User {
	return model.User{
		Username: form.Username,
		Nickname: form.Nickname,
		Pwd:      form.Pwd,
		Email:    form.Email,
		SiteUrl:  form.SiteUrl,
	}
}

/*// 定义返回错误信息
func (form LoginForm) GetError(err validator.ValidationErrors) string {
	for _, fieldError := range err {
		if fieldError.Field() == "Username" {
			switch fieldError.Tag() {
			case "required":
				return "请输入用户名"
			}
		}
		if fieldError.Field() == "Pwd" {
			switch fieldError.Tag() {
			case "required":
				return "请输入密码"
			}
		}
	}
	return err.Error()
}*/
