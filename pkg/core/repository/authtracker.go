package repository

import (
	"github.com/google/uuid"
)

// AuthoriserInMemoryTracker keeps track of authorisations.
// This struct mimics a database.
type AuthoriserInMemoryTracker struct {
	// Authorisations maps a UID to a credit card number
	Authorisations map[string]int64
}

// NewAuthoriserInMemoryTracker creates a new AuthoriserInMemoryTracker.
func NewAuthoriserInMemoryTracker() *AuthoriserInMemoryTracker {
	at := AuthoriserInMemoryTracker{Authorisations: make(map[string]int64)}
	return &at
}

// Authorise generates a new UID and returns it.
func (at *AuthoriserInMemoryTracker) Authorise(ccNumber int64) (uid string) {
	uid = uuid.NewString()
	at.Authorisations[uid] = ccNumber

	return uid
}

func (at *AuthoriserInMemoryTracker) GetAssociatedCreditCard(uid string) (ccNumber int64, ok bool) {
	ccNumber, ok = at.Authorisations[uid]
	return ccNumber, ok
}
