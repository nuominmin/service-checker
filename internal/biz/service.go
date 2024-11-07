package biz

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Repo .
type Repo interface{}

// Checker .
type Checker struct {
	Name       string
	URL        string
	MaxRetries uint32
}

// IsServiceHealthy 检查服务健康状态
func (s *Service) isServiceHealthy(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true
	}

	fmt.Printf("服务响应状态异常，状态码: %d\n", resp.StatusCode)
	return false
}

func (s *Service) healthy(checker *Checker) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		maxRetries := int(checker.MaxRetries)
		for i := 0; i < maxRetries; i++ {
			if s.isServiceHealthy(checker.URL) {
				return nil
			}
			time.Sleep(time.Second * time.Duration(1<<uint(i)))
		}

		// 如果失败次数达到最大重试次数，则触发告警
		_ = s.alert.SendMessage(context.Background(), fmt.Sprintf("警告: 服务健康检查连续失败 %d 次，请及时处理！", maxRetries))
		return nil
	}
}
