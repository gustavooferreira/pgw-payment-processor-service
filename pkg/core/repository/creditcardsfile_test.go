package repository_test

import (
	"testing"

	"github.com/gustavooferreira/pgw-payment-processor-service/pkg/core"
	"github.com/gustavooferreira/pgw-payment-processor-service/pkg/core/repository"
	"github.com/stretchr/testify/assert"
)

func TestCreditCardShouldFail(t *testing.T) {
	tests := map[string]struct {
		ccNumber       int64
		reason         core.CCFailReason
		expectedResult bool
	}{
		"authorise fail": {ccNumber: 4000000000000119, reason: core.CCFailReason_Authorise, expectedResult: true},
		"capture fail":   {ccNumber: 4000000000000259, reason: core.CCFailReason_Capture, expectedResult: true},
		"refund fail":    {ccNumber: 4000000000003238, reason: core.CCFailReason_Refund, expectedResult: true},
		"no fail":        {ccNumber: 123, expectedResult: false},
	}

	cch := createCreditCardsHolder()

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			value := cch.ShouldFail(test.ccNumber, test.reason)
			assert.Equal(t, test.expectedResult, value)
		})
	}
}

func createCreditCardsHolder() *repository.CreditCardsHolder {
	cch := repository.NewCreditCardsHolder()

	cch.CreditCards[4000000000000119] = core.CCFailReason_Authorise
	cch.CreditCards[4000000000000259] = core.CCFailReason_Capture
	cch.CreditCards[4000000000003238] = core.CCFailReason_Refund

	return &cch
}
