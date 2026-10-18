package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"cosmossdk.io/log"
	"github.com/cliffordx/opencasedigests/chain/app"
	cfg "github.com/cometbft/cometbft/config"
	cmtlog "github.com/cometbft/cometbft/libs/log"
	cmtos "github.com/cometbft/cometbft/libs/os"
	"github.com/cometbft/cometbft/node"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/privval"
	"github.com/cometbft/cometbft/proxy"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	homeDir string
	chainID string

	// RootCmd is the primary CLI command for casedigestord.
	RootCmd = &cobra.Command{
		Use:   "casedigestord",
		Short: "Case Digestor Protocol daemon and validator node",
		Long: `casedigestord is the node daemon for Case Digestor Protocol—a decentralized
legal provenance and peer-attested jurisprudence network for Philippine law.`,
	}
)

func init() {
	defaultHome := filepath.Join(os.Getenv("HOME"), ".casedigestord")
	RootCmd.PersistentFlags().StringVar(&homeDir, "home", defaultHome, "node directory for config and data")
	RootCmd.PersistentFlags().StringVar(&chainID, "chain-id", app.DefaultChainID, "blockchain network identifier")

	RootCmd.AddCommand(initCmd)
	RootCmd.AddCommand(startCmd)
	RootCmd.AddCommand(testnetCmd)
	RootCmd.AddCommand(versionCmd)
}

// Execute runs the root command.
func Execute() error {
	return RootCmd.Execute()
}

// versionCmd prints the version of casedigestord.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the casedigestord binary version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("casedigestord version %s (protocol v1, Cosmos SDK v0.50, CometBFT v0.38)\n", app.AppVersion)
	},
}

// initCmd initializes node configuration and genesis files.
var initCmd = &cobra.Command{
	Use:   "init [moniker]",
	Short: "Initialize private validator, p2p, and genesis files",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		moniker := args[0]
		config := cfg.DefaultConfig()
		config.SetRoot(homeDir)
		config.Moniker = moniker

		cfg.EnsureRoot(config.RootDir)

		// Create private validator key if not present
		privValKeyFile := config.PrivValidatorKeyFile()
		privValStateFile := config.PrivValidatorStateFile()
		var pv *privval.FilePV
		if cmtos.FileExists(privValKeyFile) {
			pv = privval.LoadFilePV(privValKeyFile, privValStateFile)
		} else {
			pv = privval.GenFilePV(privValKeyFile, privValStateFile)
			pv.Save()
		}

		// Create node key for p2p identity
		nodeKeyFile := config.NodeKeyFile()
		if !cmtos.FileExists(nodeKeyFile) {
			if _, err := p2p.LoadOrGenNodeKey(nodeKeyFile); err != nil {
				return err
			}
		}

		// Generate initial genesis document
		pubKey, err := pv.GetPubKey()
		if err != nil {
			return err
		}

		genDoc := cmttypes.GenesisDoc{
			ChainID:     chainID,
			GenesisTime: time.Now().UTC(),
			Validators: []cmttypes.GenesisValidator{
				{
					Address: pubKey.Address(),
					PubKey:  pubKey,
					Power:   10,
					Name:    moniker,
				},
			},
		}

		genFile := config.GenesisFile()
		if !cmtos.FileExists(genFile) {
			if err := genDoc.SaveAs(genFile); err != nil {
				return err
			}
		}

		cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)

		fmt.Printf("Initialized %s node at %s (chain-id: %s)\n", moniker, homeDir, chainID)
		return nil
	},
}

// startCmd boots the CometBFT consensus engine and ABCI state machine.
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Run the full node and start consensus replication",
	RunE: func(cmd *cobra.Command, args []string) error {
		config := cfg.DefaultConfig()
		config.SetRoot(homeDir)

		v := viper.New()
		v.SetConfigName("config")
		v.SetConfigType("toml")
		v.AddConfigPath(filepath.Join(homeDir, "config"))
		if err := v.ReadInConfig(); err == nil {
			_ = v.Unmarshal(config)
			config.SetRoot(homeDir)
		}

		appLogger := log.NewLogger(os.Stdout)
		casedApp := app.NewCaseDigestorApp(appLogger)

		pv := privval.LoadFilePV(config.PrivValidatorKeyFile(), config.PrivValidatorStateFile())
		nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
		if err != nil {
			return fmt.Errorf("failed to load node key: %w", err)
		}

		clientCreator := proxy.NewLocalClientCreator(casedApp)
		cmtLogger := cmtlog.NewTMLogger(os.Stdout)

		tmNode, err := node.NewNode(
			config,
			pv,
			nodeKey,
			clientCreator,
			node.DefaultGenesisDocProviderFunc(config),
			cfg.DefaultDBProvider,
			node.DefaultMetricsProvider(config.Instrumentation),
			cmtLogger,
		)
		if err != nil {
			return fmt.Errorf("failed to create node: %w", err)
		}

		if err := tmNode.Start(); err != nil {
			return fmt.Errorf("failed to start node: %w", err)
		}

		appLogger.Info("casedigestord node started successfully", "moniker", config.Moniker, "rpc", config.RPC.ListenAddress)

		// Wait for signal
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		appLogger.Info("shutting down casedigestord...")
		return tmNode.Stop()
	},
}

// testnetCmd bootstraps a multi-node testnet cluster for the institutional consortium.
var testnetCmd = &cobra.Command{
	Use:   "testnet",
	Short: "Generate configurations for a 3-node institutional validator consortium",
	RunE: func(cmd *cobra.Command, args []string) error {
		outputDir, _ := cmd.Flags().GetString("output-dir")
		if outputDir == "" {
			outputDir = "./deploy/testnet/data"
		}

		validators := []struct {
			Moniker string
			P2PPort int
			RPCPort int
		}{
			{Moniker: "UPLaw-Validator", P2PPort: 26656, RPCPort: 26657},
			{Moniker: "AteneoLaw-Validator", P2PPort: 26666, RPCPort: 26667},
			{Moniker: "SanBedaLaw-Validator", P2PPort: 26676, RPCPort: 26677},
		}

		type nodeSetup struct {
			Moniker string
			NodeKey *p2p.NodeKey
			PV      *privval.FilePV
			Dir     string
			P2PPort int
			RPCPort int
		}

		setups := make([]nodeSetup, len(validators))
		genVals := make([]cmttypes.GenesisValidator, len(validators))

		for i, v := range validators {
			nodeDir := filepath.Join(outputDir, fmt.Sprintf("node%d", i))
			configDir := filepath.Join(nodeDir, "config")
			dataDir := filepath.Join(nodeDir, "data")
			os.MkdirAll(configDir, 0o755)
			os.MkdirAll(dataDir, 0o755)

			pvFile := filepath.Join(configDir, "priv_validator_key.json")
			stateFile := filepath.Join(dataDir, "priv_validator_state.json")
			pv := privval.GenFilePV(pvFile, stateFile)
			pv.Save()

			nkFile := filepath.Join(configDir, "node_key.json")
			nk, err := p2p.LoadOrGenNodeKey(nkFile)
			if err != nil {
				return err
			}

			pubKey, err := pv.GetPubKey()
			if err != nil {
				return err
			}

			genVals[i] = cmttypes.GenesisValidator{
				Address: pubKey.Address(),
				PubKey:  pubKey,
				Power:   10,
				Name:    v.Moniker,
			}

			setups[i] = nodeSetup{
				Moniker: v.Moniker,
				NodeKey: nk,
				PV:      pv,
				Dir:     nodeDir,
				P2PPort: v.P2PPort,
				RPCPort: v.RPCPort,
			}
		}

		// Create unified genesis
		genDoc := cmttypes.GenesisDoc{
			ChainID:     chainID,
			GenesisTime: time.Now().UTC(),
			Validators:  genVals,
			AppState:    json.RawMessage(`{"casedigestor":{"status":"initialized"}}`),
		}

		isDocker, _ := cmd.Flags().GetBool("docker")

		// Write configs for each node with persistent peers
		for i, s := range setups {
			config := cfg.DefaultConfig()
			config.SetRoot(s.Dir)
			config.Moniker = s.Moniker
			config.P2P.AddrBookStrict = false
			config.P2P.AllowDuplicateIP = true
			config.RPC.CORSAllowedOrigins = []string{"*"}
			config.Consensus.TimeoutCommit = 1 * time.Second

			var peers []string
			if isDocker {
				config.P2P.ListenAddress = "tcp://0.0.0.0:26656"
				config.RPC.ListenAddress = "tcp://0.0.0.0:26657"
				for j, other := range setups {
					if i != j {
						peerAddr := fmt.Sprintf("%s@validator%d:26656", other.NodeKey.ID(), j)
						peers = append(peers, peerAddr)
					}
				}
			} else {
				config.P2P.ListenAddress = fmt.Sprintf("tcp://0.0.0.0:%d", s.P2PPort)
				config.RPC.ListenAddress = fmt.Sprintf("tcp://0.0.0.0:%d", s.RPCPort)
				for j, other := range setups {
					if i != j {
						peerAddr := fmt.Sprintf("%s@127.0.0.1:%d", other.NodeKey.ID(), other.P2PPort)
						peers = append(peers, peerAddr)
					}
				}
			}
			config.P2P.PersistentPeers = joinStrings(peers, ",")

			cfg.WriteConfigFile(filepath.Join(s.Dir, "config", "config.toml"), config)
			genDoc.SaveAs(filepath.Join(s.Dir, "config", "genesis.json"))
		}

		fmt.Printf("Generated 3-node consortium testnet cluster at %s (chain-id: %s, docker: %v)\n", outputDir, chainID, isDocker)
		for i, v := range validators {
			fmt.Printf("  Validator %d: %s (P2P: %d, RPC: %d)\n", i, v.Moniker, v.P2PPort, v.RPCPort)
		}
		return nil
	},
}

func init() {
	testnetCmd.Flags().String("output-dir", "./deploy/testnet/data", "output directory for testnet node configurations")
	testnetCmd.Flags().Bool("docker", true, "configure node peering for Docker network topology")
}

func joinStrings(elems []string, sep string) string {
	res := ""
	for i, s := range elems {
		if i > 0 {
			res += sep
		}
		res += s
	}
	return res
}
