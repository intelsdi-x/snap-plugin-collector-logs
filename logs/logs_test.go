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
	"io/ioutil"
	"os"
	"testing"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/ctypes"
	. "github.com/smartystreets/goconvey/convey"
)

const logFileContentApache = `2016-12-06 09:21:08.341 6 INFO oslo_service.service [req-cb760354-bbb0-4968-92e6-3312b8a7d223 - - - - -] Starting 5 workers
2016-12-06 09:21:08.349 6 INFO nova.network.driver [req-cb760354-bbb0-4968-92e6-3312b8a7d223 - - - - -] Loading network driver 'nova.network.linux_net'
2016-12-06 09:21:08.488 20 INFO nova.osapi_compute.wsgi.server [req-67440a41-6667-4e07-b546-fa336ab5c3af - - - - -] (20) wsgi starting up on http://10.0.0.1:8774
2016-12-06 09:21:08.493 18 INFO nova.osapi_compute.wsgi.server [req-de1ec4f5-ddb9-4726-b03b-900cd78b03ea - - - - -] (18) wsgi starting up on http://10.0.0.1:8774
2016-12-06 09:21:08.494 22 INFO nova.osapi_compute.wsgi.server [req-163c5cdb-e764-4838-93f5-8fd56b4942b2 - - - - -] (22) wsgi starting up on http://10.0.0.1:8774
2016-12-06 09:21:08.499 19 INFO nova.osapi_compute.wsgi.server [req-ee7ac782-32aa-492c-aa48-3d84825814e7 - - - - -] (19) wsgi starting up on http://10.0.0.1:8774
2016-12-06 09:21:08.500 21 INFO nova.osapi_compute.wsgi.server [req-a253f557-bee0-4d13-933c-e690a9a4c69f - - - - -] (21) wsgi starting up on http://10.0.0.1:8774
2016-12-06 09:21:09.478 6 INFO nova.wsgi [req-cb760354-bbb0-4968-92e6-3312b8a7d223 - - - - -] metadata listening on 10.0.0.1:8775
2016-12-06 09:21:09.478 6 INFO oslo_service.service [req-cb760354-bbb0-4968-92e6-3312b8a7d223 - - - - -] Starting 5 workers
2016-12-06 09:21:09.585 29 INFO nova.metadata.wsgi.server [req-8354a5d5-3c74-4ab3-826a-a69f3736f488 - - - - -] (29) wsgi starting up on http://10.0.0.1:8775
2016-12-06 09:21:09.589 30 INFO nova.metadata.wsgi.server [req-7921f263-7839-4cce-957b-f2ecc2e8814e - - - - -] (30) wsgi starting up on http://10.0.0.1:8775
2016-12-06 09:21:09.603 31 INFO nova.metadata.wsgi.server [req-773f91ec-7d02-44e1-aa0d-1cf056b799f9 - - - - -] (31) wsgi starting up on http://10.0.0.1:8775
2016-12-06 09:21:09.616 32 INFO nova.metadata.wsgi.server [req-5936911e-68e5-432d-93b6-b68c89cff3ea - - - - -] (32) wsgi starting up on http://10.0.0.1:8775
2016-12-06 09:21:09.640 33 INFO nova.metadata.wsgi.server [req-a12a4378-da99-49a7-86a5-89be0f2cee60 - - - - -] (33) wsgi starting up on http://10.0.0.1:8775
2016-12-07 03:39:17.960 18 INFO nova.wsgi [-] Stopping WSGI server.
`

const logFileContentRabbit = `
=INFO REPORT==== 7-Dec-2016::03:48:22 ===
accepting AMQP connection <0.5567.0> (10.0.0.1:40198 -> 10.0.0.1:5672)

=INFO REPORT==== 7-Dec-2016::03:48:22 ===
Connection <0.5567.0> (10.0.0.1:40198 -> 10.0.0.1:5672) has a client-provided name: nova-scheduler:6:69878ca7-bf1d-47e5-8d98-9614c0ece4c6

=INFO REPORT==== 7-Dec-2016::03:49:15 ===
accepting AMQP connection <0.6031.0> (10.0.0.1:41234 -> 10.0.0.1:5672)

=INFO REPORT==== 7-Dec-2016::03:49:15 ===
Connection <0.6031.0> (10.0.0.1:41234 -> 10.0.0.1:5672) has a client-provided name: neutron-server:26:0edf54c9-8868-4db1-bb40-60fc2b85539f

=INFO REPORT==== 7-Dec-2016::03:49:15 ===
accepting AMQP connection <0.6042.0> (10.0.0.1:41240 -> 10.0.0.1:5672)

=INFO REPORT==== 7-Dec-2016::03:49:15 ===
Connection <0.6042.0> (10.0.0.1:41240 -> 10.0.0.1:5672) has a client-provided name: neutron-server:25:317888f2-4539-4451-b2f9-3c26fdcb1c5c

=WARNING REPORT==== 7-Dec-2016::03:49:18 ===
closing AMQP connection <0.5006.0> (10.0.0.1:40004 -> 10.0.0.1:5672 - nova-compute:7:dc749d8b-15ea-4ae3-b524-22b29b810880):
client unexpectedly closed TCP connection

=WARNING REPORT==== 7-Dec-2016::03:49:18 ===
closing AMQP connection <0.4881.0> (10.0.0.1:40000 -> 10.0.0.1:5672 - nova-compute:7:89ae9d72-d74f-4509-9896-cc7e7b388b3b):
client unexpectedly closed TCP connection

=INFO REPORT==== 7-Dec-2016::03:49:20 ===
accepting AMQP connection <0.6091.0> (10.0.0.1:41314 -> 10.0.0.1:5672)

=INFO REPORT==== 7-Dec-2016::03:49:20 ===
Connection <0.6091.0> (10.0.0.1:41314 -> 10.0.0.1:5672) has a client-provided name: nova-compute:7:d377fb4a-b463-40c6-9049-8f341ffdac44

=INFO REPORT==== 7-Dec-2016::03:49:20 ===
accepting AMQP connection <0.6112.0> (10.0.0.1:41316 -> 10.0.0.1:5672)

=INFO REPORT==== 7-Dec-2016::03:49:20 ===
Connection <0.6112.0> (10.0.0.1:41316 -> 10.0.0.1:5672) has a client-provided name: nova-compute:7:8bf04dbd-4db4-4886-b16a-05db8a8c7ce1
`

func TestIsDir(t *testing.T) {
	os.Mkdir("testdir", 0755)
	os.OpenFile("testfile", os.O_CREATE, 0644)

	Convey("should not panic", t, func() {
		So(func() { isDir("testdir") }, ShouldNotPanic)
	})
	Convey("should return true for directory", t, func() {
		result, err := isDir("testdir")
		So(err, ShouldBeNil)
		So(result, ShouldBeTrue)
	})
	Convey("should return false for not existing item", t, func() {
		result, err := isDir("notexisting")
		So(err, ShouldNotBeNil)
		So(result, ShouldBeFalse)
	})
	Convey("should return false for file", t, func() {
		result, err := isDir("testfile")
		So(err, ShouldBeNil)
		So(result, ShouldBeFalse)
	})

	os.Remove("testdir")
	os.Remove("testfile")
}

func TestExpandPaths(t *testing.T) {
	os.MkdirAll("testdir1/testdir11", 0755)
	os.MkdirAll("testdir1/testdir12", 0755)
	os.MkdirAll("testdir1/testdir13", 0755)
	os.Mkdir("testdir2", 0755)
	os.Mkdir("testdir3", 0755)
	os.OpenFile("testfile1", os.O_CREATE, 0644)
	os.OpenFile("testfile2", os.O_CREATE, 0644)

	Convey("should not panic", t, func() {
		So(func() { expandPaths("", &[]string{}) }, ShouldNotPanic)
	})

	Convey("should list only matching directories (1)", t, func() {
		list := []string{}
		expandPaths("*", &list)
		So(list, ShouldResemble, []string{"testdir1", "testdir2", "testdir3"})
	})

	Convey("should list only matching directories (2)", t, func() {
		list := []string{}
		expandPaths("*/testdir11", &list)
		So(list, ShouldResemble, []string{"testdir1/testdir11"})
	})

	Convey("should list only matching directories (3)", t, func() {
		list := []string{}
		expandPaths("*/{testdir11,testdir12}", &list)
		So(list, ShouldResemble, []string{"testdir1/testdir11", "testdir1/testdir12"})

		list = []string{}
		expandPaths("*/(testdir11,testdir12)", &list)
		So(list, ShouldResemble, []string{"testdir1/testdir11", "testdir1/testdir12"})

		list = []string{}
		expandPaths("*/{testdir11|testdir12}", &list)
		So(list, ShouldResemble, []string{"testdir1/testdir11", "testdir1/testdir12"})

		list = []string{}
		expandPaths("*/(testdir11|testdir12)", &list)
		So(list, ShouldResemble, []string{"testdir1/testdir11", "testdir1/testdir12"})
	})

	os.Remove("testdir1/testdir11")
	os.Remove("testdir1/testdir12")
	os.Remove("testdir1/testdir13")
	os.Remove("testdir1")
	os.Remove("testdir2")
	os.Remove("testdir3")
	os.Remove("testfile1")
	os.Remove("testfile2")
}

func TestFilterFiles(t *testing.T) {
	os.Mkdir("testdir", 0755)
	os.OpenFile("testdir/testfile1", os.O_CREATE, 0644)
	os.OpenFile("testdir/testfile2", os.O_CREATE, 0644)
	os.OpenFile("testdir/testfile3", os.O_CREATE, 0644)

	Convey("should not panic", t, func() {
		So(func() { filterFiles(".", "", &[]string{}) }, ShouldNotPanic)
	})

	Convey("should list only matching files", t, func() {
		list := []string{}
		filterFiles("testdir", ".*2|3", &list)
		So(list, ShouldResemble, []string{"testdir/testfile2", "testdir/testfile3"})
	})

	os.Remove("testdir/testfile1")
	os.Remove("testdir/testfile2")
	os.Remove("testdir/testfile3")
	os.Remove("testdir")
}

func TestLoadConfig(t *testing.T) {
	cfg := plugin.NewPluginConfigType()
	cfg.AddItem("log_dir", ctypes.ConfigValueStr{Value: "logdir"})
	cfg.AddItem("log_file", ctypes.ConfigValueStr{Value: "testfile"})
	cfg.AddItem("log_type", ctypes.ConfigValueStr{Value: "apache"})

	Convey("should not panic", t, func() {
		So(func() {
			l := Logs{}
			l.loadConfig(cfg.Table())
		}, ShouldNotPanic)
	})

	Convey("should load config properly", t, func() {
		l := Logs{}
		l.loadConfig(cfg.Table())

		So(len(l.config), ShouldEqual, len(cfg.Table()))
		So(l.config, ShouldResemble, cfg.Table())
	})
}

func TestRefreshLogList(t *testing.T) {
	os.Mkdir("logdira", 0755)
	os.Mkdir("logdirb", 0755)
	os.Mkdir("logdirc", 0755)
	os.OpenFile("logdira/testfile1", os.O_CREATE, 0644)
	os.OpenFile("logdira/testfile2", os.O_CREATE, 0644)
	os.OpenFile("logdira/testfile3", os.O_CREATE, 0644)
	os.OpenFile("logdirb/testfile1", os.O_CREATE, 0644)
	os.OpenFile("logdirb/testfile2", os.O_CREATE, 0644)
	os.OpenFile("logdirb/testfile3", os.O_CREATE, 0644)
	os.OpenFile("logdirc/testfile1", os.O_CREATE, 0644)
	os.OpenFile("logdirc/testfile2", os.O_CREATE, 0644)
	os.OpenFile("logdirc/testfile3", os.O_CREATE, 0644)

	cfg := plugin.NewPluginConfigType()
	cfg.AddItem("log_dir", ctypes.ConfigValueStr{Value: "./(*dira|*dirb)"})
	cfg.AddItem("log_file", ctypes.ConfigValueStr{Value: ".*file(2|3)"})
	cfg.AddItem("log_type", ctypes.ConfigValueStr{Value: "apache"})

	Convey("should not panic", t, func() {
		So(
			func() {
				l := Logs{}
				l.loadConfig(cfg.Table())
				l.refreshLogList()
			},
			ShouldNotPanic)
	})

	Convey("should list only matching logs", t, func() {
		l := Logs{}
		l.loadConfig(cfg.Table())
		l.refreshLogList()
		So(l.logFiles, ShouldResemble, []string{"logdira/testfile2", "logdira/testfile3", "logdirb/testfile2", "logdirb/testfile3"})
	})

	os.Remove("logdira/testfile1")
	os.Remove("logdira/testfile2")
	os.Remove("logdira/testfile3")
	os.Remove("logdirb/testfile1")
	os.Remove("logdirb/testfile2")
	os.Remove("logdirb/testfile3")
	os.Remove("logdirc/testfile1")
	os.Remove("logdirc/testfile2")
	os.Remove("logdirc/testfile3")
	os.Remove("logdira")
	os.Remove("logdirb")
	os.Remove("logdirc")
}

func TestGetMetricTypes(t *testing.T) {
	cfg := plugin.NewPluginConfigType()
	cfg.AddItem("log_dir", ctypes.ConfigValueStr{Value: "logdir"})
	cfg.AddItem("log_file", ctypes.ConfigValueStr{Value: "testfile"})
	cfg.AddItem("log_type", ctypes.ConfigValueStr{Value: "apache"})

	Convey("should not panic", t, func() {
		So(func() {
			l := Logs{}
			l.GetMetricTypes(cfg)
		}, ShouldNotPanic)
	})

	Convey("should return valid metric type", t, func() {
		l := Logs{}
		mt, err := l.GetMetricTypes(cfg)
		So(err, ShouldBeNil)
		So(mt, ShouldNotBeEmpty)
		So(len(mt), ShouldEqual, 1)
		So(mt[0].Namespace().Strings(), ShouldResemble, []string{"intel", "logs", "*", "*", "message"})
	})
}

func TestCollectMetrics(t *testing.T) {
	os.Mkdir("logcache", 0755)
	os.Mkdir("logdir", 0755)
	if file, err := os.Create("logdir/testapache.log"); err != nil {
		t.Errorf("Test log creation error: %s", err)
		t.Fail()
	} else {
		file.WriteString(logFileContentApache)
		file.Close()
	}
	if file, err := os.Create("logdir/testrabbit.log"); err != nil {
		t.Errorf("Test log creation error: %s", err)
		t.Fail()
	} else {
		file.WriteString(logFileContentRabbit)
		file.Close()
	}

	cfgApache := plugin.NewPluginConfigType()
	cfgApache.AddItem("log_dir", ctypes.ConfigValueStr{Value: "logdir"})
	cfgApache.AddItem("log_file", ctypes.ConfigValueStr{Value: "testapache.log"})
	cfgApache.AddItem("log_type", ctypes.ConfigValueStr{Value: "apache"})
	cfgApache.AddItem("cache_dir", ctypes.ConfigValueStr{Value: "logcache"})
	cfgApache.AddItem("metric_name", ctypes.ConfigValueStr{Value: "nova"})
	cfgApache.AddItem("collection_time", ctypes.ConfigValueInt{Value: 300})
	cfgApache.AddItem("scanning_dir_counter", ctypes.ConfigValueInt{Value: 0})

	cfgRabbit := plugin.NewPluginConfigType()
	cfgRabbit.AddItem("log_dir", ctypes.ConfigValueStr{Value: "logdir"})
	cfgRabbit.AddItem("log_file", ctypes.ConfigValueStr{Value: "testrabbit.log"})
	cfgRabbit.AddItem("log_type", ctypes.ConfigValueStr{Value: "rabbit"})
	cfgRabbit.AddItem("cache_dir", ctypes.ConfigValueStr{Value: "logcache"})
	cfgRabbit.AddItem("metric_name", ctypes.ConfigValueStr{Value: "rabbitmq"})
	cfgRabbit.AddItem("collection_time", ctypes.ConfigValueInt{Value: 300})
	cfgRabbit.AddItem("scanning_dir_counter", ctypes.ConfigValueInt{Value: 3})

	mtsApache := []plugin.MetricType{
		plugin.MetricType{
			Namespace_: core.NewNamespace("intel", "logs", "nova", "testapache.log", "message"),
			Config_:    cfgApache.ConfigDataNode,
		},
	}
	mtsRabbit := []plugin.MetricType{
		plugin.MetricType{
			Namespace_: core.NewNamespace("intel", "logs", "rabbitmq", "testrabbit.log", "message"),
			Config_:    cfgRabbit.ConfigDataNode,
		},
	}

	Convey("should not panic", t, func() {
		So(func() {
			l := Logs{}
			l.CollectMetrics(append(mtsApache, mtsRabbit...))
		}, ShouldNotPanic)
	})

	Convey("should return error with empty metric list", t, func() {
		l := Logs{}
		_, err := l.CollectMetrics([]plugin.MetricType{})
		So(err, ShouldNotBeNil)
	})

	Convey("should collect valid metrics and create valid cache file (Apache)", t, func() {
		os.Remove("logcache/nova_testapache.log.json")

		l := Logs{}
		m, err := l.CollectMetrics(mtsApache)
		So(err, ShouldBeNil)
		So(len(m), ShouldEqual, 15)

		allData := ""
		for _, v := range m {
			allData += v.Data().(string) + "\n"
		}
		So(allData, ShouldEqual, logFileContentApache)

		positionCache := positionCache{}
		posData, err := ioutil.ReadFile("logcache/nova_testapache.log.json")
		So(err, ShouldBeNil)
		err = json.Unmarshal(posData, &positionCache)
		So(err, ShouldBeNil)
		So(positionCache.Offset, ShouldEqual, 2193)

		// Should refresh log files list after each collection cycle
		os.Remove("logdir/testapache.log")
		l.CollectMetrics(mtsApache)
		So(l.logFiles, ShouldBeEmpty)
	})

	Convey("should collect valid metrics and create valid cache file (Rabbit)", t, func() {
		os.Remove("logcache/rabbitmq_testrabbit.log.json")

		l := Logs{}
		m, err := l.CollectMetrics(mtsRabbit)
		So(err, ShouldBeNil)
		So(len(m), ShouldEqual, 12)

		allData := ""
		for _, v := range m {
			allData += "\n" + v.Data().(string) + "\n"
		}
		So(allData, ShouldEqual, logFileContentRabbit)

		positionCache := positionCache{}
		posData, err := ioutil.ReadFile("logcache/rabbitmq_testrabbit.log.json")
		So(err, ShouldBeNil)
		err = json.Unmarshal(posData, &positionCache)
		So(err, ShouldBeNil)
		So(positionCache.Offset, ShouldEqual, 1897)

		// Should refresh log files list after 3 collection cycles
		os.Remove("logdir/testrabbit.log")
		l.CollectMetrics(mtsRabbit)
		So(l.logFiles, ShouldNotBeEmpty)
		l.CollectMetrics(mtsRabbit)
		So(l.logFiles, ShouldNotBeEmpty)
		l.CollectMetrics(mtsRabbit)
		So(l.logFiles, ShouldNotBeEmpty)
		l.CollectMetrics(mtsRabbit)
		So(l.logFiles, ShouldBeEmpty) // <- 4th collection - list should be updated now
	})

	os.Remove("logdir")
	os.Remove("logcache/nova_testapache.log.json")
	os.Remove("logcache/rabbitmq_testrabbit.log.json")
	os.Remove("logcache")
}
