package service

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/dengdeng-harmenyos/server/internal/config"
	"github.com/dengdeng-harmenyos/server/internal/logger"
)

// HuaweiPushService 华为Push Kit v3推送服务
type HuaweiPushService struct {
	config      config.HuaweiPushConfig
	httpClient  *http.Client
	accessToken string
	tokenExpiry time.Time
	tokenMutex  sync.Mutex
	privateKey  *rsa.PrivateKey
	keyID       string
	subAccount  string
	projectID   string
}

// ServiceAccountConfig 服务账号密钥文件结构
type ServiceAccountConfig struct {
	ProjectID  string `json:"project_id"`
	KeyID      string `json:"key_id"`
	PrivateKey string `json:"private_key"`
	SubAccount string `json:"sub_account"`
	TokenURI   string `json:"token_uri"`
}

// JWTClaims JWT声明
type JWTClaims struct {
	Iss string `json:"iss"` // sub_account
	Aud string `json:"aud"` // token_uri
	Iat int64  `json:"iat"` // 签发时间
	Exp int64  `json:"exp"` // 过期时间
	jwt.RegisteredClaims
}

// NewHuaweiPushService 创建推送服务
func NewHuaweiPushService(cfg config.HuaweiPushConfig) (*HuaweiPushService, error) {
	logger.Info("Initializing Huawei Push Service...")
	logger.Debug("  Service Account File: %s", cfg.ServiceAccountFile)

	// 从服务账号密钥文件读取配置
	privateKey, keyID, subAccount, projectID, err := loadServiceAccount(cfg.ServiceAccountFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load service account: %w", err)
	}

	logger.Info("✓ Huawei Push service account loaded")
	logger.Debug("  Key ID: %s", keyID)
	logger.Debug("  Sub Account: %s", subAccount)
	logger.Debug("  Project ID: %s", projectID)

	return &HuaweiPushService{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		privateKey: privateKey,
		keyID:      keyID,
		subAccount: subAccount,
		projectID:  projectID,
	}, nil
}

// loadServiceAccount 从嵌入的配置或文件加载服务账号
func loadServiceAccount(filePath string) (*rsa.PrivateKey, string, string, string, error) {
	logger.Debug("Loading service account...")

	var data []byte
	var err error

	// 优先使用嵌入的配置
	embeddedJSON := config.GetEmbeddedPrivateJSON()
	if embeddedJSON != "" {
		logger.Debug("Using embedded service account configuration")
		data = []byte(embeddedJSON)
	} else {
		// 如果嵌入配置为空，从文件读取（用于开发环境）
		logger.Debug("Loading service account from file: %s", filePath)
		data, err = os.ReadFile(filePath)
		if err != nil {
			return nil, "", "", "", fmt.Errorf("failed to read service account file: %w", err)
		}
	}

	var config ServiceAccountConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, "", "", "", fmt.Errorf("failed to parse service account file: %w", err)
	}

	// 解析私钥
	privateKey, err := parsePrivateKey(config.PrivateKey)
	if err != nil {
		return nil, "", "", "", fmt.Errorf("failed to parse private key: %w", err)
	}

	if config.KeyID == "" || config.SubAccount == "" || config.ProjectID == "" {
		return nil, "", "", "", fmt.Errorf("key_id, sub_account or project_id is empty")
	}

	return privateKey, config.KeyID, config.SubAccount, config.ProjectID, nil
}

// parsePrivateKey 解析PEM格式的私钥
func parsePrivateKey(pemString string) (*rsa.PrivateKey, error) {
	// 清理PEM字符串
	pemString = strings.ReplaceAll(pemString, "\\n", "\n")

	block, _ := pem.Decode([]byte(pemString))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PKCS8 private key: %w", err)
	}

	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA private key")
	}

	return rsaKey, nil
}

// V3推送请求结构
type V3PushRequest struct {
	Payload     interface{}  `json:"payload"`
	Target      V3Target     `json:"target"`
	PushOptions *PushOptions `json:"pushOptions,omitempty"`
}

type V3Target struct {
	Token []string `json:"token"` // 文档要求使用token字段
}

type PushOptions struct {
	TestMessage bool   `json:"testMessage,omitempty"` // 是否为测试消息
	TTL         int    `json:"ttl,omitempty"`         // 消息缓存时间(秒)，默认86400(1天)
	BiTag       string `json:"biTag,omitempty"`       // 批量任务消息标识
}

// Alert消息Payload（通知消息，push-type=0）
type AlertPayload struct {
	Notification Notification `json:"notification"`
}

// Notification 通知消息结构体
type Notification struct {
	Category    string      `json:"category"`        // 通知消息类型（必填）
	Title       string      `json:"title"`           // 通知标题（必填）
	Body        string      `json:"body"`            // 通知内容（必填）
	Image       string      `json:"image,omitempty"` // 右侧大图标URL
	ClickAction ClickAction `json:"clickAction"`     // 点击消息动作（必填）
	Badge       *Badge      `json:"badge,omitempty"` // 通知消息角标
	Sound       string      `json:"sound,omitempty"` // 自定义铃声
}

// ClickAction 点击行为
type ClickAction struct {
	ActionType int                    `json:"actionType"`       // 0:打开首页, 1:打开自定义页面, 3:清除通知
	Action     string                 `json:"action,omitempty"` // 应用内置页面action
	URI        string                 `json:"uri,omitempty"`    // 应用内置页面uri
	Data       map[string]interface{} `json:"data,omitempty"`   // 传递给应用的数据
}

// Badge 角标
type Badge struct {
	AddNum int `json:"addNum,omitempty"` // 角标累加数字(1-99)
	SetNum int `json:"setNum,omitempty"` // 角标设置数字(0-99)
}

// 卡片刷新Payload（push-type=1）
type FormUpdatePayload struct {
	FormID      int64                  `json:"formId"`           // 服务卡片实例ID（必填）
	Version     int                    `json:"version"`          // 卡片刷新版本号（必填）
	ModuleName  string                 `json:"moduleName"`       // 服务卡片模块名称（必填）
	FormName    string                 `json:"formName"`         // 服务卡片名称（必填）
	AbilityName string                 `json:"abilityName"`      // 服务卡片Ability名称（必填）
	FormData    map[string]interface{} `json:"formData"`         // 待刷新卡片数据（必填）
	Images      []FormImage            `json:"images,omitempty"` // 卡片图片数据
}

// FormImage 卡片图片
type FormImage struct {
	KeyName string `json:"keyName"` // 图片对应的key
	URL     string `json:"url"`     // 图片下载地址（HTTPS）
	Require int    `json:"require"` // 图片刷新策略：1=失败不刷新，0=失败仅刷新文字
}

// 后台消息Payload（push-type=6）
type BackgroundPayload struct {
	ExtraData string `json:"extraData"`           // 传递给应用的数据（必填）
	ProxyData string `json:"proxyData,omitempty"` // 数据代理："ENABLE"
}

// 语音播报消息Payload（push-type=2）
type ExtensionPayload struct {
	Notification Notification `json:"notification"` // category必须为"PLAY_VOICE"
	ExtraData    string       `json:"extraData"`    // 语音播报额外数据（必填）
}

// 应用内通话消息Payload（push-type=10）
type VoIPCallPayload struct {
	ExtraData string `json:"extraData"` // 传递给应用的数据（必填）
}

// 推送响应
type PushResponse struct {
	Code      string `json:"code"`
	Msg       string `json:"msg"`
	RequestID string `json:"requestId"`
}

// SendNotification 发送通知消息（Alert）
func (s *HuaweiPushService) SendNotification(pushToken, title, body string, data map[string]interface{}) error {
	logger.Debug("Sending notification: title=%s, body=%s, token=%s...", title, body, pushToken[:20])

	// 构建点击行为
	clickAction := ClickAction{
		ActionType: 0, // 0: 打开应用首页
	}

	// 如果有额外数据，设置为点击时传递的数据
	if data != nil {
		clickAction.Data = data
		logger.Debug("  Extra data: %+v", data)
	}

	// 构建通知消息payload
	payload := AlertPayload{
		Notification: Notification{
			Category:    "WORK", // 默认使用工作提醒类型，可根据业务需求修改
			Title:       title,
			Body:        body,
			ClickAction: clickAction,
			Badge:       &Badge{AddNum: 1}, // 默认角标加1
		},
	}

	// 默认使用测试消息选项
	options := &PushOptions{
		TestMessage: false,
		TTL:         86400, // 1天
	}

	return s.sendPush(0, []string{pushToken}, payload, options)
}

// SendFormUpdate 发送卡片刷新消息
func (s *HuaweiPushService) SendFormUpdate(pushToken string, formID int64, version int, moduleName, formName, abilityName string, formData map[string]interface{}) error {
	payload := FormUpdatePayload{
		FormID:      formID,
		Version:     version,
		ModuleName:  moduleName,
		FormName:    formName,
		AbilityName: abilityName,
		FormData:    formData,
	}

	options := &PushOptions{
		TTL: 666,
	}

	return s.sendPush(1, []string{pushToken}, payload, options)
}

// SendFormUpdateSimple 发送卡片刷新消息（简化版，使用默认参数）
func (s *HuaweiPushService) SendFormUpdateSimple(pushToken string, formID int64, formData map[string]interface{}) error {
	// 使用默认值，实际使用时应该从配置或数据库中获取这些参数
	return s.SendFormUpdate(pushToken, formID, 0, "entry", "widget", "EntryAbility", formData)
}

// SendBackgroundMessage 发送后台消息
func (s *HuaweiPushService) SendBackgroundMessage(pushToken string, extraData string) error {
	payload := BackgroundPayload{
		ExtraData: extraData,
	}

	return s.sendPush(6, []string{pushToken}, payload, nil)
}

// SendVoIPCall 发送应用内通话消息
func (s *HuaweiPushService) SendVoIPCall(pushToken string, extraData string) error {
	payload := VoIPCallPayload{
		ExtraData: extraData,
	}

	// 应用内通话消息建议TTL为30-60秒
	options := &PushOptions{
		TTL: 30,
	}

	return s.sendPush(10, []string{pushToken}, payload, options)
}

// SendBatchNotification 批量发送通知消息
func (s *HuaweiPushService) SendBatchNotification(pushTokens []string, title, body string, data map[string]interface{}) error {
	if len(pushTokens) > 1000 {
		return fmt.Errorf("batch size exceeds limit: %d (max 1000)", len(pushTokens))
	}

	// 构建点击行为
	clickAction := ClickAction{
		ActionType: 0, // 0: 打开应用首页
	}

	if data != nil {
		clickAction.Data = data
	}

	// 构建通知消息payload
	payload := AlertPayload{
		Notification: Notification{
			Category:    "WORK",
			Title:       title,
			Body:        body,
			ClickAction: clickAction,
			Badge:       &Badge{AddNum: 1}, // 默认角标加1
		},
	}

	options := &PushOptions{
		TestMessage: false,
		TTL:         86400,
	}

	return s.sendPush(0, pushTokens, payload, options)
}

// sendPush 通用推送方法
func (s *HuaweiPushService) sendPush(pushType int, tokens []string, payload interface{}, options *PushOptions) error {
	logger.Debug("sendPush: type=%d, tokens=%d", pushType, len(tokens))

	// 获取JWT token
	jwtToken, err := s.getAccessToken()
	if err != nil {
		logger.Error("Failed to get JWT token: %v", err)
		return fmt.Errorf("failed to get JWT token: %w", err)
	}
	logger.Debug("✓ JWT token obtained")

	// 构建请求体
	requestBody := V3PushRequest{
		Payload: payload,
		Target: V3Target{
			Token: tokens, // 使用token字段
		},
		PushOptions: options,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	logger.Debug("Request payload: %s", string(jsonData))

	// 构建请求
	url := fmt.Sprintf("%s/%s/messages:send", s.config.PushAPIURL, s.projectID)
	logger.Debug("Push URL: %s", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("push-type", strconv.Itoa(pushType))

	// 发送请求
	logger.Debug("Sending push request to Huawei...")
	resp, err := s.httpClient.Do(req)
	if err != nil {
		logger.Error("Failed to send request: %v", err)
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}
	logger.Debug("Push response (status=%d): %s", resp.StatusCode, string(body))

	// 解析响应
	var pushResp PushResponse
	if err := json.Unmarshal(body, &pushResp); err != nil {
		logger.Error("Failed to parse response: %v", err)
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// 检查响应状态
	if pushResp.Code != "80000000" {
		logger.Error("Push failed: code=%s, msg=%s, requestId=%s", pushResp.Code, pushResp.Msg, pushResp.RequestID)
		return fmt.Errorf("push failed: code=%s, msg=%s", pushResp.Code, pushResp.Msg)
	}

	logger.Info("✓ Push sent successfully (requestId=%s)", pushResp.RequestID)
	return nil
}

// getAccessToken 生成JWT token作为访问令牌
func (s *HuaweiPushService) getAccessToken() (string, error) {
	s.tokenMutex.Lock()
	defer s.tokenMutex.Unlock()

	// 如果token还有效（提前5分钟刷新），直接返回
	if s.accessToken != "" && time.Now().Add(5*time.Minute).Before(s.tokenExpiry) {
		logger.Debug("Using cached JWT token (expires in %v)", time.Until(s.tokenExpiry))
		return s.accessToken, nil
	}

	logger.Debug("Generating new JWT token...")

	// 签发时间和过期时间
	iat := time.Now().Unix()
	exp := iat + 3600 // 1小时后过期

	// 创建claims
	claims := JWTClaims{
		Iss: s.subAccount,
		Aud: "https://oauth-login.cloud.huawei.com/oauth2/v3/token",
		Iat: iat,
		Exp: exp,
	}

	// 创建token（使用PS256算法）
	token := jwt.NewWithClaims(&jwt.SigningMethodRSAPSS{
		SigningMethodRSA: jwt.SigningMethodRS256,
		Options: &rsa.PSSOptions{
			SaltLength: rsa.PSSSaltLengthEqualsHash,
			Hash:       0,
		},
	}, claims)

	// 设置header
	token.Header["kid"] = s.keyID
	token.Header["typ"] = "JWT"
	token.Header["alg"] = "PS256"

	// 签名
	jwtToken, err := token.SignedString(s.privateKey)
	if err != nil {
		logger.Error("Failed to sign JWT token: %v", err)
		return "", fmt.Errorf("failed to sign JWT token: %w", err)
	}

	// 保存token和过期时间
	s.accessToken = jwtToken
	s.tokenExpiry = time.Unix(exp, 0)

	logger.Info("✓ New JWT token generated (expires in 3600s)")
	logger.Debug("JWT token: %s...", jwtToken[:min(50, len(jwtToken))])

	return s.accessToken, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
