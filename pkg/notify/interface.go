package notify

import (
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/model"
)

type Service interface {
	Notify(transaction *model.Transaction) error
}
