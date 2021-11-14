package cmd

type Config struct {
	WorkerPoolSize uint16
	DevLogLevel    bool
	CfgFile        string
	SocketURL      string
	ProductIDs     []string // "BTC-USD","ETH-USD","ETH-BTC"
}
