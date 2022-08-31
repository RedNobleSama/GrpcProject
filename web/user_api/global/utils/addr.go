/**
    @auther: oreki
    @date: 2022/4/26
    @note: 图灵老祖保佑,永无BUG
**/

package utils

import (
	"fmt"
	"net"
)

// GetAvailablePort 获取可用端口
func GetAvailablePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", "localhost"))
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}

	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

// IsPortAvailable 判断端口是否可以(未被占用)
func IsPortAvailable(port int) bool {
	addr := fmt.Sprintf("%s:%d", "0.0.0.0", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}

	defer listener.Close()
	return true
}
