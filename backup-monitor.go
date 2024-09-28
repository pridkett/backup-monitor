package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/naoina/toml"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Server struct {
		Host string
		Port int
	}
	Directories []DirectoryConfig
}

type Condition string

const (
	ConditionAny Condition = "any"
	ConditionAll Condition = "all"
)

type DirectoryConfig struct {
	Id        string
	Globs     []FileGlobConfig
	Condition Condition
}

type FileGlobConfig struct {
	Path string
	Glob string
	Time string
}

var config Config

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func main() {

	var configFile = flag.String("config", "config.toml", "Path to the configuration file")
	flag.Parse()

	// os.Args provides access to raw command-line arguments
	args := os.Args

	// Print the arguments
	log.Info("Command-line arguments:")
	for i, arg := range args {
		log.WithFields(log.Fields{
			"arg_num": i,
			"arg":     arg,
		}).Info("Argument")
	}

	if *configFile != "" {
		var f, err = os.Open(*configFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		if err := toml.NewDecoder(f).Decode(&config); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("No configuration file specified")
	}

	startWebService(config.Server.Host, config.Server.Port)
}

func getDirectoryConfig(id string) (DirectoryConfig, error) {
	for _, directoryConfig := range config.Directories {
		if directoryConfig.Id == id {
			return directoryConfig, nil
		}
	}
	return DirectoryConfig{}, fmt.Errorf("Directory config not found for id: %s", id)
}

func startWebService(host string, port int) {
	r := gin.Default()
	address := fmt.Sprintf("%s:%d", host, port)
	log.Infof("Starting web service on %s", address)

	r.GET("/config", getConfig)
	r.GET("/getUpdated/:configSet", getUpdated)

	err := r.Run(address)
	if err != nil {
		log.Fatalf("Failed to start web service: %v", err)
	}

}

func convertTimeDelta(timestr string) (time.Time, error) {
	// delta parses the given duration string (timestr) into a time.Duration object.
	// If the parsing fails, it returns an error.

	delta, err := time.ParseDuration(timestr)
	if err != nil {
		return time.Time{}, err
	}

	// get the current time and subtract out the delta
	now := time.Now()
	return now.Add(-delta), nil
}

func getUpdated(c *gin.Context) {
	configSet := c.Param("configSet")
	if configSet == "" {
		c.JSON(400, gin.H{"error": "configSet parameter is required"})
		return
	}
	directoryConfig, err := getDirectoryConfig(configSet)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var filteredFiles []string

	conditionsMet := make([]bool, len(directoryConfig.Globs))
	for i := range conditionsMet {
		conditionsMet[i] = false
	}

	// iterate over the globs
	for i, glob := range directoryConfig.Globs {
		log.Infof("Processing glob: %s", glob.Glob)
		checkTime, err := convertTimeDelta(glob.Time)
		if err != nil {
			log.Errorf("Error parsing time delta: %v", err)
		}
		log.Infof("Checking files modified after: %s", checkTime)

		// list all files in glob.Path that match glob.Glob
		matches, err := filepath.Glob(filepath.Join(glob.Path, glob.Glob))
		if err != nil {
			log.Errorf("Error processing glob %s: %v", glob.Glob, err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		for _, match := range matches {
			fileinfo, err := os.Stat(match)
			if err != nil {
				log.Errorf("Error getting file info for %s: %v", match, err)
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			log.Infof("File: %s, Modified: %s, CheckTime: %s", match, fileinfo.ModTime(), checkTime)
			if fileinfo.ModTime().After(checkTime) {
				filteredFiles = append(filteredFiles, match)
				conditionsMet[i] = true
			}
		}
	}

	rv := false
	if directoryConfig.Condition == ConditionAny {
		for _, conditionMet := range conditionsMet {
			if conditionMet {
				rv = true
				break
			}
		}
	} else if directoryConfig.Condition == ConditionAll {
		rv = true
		for _, conditionMet := range conditionsMet {
			if !conditionMet {
				rv = false
				break
			}
		}
	}

	c.JSON(200, gin.H{"status": rv, "files": filteredFiles})
	return
}

func getConfig(c *gin.Context) {
	c.JSON(200, config)
}
