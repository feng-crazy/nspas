package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/nspas/go-service/config"
	"github.com/nspas/go-service/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WeChatService struct {
	cfg *config.Config
}

func NewWeChatService(cfg *config.Config) *WeChatService {
	return &WeChatService{cfg: cfg}
}

// WeChatOAuthURLResponse 微信授权URL响应
type WeChatOAuthURLResponse struct {
	URL string `json:"url"`
}

// WeChatTokenResponse 微信令牌响应
type WeChatTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid,omitempty"`
}

// WeChatUserResponse 微信用户信息响应
type WeChatUserResponse struct {
	OpenID     string `json:"openid"`
	Nickname   string `json:"nickname"`
	Sex        int    `json:"sex"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Country    string `json:"country"`
	HeadImgURL string `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	UnionID    string `json:"unionid,omitempty"`
}

// GetWeChatOAuthURL 获取微信授权URL
func (s *WeChatService) GetWeChatOAuthURL(ctx context.Context, state string) string {
	wechatCfg := s.cfg.WeChat
	authURL := fmt.Sprintf(
		"https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect",
		wechatCfg.AppID,
		url.QueryEscape(wechatCfg.RedirectURI),
		wechatCfg.Scope,
		state,
	)

	return authURL
}

// GetWeChatToken 根据code获取微信令牌
func (s *WeChatService) GetWeChatToken(ctx context.Context, code string) (*WeChatTokenResponse, error) {
	wechatCfg := s.cfg.WeChat
	tokenURL := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		wechatCfg.AppID,
		wechatCfg.AppSecret,
		code,
	)

	// 发送请求
	resp, err := http.Get(tokenURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var tokenResp WeChatTokenResponse
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// GetWeChatUserInfo 获取微信用户信息
func (s *WeChatService) GetWeChatUserInfo(ctx context.Context, accessToken, openID string) (*WeChatUserResponse, error) {
	userInfoURL := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN",
		accessToken,
		openID,
	)

	// 发送请求
	resp, err := http.Get(userInfoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var userResp WeChatUserResponse
	err = json.Unmarshal(body, &userResp)
	if err != nil {
		return nil, err
	}

	return &userResp, nil
}

// WeChatLogin 微信登录
func (s *WeChatService) WeChatLogin(ctx context.Context, code string) (*models.User, string, error) {
	// 1. 根据code获取微信令牌
	token, err := s.GetWeChatToken(ctx, code)
	if err != nil {
		return nil, "", err
	}

	// 2. 根据令牌获取微信用户信息
	wechatUser, err := s.GetWeChatUserInfo(ctx, token.AccessToken, token.OpenID)
	if err != nil {
		return nil, "", err
	}

	// 3. 根据微信用户信息创建或获取系统用户
	// 这里需要实现根据openid查找用户，如果不存在则创建
	// 暂时返回mock数据
	user := &models.User{
		ID:        primitive.NewObjectID(),
		Email:     fmt.Sprintf("%s@wechat.com", wechatUser.OpenID),
		Phone:     "",
		Password:  "", // 微信登录不需要密码
		Role:      models.RoleUser,
	}

	// 4. 生成JWT token
	// 这里需要调用现有的generateToken方法，暂时返回空字符串
	return user, "", nil
}
