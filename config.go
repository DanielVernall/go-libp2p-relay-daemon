package relaydaemon

import (
	"encoding/json"
	"os"
	"time"

	relayv2 "github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
)

// Config stores the full configuration of the relays, ACLs and other settings
// that influence behaviour of a relay daemon.
type Config struct {
	Network    NetworkConfig
	ConnMgr    ConnMgrConfig
	Relay      RelayConfig
	NATService NATServiceConfig
	ACL        ACLConfig
	Daemon     DaemonConfig
}

// DaemonConfig controls settings for the relay-daemon itself.
type DaemonConfig struct {
	PprofPort int
}

// NetworkConfig controls listen and annouce settings for the libp2p host.
type NetworkConfig struct {
	ListenAddrs   []string
	AnnounceAddrs []string
}

// ConnMgrConfig controls the libp2p connection manager settings.
type ConnMgrConfig struct {
	ConnMgrLo    int
	ConnMgrHi    int
	ConnMgrGrace time.Duration
}

// RelayConfig controls activation of V2 circuits and resource configuration
// for them.
type RelayConfig struct {
	Enabled   bool
	Resources relayv2.Resources
}

// NATServiceConfig controls activation of the AutoNAT service
type NATServiceConfig struct {
	Enabled bool
}

// ACLConfig provides filtering configuration to allow specific peers or
// subnets to be fronted by relays. This specifies the peers/subnets
// that are able to make reservations on the relay.
type ACLConfig struct {
	AllowPeers   []string
	AllowSubnets []string
}

// DefaultConfig returns a default relay configuration using default resource
// settings and no ACLs.
func DefaultConfig() Config {
	return Config{
		Network: NetworkConfig{
			ListenAddrs: []string{
				"/ip4/0.0.0.0/udp/4001/quic",
				"/ip6/::/udp/4001/quic",
				"/ip4/0.0.0.0/tcp/4001",
				"/ip6/::/tcp/4001",
			},
		},
		ConnMgr: ConnMgrConfig{
			ConnMgrLo:    512,
			ConnMgrHi:    768,
			ConnMgrGrace: 2 * time.Minute,
		},
		Relay: RelayConfig{
			Enabled:   true,
			Resources: relayv2.DefaultResources(),
		},
		NATService: NATServiceConfig{
			Enabled: true,
		},
		Daemon: DaemonConfig{
			PprofPort: 6060,
		},
	}
}

// LoadConfig reads a relay daemon JSON configuration from the given path.
// The configuration is first initialized with DefaultConfig, so all unset
// fields will take defaults from there.
func LoadConfig(cfgPath string) (Config, error) {
	cfg := DefaultConfig()

	if cfgPath != "" {
		cfgFile, err := os.Open(cfgPath)
		if err != nil {
			return Config{}, err
		}
		defer cfgFile.Close()

		decoder := json.NewDecoder(cfgFile)
		err = decoder.Decode(&cfg)
		if err != nil {
			return Config{}, err
		}
	}

	return cfg, nil
}
