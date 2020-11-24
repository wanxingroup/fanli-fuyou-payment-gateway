package transaction_test

import (
	"os"
	"testing"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/utils/test"
)

func TestMain(m *testing.M) {

	test.Init()

	code := m.Run()

	test.Release()

	os.Exit(code)
}
