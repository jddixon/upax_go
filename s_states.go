package upax_go

// upax_go/s_states.go

// For server-server connections, there is little state to track.
const (
	// States through which a intra-cluster input cnx may pass
	S_HELLO_RCVD = iota

	// After the peer has sent a message containing its nodeID and a
	// digital signature, we can determine which peer we are speaking to.
	S_ID_VERIFIED

	// Once the connection has reached this state, no more messages
	// can be accepted.
	S_BYE_RCVD

	// When we reach this state, the connection must be closed.
	S_IN_CLOSED
)
