package repository_test

import (
	"testing"

	"github.com/gustavooferreira/pgw-payment-processor-service/pkg/core/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthTrackerAuthorise(t *testing.T) {
	at := repository.NewAuthTracker()

	var ccNumber int64 = 4000000000000119

	uid := at.Authorise(ccNumber)

	number, ok := at.GetAssociatedCreditCard(uid)
	require.Equal(t, true, ok)
	assert.Equal(t, ccNumber, number)
}
