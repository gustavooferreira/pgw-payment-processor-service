package repository

import (
	"github.com/google/uuid"
)

// AuthTracker keeps track of authorisations.
// This struct mimics a database.
type AuthTracker struct {
	// Authorisations maps a UID to a credit card number
	Authorisations map[string]int64
}

// NewAuthTracker creates a new AuthTracker.
func NewAuthTracker() AuthTracker {
	at := AuthTracker{Authorisations: make(map[string]int64)}
	return at
}

// Authorise generates a new UID and returns it.
func (at *AuthTracker) Authorise(ccNumber int64) (uid string) {
	uid = uuid.NewString()
	at.Authorisations[uid] = ccNumber

	return uid
}

func (at *AuthTracker) GetAssociatedCreditCard(uid string) (ccNumber int64, ok bool) {
	ccNumber, ok = at.Authorisations[uid]
	return ccNumber, ok
}
