package repository

import (
	"github.com/gustavooferreira/pgw-payment-processor-service/pkg/core"
	"gopkg.in/yaml.v2"
)

// CreditCardsHolder holds the credit cards number and the reason to fail.
// This struct mimics a database.
type CreditCardsHolder struct {
	CreditCards map[int64]core.CCFailReason `yaml:"creditCards"`
}

// NewCreditCardsHolder creates a new CreditCardsHolder.
func NewCreditCardsHolder() CreditCardsHolder {
	ch := CreditCardsHolder{CreditCards: make(map[int64]core.CCFailReason)}
	return ch
}

// Load loads data into the CreditCardsHolder.
func (cch *CreditCardsHolder) Load(data []byte) error {
	err := yaml.Unmarshal([]byte(data), &cch)
	return err
}

// Fail checks whether the provided credit card number should fail.
// If yes, returns the reason.
func (cch *CreditCardsHolder) ShouldFail(ccNumber int64, reason core.CCFailReason) bool {
	if v, ok := cch.CreditCards[ccNumber]; ok {
		if reason == v {
			return true
		}
	}
	return false
}
