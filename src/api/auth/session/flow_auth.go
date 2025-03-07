package session

const (
	AuthFlowStateNone      = 0
	AuthFlowStateCreated   = 1
	AuthFlowStateConnected = 2
	AuthFlowStateVerified  = 3

	FlowStatePending = "pending"
	FlowStateSuccess = "success"
	FlowStateFailed  = "failed"
)

type AuthFlow struct {
	slyWalletAddress string
	eoa              string
	accountId        string
	audiences        []string
	domain           string
	state            int
}

func NewAuthFlow() *AuthFlow {
	return &AuthFlow{
		state: AuthFlowStateCreated,
	}
}

func (a *AuthFlow) isCreatedState() bool {
	return a.state == AuthFlowStateCreated
}
func (a *AuthFlow) isConnectedState() bool {
	return a.state == AuthFlowStateConnected
}
func (a *AuthFlow) isVerified() bool {
	return a.state == AuthFlowStateVerified
}
func (a *AuthFlow) SessionState() string {

	switch a.state {
	case AuthFlowStateNone:
		return FlowStatePending
	case AuthFlowStateCreated:
		return FlowStatePending
	case AuthFlowStateConnected:
		return FlowStatePending
	case AuthFlowStateVerified:
		return FlowStateSuccess
	default:
		return FlowStateFailed
	}
}

func (a *AuthFlow) setPayload(payload *PayloadAccountsResponse) {
	a.eoa = payload.EOA
	a.slyWalletAddress = payload.SLYWalletAddress
	a.state = AuthFlowStateConnected
}

func (a *AuthFlow) setVerified(accountId string) {
	a.state = AuthFlowStateVerified
	a.accountId = accountId
}
