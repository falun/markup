package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/falun/markup/app"
	"github.com/falun/markup/bundle"
)

func main() {

	var (
		rootDir             = mustGetPwd()
		port                = 8080
		host                = "localhost"
		index               = "README.md"
		excludeDirs         = ""
		noUseExistingConfig = false
		noWriteConfig       = false
		configDir           = filepath.Join(rootDir, ".markup")
		assetDir            = filepath.Join(configDir, "assets")
		configFile          = "conf.json"
		scanFreq            = 0 * time.Second
		dumpConfig          = false
		forceInstallAssets  = false
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
	flag.DurationVar(
		&scanFreq, "scan-freq", scanFreq,
		"If > 0s this will control how long markup waits between each indexing pass over root directory")
	flag.StringVar(
		&excludeDirs, "exclude-dirs", "",
		"Comma-separated set of directories that will be excluded when building an index")
	flag.BoolVar(
		&forceInstallAssets, "force-install-assets", false,
		"Installs default assets into the asset directory. Will overwrite existing content if the directory in not empty.")

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
	if useDefault || flagSet["scan-freq"] {
		flagged.RootScanFreq = &scanFreq
	}
	if useDefault || flagSet["exclude-dirs"] {
		flagged.ExcludeDirs = strings.Split(excludeDirs, ",")
	}

	cfg.Merge(flagged)

	if !dumpConfig && !noWriteConfig {
		cfg.SaveToFile(configFilePath)
	}

	shouldInstallAssets := !exists(cfg.AssetDir)

	installAssets := shouldInstallAssets || forceInstallAssets
	if dumpConfig {
		if noUseExistingConfig {
			configFilePath = ""
		}
		if installAssets {
			fmt.Printf("Wolud install assets to %v\n", cfg.AssetDir)
		}
		fmt.Printf("Would use mime map: %v\n", filepath.Join(cfg.AssetDir, "mime.json"))
		fmt.Printf("Would use config (%s): ", configFilePath)
		fmt.Printf("%s\n", cfg.JSON())
		return
	}

	if installAssets {
		bundle.Install(cfg.AssetDir)
	}

	app.Main(*cfg, getMimeMap(cfg.AssetDir))
}

func getMimeMap(cfgDir string) map[string]string {
	p := filepath.Join(cfgDir, "mime.json")
	b, err := ioutil.ReadFile(p)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		log.Fatalf("Error attempting to read MIME config (%v): %v\n", p, err.Error())
	}

	// config stores mime -> [extenion]
	config := map[string][]string{}
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatalf("Error parsing MIME config (%v): %v", p, err.Error())
	}

	// returns in format extension -> MIME type
	result := map[string]string{}
	for mt, exts := range config {
		for _, e := range exts {
			// ensure extension also includes the period
			if e[0] != '.' {
				e = "." + e
			}
			result[e] = mt
		}
	}

	return result
}

func mustGetPwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not get working directory: %v", err)
	}
	return dir
}

func exists(path string) bool {
	fp, err := os.Open(path)
	defer fp.Close()
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		log.Fatalf("failed check for existence of %v: %v", path, err.Error())
	}

	return true
}

func isDirOrDirSymlink(path string) (bool, error) {
	fp, err := os.Open(path)
	if err != nil {
		return false, fmt.Errorf("could not open %v: %v", path, err)
	}

	i, err := fp.Stat()
	if err != nil {
		return false, fmt.Errorf("could not Stat %v: %v", path, err)
	}

	if i.IsDir() {
		return true, nil
	}

	if 0 == (os.ModeSymlink & i.Mode()) {
		return false, nil
	}

	newPath, err := os.Readlink(path)
	if err != nil {
		return false, fmt.Errorf("could not resolve symlink %v: %v", path, err)
	}

	fp, err = os.Open(newPath)
	if err != nil {
		return false, fmt.Errorf("could not open resolved symlink %v: %v", newPath, err)
	}

	i, err = fp.Stat()
	if err != nil {
		return false, fmt.Errorf("could not Stat resolved symlink %v: %v", newPath, err)
	}

	return i.IsDir(), nil
}

func mustIsDirOrDirSymlink(path string) bool {
	isDir, err := isDirOrDirSymlink(path)
	if err != nil {
		log.Fatal(err)
	}

	return isDir
}

func getConfigPath(args []string, configDir, configFile string) string {
	if len(args) != 0 {
		return args[0]
	}

	if !mustIsDirOrDirSymlink(configDir) {
		log.Fatalf("configDir '%v' must be a directory")
	}

	return filepath.Join(configDir, configFile)
}

func ensureConfigDir(d string) error {
	err := os.MkdirAll(d, 0744)
	if err != nil && os.IsExist(err) {
		return nil
	}
	return err
}
