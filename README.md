# Snap collector plugin - logs
This plugin collects log messages partially for each collection run. Log file reading is limited by time.

It's used in the [Snap framework](http://github.com:intelsdi-x/snap).

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license-and-authors)
6. [Acknowledgements](#acknowledgements)

## Getting Started
### System Requirements 
* [golang 1.6+](https://golang.org/dl/) (needed only for building)

### Operating systems
All OSs currently supported by snap:
* Linux/amd64
* Darwin/amd64

### Installation

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-collector-logs

Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-logs.git
```

Build the plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `./build/`

### Configuration and Usage
* Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)

#### Task manifest configuration options
Option|Description|Default value
------|-----------|-------------
"metric_name"|Declaration of metric name, the first dynamic part of namespace|all
"scanning_dir_counter"|Defines when directory should be scanned (per collection or after several collections), for checking if there are new files in logs directory (if not defined|logs directory is scanned per metrics collection)|0
"log_dir"|Filepath expression to get logs directory, e.g.:"/var/log/kolla/(neutron\|nova\|cinder)". Filepath expressions available: "(dir1\|dir2\|dirn)", "{dir1,dir2,dirn}", "*"|/var/log
"log_file"|Regular expression to get file/files which logs in directory defined as a “log_dir”, e.g.: "keystone_apache_\S{1,}"|.*
"splitter_type"|Predefined splitter type. Available options: new-line, empty-line, custom. If custom, you can set "splitter" option manually.|new-line
"splitter_pos"|Position of splitter. Available options: before, after.|after
"splitter"|Characteristic character/characters to split logs (on default logs are splitted per lines)|\\n
"cache_dir"|Directory in which offsets for next reading of logs are saved, e.g: "/var/cache/snap/". Directory must already exist before collection, otherwise last log read position will not be saved.|/var/cache/snap
"collection_time"|Maximum time for continuous reading of logs, it should be lower than task manifest|300ms

## Documentation

### Collected Metrics
This plugin has the ability to gather the following metrics:

Namespace | Description
----------|-----------------------
/intel/logs/[metric_name]/[log_file]/message|Single log message


### Examples
This is an example running logs collector and writing data to a file. It is assumed that you are using the latest Snap binary and plugins.

The example is run from a directory which includes snaptel, snapteld, along with the plugins and task file.

In one terminal window, open the Snap daemon (in this case with logging set to 1 and trust disabled):
```
$ snapteld -l 1 -t 0
```

In another terminal window:
Load logs plugin
```
$ snaptel plugin load snap-plugin-collector-logs
```
See available metrics for your system
```
$ snaptel metric list
```

Create a task manifest file (e.g. `task-logs.json`):    
```json
{
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "3s"
    },
    "workflow": {
        "collect": {
            "metrics": {
                "/intel/logs/*": {}
            },
            "config": {
                "/intel/logs": {
                    "metric_name": "rabbit_logs",
                    "cache_dir": "/home/test/cache/snap",
                    "log_dir": "/home/test/kolla/*",
                    "log_file": ".*rabbit.*",
                    "splitter_type": "empty-line"
                }
            },
            "publish": [
                {
                    "plugin_name": "file",
                    "config": {
                        "file": "/tmp/published_logs"
                    }
                }
            ]
        }
    }
}
```

Load file plugin for publishing:
```
$ snaptel plugin load snap-plugin-publisher-file
Plugin loaded
Name: file
Version: 3
Type: publisher
Signed: false
Loaded Time: Fri, 20 Nov 2015 11:41:39 PST
```

Create task:
```
$ snaptel task create -t task-logs.json
Using task manifest to create task
Task created
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
Name: Task-02dd7ff4-8106-47e9-8b86-70067cd0a850
State: Running
```

See file output (this is just single collection output with default collection_time of 300ms): [EXAMPLE_OUTPUT.md](EXAMPLE_OUTPUT.md)

Stop task:
```
$ snaptel task stop 02dd7ff4-8106-47e9-8b86-70067cd0a850
Task stopped:
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release. If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-logs/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-logs/pulls).

## Community Support
This repository is one of **many** plugins in **Snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[Snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [@mkleina](https://github.com/mkleina)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.
