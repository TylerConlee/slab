package plugins

import (
	"bytes"
	"os"

	"github.com/BurntSushi/toml"
)

// LoadPlugins grabs the command line argument for where the configuration file
// is located and loads that into memory.
func LoadPlugins() (config Plugins) {
	if _, err := os.Stat("plugins.toml"); err == nil {
		if _, err := toml.DecodeFile("plugins.toml", &config); err != nil {
			log.Error("Configuration file not found.", map[string]interface{}{
				"module": "plugins",
				"error":  err,
			})
			return
		}
		log.Info("Configuration loaded successfully", map[string]interface{}{
			"module": "plugins",
			"file":   "plugins.toml",
		})
		return
	}
	return
}

// SavePlugins takes a config and saves it to the local file, config.toml.
func SavePlugins(config Plugins) bool {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(config); err != nil {
		log.Error("Error creating new buffer for config", map[string]interface{}{
			"module": "plugins",
			"config": config,
			"error":  err,
		})
		return false
	}

	f, err := os.Create("plugins.toml")
	if nil != err {
		log.Error("error saving file", map[string]interface{}{
			"module": "plugins",
			"error":  err,
		})
	}
	defer f.Close()
	n, err := f.WriteString(buf.String())
	if nil != err {
		log.Error("error saving file", map[string]interface{}{
			"module": "plugins",
		})
	}
	f.Sync()
	log.Debug("Saved configuration file", map[string]interface{}{
		"module": "plugins",
		"output": n,
	})
	return true

}
