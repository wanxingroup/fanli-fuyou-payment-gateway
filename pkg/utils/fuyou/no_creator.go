package fuyou

import (
	"fmt"
	"math/rand"
	"strings"
	"sync/atomic"
	"time"
)

var intNum uint32 = 0

// 富友订单号生成规则
// 1419+yyyyMMddHHmmss+三位随机数+五位自增序列
// 共26位字符串长度
func GeneralFuYouOrderId() string {
	var n, v uint32
	for {
		v = atomic.LoadUint32(&intNum)
		n = v + 1
		if atomic.CompareAndSwapUint32(&intNum, v, n) {
			break
		}
	}
	s1 := fmt.Sprintf("%d", rand.Intn(1000))
	s2 := fmt.Sprintf("%d", n%100000)
	return "1419" + time.Now().Format("20060102150405") + strings.Repeat("0", 3-len(s1)) + s1 + strings.Repeat("0", 5-len(s2)) + s2
}
