// Package trace 提供进程内链路 ID 生成。
package trace

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// seq 进程内递增计数，保证同一毫秒内生成的 ID 不重复。
var seq uint64

// ipHexOnce 缓存本机 IP 的十六进制值，避免每个请求都枚举网卡。
var (
	ipHexOnce sync.Once
	ipHexVal  string
)

// GenerateTraceID 生成 traceID：IP 十六进制 + 进程号 + 毫秒时间戳 + 递增计数。
// 例：c0a80101-12345-1760000000000-7。
func GenerateTraceID() string {
	return fmt.Sprintf("%s-%d-%d-%d",
		localIPHex(),
		os.Getpid(),
		time.Now().UnixMilli(),
		atomic.AddUint64(&seq, 1),
	)
}

// localIPHex 返回本机 IP 的十六进制表示；取不到时返回 00000000。
func localIPHex() string {
	ipHexOnce.Do(func() {
		ipHexVal = computeLocalIPHex()
	})
	return ipHexVal
}

func computeLocalIPHex() string {
	ip := localIP()
	if ip == nil {
		return "00000000"
	}
	if v4 := ip.To4(); v4 != nil {
		return fmt.Sprintf("%02x%02x%02x%02x", v4[0], v4[1], v4[2], v4[3])
	}
	return strings.ReplaceAll(ip.String(), ":", "")
}

// localIP 优先取第一块非回环的网卡地址，取不到时用 UDP 拨号获取出口地址。
func localIP() net.IP {
	ifaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range ifaces {
			if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
				continue
			}
			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				ipNet, ok := addr.(*net.IPNet)
				if !ok || ipNet.IP.IsLoopback() {
					continue
				}
				if v4 := ipNet.IP.To4(); v4 != nil {
					return v4
				}
			}
		}
	}
	// UDP 拨号不会真正发包，仅用于让内核选择出口地址。
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil
	}
	defer conn.Close()
	if addr, ok := conn.LocalAddr().(*net.UDPAddr); ok {
		return addr.IP
	}
	return nil
}
