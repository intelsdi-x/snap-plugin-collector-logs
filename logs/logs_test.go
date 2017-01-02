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
	"io/ioutil"
	"os"
	"testing"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
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
		So(func() { filterFiles(".", "") }, ShouldNotPanic)
	})

	Convey("should list only matching files", t, func() {
		list := filterFiles("testdir", ".*2|3")
		So(list, ShouldResemble, []string{"testdir/testfile2", "testdir/testfile3"})
	})

	os.Remove("testdir/testfile1")
	os.Remove("testdir/testfile2")
	os.Remove("testdir/testfile3")
	os.Remove("testdir")
}

func TestLoadConfig(t *testing.T) {
	cfg := make(plugin.Config)
	cfg["metric_name"] = "all"
	cfg["log_dir"] = "/var/log"
	cfg["log_file"] = ".*"
	cfg["splitter_type"] = "new-line"
	cfg["splitter"] = splitterTypes["apache"]
	cfg["cache_dir"] = "/var/cache/snap"
	cfg["scanning_dir_counter"] = int64(2)
	cfg["collection_time"] = int64(321)

	cfgBad1 := make(plugin.Config)
	cfgBad1["splitter_type"] = "new-line"
	cfgBad1["collection_time"] = "abcd"

	cfgBad2 := make(plugin.Config)
	cfgBad2["splitter_type"] = "abcd"
	cfgBad2["collection_time"] = "300ms"

	cfgBad3 := make(plugin.Config)
	cfgBad3["splitter_type"] = "custom"
	cfgBad3["splitter"] = "bad(splitter"
	cfgBad3["collection_time"] = "300ms"

	cfgBad4 := make(plugin.Config)
	cfgBad4["splitter_type"] = "new-line"
	cfgBad4["collection_time"] = "300ms"
	cfgBad4["log_file"] = "log(file"

	Convey("should not panic", t, func() {
		So(func() {
			l := Logs{}
			l.loadConfig(cfg)
		}, ShouldNotPanic)
	})

	Convey("should load config properly", t, func() {
		l := Logs{}
		l.loadConfig(cfg)

		So(len(l.configInt), ShouldEqual, 2)
		So(len(l.configStr), ShouldEqual, 6)
	})

	Convey("should return error on invalid collection_time value", t, func() {
		l := Logs{}
		err := l.loadConfig(cfgBad1)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldResemble, "collection time value (collection_time) is invalid")
	})

	Convey("should return error on invalid splitter_type value", t, func() {
		l := Logs{}
		err := l.loadConfig(cfgBad2)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldResemble, "splitter type \"abcd\" is not supported")
	})

	Convey("should return error on invalid splitter value", t, func() {
		l := Logs{}
		err := l.loadConfig(cfgBad3)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldResemble, "splitter value is invalid")
	})

	Convey("should return error on invalid log_file value", t, func() {
		l := Logs{}
		err := l.loadConfig(cfgBad4)

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldResemble, "log_file value is invalid")
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

	cfg := make(plugin.Config)
	cfg["log_dir"] = "./(*dira|*dirb)"
	cfg["log_file"] = ".*file(2|3)"
	cfg["splitter_type"] = "new-line"

	Convey("should not panic", t, func() {
		So(
			func() {
				l := Logs{}
				l.loadConfig(cfg)
				l.refreshLogList()
			},
			ShouldNotPanic)
	})

	Convey("should list only matching logs", t, func() {
		l := Logs{}
		l.loadConfig(cfg)
		l.refreshLogList()
		So(logFiles, ShouldResemble, []string{"logdira/testfile2", "logdira/testfile3", "logdirb/testfile2", "logdirb/testfile3"})
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
	cfg := make(plugin.Config)
	cfg["log_dir"] = "logdir"
	cfg["log_file"] = "testfile"
	cfg["splitter_type"] = "new-line"
	cfg["collection_time"] = "300ms"

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
		So(mt[0].Namespace.Strings(), ShouldResemble, []string{"intel", "logs", "*", "*", "message"})
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

	cfgApache := make(plugin.Config)
	cfgApache["log_dir"] = "logdir"
	cfgApache["log_file"] = "testapache.log"
	cfgApache["splitter_type"] = "new-line"
	cfgApache["cache_dir"] = "logcache"
	cfgApache["metric_name"] = "nova"
	cfgApache["collection_time"] = "300ms"
	cfgApache["scanning_dir_counter"] = int64(0)

	cfgRabbit := make(plugin.Config)
	cfgRabbit["log_dir"] = "logdir"
	cfgRabbit["log_file"] = "testrabbit.log"
	cfgRabbit["splitter_type"] = "empty-line"
	cfgRabbit["cache_dir"] = "logcache"
	cfgRabbit["metric_name"] = "rabbitmq"
	cfgRabbit["collection_time"] = "300ms"
	cfgRabbit["scanning_dir_counter"] = int64(3)

	mtsApache := []plugin.Metric{
		plugin.Metric{
			Namespace: plugin.NewNamespace("intel", "logs").AddDynamicElement("metric_name", "Metric name defined in config file").
				AddDynamicElement("log_file", "Log file name").AddStaticElement("message"),
			Config: cfgApache,
		},
	}
	mtsApache[0].Namespace[2].Value = "nova"
	mtsApache[0].Namespace[3].Value = "testapache.log"

	mtsRabbit := []plugin.Metric{
		plugin.Metric{
			Namespace: plugin.NewNamespace("intel", "logs").AddDynamicElement("metric_name", "Metric name defined in config file").
				AddDynamicElement("log_file", "Log file name").AddStaticElement("message"),
			Config: cfgRabbit,
		},
	}
	mtsRabbit[0].Namespace[2].Value = "rabbitmq"
	mtsRabbit[0].Namespace[3].Value = "testrabbit.log"

	Convey("should not panic", t, func() {
		So(func() {
			l := Logs{}
			l.CollectMetrics(append(mtsApache, mtsRabbit...))
		}, ShouldNotPanic)
	})

	Convey("should collect valid metrics and create valid cache file (Apache)", t, func() {
		os.Remove("logcache/nova_testapache.log.json")

		l := Logs{}
		m, err := l.CollectMetrics(mtsApache)
		So(err, ShouldBeNil)
		So(len(m), ShouldEqual, 15)

		allData := ""
		for _, v := range m {
			So(v.Namespace[2].IsDynamic(), ShouldBeTrue)
			So(v.Namespace[2].Description, ShouldEqual, "Metric name defined in config file")
			So(v.Namespace[3].IsDynamic(), ShouldBeTrue)
			So(v.Namespace[3].Description, ShouldEqual, "Log file name")
			So(v.Namespace[4].IsDynamic(), ShouldBeFalse)
			allData += v.Data.(string) + "\n"
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
		So(logFiles, ShouldBeEmpty)
	})

	Convey("should collect valid metrics and create valid cache file (Rabbit)", t, func() {
		os.Remove("logcache/rabbitmq_testrabbit.log.json")

		l := Logs{}
		m, err := l.CollectMetrics(mtsRabbit)
		So(err, ShouldBeNil)
		So(len(m), ShouldEqual, 12)

		allData := ""
		for _, v := range m {
			So(v.Namespace[2].IsDynamic(), ShouldBeTrue)
			So(v.Namespace[2].Description, ShouldEqual, "Metric name defined in config file")
			So(v.Namespace[3].IsDynamic(), ShouldBeTrue)
			So(v.Namespace[3].Description, ShouldEqual, "Log file name")
			So(v.Namespace[4].IsDynamic(), ShouldBeFalse)
			allData += "\n" + v.Data.(string) + "\n"
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
		So(logFiles, ShouldNotBeEmpty)
		l.CollectMetrics(mtsRabbit)
		So(logFiles, ShouldNotBeEmpty)
		l.CollectMetrics(mtsRabbit)
		So(logFiles, ShouldNotBeEmpty)
		l.CollectMetrics(mtsRabbit)
		So(logFiles, ShouldBeEmpty) // <- 4th collection - list should be updated now
	})

	os.Remove("logdir/testapache.log")
	os.Remove("logdir/testrabbit.log")
	os.Remove("logdir")
	os.Remove("logcache/nova_testapache.log.json")
	os.Remove("logcache/rabbitmq_testrabbit.log.json")
	os.Remove("logcache")
}
