package core

import "context"

// Repository represents a database holding credentials
type Repository interface {
	ShouldFail(ccNumber int64) (reason CCFailReason, fail bool)
}

type Authoriser interface {
	Authorise(ccNumber int64) (uid string)
	GetAssociatedCreditCard(uid string) (ccNumber int64, ok bool)
}

// ShutDowner represents anything that can be shutdown like an HTTP server.
type ShutDowner interface {
	ShutDown(ctx context.Context) error
}
