package wechat

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"app/global"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	wxapp "github.com/wechatpay-apiv3/wechatpay-go/services/payments/app"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

// Config 微信支付配置
type Config struct {
	AppID          string
	MchID          string
	APIv3Key       string
	CertSerialNo   string
	PrivateKeyPath string
	NotifyURL      string
}

// AppPayParams App 支付所需参数
type AppPayParams struct {
	Appid     string `json:"appid"`
	Partnerid string `json:"partnerid"`
	Prepayid  string `json:"prepayid"`
	Package   string `json:"package"`
	Noncestr  string `json:"noncestr"`
	Timestamp string `json:"timestamp"`
	Sign      string `json:"sign"`
}

var (
	cfg       *Config
	client    *core.Client
	notifyHdl *notify.Handler
)

// Init 初始化微信支付客户端。配置不完整时仅记录警告，不阻断服务启动。
func Init() (err error) {
	cfg = &Config{
		AppID:          global.Global.Viper.GetString("wechat.appid"),
		MchID:          global.Global.Viper.GetString("wechat.mchID"),
		APIv3Key:       global.Global.Viper.GetString("wechat.apiV3Key"),
		CertSerialNo:   global.Global.Viper.GetString("wechat.certSerialNo"),
		PrivateKeyPath: global.Global.Viper.GetString("wechat.privateKeyPath"),
		NotifyURL:      global.Global.Viper.GetString("wechat.notifyURL"),
	}

	if cfg.AppID == "" || cfg.MchID == "" || cfg.APIv3Key == "" || cfg.PrivateKeyPath == "" {
		global.Global.Logger.Warn("微信支付配置不完整，会员充值功能不可用，请配置 wechat 节")
		return nil
	}

	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(cfg.PrivateKeyPath)
	if err != nil {
		return fmt.Errorf("加载商户私钥失败: %w", err)
	}

	ctx := context.Background()
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(cfg.MchID, cfg.CertSerialNo, mchPrivateKey, cfg.APIv3Key),
	}
	client, err = core.NewClient(ctx, opts...)
	if err != nil {
		return fmt.Errorf("初始化微信支付客户端失败: %w", err)
	}

	// 构建回调通知处理器（用于验签和解密）
	certVisitor := downloader.MgrInstance().GetCertificateVisitor(cfg.MchID)
	verifier := verifiers.NewSHA256WithRSAVerifier(certVisitor)
	notifyHdl, err = notify.NewRSANotifyHandler(cfg.APIv3Key, verifier)
	if err != nil {
		return fmt.Errorf("初始化微信支付回调处理器失败: %w", err)
	}

	return nil
}

// AppPay 创建 App 支付订单，返回 App 端调起支付所需参数
func AppPay(ctx context.Context, orderNo, description string, amountFen int64) (*AppPayParams, error) {
	if client == nil {
		return nil, fmt.Errorf("微信支付客户端未初始化")
	}

	svc := wxapp.AppApiService{Client: client}
	resp, _, err := svc.Prepay(ctx, wxapp.PrepayRequest{
		Appid:       core.String(cfg.AppID),
		Mchid:       core.String(cfg.MchID),
		Description: core.String(description),
		OutTradeNo:  core.String(orderNo),
		NotifyUrl:   core.String(cfg.NotifyURL),
		Amount: &wxapp.Amount{
			Currency: core.String("CNY"),
			Total:    core.Int64(amountFen),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("微信支付下单失败: %w", err)
	}
	if resp.PrepayId == nil || *resp.PrepayId == "" {
		return nil, fmt.Errorf("微信支付下单返回 prepay_id 为空")
	}

	// 生成 App 支付签名参数
	nonce := randomNonce()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	prepayID := *resp.PrepayId
	message := fmt.Sprintf("%s\n%s\n%s\n%s\n", cfg.AppID, timestamp, nonce, prepayID)

	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(cfg.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("加载商户私钥失败: %w", err)
	}
	sign, err := utils.SignSHA256WithRSA(message, mchPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("生成支付签名失败: %w", err)
	}

	return &AppPayParams{
		Appid:     cfg.AppID,
		Partnerid: cfg.MchID,
		Prepayid:  prepayID,
		Package:   "Sign=WXPay",
		Noncestr:  nonce,
		Timestamp: timestamp,
		Sign:      sign,
	}, nil
}

// ParseNotify 解析并验签微信支付回调
func ParseNotify(ctx context.Context, req *http.Request) (*payments.Transaction, error) {
	if notifyHdl == nil {
		return nil, fmt.Errorf("微信支付回调处理器未初始化")
	}
	content := new(payments.Transaction)
	_, err := notifyHdl.ParseNotifyRequest(ctx, req, content)
	if err != nil {
		return nil, fmt.Errorf("解析微信支付回调失败: %w", err)
	}
	return content, nil
}

// QueryOrderByOutTradeNo 根据商户订单号查询支付订单
func QueryOrderByOutTradeNo(ctx context.Context, outTradeNo string) (*payments.Transaction, error) {
	if client == nil {
		return nil, fmt.Errorf("微信支付客户端未初始化")
	}
	svc := wxapp.AppApiService{Client: client}
	resp, _, err := svc.QueryOrderByOutTradeNo(ctx, wxapp.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(outTradeNo),
		Mchid:      core.String(cfg.MchID),
	})
	if err != nil {
		return nil, fmt.Errorf("查询微信支付订单失败: %w", err)
	}
	return resp, nil
}

// randomNonce 生成 32 位随机字符串
func randomNonce() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
