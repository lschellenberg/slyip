package info

type ChainInfoDTO struct {
	WalletAddress string `json:"walletAddress"`
	RPCUrl        string `json:"rpcURL"`
	ChainId       string `json:"chainId"`
	WalletName    string `json:"walletName"`
	WalletValue   string `json:"walletValue"`
}

type InvitationCodesDTO struct {
	Code      string `json:"code"`
	Status    string `json:"status"`
	ExpiresAt int64  `json:"expiresAt"`
}
