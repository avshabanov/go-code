package trivial

// AccountState is a state of account
type AccountState int

const (
	// Suspended is a state for inactive users
	Suspended = iota
	// Stopped is a state, synonymous to Suspended
	Stopped = Suspended
	// Active is a state for currently active users
	Active
	// RequiresVerification is a state for created accounts that require user verification
	RequiresVerification
)

// NB: comment below requires stringer package: go get golang.org/x/tools/cmd/stringer

//go:generate stringer -type=AccountState
