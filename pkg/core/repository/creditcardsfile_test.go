package repository_test

import (
	"testing"

	"github.com/gustavooferreira/pgw-payment-processor-service/pkg/core"
	"github.com/gustavooferreira/pgw-payment-processor-service/pkg/core/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreditCardShouldFail(t *testing.T) {
	tests := map[string]struct {
		ccNumber       int64
		expectedFail   bool
		expectedReason core.CCFailReason
	}{
		"authorise fail": {ccNumber: 4000000000000119, expectedFail: true, expectedReason: core.CCFailReason_Authorise},
		"capture fail":   {ccNumber: 4000000000000259, expectedFail: true, expectedReason: core.CCFailReason_Capture},
		"refund fail":    {ccNumber: 4000000000003238, expectedFail: true, expectedReason: core.CCFailReason_Refund},
		"no fail":        {ccNumber: 123, expectedFail: false},
	}

	cch := createCreditCardsHolder()

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			reason, fail := cch.ShouldFail(test.ccNumber)
			require.Equal(t, test.expectedFail, fail)

			if test.expectedFail {
				assert.Equal(t, test.expectedReason, reason)
			}
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
