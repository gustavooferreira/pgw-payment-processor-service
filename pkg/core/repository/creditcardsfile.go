package repository

import (
	"github.com/gustavooferreira/pgw-payment-processor-service/pkg/core"
	"gopkg.in/yaml.v2"
)

// CreditCardFileChecker holds the credit cards number and the reason to fail.
// This struct mimics a database.
type CreditCardFileChecker struct {
	CreditCards map[int64]core.CCFailReason `yaml:"creditCards"`
}

// NewCreditCardFileChecker creates a new CreditCardsHolder.
func NewCreditCardFileChecker() *CreditCardFileChecker {
	ccfc := CreditCardFileChecker{CreditCards: make(map[int64]core.CCFailReason)}
	return &ccfc
}

// Load loads data into the CreditCardFileChecker.
func (ccfc *CreditCardFileChecker) Load(data []byte) error {
	err := yaml.Unmarshal([]byte(data), &ccfc)
	return err
}

// ShouldFail checks whether the provided credit card number should fail for the provided reason.
func (ccfc *CreditCardFileChecker) ShouldFail(ccNumber int64, reason core.CCFailReason) bool {
	if v, ok := ccfc.CreditCards[ccNumber]; ok {
		if reason == v {
			return true
		}
	}
	return false
}
