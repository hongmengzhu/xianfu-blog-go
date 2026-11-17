package loginLimitPg

// 登录失败限流（内存实现）：按 登录域 + 账号 + 客户端 IP + 租户 计数，
// 达到阈值后短期锁定，防止撞库/暴力破解。

import (
	"strings"
	"sync"
	"time"
)

const (
	// failWindow 统计窗口
	failWindow = 5 * time.Minute
	// maxFailures 窗口内允许的失败次数，超过即锁定
	maxFailures = 5
	// lockDuration 锁定时间
	lockDuration = 15 * time.Minute
	// maxKeys 防止 map 无限增长
	maxKeys = 100000
)

type entry struct {
	fails       int
	firstFail   time.Time
	lockedUntil time.Time
}

var (
	mu      sync.Mutex
	entries = make(map[string]*entry)
	lastGC  time.Time
)

// Key 生成限流键
func Key(parts ...string) string {
	return strings.Join(parts, "|")
}

// Allowed 是否允许本次登录尝试；不允许时返回建议的等待秒数
func Allowed(key string) (bool, int) {
	now := time.Now()
	mu.Lock()
	defer mu.Unlock()
	c := entries[key]
	if nil == c {
		return true, 0
	}
	if now.Before(c.lockedUntil) {
		return false, int(c.lockedUntil.Sub(now).Seconds()) + 1
	}
	// 窗口外清零
	if c.fails > 0 && now.Sub(c.firstFail) > failWindow {
		delete(entries, key)
	}
	return true, 0
}

// Failure 记录一次失败；返回当前是否已被锁定
func Failure(key string) bool {
	now := time.Now()
	mu.Lock()
	defer mu.Unlock()
	c := entries[key]
	if nil == c || now.Sub(c.firstFail) > failWindow {
		c = &entry{firstFail: now}
		entries[key] = c
	}
	c.fails++
	if c.fails >= maxFailures {
		c.lockedUntil = now.Add(lockDuration)
		c.fails = 0
		c.firstFail = now
		return true
	}
	gc(now)
	return false
}

// Reset 登录成功后清除计数
func Reset(key string) {
	mu.Lock()
	defer mu.Unlock()
	delete(entries, key)
}

// gc 定期清理过期条目与超量条目
func gc(now time.Time) {
	if now.Sub(lastGC) < failWindow {
		return
	}
	lastGC = now
	for k, c := range entries {
		if now.After(c.lockedUntil) && now.Sub(c.firstFail) > failWindow {
			delete(entries, k)
		}
	}
	if len(entries) > maxKeys {
		for k := range entries {
			delete(entries, k)
			if len(entries) <= maxKeys/2 {
				break
			}
		}
	}
}
