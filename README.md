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
"log_type"|Predefined log type. Available options: apache, rabbit, custom. If custom, you can set "splitter" option manually.|apache
"splitter"|Characteristic character/characters to split logs (on default logs are splitted per lines)|\\n
"cache_dir"|Directory in which offsets for next reading of logs are saved, e.g: "/var/cache/snap/"|/var/cache/snap
"collection_time"|Maximum time for continuous reading of logs in milliseconds, it should be lower than task manifest|300

## Documentation

### Collected Metrics
This plugin has the ability to gather the following metrics:

Namespace | Description (optional)
----------|-----------------------
/intel/logs/[metric_name]/[log_file]/message|Single log message


### Examples
This is an example running psutil and writing data to a file. It is assumed that you are using the latest Snap binary and plugins.

The example is run from a directory which includes snaptel, snapteld, along with the plugins and task file.

In one terminal window, open the Snap daemon (in this case with logging set to 1 and trust disabled):
```
$ snapteld -l 1 -t 0
```

In another terminal window:
Load logs plugin
```
$ snaptel plugin load snap-plugin-collector-logss
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
                    "log_type": "rabbit"
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

See file output (this is just single collection output with default collection_time of 300ms):
```json
[
    {
        "timestamp": "2016-12-28T12:17:21.949648175+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:15 ===\naccepting AMQP connection <0.4612.0> (10.0.0.1:39908 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251356733+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:21.954813257+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:15 ===\nConnection <0.4612.0> (10.0.0.1:39908 -> 10.0.0.1:5672) has a client-provided name: neutron-dhcp-agent:10:1f5aeddf-66d7-43e2-bc03-b71c002cf180",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251357266+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:21.957130655+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:15 ===\naccepting AMQP connection <0.4633.0> (10.0.0.1:39910 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251365284+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:21.95962844+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:15 ===\naccepting AMQP connection <0.4636.0> (10.0.0.1:39912 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251365456+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:21.964532387+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:15 ===\nConnection <0.4633.0> (10.0.0.1:39910 -> 10.0.0.1:5672) has a client-provided name: neutron-dhcp-agent:10:7db5581d-3dd5-4db2-a15e-e4a684927470",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251365638+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:21.970014508+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:15 ===\nConnection <0.4636.0> (10.0.0.1:39912 -> 10.0.0.1:5672) has a client-provided name: neutron-dhcp-agent:10:bd1351f0-0095-4439-a126-96569dbb8e47",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251365811+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:21.972730693+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\naccepting AMQP connection <0.4672.0> (10.0.0.1:39980 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251365982+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:21.977676051+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\nConnection <0.4672.0> (10.0.0.1:39980 -> 10.0.0.1:5672) has a client-provided name: heat-engine:23:e56a3c56-8939-4df4-92f3-69cdf66da68f",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251366164+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:21.980087857+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\naccepting AMQP connection <0.4712.0> (10.0.0.1:39982 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251366336+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:21.982826629+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\naccepting AMQP connection <0.4715.0> (10.0.0.1:39984 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251366517+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:21.985610118+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\naccepting AMQP connection <0.4718.0> (10.0.0.1:39988 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251366688+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:21.988334773+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\naccepting AMQP connection <0.4724.0> (10.0.0.1:39990 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.25136686+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:21.991164509+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\naccepting AMQP connection <0.4721.0> (10.0.0.1:39986 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251367043+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:21.996238704+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\nConnection <0.4712.0> (10.0.0.1:39982 -> 10.0.0.1:5672) has a client-provided name: heat-engine:20:6ad18750-80f6-49f1-93e6-824ad0569b31",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251367212+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.001032325+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\nConnection <0.4715.0> (10.0.0.1:39984 -> 10.0.0.1:5672) has a client-provided name: heat-engine:22:1ae8745d-7621-4bd3-bcef-33c2b0b59e92",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.2513674+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.006767305+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\nConnection <0.4718.0> (10.0.0.1:39988 -> 10.0.0.1:5672) has a client-provided name: heat-engine:23:5a4be275-54ca-42a3-9f93-0d82da1c8ba4",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251367574+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.011168225+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\nConnection <0.4721.0> (10.0.0.1:39986 -> 10.0.0.1:5672) has a client-provided name: heat-engine:21:f318f00f-2cd1-4288-9869-9b71449e8eb4",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.25136776+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.016666735+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\nConnection <0.4724.0> (10.0.0.1:39990 -> 10.0.0.1:5672) has a client-provided name: heat-engine:19:9f2c22bc-2777-4ad8-8d21-6bad76d6dbdd",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251367933+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.019465003+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\naccepting AMQP connection <0.4866.0> (10.0.0.1:39992 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251368103+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.022797272+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\naccepting AMQP connection <0.4869.0> (10.0.0.1:39994 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251368284+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.025154897+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\naccepting AMQP connection <0.4872.0> (10.0.0.1:39996 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251368458+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.028546362+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\naccepting AMQP connection <0.4877.0> (10.0.0.1:39998 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251368639+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.033673059+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\nConnection <0.4872.0> (10.0.0.1:39996 -> 10.0.0.1:5672) has a client-provided name: heat-engine:20:8cc1e8b6-ca4b-427c-991b-f0b55cb0a6f8",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251368808+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.03602268+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\naccepting AMQP connection <0.4881.0> (10.0.0.1:40000 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.25136898+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.040765584+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\nConnection <0.4866.0> (10.0.0.1:39992 -> 10.0.0.1:5672) has a client-provided name: heat-engine:22:7c0e1076-a589-4d42-b573-cd3553343037",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251369163+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.045311617+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\nConnection <0.4881.0> (10.0.0.1:40000 -> 10.0.0.1:5672) has a client-provided name: nova-compute:7:89ae9d72-d74f-4509-9896-cc7e7b388b3b",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251369369+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.049710171+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\nConnection <0.4869.0> (10.0.0.1:39994 -> 10.0.0.1:5672) has a client-provided name: heat-engine:23:a5d27411-528e-4290-b645-5501409cbb1d",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251369551+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.054677197+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\nConnection <0.4877.0> (10.0.0.1:39998 -> 10.0.0.1:5672) has a client-provided name: heat-engine:21:a69c052b-3372-4822-b55c-c92d20fdb66b",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.25136973+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.057850248+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\naccepting AMQP connection <0.5000.0> (10.0.0.1:40002 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251369901+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.060366205+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\naccepting AMQP connection <0.5006.0> (10.0.0.1:40004 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251370072+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.062711832+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\naccepting AMQP connection <0.5015.0> (10.0.0.1:40006 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251370248+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.067217239+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\nConnection <0.5000.0> (10.0.0.1:40002 -> 10.0.0.1:5672) has a client-provided name: heat-engine:19:0e58ee85-9686-404d-bbee-96fee9b28da8",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251370422+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.07254339+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\nConnection <0.5006.0> (10.0.0.1:40004 -> 10.0.0.1:5672) has a client-provided name: nova-compute:7:dc749d8b-15ea-4ae3-b524-22b29b810880",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251370605+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.075856599+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\naccepting AMQP connection <0.5044.0> (10.0.0.1:40008 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251370775+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.081711333+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\nConnection <0.5044.0> (10.0.0.1:40008 -> 10.0.0.1:5672) has a client-provided name: heat-engine:21:78d0a0e3-091d-4167-af73-572ef5804996",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251370969+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.087111302+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\nConnection <0.5015.0> (10.0.0.1:40006 -> 10.0.0.1:5672) has a client-provided name: heat-engine:20:0d5ea932-de99-4f68-857a-0da1480b6715",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251371142+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.089583767+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\naccepting AMQP connection <0.5096.0> (10.0.0.1:40010 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.25137133+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.094390276+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\nConnection <0.5096.0> (10.0.0.1:40010 -> 10.0.0.1:5672) has a client-provided name: heat-engine:22:34b8213b-2c75-46cf-90e9-dca50eb36ae5",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251371511+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.100626606+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\naccepting AMQP connection <0.5128.0> (10.0.0.1:40016 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251371691+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.105555652+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:18 ===\nConnection <0.5128.0> (10.0.0.1:40016 -> 10.0.0.1:5672) has a client-provided name: heat-engine:19:8232ffc3-9a9a-4ae2-b31d-96df2e3c64f7",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251371864+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.107850697+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:22 ===\naccepting AMQP connection <0.5185.0> (10.0.0.1:40130 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251372035+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.112250736+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:22 ===\nConnection <0.5185.0> (10.0.0.1:40130 -> 10.0.0.1:5672) has a client-provided name: neutron-server:25:9981b0ec-8b61-4cd4-9f08-f3ec3b36d6b6",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251372247+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.11503375+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:22 ===\naccepting AMQP connection <0.5215.0> (10.0.0.1:40132 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251372417+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.119189709+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:22 ===\nConnection <0.5215.0> (10.0.0.1:40132 -> 10.0.0.1:5672) has a client-provided name: neutron-server:26:8d827ca6-fcee-4234-96bd-64ea5f0f918f",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251372618+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.12115772+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:22 ===\naccepting AMQP connection <0.5235.0> (10.0.0.1:40138 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251372788+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.125061027+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:22 ===\nConnection <0.5235.0> (10.0.0.1:40138 -> 10.0.0.1:5672) has a client-provided name: neutron-server:27:f7fb6b93-8d52-447b-a7b0-f0e3439b8455",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251372959+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.127687636+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:22 ===\naccepting AMQP connection <0.5302.0> (10.0.0.1:40140 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251373139+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.132119977+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:22 ===\nConnection <0.5302.0> (10.0.0.1:40140 -> 10.0.0.1:5672) has a client-provided name: neutron-server:25:d96caa7f-2f06-4096-aaa1-4ecb8e2ca42e",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251373309+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.134077481+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:22 ===\naccepting AMQP connection <0.5341.0> (10.0.0.1:40142 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251373491+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.138073992+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:22 ===\nConnection <0.5341.0> (10.0.0.1:40142 -> 10.0.0.1:5672) has a client-provided name: neutron-server:25:a51bc61f-2d1d-4ab8-a0c2-28c5e20d4058",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251373661+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.141610724+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:22 ===\naccepting AMQP connection <0.5363.0> (10.0.0.1:40144 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251373829+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.145718292+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:22 ===\nConnection <0.5363.0> (10.0.0.1:40144 -> 10.0.0.1:5672) has a client-provided name: neutron-server:25:7f07e97e-c450-4f8b-b372-f224a2954214",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251373999+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.148179935+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:22 ===\naccepting AMQP connection <0.5385.0> (10.0.0.1:40146 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251374182+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.152767091+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:22 ===\nConnection <0.5385.0> (10.0.0.1:40146 -> 10.0.0.1:5672) has a client-provided name: nova-consoleauth:6:02ea1696-f427-4a94-b287-a0aa3cb3a098",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251374353+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.155131288+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:22 ===\naccepting AMQP connection <0.5424.0> (10.0.0.1:40172 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251374522+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.159602875+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:22 ===\nConnection <0.5424.0> (10.0.0.1:40172 -> 10.0.0.1:5672) has a client-provided name: nova-conductor:23:0f4d91b0-7636-4e29-8e94-ebf7e3aba0af",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251374723+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.16169879+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:22 ===\naccepting AMQP connection <0.5439.0> (10.0.0.1:40182 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251374902+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.165503761+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:22 ===\nConnection <0.5439.0> (10.0.0.1:40182 -> 10.0.0.1:5672) has a client-provided name: nova-conductor:21:371f43f3-4c55-4154-8247-ddb109cec75a",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251375071+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.167995166+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:22 ===\naccepting AMQP connection <0.5464.0> (10.0.0.1:40188 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.25137524+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.172426497+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:22 ===\nConnection <0.5464.0> (10.0.0.1:40188 -> 10.0.0.1:5672) has a client-provided name: nova-conductor:22:af94c80c-04ea-4efd-8974-d0f39bc4ba55",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251375419+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.174248504+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:22 ===\naccepting AMQP connection <0.5507.0> (10.0.0.1:40190 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251375588+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.177982128+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:22 ===\nConnection <0.5507.0> (10.0.0.1:40190 -> 10.0.0.1:5672) has a client-provided name: nova-conductor:20:896562fc-baf2-4bf7-8479-d476b21bba47",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251375758+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.1807925+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:22 ===\naccepting AMQP connection <0.5518.0> (10.0.0.1:40196 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251375938+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.18478799+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:22 ===\nConnection <0.5518.0> (10.0.0.1:40196 -> 10.0.0.1:5672) has a client-provided name: nova-conductor:19:55ec74d8-e5c8-42eb-a759-7ed3dc26365a",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251376114+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.186922446+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:22 ===\naccepting AMQP connection <0.5567.0> (10.0.0.1:40198 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251376301+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.190999796+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:48:22 ===\nConnection <0.5567.0> (10.0.0.1:40198 -> 10.0.0.1:5672) has a client-provided name: nova-scheduler:6:69878ca7-bf1d-47e5-8d98-9614c0ece4c6",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251376472+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.193595564+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:49:15 ===\naccepting AMQP connection <0.6031.0> (10.0.0.1:41234 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251376642+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.197401213+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:49:15 ===\nConnection <0.6031.0> (10.0.0.1:41234 -> 10.0.0.1:5672) has a client-provided name: neutron-server:26:0edf54c9-8868-4db1-bb40-60fc2b85539f",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251376821+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.202282702+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:49:15 ===\naccepting AMQP connection <0.6042.0> (10.0.0.1:41240 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251376989+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.205862585+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:49:15 ===\nConnection <0.6042.0> (10.0.0.1:41240 -> 10.0.0.1:5672) has a client-provided name: neutron-server:25:317888f2-4539-4451-b2f9-3c26fdcb1c5c",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.25137716+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.211708268+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=WARNING REPORT==== 7-Dec-2016::03:49:18 ===\nclosing AMQP connection <0.5006.0> (10.0.0.1:40004 -> 10.0.0.1:5672 - nova-compute:7:dc749d8b-15ea-4ae3-b524-22b29b810880):\nclient unexpectedly closed TCP connection",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251377339+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.216067866+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=WARNING REPORT==== 7-Dec-2016::03:49:18 ===\nclosing AMQP connection <0.4881.0> (10.0.0.1:40000 -> 10.0.0.1:5672 - nova-compute:7:89ae9d72-d74f-4509-9896-cc7e7b388b3b):\nclient unexpectedly closed TCP connection",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.25137751+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.21852199+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:49:20 ===\naccepting AMQP connection <0.6091.0> (10.0.0.1:41314 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251377685+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.222162549+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:49:20 ===\nConnection <0.6091.0> (10.0.0.1:41314 -> 10.0.0.1:5672) has a client-provided name: nova-compute:7:d377fb4a-b463-40c6-9049-8f341ffdac44",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251377865+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.224179949+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:49:20 ===\naccepting AMQP connection <0.6112.0> (10.0.0.1:41316 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251378035+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.227968151+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:49:20 ===\nConnection <0.6112.0> (10.0.0.1:41316 -> 10.0.0.1:5672) has a client-provided name: nova-compute:7:8bf04dbd-4db4-4886-b16a-05db8a8c7ce1",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251378205+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.230082952+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:49:20 ===\naccepting AMQP connection <0.6123.0> (10.0.0.1:41318 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251378388+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.234426372+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:49:20 ===\nConnection <0.6123.0> (10.0.0.1:41318 -> 10.0.0.1:5672) has a client-provided name: nova-conductor:23:4f4cc978-87e4-4fea-88fa-2d13147e857b",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251378559+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.236259963+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:49:20 ===\naccepting AMQP connection <0.6150.0> (10.0.0.1:41326 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251378729+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.239823975+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:49:20 ===\nConnection <0.6150.0> (10.0.0.1:41326 -> 10.0.0.1:5672) has a client-provided name: nova-conductor:21:86291add-c380-4101-baa5-10e6644b8956",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251378906+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.241690078+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:49:20 ===\naccepting AMQP connection <0.6161.0> (10.0.0.1:41338 -> 10.0.0.1:5672)",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251379093+01:00"
    },
    {
        "timestamp": "2016-12-28T12:17:22.246863594+01:00",
        "namespace": "/intel/logs/rabbit_logs/rabbit.log/message",
        "data": "=INFO REPORT==== 7-Dec-2016::03:49:20 ===\nConnection <0.6161.0> (10.0.0.1:41338 -> 10.0.0.1:5672) has a client-provided name: nova-conductor:22:569628a7-f04a-469a-9b99-3cdb13572540",
        "unit": "",
        "tags": {
            "plugin_running_on": "mkleina-dev"
        },
        "version": 0,
        "last_advertised_time": "2016-12-28T12:17:22.251379274+01:00"
    }
]
```

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
