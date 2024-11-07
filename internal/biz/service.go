package biz

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// 定义健康状态等级
type HealthStatus int

const (
	Unknown  HealthStatus = iota // 0: 未知状态
	Healthy                      // 1: 健康
	Degraded                     // 2: 降级
	Unstable                     // 3: 不稳定
	Critical                     // 4: 严重
	Down                         // 5: 不可用
)

// Repo .
type Repo interface{}

// Checker 代表服务的健康检查器，并记录当前健康状态等级
type Checker struct {
	Name          string
	URL           string
	Status        HealthStatus // 当前健康状态等级
	LastCheckedAt time.Time    // 最近一次检查时间
	mu            sync.RWMutex // 用于并发控制
}

// SetStatus 更新健康状态等级及检查时间
func (c *Checker) SetStatus(status HealthStatus) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Status = status
	c.LastCheckedAt = time.Now()
}

// GetStatus 获取当前健康状态等级
func (c *Checker) GetStatus() HealthStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Status
}

// IsServiceHealthy 检查服务健康状态
func (s *Service) isServiceHealthy(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// UpdateStatusOnFailure 更新健康状态等级，逐步降级
func (s *Service) updateStatusOnFailure(checker *Checker) {
	switch checker.GetStatus() {
	case Healthy:
		checker.SetStatus(Degraded)
	case Degraded:
		checker.SetStatus(Unstable)
	case Unstable:
		checker.SetStatus(Critical)
	case Critical:
		checker.SetStatus(Down)
	default:
		// 如果已是 Down 状态，不再进一步降级
	}
}

// UpdateStatusOnSuccess 更新健康状态等级，逐步恢复
func (s *Service) updateStatusOnSuccess(checker *Checker) {
	switch checker.GetStatus() {
	case Down:
		checker.SetStatus(Critical)
	case Critical:
		checker.SetStatus(Unstable)
	case Unstable:
		checker.SetStatus(Degraded)
	case Degraded:
		checker.SetStatus(Healthy)
	default:
		// 如果已是 Healthy 状态，无需进一步提升
	}
}

// Healthy 对服务进行健康检查，根据检查结果动态调整健康状态
func (s *Service) healthy(checker *Checker) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		// 检查服务的健康状态
		if s.isServiceHealthy(checker.URL) {
			if checker.GetStatus() == Healthy {
				// 当前状态是正常的。延长下一次检查间隔
				time.Sleep(10 * time.Second)
				return nil
			}

			// 成功时逐步恢复健康状态
			s.updateStatusOnSuccess(checker)

			return nil
		}

		// 失败时逐步降级健康状态
		s.updateStatusOnFailure(checker)

		// 如果降级到不可用状态，发送告警
		if checker.GetStatus() == Down {
			_ = s.alert.SendMessage(context.Background(), fmt.Sprintf("警告: 服务健康检查连续失败，当前状态: 不可用，请及时处理！"))
		}

		return nil
	}
}
