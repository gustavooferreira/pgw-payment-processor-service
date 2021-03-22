package repository_test

import (
	"testing"

	"github.com/gustavooferreira/pgw-payment-processor-service/pkg/core/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthorisation(t *testing.T) {
	auth := repository.NewAuthoriserInMemoryTracker()

	var ccNumber int64 = 4000000000000119

	uid := auth.Authorise(ccNumber)

	number, ok := auth.GetAssociatedCreditCard(uid)
	require.Equal(t, true, ok)
	assert.Equal(t, ccNumber, number)
}
