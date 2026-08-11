package grok

import (
	"fmt"

	"github.com/goloop/ai"
)

// checkHosted refuses a request that asks the provider to run a capability on
// its own side.
//
// This provider does run searches of its own, but it reaches them through
// request parameters this driver has not verified against the live API, and
// two different mechanisms are plausible from its documentation. Declaring
// support on a guess would mean sending a shape nobody has seen accepted, and
// a search that silently does not happen is the failure this whole contract
// exists to prevent. A later version adds it once the wire shape is confirmed.
//
// Returning [ai.ErrNoHosted] before the request leaves is the documented
// behavior, not a placeholder: an answer produced without the search that was
// asked for looks exactly like one produced with it, so failing loudly is the
// only way a caller can tell the difference.
func checkHosted(req *ai.Request) error {
	if len(req.Hosted) == 0 {
		return nil
	}
	return fmt.Errorf("%w: this provider's search is not reachable from the chat endpoint", ai.ErrNoHosted)
}
