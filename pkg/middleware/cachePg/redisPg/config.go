package redisPg

import (
	"fmt"
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg/pg"
	"github.com/redis/go-redis/v9"
)

var driverRegistry = map[string]Driver{}

func init() {
	RegisterDriver("DefaultDriver", DefaultDriver{})
}

// Driver interface defines how to create a Redis client.
type Driver interface {
	CreateClient(c pg.Redis) (*redis.Client, error)
}

// RegisterDriver registers a Redis driver with the given name.
// It panics if the driver name has already been registered.
func RegisterDriver(name string, driver Driver) {
	if _, ok := driverRegistry[name]; ok {
		panic("redis driver already registered: " + name)
	}
	driverRegistry[name] = driver
}

// DefaultDriver is the default implementation of the Driver interface.
type DefaultDriver struct{}

// CreateClient creates a new Redis client based on the provided configuration.
func (DefaultDriver) CreateClient(c pg.Redis) (*redis.Client, error) {
	address := fmt.Sprintf("%s:%d", c.Host, c.Port)
	return redis.NewClient(&redis.Options{
		Addr:            address,
		Username:        c.Username,
		Password:        c.Password,
		DB:              c.Database,
		DialTimeout:     time.Duration(c.ConnectTimeout) * time.Millisecond,
		ReadTimeout:     time.Duration(c.ReadTimeout) * time.Millisecond,
		WriteTimeout:    time.Duration(c.WriteTimeout) * time.Millisecond,
		ConnMaxIdleTime: time.Duration(c.IdleTimeout) * time.Millisecond,
	}), nil
}
