package config

type EthConfig struct {
	Wallet      WalletConfig `json:"wallet"`
	Contracts   Contracts    `json:"contracts"`
	Chain       Chain        `json:"chain"`
	SyncService SyncService  `json:"syncService"`
	Sly         Sly          `json:"sly"`
}

type WalletConfig struct {
	Name    string `json:"name"`
	Private string `json:"private"`
}

type Chain struct {
	RPCUrl string `json:"rpc_url"`
	ID     string `json:"id"`
}

type SyncService struct {
	SyncInterval int  `json:"syncInterval"`
	On           bool `json:"on"`
}

type Contracts struct {
	HubAddress       string `json:"hubAddress"`
	HubCreationBlock int64  `json:"hubCreationBlock"`
}

type Sly struct {
	FactoryAddress string `json:"factoryAddress"`
}
