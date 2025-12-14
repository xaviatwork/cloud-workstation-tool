package cloudworkstation

import (
	"encoding/json"
	"os"
	"path"
)

type Config struct {
	Path string
}

func useDefault() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = ""
	}
	return &Config{
		Path: path.Join(homeDir, ".config", "cloud-workstation", "cw.cfg"),
	}
}

func NewCloudWorkstation() *CloudWorkstation {
	return &CloudWorkstation{
		TunnelPorts: TunnelPorts{
			Remote: "22",
		}}
}

func Load(cfg *Config) (*CloudWorkstation, error) {
	if cfg == nil {
		cfg = useDefault()
	}
	cfgBytes, err := os.ReadFile(cfg.Path)
	if err != nil {
		return &CloudWorkstation{}, err
	}
	cw := NewCloudWorkstation()
	if err := json.Unmarshal(cfgBytes, cw); err != nil {
		return &CloudWorkstation{}, err
	}
	return cw, nil
}
