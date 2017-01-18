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

	"bufio"

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

type splitterType struct {
	regex  string
	length int64
}

// Log files list related variables that should persist between CollectMetrics calls
var (
	logFiles        []string
	logFilesScanned bool
	scanningCounter int64
)

// splitterTypes is a map with predefined regexp separators for different log types
var splitterTypes = map[string]splitterType{"new-line": splitterType{"\n", 1}, "empty-line": splitterType{"\n\n", 2}, "date-time": splitterType{"(^|\n)[0-9]{4}-[0-1][0-2]-[0-3][0-9] [0-2][0-9]:[0-5][0-9]:[0-5][0-9].[0-9]{3}$", 24}}

// positionCache is log file seek position in bytes
type positionCache struct {
	Offset int64 `json:"offset"`
}

// CollectMetrics collects metrics for testing
func (l Logs) CollectMetrics(mts []plugin.Metric) ([]plugin.Metric, error) {
	if err := l.loadConfig(mts[0].Config); err != nil {
		return nil, err
	}
	l.refreshLogList(mts)

	// Automatically create cache dir if needed
	if err := os.MkdirAll(l.configStr["cache_dir"], 0755); err != nil {
		return nil, fmt.Errorf("cannot create offset cache directory")
	}

	// Set log file splitter
	splitter, err := regexp.Compile(l.configStr["splitter"])
	if err != nil {
		return nil, fmt.Errorf("splitter value is invalid")
	}

	// Set collection time limit
	collectionTime, err := time.ParseDuration(l.configStr["collection_time"])
	if err != nil {
		return nil, fmt.Errorf("collection time value (collection_time) is invalid")
	}
	if len(logFiles) > 0 {
		collectionTime /= time.Duration(len(logFiles))
	}
	logrus.WithFields(logrus.Fields{"files_count": len(logFiles), "collection_time_per_file": collectionTime}).Info("Collection time per file calculated")

	metrics := []plugin.Metric{}

	for _, logFilePath := range logFiles {
		_, logFileName := filepath.Split(logFilePath)

		// Load log file
		logFile, err := os.OpenFile(logFilePath, os.O_RDONLY, 0)
		if err != nil {
			logrus.WithFields(logrus.Fields{"filename": logFilePath, "error": err}).Error("Error while opening log file")
			logFile.Close()
			continue
		}

		// Load last log file offset
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

		// Seek loaded file offset or reset offset if file is smaller (log file retention)
		fi, err := logFile.Stat()
		if err != nil {
			logrus.WithFields(logrus.Fields{"filename": logFilePath, "error": err}).Error("Cannot get log file size")
		}
		if positionCache.Offset <= fi.Size() {
			if _, err := logFile.Seek(positionCache.Offset, os.SEEK_SET); err != nil {
				logrus.WithFields(logrus.Fields{"filename": logFilePath, "error": err}).Error("Cannot seek log file offset")
				logFile.Close()
				continue
			}
		} else {
			logrus.WithField("filename", logFilePath).Warning("Offset outside file content, reading from beginning")
			positionCache.Offset = 0
		}

		// Create buffered reader for opened log file
		logFileReader := bufio.NewReader(logFile)

		// Create log entry helper buffers
		logEntry := ""
		prevSplitter := ""

		// Collect as many data as it is possible during configured collection time limit
		collectStart := time.Now()
		bytesReadBefore := positionCache.Offset
		bytesRead := 0
		fileMetrics := []plugin.Metric{}
		var logFileErr error
		for time.Since(collectStart) < collectionTime && int64(len(metrics)+len(fileMetrics)) < l.configInt["metrics_limit"] {
			// Read data from log file into memory buffer
			var b byte
			b, logFileErr = logFileReader.ReadByte()
			if logFileErr != nil {
				if logFileErr != io.EOF {
					logrus.WithFields(logrus.Fields{"filename": logFilePath, "error": logFileErr}).Error("Error while reading log file data")
					break
				}
			} else {
				logEntry += string(b)
				bytesRead++
			}

			// Get short sample from end of logEntry for quick matching
			splitterLookahead := copyFromEnd(logEntry, l.configInt["splitter_length"])

			if splitter.MatchString(splitterLookahead) || (logFileErr == io.EOF && len(logEntry) > 0) {
				// Trim splitter
				data := splitter.ReplaceAllString(logEntry, "")

				// Save current file position
				positionCache.Offset = bytesReadBefore + int64(bytesRead)

				currentSplitter := splitter.FindString(logEntry)
				if len(data) > 0 {
					// Add splitter back based on splitter_pos config
					switch strings.ToLower(l.configStr["splitter_pos"]) {
					case "before":
						data = fmt.Sprintf("%s%s", prevSplitter, data)

						// Rewind offset before current splitter if "before" option enabled, because previous splitter is used
						// Do not apply if EOF found
						if !(logFileErr == io.EOF) {
							positionCache.Offset -= int64(len(currentSplitter))
						}
					case "after":
						data = fmt.Sprintf("%s%s", data, currentSplitter)
					default:
						return nil, fmt.Errorf("splitter_pos is invalid")
					}

					// Add new metric
					ns := make([]plugin.NamespaceElement, len(mts[0].Namespace))
					copy(ns, mts[0].Namespace)
					ns[2].Value = l.configStr["metric_name"]
					ns[3].Value = logFileName
					mt := plugin.Metric{
						Data:      data,
						Namespace: ns,
						Timestamp: time.Now(),
						Version:   Version,
					}
					fileMetrics = append(fileMetrics, mt)
				}
				prevSplitter = currentSplitter

				// Clear log entry buffer
				logEntry = ""

				if logFileErr == io.EOF {
					break
				}
			}
		}
		logFile.Close()
		logrus.WithFields(logrus.Fields{"metric_count": len(fileMetrics), "filename": logFilePath}).Info("Metric read count during collection time")

		if logFileErr != io.EOF && len(fileMetrics) == 0 {
			logrus.WithFields(logrus.Fields{
				"file_path":                logFilePath,
				"collection_time_per_file": collectionTime,
				"collection_time":          l.configStr["collection_time"],
				"metrics_limit":            l.configInt["metrics_limit"],
			}).Warn("Logs to read stay in log file but there were not read, check task configuration")
		}

		if len(fileMetrics) > 0 {
			posData, err := json.Marshal(positionCache)
			if err != nil {
				logrus.WithError(err).Error("Cannot marshal offset cache JSON data")
				continue
			}
			if err := ioutil.WriteFile(posFilePath, posData, 0644); err != nil {
				logrus.WithField("filename", logFilePath).Error("Cannot save log offset cache file")
				continue
			}

			// Return file metrics only if cache file successfully saved
			metrics = append(metrics, fileMetrics...)
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
	policy.AddNewStringRule([]string{"intel", Name}, "splitter_pos", false, plugin.SetDefaultString("after"))
	policy.AddNewStringRule([]string{"intel", Name}, "splitter_type", false, plugin.SetDefaultString("new-line"))
	policy.AddNewStringRule([]string{"intel", Name}, "splitter", false, plugin.SetDefaultString(splitterTypes["new-line"].regex))
	policy.AddNewStringRule([]string{"intel", Name}, "cache_dir", false, plugin.SetDefaultString("/var/cache/snap"))
	policy.AddNewStringRule([]string{"intel", Name}, "collection_time", false, plugin.SetDefaultString("300ms"))
	policy.AddNewIntRule([]string{"intel", Name}, "scanning_dir_counter", false, plugin.SetDefaultInt(0))
	policy.AddNewIntRule([]string{"intel", Name}, "splitter_length", false, plugin.SetDefaultInt(1))
	policy.AddNewIntRule([]string{"intel", Name}, "metrics_limit", false, plugin.SetDefaultInt(300))
	return *policy, nil
}

// copyFromEnd returns n characters from end of string s
func copyFromEnd(s string, n int64) string {
	if n <= int64(len(s)) {
		return s[int64(len(s))-n:]
	}
	return s
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
				expandPaths(strings.Join(append(append(patternElements[:i], d), patternElements[i+1:]...), string(os.PathSeparator)), collected)
			}
			return
		}
	}

	expandedPath, err := filepath.Glob(pattern)
	if err != nil || len(expandedPath) == 0 {
		logrus.WithFields(logrus.Fields{"path_pattern": pattern, "expanded_path": expandedPath, "error": err}).Error("Cannot expand path, check if path pattern is correct")
	}
	for _, path := range expandedPath {
		if result, err := isDir(path); result && err == nil {
			*collected = append(*collected, path)
		}
	}
}

// Check if file name matches log_file namespace entry
func matchLogFileNamespaceEntry(mts []plugin.Metric, name string) bool {
	for _, m := range mts {
		if m.Namespace[3].Value == "*" || m.Namespace[3].Value == name {
			return true
		}
	}
	return false
}

// List all files that matches the specified regexp
func filterFiles(path string, filePattern string, mts []plugin.Metric) []string {
	// Read dir contents
	files, err := ioutil.ReadDir(path)
	if err != nil {
		logrus.WithField("path", path).Error("Cannot access path! Log file list generation failed.")
		return []string{}
	}

	// Filter files inside dir
	fp, err := regexp.Compile(filePattern)
	if err != nil {
		logrus.WithError(err).Error("File pattern must be valid regular expression!")
		return []string{}
	}
	result := []string{}
	for _, file := range files {
		if !file.IsDir() && fp.MatchString(file.Name()) && matchLogFileNamespaceEntry(mts, file.Name()) {
			result = append(result, filepath.Join(path, file.Name()))
		}
	}
	return result
}

// Initialize plugin configuration
func (l *Logs) refreshLogList(mts []plugin.Metric) {
	if scanningCounter <= 0 || !logFilesScanned {
		scanningCounter = l.configInt["scanning_dir_counter"]

		logDir := l.configStr["log_dir"]
		logFile := l.configStr["log_file"]

		allPaths := []string{}
		expandPaths(logDir, &allPaths)

		logFiles = []string{}
		for _, path := range allPaths {
			logFiles = append(logFiles, filterFiles(path, fmt.Sprintf("^%s$", logFile), mts)...)
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
			l.configStr["splitter"] = val.regex
			l.configInt["splitter_length"] = val.length
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
