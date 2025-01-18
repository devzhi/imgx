package util

import (
	"bufio"
	"golang.org/x/term"
	"os"
	"strings"
)

// ReadPassword 读取密码，支持隐藏输入
func ReadPassword() (string, error) {
	// 首先尝试使用 term.ReadPassword
	if password, err := term.ReadPassword(int(os.Stdin.Fd())); err == nil {
		return string(password), nil
	}

	// 如果 term.ReadPassword 失败，回退到普通输入（会显示字符）
	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	// 去除末尾的换行符
	return strings.TrimSpace(password), nil
}
