package cmd

type Config struct {
	// WorkerPoolSize is the number of go routines VWAP producers
	WorkerPoolSize uint16
	// WindowsSize is the moving window size of VWAP data points i.e. 200
	WindowsSize    uint16
	// True for development level logging, false for production.
	DevLogLevel    bool
	CfgFile        string
	// SocketURL is the host URL for the matches channel
	SocketURL      string
	// Products IDs to subscribe trades "BTC-USD","ETH-USD","ETH-BTC"
	ProductIDs     []string
}
