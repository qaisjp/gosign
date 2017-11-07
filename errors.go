package gosign

import "errors"

// ErrLoggedOut is the error for all errors related to being logged out.
var ErrLoggedOut = errors.New("User is already logged out")
