package deviceauth

import (
	"sync"
	"time"
)

// Store 设备授权存储接口
type Store interface {
	// Save 保存缓存条目（key 可以是 deviceCode 或 userCode）
	Save(key string, value DeviceAuthCache)
	// LoadByDeviceCode 按设备码加载
	LoadByDeviceCode(deviceCode string) (DeviceAuthCache, bool)
	// LoadByUserCode 按用户码加载
	LoadByUserCode(userCode string) (DeviceAuthCache, bool)
	// Update 更新缓存条目（同时更新 deviceCode 和 userCode 两个 key）
	Update(cache DeviceAuthCache)
	// Delete 按设备码删除（同时删除 deviceCode 和 userCode 两个 key）
	Delete(deviceCode string)
}

// MemoryStore 基于 sync.Map 的内存存储实现
type MemoryStore struct {
	m sync.Map
}

// NewMemoryStore 创建内存存储并启动过期清理协程
func NewMemoryStore() *MemoryStore {
	s := &MemoryStore{}
	go s.cleanupLoop()
	return s
}

func (s *MemoryStore) Save(key string, value DeviceAuthCache) {
	s.m.Store(key, value)
}

func (s *MemoryStore) LoadByDeviceCode(deviceCode string) (DeviceAuthCache, bool) {
	v, ok := s.m.Load(deviceCode)
	if !ok {
		return DeviceAuthCache{}, false
	}
	cache, ok := v.(DeviceAuthCache)
	return cache, ok
}

func (s *MemoryStore) LoadByUserCode(userCode string) (DeviceAuthCache, bool) {
	v, ok := s.m.Load(userCode)
	if !ok {
		return DeviceAuthCache{}, false
	}
	cache, ok := v.(DeviceAuthCache)
	return cache, ok
}

func (s *MemoryStore) Update(cache DeviceAuthCache) {
	if cache.DeviceCode != "" {
		s.m.Store(cache.DeviceCode, cache)
	}
	if cache.UserCode != "" {
		s.m.Store(cache.UserCode, cache)
	}
}

func (s *MemoryStore) Delete(deviceCode string) {
	if cache, ok := s.LoadByDeviceCode(deviceCode); ok {
		if cache.UserCode != "" {
			s.m.Delete(cache.UserCode)
		}
		s.m.Delete(deviceCode)
	}
}

// cleanupLoop 定期清理过期条目
func (s *MemoryStore) cleanupLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for now := range ticker.C {
		_ = now
		s.m.Range(func(key, value any) bool {
			cache, ok := value.(DeviceAuthCache)
			if !ok {
				return true
			}
			if cache.IsExpired() {
				s.m.Delete(key)
			}
			return true
		})
	}
}
