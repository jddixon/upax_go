package upax_go

// xlattice_go/upax_go/c_var.go

func init() {
	// cMsgHandlers = make([][]interface{}, C_BYE_RCVD, C_MSG_HANDLER_COUNT)

	cMsgHandlers = [][]interface{}{
		// client messages permitted in C_HELLO_RCVD state
		{doCIntroMsg, doCItsMeMsg, badCCombo, badCCombo, badCCombo,
			badCCombo, badCCombo},
		// client messages permitted in C_ID_VERIFIED state
		{badCCombo, badCCombo, doCKeepAliveMsg, doCQueryMsg, doCGetMsg,
			doCPutMsg, doCByeMsg},
	}
}
