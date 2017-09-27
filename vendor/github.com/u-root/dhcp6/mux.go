package dhcp6

import (
	"sync"
)

// DefaultServeMux is the default ServeMux used by Serve.  When the Handle and
// HandleFunc functions are called, handlers are applied to DefaultServeMux.
var DefaultServeMux = NewServeMux()

// ServeMux is a DHCP request multiplexer, which implements Handler.  ServeMux
// matches handlers based on their MessageType, enabling different handlers
// to be used for different types of DHCP messages.  ServeMux can be helpful
// for structuring your application, but may not be needed for very simple
// DHCP servers.
type ServeMux struct {
	mu sync.RWMutex
	m  map[MessageType]Handler
}

// NewServeMux creates a new ServeMux which is ready to accept Handlers.
func NewServeMux() *ServeMux {
	return &ServeMux{
		m: make(map[MessageType]Handler),
	}
}

// ServeDHCP implements Handler for ServeMux, and serves a DHCP request using
// the appropriate handler for an input Request's MessageType.  If the
// MessageType does not match a valid Handler, ServeDHCP does not invoke any
// handlers, ignoring a client's request.
func (mux *ServeMux) ServeDHCP(w ResponseSender, r *Request) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()
	h, ok := mux.m[r.MessageType]
	if !ok {
		return
	}

	h.ServeDHCP(w, r)
}

// Handle registers a MessageType and Handler with a ServeMux, so that
// future requests with that MessageType will invoke the Handler.
func (mux *ServeMux) Handle(mt MessageType, handler Handler) {
	mux.mu.Lock()
	mux.m[mt] = handler
	mux.mu.Unlock()
}

// Handle registers a MessageType and Handler with the DefaultServeMux,
// so that future requests with that MessageType will invoke the Handler.
func Handle(mt MessageType, handler Handler) {
	DefaultServeMux.Handle(mt, handler)
}

// HandleFunc registers a MessageType and function as a HandlerFunc with a
// ServeMux, so that future requests with that MessageType will invoke the
// HandlerFunc.
func (mux *ServeMux) HandleFunc(mt MessageType, handler func(ResponseSender, *Request)) {
	mux.Handle(mt, HandlerFunc(handler))
}

// HandleFunc registers a MessageType and function as a HandlerFunc with the
// DefaultServeMux, so that future requests with that MessageType will invoke
// the HandlerFunc.
func HandleFunc(mt MessageType, handler func(ResponseSender, *Request)) {
	DefaultServeMux.HandleFunc(mt, handler)
}
