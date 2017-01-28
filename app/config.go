package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Config struct {
	Host        string   `json:"host"`
	Port        int      `json:"port"`
	RootDir     string   `json:"root_dir"`
	AssetDir    string   `json:"asset_dir"`
	Index       string   `json:"index_file"`
	BrowseToken string   `json:"browse_token"`
	IndexToken  string   `json:"index_token"`
	AssetToken  string   `json:"asset_token"`
	ExcludeDirs []string `json:"exclude_dirs"`
}

// Merge will use the Config as a set of defaults and
func (c *Config) Merge(other Config) *Config {
	ck := func(v string) bool {
		return strings.TrimSpace(v) != ""
	}

	if ck(other.Host) {
		c.Host = other.Host
	}

	if other.Port != 0 {
		c.Port = other.Port
	}

	if ck(other.RootDir) {
		c.RootDir = other.RootDir
	}

	if ck(other.AssetDir) {
		c.AssetDir = other.AssetDir
	}

	if ck(other.Index) {
		c.Index = other.Index
	}

	if ck(other.BrowseToken) {
		c.BrowseToken = other.BrowseToken
	}

	if ck(other.IndexToken) {
		c.IndexToken = other.IndexToken
	}

	if ck(other.AssetToken) {
		c.AssetToken = other.AssetToken
	}

	if len(other.ExcludeDirs) != 0 {
		c.ExcludeDirs = other.ExcludeDirs
	}

	return c
}

// ConfigFromFile reads a config object from the specified file path and a bool
// indicating if a file was found.
func ConfigFromFile(path string) (Config, bool) {
	b, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		return Config{}, false
	}

	if err != nil {
		log.Fatalf("could not read config from %v: %v", path, err)
	}

	c := Config{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		log.Fatalf("colud not read config from %v: %v", path, err)
	}

	return c, true
}

func (c Config) SaveToFile(path string) {
	err := ioutil.WriteFile(path, []byte(c.JSON()), 0644)
	if err != nil {
		log.Fatalf("could not save config to %v: %v", path, err)
	}
}

func (c Config) JSON() string {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Fatalf("could not encode config object for save: %v", err)
	}
	return string(b)
}
