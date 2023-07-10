package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/SUSE/connect-ng/internal/connect"
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)


const intro = `
<node>
	<interface name="com.suse.Connect">
		<method name="Version">
			<arg direction="in" type="b"/>
			<arg direction="out" type="s"/>
		</method>
		<method name="Status">
			<arg direction="in" type="s" />
			<arg direction="out" type="s" />
		</method>
		<method name="DeactivateSytem">
			<arg direction="in" type="s" />
			<arg direction="out" type="s" />
		</method>
	</interface>` + introspect.IntrospectDataString + `</node>`

type Connectd string

func loadConfig(clientParams string) {
	// unmarshal extra config fields only for local use
	var extConfig struct {
		Debug string `json:"debug"`
	}
	json.Unmarshal([]byte(clientParams), &extConfig)
	connect.CFG.Load()
	connect.CFG.MergeJSON(clientParams)
}


func (f Connectd) Version(fullVersion bool) (string, *dbus.Error) {
	var version string
	if fullVersion {
		version = connect.GetFullVersion()
	} else {
		version = connect.GetShortenedVersion()
	}
	return version, nil
}

func (f Connectd) Status(format string) (string, *dbus.Error) {
	connect.CFG.Load()
	output, err := connect.GetProductStatuses(format)

	if err != nil {
		return output, dbus.MakeFailedError(err)
	}
	return output, nil
}

func DeactivateSytem(clientParams string) (string) {
	loadConfig(clientParams)
	err := connect.DeregisterSystem()

	if err != nil {
		return errorToJSON(err)
	}
	return "{}"
}

func main() {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		panic(err)
	}

	f := Connectd("Connectd")
	conn.Export(f, "/com/suse/Connect", "com.suse.Connect")
	conn.Export(introspect.Introspectable(intro), "/com/suse/Connect",
	"org.freedesktop.DBus.Introspectable")

	reply, err := conn.RequestName("com.suse.Connect", dbus.NameFlagDoNotQueue)

	if err != nil {
		panic(err)

	}

	if reply !=  dbus.RequestNameReplyPrimaryOwner {
		fmt.Fprintln(os.Stderr, "name already taken")
		os.Exit(1)
	}
	fmt.Println("Listening on com.suse.Connect /com/suse/Connect ...")
	select {}
}