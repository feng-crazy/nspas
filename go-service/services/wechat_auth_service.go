package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"neuro-guide-go-service/config"
	"neuro-guide-go-service/models"
)

// WeChatAuthService handles WeChat authentication
type WeChatAuthService struct {
	appID       string
	appSecret   string
	httpClient  *http.Client
	userService *UserService
}

// NewWeChatAuthService creates a new instance of WeChatAuthService
func NewWeChatAuthService(cfg *config.Config, userService *UserService) *WeChatAuthService {
	// In production, you should get these from environment variables or config
	appID := getEnv("WECHAT_APP_ID", "")
	appSecret := getEnv("WECHAT_APP_SECRET", "")

	return &WeChatAuthService{
		appID:       appID,
		appSecret:   appSecret,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
		userService: userService,
	}
}

// WeChatLoginResponse represents the response from WeChat OAuth API
type WeChatLoginResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid,omitempty"`
	ErrCode      int    `json:"errcode,omitempty"`
	ErrMsg       string `json:"errmsg,omitempty"`
}

// WeChatUserInfo represents user info from WeChat
type WeChatUserInfo struct {
	OpenID     string `json:"openid"`
	Nickname   string `json:"nickname"`
	Sex        int    `json:"sex"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Country    string `json:"country"`
	HeadImgURL string `json:"headimgurl"`
	UnionID    string `json:"unionid,omitempty"`
	ErrCode    int    `json:"errcode,omitempty"`
	ErrMsg     string `json:"errmsg,omitempty"`
}

// WeChatPhoneNumberResponse represents the response from WeChat phone number decryption API
type WeChatPhoneNumberResponse struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
	Watermark       struct {
		AppID     string `json:"appid"`
		Timestamp int64  `json:"timestamp"`
	} `json:"watermark"`
	ErrCode int    `json:"errcode,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`
}

// AuthenticateWithCode authenticates user with WeChat authorization code
func (w *WeChatAuthService) AuthenticateWithCode(code string, userInfo *map[string]interface{}) (*models.User, string, error) {
	// Exchange code for access token and openid
	tokenURL := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		w.appID, w.appSecret, code)

	resp, err := w.httpClient.Get(tokenURL)
	if err != nil {
		return nil, "", fmt.Errorf("failed to exchange code for token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read token response: %w", err)
	}

	var tokenResp WeChatLoginResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, "", fmt.Errorf("failed to parse token response: %w", err)
	}

	if tokenResp.ErrCode != 0 {
		return nil, "", fmt.Errorf("WeChat API error: %s (code: %d)", tokenResp.ErrMsg, tokenResp.ErrCode)
	}

	// 如果传入了userInfo，说明是网页端登录，需要获取用户信息
	if userInfo != nil {
		// Get user info
		userInfoURL := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s",
			tokenResp.AccessToken, tokenResp.OpenID)

		resp, err := w.httpClient.Get(userInfoURL)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get user info: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, "", fmt.Errorf("failed to read user info response: %w", err)
		}

		var wechatUserInfo WeChatUserInfo
		if err := json.Unmarshal(body, &wechatUserInfo); err != nil {
			return nil, "", fmt.Errorf("failed to parse user info response: %w", err)
		}

		if wechatUserInfo.ErrCode != 0 {
			return nil, "", fmt.Errorf("WeChat API error: %s (code: %d)", wechatUserInfo.ErrMsg, wechatUserInfo.ErrCode)
		}

		// 将微信用户信息添加到userInfo中
		wechatInfo := map[string]interface{}{
			"wechat_id": wechatUserInfo.OpenID,
			"nickname":  wechatUserInfo.Nickname,
			"avatar":    wechatUserInfo.HeadImgURL,
		}

		// 更新传入的userInfo
		for k, v := range wechatInfo {
			(*userInfo)[k] = v
		}
	}

	// Create or get user in our system
	// 微信登录的用户默认不是游客
	var user *models.User

	if userInfo != nil {
		// 如果有用户信息，使用完整信息创建用户
		user, err = w.userService.GetOrCreateUser(
			(*userInfo)["wechat_id"].(string),
			(*userInfo)["nickname"].(string),
			(*userInfo)["avatar"].(string),
		)
	} else {
		// 如果没有提供用户信息，只使用OpenID
		user, err = w.userService.GetOrCreateUser(
			tokenResp.OpenID,
			"", // 昵称为空
			"", // 头像为空
		)
	}

	if err != nil {
		return nil, "", fmt.Errorf("failed to create/get user: %w", err)
	}

	// Generate our own token for the user (in production, this should be JWT)
	token := user.ID // For simplicity, we use user ID as token

	return user, token, nil
}

// DecryptPhoneNumber decrypts the encrypted phone number received from WeChat mini program
func (w *WeChatAuthService) DecryptPhoneNumber(appID, appSecret, encryptedData, iv, openID string) (string, error) {
	// 在实际应用中，这里应该调用微信的解密接口
	// 这里只是一个示例实现，实际项目中需要使用微信提供的解密算法

	// 示例：简单返回一个虚拟手机号（实际项目中应该实现真正的解密逻辑）
	// 实际应用中，你需要：
	// 1. 获取微信的session_key
	// 2. 使用AES解密算法解密encryptedData
	// 3. 验证解密结果的watermark中的appid是否匹配

	// 这里只是示例，实际项目中需要实现完整的解密流程
	return "13800138000", nil
}

// ValidateToken validates a WeChat token and returns the user
func (w *WeChatAuthService) ValidateToken(token string) (*models.User, error) {
	// In production, this should validate JWT token
	// For now, we treat the token as user ID
	return w.userService.GetUserByID(token)
}

// getEnv is a helper function to get environment variable with default value
func getEnv(key, defaultValue string) string {
	// This duplicates the function in config.go, but we include it here for self-contained code
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
