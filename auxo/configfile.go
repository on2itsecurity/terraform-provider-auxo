// Description: This file contains functions for reading the (ztctl) config file

package auxo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
)

type ConfigFile struct {
	Configs []ConfigEntry `json:"configs"`
}

type ConfigEntry struct {
	Alias       string `json:"alias"`
	Description string `json:"description"`
	Token       string `json:"token"`
	APIAddress  string `json:"apiaddress"`
	Debug       bool   `json:"debug"`
}

// Will return the default config file location
func getDefaultConfigLocation() string {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	location := home + "/.ztctl/" + "config.json"

	return location
}

// GetConfigs will get all config entries from the specified file
// Specify filename for different location, leave empty for default.
// returns ConfigFile and error
func getConfigs(fileName string) (ConfigFile, error) {

	//Read full config and select the Alias
	cfgFile := ConfigFile{}
	cfgFileAsByte, err := readFile(fileName)
	if err != nil {
		return ConfigFile{}, err
	}

	err = json.Unmarshal(cfgFileAsByte, &cfgFile)
	if err != nil {
		return ConfigFile{}, err
	}

	return cfgFile, nil
}

// Read file and return byte array or error
func readFile(fileName string) ([]byte, error) {
	f, err := os.Open(fileName)

	if err != nil {
		return nil, err
	}

	defer f.Close()
	fileContent, err := ioutil.ReadAll(f)

	if err != nil {
		return nil, err
	}

	return fileContent, nil
}

// GetConfig will get the config entry with the specified alias
// Specify filename for different location, leave empty for default.
// returns ConfigEntry and error
func getConfig(filename, alias string) (ConfigEntry, error) {
	cfg, err := getConfigs(filename)

	if err != nil {
		return ConfigEntry{}, err
	}

	//Find the config entry with the alias
	for _, ce := range cfg.Configs {
		if ce.Alias == alias {
			return ce, nil
		}
	}

	return ConfigEntry{}, fmt.Errorf("Could not find config entry with alias %s", alias)
}
