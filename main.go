package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/falun/markup/app"
)

func main() {

	var (
		rootDir             = mustGetPwd()
		port                = 8080
		host                = "localhost"
		index               = "README.md"
		noUseExistingConfig = false
		noWriteConfig       = false
		configDir           = path.Join(rootDir, ".markup")
		assetDir            = path.Join(configDir, "assets")
		configFile          = "conf.json"
		dumpConfig          = false
	)

	// config around how to run markup
	flag.StringVar(&rootDir, "root", rootDir, "the root serving directory")
	flag.StringVar(&host, "ip", host, "the interface we should be listening on")
	flag.IntVar(&port, "port", port, "the port markup will listen on")
	flag.StringVar(
		&index, "default-file", index,
		"the file returned if '/' is requested; resolved relative to root")

	// config around how to handle config
	flag.StringVar(
		&configFile, "config-file", configFile,
		"determines the file name to load config from and write config to")
	flag.BoolVar(
		&noUseExistingConfig, "no-use-config", noUseExistingConfig,
		"If set markup will not use an existing config from disk relying only on command line flags")
	flag.BoolVar(
		&noWriteConfig, "no-write-config", noWriteConfig,
		"If set this will cause the current config to not be saved to disk")
	flag.StringVar(&configDir, "config-dir", configDir,
		"where the config should be written")
	flag.BoolVar(
		&dumpConfig, "dc", dumpConfig,
		"If set this will print the config values that will be used and exit. A new config will not be written to disk.")
	flag.StringVar(
		&assetDir, "asset-dir", assetDir,
		"If set markup will attempt to load templates and other app assets from this directory.")

	flag.Parse()

	flagSet := map[string]bool{}
	flag.Visit(func(f *flag.Flag) { flagSet[f.Name] = true })

	if !mustIsDirOrDirSymlink(rootDir) {
		log.Fatalf("rootDir '%v' must be a directory")
	}

	if !noWriteConfig {
		if err := ensureConfigDir(configDir); err != nil {
			log.Fatal("could not find or create config dir %v: %v", configDir, err)
		}
	}
	configFilePath := getConfigPath(flag.Args(), configDir, configFile)

	cfg := &app.Config{
		Port:        port,
		BrowseToken: "browse",
		IndexToken:  "index",
		AssetToken:  "asset",
	}

	cfgFile, foundExisting := app.ConfigFromFile(configFilePath)
	if !noUseExistingConfig {
		cfg.Merge(cfgFile)
	}

	useDefault := !foundExisting || noUseExistingConfig

	flagged := app.Config{}
	if useDefault || flagSet["root"] {
		flagged.RootDir = rootDir
	}
	if useDefault || flagSet["ip"] {
		flagged.Host = host
	}
	if useDefault || flagSet["port"] {
		flagged.Port = port
	}
	if useDefault || flagSet["default-file"] {
		flagged.Index = index
	}
	if useDefault || flagSet["asset-dir"] {
		flagged.AssetDir = assetDir
	}
	cfg.Merge(flagged)

	if !dumpConfig && !noWriteConfig {
		cfg.SaveToFile(configFilePath)
	}

	if dumpConfig {
		if noUseExistingConfig {
			configFilePath = ""
		}
		fmt.Printf("Would use config (%s):", configFilePath)
		fmt.Printf("%s\n", cfg.JSON())
		return
	}
	app.Main(*cfg)
}

func mustGetPwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not get working directory: %v", err)
	}
	return dir
}

func mustIsDirOrDirSymlink(path string) bool {
	fp, err := os.Open(path)
	if err != nil {
		log.Fatalf("could not open %v: %v", path, err)
	}

	i, err := fp.Stat()
	if err != nil {
		log.Fatalf("could not Stat %v: %v", path, err)
	}

	if i.IsDir() {
		return true
	}

	if 0 == (os.ModeSymlink & i.Mode()) {
		return false
	}

	newPath, err := os.Readlink(path)
	if err != nil {
		log.Fatalf("could not resolve symlink %v: %v", path, err)
	}

	fp, err = os.Open(newPath)
	if err != nil {
		log.Fatalf("could not open resolved symlink %v: %v", newPath, err)
	}

	i, err = fp.Stat()
	if err != nil {
		log.Fatalf("could not Stat resolved symlink %v: %v", newPath, err)
	}

	return i.IsDir()
}

func getConfigPath(args []string, configDir, configFile string) string {
	if len(args) != 0 {
		return args[0]
	}

	if !mustIsDirOrDirSymlink(configDir) {
		log.Fatalf("configDir '%v' must be a directory")
	}

	return path.Join(configDir, configFile)
}

func ensureConfigDir(d string) error {
	err := os.MkdirAll(d, 0744)
	if err != nil && os.IsExist(err) {
		return nil
	}
	return err
}
