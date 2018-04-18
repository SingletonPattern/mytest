package wechat

import (
	"os"

	"weixin/initconfig"
)

var WxConfig = initconfig.WxConfig

var WxFunction = NewWechat(WxConfig)

type WechatFunc interface {
	//获取开放平台接口令牌
	GetComponentAccessToken() string
	//指定是否强制刷新开放平台接口令牌
	RefreshComponentAccessToken(forceRefresh bool) string
	//获取公众号管理员授权码
	GetPreAuthCode() (string, []error)
	//获取公众号授权相关信息及初始化公众号接口令牌
	GetAuthorization(authCode string) (*ApiQueryAuth, []error)
	//获取公众号接口令牌
	GetAuthAccessToken() string
	//指定是否强制刷新公众号接口令牌
	RefreshAccessToken(forceRefresh bool) string
	//获取授权信息
	GetAuthorizerInfo() (*ApiGetAuthorizerInfo, []error)
	//获取选项设置值
	GetAuthorizerOption(option string) (*ApiGetAuthorizerOption, []error)
	//设置选项值
	SetAuthorizerOption(optionName, optionValue string) []error
	//公众号管理员授权登录页面
	CreateComponentLoginPageUrl(redirectUrl, preAuthCode string) string

	GetJsapiTicket() string

	RefreshJsapiTicket(forceRefresh bool) string

	CreateJsapiSignature(uri string) *JsApiSignature

	BuildOauth2AuthorizationUrl(redirectUrl, scope, state string) string

	GetOAuth2AccessToken(code string) (*OAuth2AccessToken, []error)

	OAuth2GwtUserInfo(oAuth2AccessToken *OAuth2AccessToken, lang string) (*UserInfo, []error)

	MaterialImageOrVoiceDownload(mediaId string) (*os.File, error)

	MaterialFileBatchGet(mediaType string, offset int64, count int64) (*NewsList, []error)

	MiniAppLogin(jsCode string) (*MiniAppSessionKey, []error)
}

type wechatConfig struct {
	config initconfig.WxOpenConfigInterface
}

//初始化微信开放平台功能函数
func NewWechat(wxConfig initconfig.WxOpenConfigInterface) WechatFunc {
	return &wechatConfig{
		config: wxConfig,
	}
}