/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2016 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package logs

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

const (
	// Name of plugin
	Name = "logs"
	// Version of plugin
	Version = 1
)

// Logs collector implementation used for testing
type Logs struct {
	configStr map[string]string
	configInt map[string]int64
}

// Log files list related variables that should persist between CollectMetrics calls
var (
	logFiles        []string
	logFilesScanned bool
	scanningCounter int64
)

// splitterTypes is a map with predefined regexp separators for different log types
var splitterTypes = map[string]string{"new-line": "\n", "empty-line": "^\n|\n\n"}

// positionCache is log file seek position in bytes
type positionCache struct {
	Offset int64 `json:"offset,omitempty"`
}

// CollectMetrics collects metrics for testing
func (l Logs) CollectMetrics(mts []plugin.Metric) ([]plugin.Metric, error) {
	if err := l.loadConfig(mts[0].Config); err != nil {
		return nil, err
	}
	l.refreshLogList()
	metrics := []plugin.Metric{}

	// Move to last known file position
	for _, logFilePath := range logFiles {
		_, logFileName := filepath.Split(logFilePath)

		// Load log file
		logFile, err := os.OpenFile(logFilePath, os.O_RDONLY, 0)
		if err != nil {
			logrus.WithField("filename", logFilePath).Error("Error while opening log file")
		}

		buffer := make([]byte, 1)
		logEntry := ""

		// Go to last log file position
		posFilePath := filepath.Join(l.configStr["cache_dir"], fmt.Sprintf("%s_%s.json", l.configStr["metric_name"], logFileName))
		positionCache := positionCache{}
		posData, err := ioutil.ReadFile(posFilePath)
		if err != nil {
			logrus.WithFields(logrus.Fields{"filename": posFilePath, "error": err}).Warning("Cannot read log offset cache file. This warning may appear when new log file found.")
		} else {
			if err := json.Unmarshal(posData, &positionCache); err != nil {
				logrus.WithFields(logrus.Fields{"filename": posFilePath, "error": err}).Error("Cannot parse log offset cache file")
			}
		}
		if _, err := logFile.Seek(positionCache.Offset, os.SEEK_SET); err != nil {
			return nil, err
		}

		// Set collection time limit
		collectStart := time.Now()
		collectionTime, err := time.ParseDuration(l.configStr["collection_time"])
		if err != nil {
			return nil, fmt.Errorf("collection time value (collection_time) is invalid")
		}

		// Collect as many data as it is possible during configured collection time limit
		for time.Since(collectStart) < collectionTime {
			// Read 1 byte from file
			_, logFileErr := logFile.Read(buffer)
			if logFileErr != nil {
				if logFileErr != io.EOF {
					return nil, err
				}
			} else {
				logEntry += string(buffer)
			}

			// Return log metric if splitter appeared or end of file reached
			splitter, err := regexp.Compile(l.configStr["splitter"])
			if err != nil {
				return nil, fmt.Errorf("splitter value is invalid")
			}
			if splitter.MatchString(logEntry) || (logFileErr == io.EOF && len(logEntry) > 0) {
				// Remove splitter string and trim whitespaces
				data := strings.TrimSpace(splitter.ReplaceAllString(logEntry, ""))

				// Clear log entry buffer and save current file position
				logEntry = ""
				if positionCache.Offset, err = logFile.Seek(0, os.SEEK_CUR); err != nil {
					logrus.WithField("filename", logFilePath).Error("Cannot get current offset in log file")
				}

				if len(data) > 0 {
					mt := plugin.Metric{
						Data:      data,
						Namespace: mts[0].Namespace,
						Timestamp: time.Now(),
						Version:   Version,
					}
					mt.Namespace[2].Value = l.configStr["metric_name"]
					mt.Namespace[3].Value = logFileName
					metrics = append(metrics, mt)
				}

				if logFileErr == io.EOF {
					break
				}
			}
		}
		logFile.Close()

		if len(metrics) > 0 {
			posData, err := json.Marshal(positionCache)
			if err != nil {
				logrus.WithField("error", err).Error("Cannot marshal offset cache JSON data")
			}
			if err := ioutil.WriteFile(posFilePath, posData, 0644); err != nil {
				logrus.WithField("filename", logFilePath).Error("Cannot save log offset cache file")
			}
		}
	}

	return metrics, nil
}

//GetMetricTypes returns metric types for testing
func (l Logs) GetMetricTypes(cfg plugin.Config) ([]plugin.Metric, error) {
	if err := l.loadConfig(cfg); err != nil {
		return nil, err
	}

	mts := []plugin.Metric{}
	mts = append(mts, plugin.Metric{
		Namespace:   plugin.NewNamespace("intel", Name).AddDynamicElement("metric_name", "Metric name defined in config file").AddDynamicElement("log_file", "Log file name").AddStaticElement("message"),
		Description: "Single log message",
		Unit:        "string",
		Version:     Version,
	})

	return mts, nil
}

//GetConfigPolicy returns a ConfigPolicy for testing
func (l Logs) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	policy := plugin.NewConfigPolicy()
	policy.AddNewStringRule([]string{"intel", Name}, "metric_name", false, plugin.SetDefaultString("all"))
	policy.AddNewStringRule([]string{"intel", Name}, "log_dir", false, plugin.SetDefaultString("/var/log"))
	policy.AddNewStringRule([]string{"intel", Name}, "log_file", false, plugin.SetDefaultString(".*"))
	policy.AddNewStringRule([]string{"intel", Name}, "splitter_type", false, plugin.SetDefaultString("new-line"))
	policy.AddNewStringRule([]string{"intel", Name}, "splitter", false, plugin.SetDefaultString(splitterTypes["new-line"]))
	policy.AddNewStringRule([]string{"intel", Name}, "cache_dir", false, plugin.SetDefaultString("/var/cache/snap"))
	policy.AddNewStringRule([]string{"intel", Name}, "collection_time", false, plugin.SetDefaultString("300ms"))
	policy.AddNewIntRule([]string{"intel", Name}, "scanning_dir_counter", false, plugin.SetDefaultInt(0))
	return *policy, nil
}

// isDir returns true if specified path is dir and false otherwise
func isDir(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return f.IsDir(), nil
}

// expandPaths converts expressions like /home/*/(Downloads|Desktop) to list of real paths
// Supported patterns: (dir1|dir2|dirn), (dir1,dir2,dirn), {dir1|dir2|dirn}, {dir1,dir2,dirn}
// and all OS filesystem patterns like *, *.*, .., ~ etc.
func expandPaths(pattern string, collected *[]string) {
	patternElements := strings.Split(pattern, string(os.PathSeparator))

	separators := regexp.MustCompile(`\,|\|`)
	brackets := regexp.MustCompile(`\{|\}|\(|\)`)

	for i, pe := range patternElements {
		if brackets.MatchString(pe) {
			dirs := separators.Split(brackets.ReplaceAllString(pe, ""), -1)
			for _, d := range dirs {
				expandPaths(filepath.Join(append(append(patternElements[:i], d), patternElements[i+1:]...)...), collected)
			}
			return
		}
	}

	expandedPath, _ := filepath.Glob(pattern)
	for _, path := range expandedPath {
		if result, err := isDir(path); result && err == nil {
			*collected = append(*collected, path)
		}
	}
}

// List all files that matches the specified regexp
func filterFiles(path string, filePattern string) []string {
	// Read dir contents
	files, err := ioutil.ReadDir(path)
	if err != nil {
		logrus.WithField("path", path).Error("Cannot access path! Log file list generation failed.")
		return []string{}
	}

	// Filter files inside dir
	fp, err := regexp.Compile(filePattern)
	if err != nil {
		logrus.WithField("error", err).Error("File pattern must be valid regular expression!")
		return []string{}
	}
	result := []string{}
	for _, file := range files {
		if !file.IsDir() && fp.MatchString(file.Name()) {
			result = append(result, filepath.Join(path, file.Name()))
		}
	}
	return result
}

// Initialize plugin configuration
func (l *Logs) refreshLogList() {
	if scanningCounter <= 0 || !logFilesScanned {
		scanningCounter = l.configInt["scanning_dir_counter"]

		logDir := l.configStr["log_dir"]
		logFile := l.configStr["log_file"]

		allPaths := []string{}
		expandPaths(logDir, &allPaths)

		logFiles = []string{}
		for _, path := range allPaths {
			logFiles = append(logFiles, filterFiles(path, fmt.Sprintf("^%s$", logFile))...)
		}
		logFilesScanned = true
	} else {
		scanningCounter--
	}
}

// Load config values
func (l *Logs) loadConfig(cfg plugin.Config) error {
	l.configStr = make(map[string]string)
	l.configInt = make(map[string]int64)

	for key := range cfg {
		if val, err := cfg.GetInt(key); err == nil {
			l.configInt[key] = val
		}
		if val, err := cfg.GetString(key); err == nil {
			l.configStr[key] = val
		}
	}

	// Configure splitter for selected preset
	if !strings.EqualFold(l.configStr["splitter_type"], "custom") {
		key := strings.ToLower(l.configStr["splitter_type"])
		if val, ok := splitterTypes[key]; ok {
			l.configStr["splitter"] = val
		} else {
			return fmt.Errorf("splitter type \"%s\" is not supported", key)
		}
	}
	if len(l.configStr["splitter"]) == 0 {
		return fmt.Errorf("please configure \"splitter\" option in your task manifest")
	}
	if _, err := regexp.Compile(l.configStr["splitter"]); err != nil {
		return fmt.Errorf("splitter value is invalid")
	}
	if _, err := regexp.Compile(l.configStr["log_file"]); err != nil {
		return fmt.Errorf("log_file value is invalid")
	}
	if _, err := time.ParseDuration(l.configStr["collection_time"]); err != nil {
		return fmt.Errorf("collection time value (collection_time) is invalid")
	}

	return nil
}
