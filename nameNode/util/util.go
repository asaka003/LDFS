package util

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

//生成指定长度的数字验证码
func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func IsMD5(s string) bool {
	// 匹配MD5值的正则表达式
	r := regexp.MustCompile(`^[a-f0-9]{32}$`)

	// 返回字符串是否符合正则表达式
	return r.MatchString(s)
}
