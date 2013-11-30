package upax_go

//upax_go/c_states.go

// XXX This is just copied from what is now s_states.go and is known to
// be WRONG.
//
const (
	// States through which client-server input cnx may pass
	C_HELLO_RCVD = iota

	// After the peer has sent a message containing its nodeID and a
	// digital signature, we can determine which peer we are speaking to.
	C_ID_VERIFIED

	// Once the connection has reached this state, no more messages
	// can be accepted.
	C_BYE_RCVD

	// When we reach this state, the connection must be closed.
	C_IN_CLOSED
)
