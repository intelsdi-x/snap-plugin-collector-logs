/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2015 Intel Corporation

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

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/ctypes"
)

const (
	// Name of plugin
	Name = "logs"
	// Version of plugin
	Version = 1
	// Type of plugin
	Type = plugin.CollectorPluginType
)

// Logs collector implementation used for testing
type Logs struct {
	logFiles []string
	config   map[string]ctypes.ConfigValue
}

// LogTypes is a map with predefined regexp separators for different log types
var logTypes = map[string]string{"apache": "\n", "rabbit": "^\n|\n\n"}

// positionCache is log file seek position in bytes
type positionCache struct {
	Position int64 `json:"position,omitempty"`
}

// CollectMetrics collects metrics for testing
func (l *Logs) CollectMetrics(mts []plugin.MetricType) ([]plugin.MetricType, error) {
	if len(mts) == 0 {
		return nil, fmt.Errorf("no metrics to collect")
	}
	if err := l.loadConfig(mts[0].Config().Table()); err != nil {
		return nil, err
	}
	l.refreshLogList()
	metrics := []plugin.MetricType{}

	// Move to last known file position
	for _, logFilePath := range l.logFiles {
		_, logFileName := filepath.Split(logFilePath)

		// Load log file
		logFile, err := os.OpenFile(logFilePath, os.O_RDONLY, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while opening log file: %s\n", err)
		}

		buffer := make([]byte, 1)
		logEntry := ""

		// Go to last log file position
		posFilePath := filepath.Join(l.config["cache_dir"].(ctypes.ConfigValueStr).Value, l.config["metric_name"].(ctypes.ConfigValueStr).Value+"_"+logFileName+".json")
		positionCache := positionCache{}
		posData, err := ioutil.ReadFile(posFilePath)
		if err == nil {
			if err := json.Unmarshal(posData, &positionCache); err != nil {
				return nil, err
			}
		}

		if _, err := logFile.Seek(positionCache.Position, os.SEEK_SET); err != nil {
			return nil, err
		}

		// Set collection time limit
		collectStart := time.Now()
		collectionTime := time.Duration(l.config["collection_time"].(ctypes.ConfigValueInt).Value) * time.Millisecond

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
			splitter := regexp.MustCompile(l.config["splitter"].(ctypes.ConfigValueStr).Value)
			if splitter.MatchString(logEntry) || (logFileErr == io.EOF && len(logEntry) > 0) {
				// Remove splitter string and trim whitespaces
				data := strings.TrimSpace(splitter.ReplaceAllString(logEntry, ""))

				// Clear log entry buffer and save current file position
				logEntry = ""
				positionCache.Position, _ = logFile.Seek(0, os.SEEK_CUR)

				if len(data) > 0 {
					mt := plugin.MetricType{
						Data_:      data,
						Namespace_: core.NewNamespace("intel", Name, l.config["metric_name"].(ctypes.ConfigValueStr).Value, logFileName, "message"),
						Timestamp_: time.Now(),
					}
					metrics = append(metrics, mt)
				}

				if logFileErr == io.EOF {
					break
				}
			}
		}
		logFile.Close()

		if len(metrics) > 0 {
			posData, _ := json.Marshal(positionCache)
			ioutil.WriteFile(posFilePath, posData, 0644)
		}
	}

	return metrics, nil
}

//GetMetricTypes returns metric types for testing
func (l *Logs) GetMetricTypes(cfg plugin.ConfigType) ([]plugin.MetricType, error) {
	if err := l.loadConfig(cfg.Table()); err != nil {
		return nil, err
	}

	mts := []plugin.MetricType{}
	mts = append(mts, plugin.MetricType{
		Namespace_:   core.NewNamespace("intel", Name).AddDynamicElement("metric_name", "Metric name defined in config file").AddDynamicElement("log_file", "Log file name").AddStaticElement("message"),
		Description_: "Single log message",
		Unit_:        "string",
	})

	return mts, nil
}

//GetConfigPolicy returns a ConfigPolicy for testing
func (l *Logs) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	rule1, _ := cpolicy.NewStringRule("metric_name", false, "all")
	rule2, _ := cpolicy.NewStringRule("log_dir", false, "/var/log")
	rule3, _ := cpolicy.NewStringRule("log_file", false, ".*")
	rule4, _ := cpolicy.NewStringRule("log_type", false, "apache")
	rule5, _ := cpolicy.NewStringRule("splitter", false, logTypes["apache"])
	rule6, _ := cpolicy.NewStringRule("cache_dir", false, "/var/cache/snap")
	rule7, _ := cpolicy.NewIntegerRule("scanning_dir_counter", false, 0)
	rule8, _ := cpolicy.NewIntegerRule("collection_time", false, 300)
	p := cpolicy.NewPolicyNode()
	p.Add(rule1, rule2, rule3, rule4, rule5, rule6, rule7, rule8)
	c.Add([]string{"intel", "logs"}, p)
	return c, nil
}

//Meta returns meta data for testing
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(
		Name,
		Version,
		Type,
		[]string{plugin.SnapGOBContentType},
		[]string{plugin.SnapGOBContentType},
		plugin.CacheTTL(100*time.Millisecond),
		plugin.RoutingStrategy(plugin.StickyRouting),
	)
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
func filterFiles(path string, filePattern string, collected *[]string) {
	// Read dir contents
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot access %s! Log file list generation failed.\n", path)
		return
	}

	// Filter files inside dir
	fp := regexp.MustCompile(filePattern)
	for _, file := range files {
		if !file.IsDir() && fp.MatchString(file.Name()) {
			*collected = append(*collected, filepath.Join(path, file.Name()))
		}
	}
}

// Initialize plugin configuration
func (l *Logs) refreshLogList() {
	logDir := l.config["log_dir"].(ctypes.ConfigValueStr).Value
	logFile := l.config["log_file"].(ctypes.ConfigValueStr).Value

	allPaths := []string{}
	expandPaths(logDir, &allPaths)

	l.logFiles = []string{}
	for _, path := range allPaths {
		filterFiles(path, fmt.Sprintf("^%s$", logFile), &l.logFiles)
	}
}

// Load config values
func (l *Logs) loadConfig(cfg map[string]ctypes.ConfigValue) error {
	if l.config == nil {
		l.config = cfg

		// Configure splitter for selected preset
		if !strings.EqualFold(l.config["log_type"].(ctypes.ConfigValueStr).Value, "custom") {
			key := strings.ToLower(l.config["log_type"].(ctypes.ConfigValueStr).Value)
			if val, ok := logTypes[key]; ok {
				l.config["splitter"] = ctypes.ConfigValueStr{Value: val}
			} else {
				return fmt.Errorf("log type \"%s\" is not supported", key)
			}
		}
		if len(l.config["splitter"].(ctypes.ConfigValueStr).Value) == 0 {
			return fmt.Errorf("please configure \"splitter\" option in your task manifest")
		}
	}
	return nil
}
