package biz

import (
	"context"
	"fmt"
	"golang.org/x/crypto/ssh"
	"net/http"
	"os"
	"sync"
	"time"
)

// 定义健康状态等级
type HealthStatus uint32

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

type Server struct {
	Host           string
	User           string
	Port           int32
	PrivateKeyPath string
	Command        string
}

// Checker 代表服务的健康检查器，并记录当前健康状态等级
type Checker struct {
	Name          string
	URL           string
	Server        Server
	Status        HealthStatus // 当前健康状态等级
	lastCheckedAt time.Time    // 最近一次检查时间
	mu            sync.RWMutex // 用于并发控制
}

// SetStatus 更新健康状态等级及检查时间
func (c *Checker) SetStatus(status HealthStatus) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Status = status
	c.lastCheckedAt = time.Now()
}

// GetStatus 获取当前健康状态等级
func (c *Checker) GetStatus() HealthStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Status
}

func (c *Checker) runServiceServerCommand() error {
	if c.Server.Host == "" || c.Server.User == "" || c.Server.Port == 0 || c.Server.PrivateKeyPath == "" || c.Server.Command == "" {
		return nil
	}

	// 读取私钥文件
	key, err := os.ReadFile(c.Server.PrivateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read file, error: %v", err)
	}

	// 解析私钥
	var signer ssh.Signer
	if signer, err = ssh.ParsePrivateKey(key); err != nil {
		return fmt.Errorf("failed to parse private key, error: %v", err)
	}

	// 配置SSH客户端
	sshConfig := &ssh.ClientConfig{
		User: c.Server.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 忽略主机密钥验证（不推荐用于生产环境）
	}

	address := fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
	var client *ssh.Client
	if client, err = ssh.Dial("tcp", address, sshConfig); err != nil {
		return fmt.Errorf("failed to ssh dial, error: %v", err)
	}
	defer client.Close()

	// 创建一个新的SSH会话
	var session *ssh.Session
	if session, err = client.NewSession(); err != nil {
		return fmt.Errorf("failed to new ssh session, error: %v", err)
	}
	defer session.Close()

	// 执行命令并获取输出
	var output []byte
	if output, err = session.CombinedOutput(c.Server.Command); err != nil {
		return fmt.Errorf("failed to conbined output, error: %v", err)
	}

	// 输出命令结果
	fmt.Printf("命令输出:\n%s\n", string(output))
	return nil
}

// addChecker 将新的 Checker 添加到 checkers 列表中.
func (s *Service) addChecker(checker *Checker) {
	s.checkersMu.Lock()
	defer s.checkersMu.Unlock()
	s.checkers = append(s.checkers, checker)
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
	case Unknown, Healthy:
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
	case Degraded, Unknown:
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
			_ = s.alert.SendMessage(context.Background(), fmt.Sprintf("警告: %s服务健康检查连续失败，当前状态: 不可用，请及时处理！", checker.Name))

			// 连接ssh执行命令
			if err := checker.runServiceServerCommand(); err != nil {
				fmt.Println("runServiceServerCommand error:", err)
			}

			// 当前状态是不可用的。延长下一次检查间隔
			time.Sleep(5 * time.Second)
		}

		return nil
	}
}

// GetAllCheckers 返回所有的健康检查器列表.
func (s *Service) GetAllCheckers() []*Checker {
	s.checkersMu.RLock()
	defer s.checkersMu.RUnlock()
	return s.checkers
}
