package upax_go

// upax_go/s_var.go

func init() {
	// sMsgHandlers = make([][]interface{}, S_BYE_RCVD, S_MSG_HANDLER_COUNT)

	sMsgHandlers = [][]interface{}{
		// client messages permitted in S_HELLO_RCVD state
		{doSItsMeMsg, badSCombo, badSCombo, badSCombo, badSCombo, badSCombo},
		// client messages permitted in S_ID_VERIFIED state
		{badSCombo, doSKeepAliveMsg, doSGetMsg, doSIHaveMsg, doSPutMsg, doSByeMsg},
	}
}
