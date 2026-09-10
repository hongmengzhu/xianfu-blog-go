package faceid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// AliyunFaceIdProvider 阿里云人脸对比提供商
// 通过阿里云 facebody API 进行人脸对比验证
// 注意：完整实现需要引入 github.com/alibabacloud-go/facebody-20191230/v5 等重型依赖
// 当前使用接口抽象，后续可按需引入阿里云 SDK
type AliyunFaceIdProvider struct {
	AccessKey             string
	AccessSecret          string
	Endpoint              string
	QualityScoreThreshold float32
	Client                *http.Client
}

// NewAliyunFaceIdProvider 创建阿里云 Face ID 提供商
func NewAliyunFaceIdProvider(accessKey string, accessSecret string, endPoint string) *AliyunFaceIdProvider {
	return &AliyunFaceIdProvider{
		AccessKey:             accessKey,
		AccessSecret:          accessSecret,
		Endpoint:              endPoint,
		QualityScoreThreshold: 0.65,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// aliyunCompareRequest 阿里云人脸对比请求
type aliyunCompareRequest struct {
	Action                string  `json:"Action"`
	AccessKeyId           string  `json:"AccessKeyId"`
	ImageDataA            string  `json:"ImageDataA"`
	ImageDataB            string  `json:"ImageDataB"`
	QualityScoreThreshold float32 `json:"QualityScoreThreshold"`
}

// aliyunCompareResponse 阿里云人脸对比响应
type aliyunCompareResponse struct {
	Data struct {
		Confidence float64   `json:"Confidence"`
		Thresholds []float64 `json:"Thresholds"`
	} `json:"Data"`
}

// Check 对比两张人脸图片
func (p *AliyunFaceIdProvider) Check(base64ImageA string, base64ImageB string) (bool, error) {
	if p.Endpoint == "" {
		return false, fmt.Errorf("阿里云 Face ID endpoint 未配置")
	}

	// 去除 data URI 前缀
	imageA := strings.Replace(base64ImageA, "data:image/png;base64,", "", -1)
	imageB := strings.Replace(base64ImageB, "data:image/png;base64,", "", -1)

	reqBody := aliyunCompareRequest{
		Action:                "CompareFace",
		AccessKeyId:           p.AccessKey,
		ImageDataA:            imageA,
		ImageDataB:            imageB,
		QualityScoreThreshold: p.QualityScoreThreshold,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return false, fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, p.Endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return false, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := p.Client
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("请求阿里云人脸对比失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, fmt.Errorf("阿里云人脸对比 HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var compareResp aliyunCompareResponse
	if err = json.Unmarshal(respBody, &compareResp); err != nil {
		return false, fmt.Errorf("解析响应失败: %w", err)
	}

	// 判断置信度是否超过阈值
	if len(compareResp.Data.Thresholds) > 0 && compareResp.Data.Thresholds[0] < compareResp.Data.Confidence {
		return true, nil
	}

	return false, nil
}
