// nolint:errcheck
package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//{"type":"subscribe","product_ids":["BTC-USD","ETH-USD","ETH-BTC"],"channels":["matches"]}
// Pair i.e, BTC-USD

type Config struct {
	WorkerPoolSize uint16
	DevLogLevel    bool
	cfgFile        string
	SocketURL      string
	ProductIDs     []string // "BTC-USD","ETH-USD","ETH-BTC"
}

var flags Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vwap",
	Short: "A real time VWAP trades calculator",
	Long: `A real-time VWAP (volume-weighted average price) calculation engine. Subscribes 
to the coinbase websocket feed to stream in trade executions and update the VWAP 
for each trading pair as updates become available.`,
	Run: func(cmd *cobra.Command, args []string) {
		if flags.WorkerPoolSize == 0 {
			_, _ = fmt.Fprintln(os.Stderr, "Please supply a positive workers pool number")
			os.Exit(1)
		}
		if len(flags.ProductIDs) == 0 {
			_, _ = fmt.Fprintln(os.Stderr, "Please supply products IDs")
			os.Exit(1)
		}
		regexc, err := regexp.Compile("[A-Z]{3}-[A-Z]{3}")
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		for i, product := range flags.ProductIDs {
			product = strings.TrimSpace(product)
			match := regexc.Match([]byte(product))
			if !match {
				_, _ = fmt.Fprintf(os.Stderr, "Invalid product ID position:%d, %s %s\n", i, product, err)
				os.Exit(1)
			}
			flags.ProductIDs[i] = product
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() Config {
	cobra.CheckErr(rootCmd.Execute())

	if rootCmd.Flags().Lookup("help").Value.String() == "true" {
		os.Exit(0)
	}

	return flags
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&flags.cfgFile, "config", "c", "", "config file (default is $HOME/.vwap.yaml)")
	rootCmd.PersistentFlags().StringSliceVarP(&flags.ProductIDs, "productids", "p", []string{"BTC-USD", "ETH-USD", "ETH-BTC"}, "The comma separated trading product ID pairs to calculate the current 200 VWAP data points e.g. BTC-USD, ETH-USD, ETH-BTC")
	rootCmd.PersistentFlags().StringVarP(&flags.SocketURL, "url", "u", "wss://ws-feed.exchange.coinbase.com", "The Coinbase URL with two choices: wss://ws-feed.exchange.coinbase.com --OR-- wss://ws-feed-public.sandbox.exchange.coinbase.com")
	rootCmd.PersistentFlags().Uint16VarP(&flags.WorkerPoolSize, "workers", "w", 200, "The workers pool size for processing the ingested trades")
	rootCmd.PersistentFlags().BoolVarP(
		&flags.DevLogLevel, "devlogging", "d", false,
		`by default logging is set to production level generating structured log entries suitable for machine processing i.e. Kafka. This offers the chance to override this to development level for human friendly log output`,
	)
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("Mario Karagiorgas"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "Mario Karagiorgas salem8@gmail.com")
	viper.SetDefault("license", "MIT")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if flags.cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(flags.cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".vwap" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".vwap")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
