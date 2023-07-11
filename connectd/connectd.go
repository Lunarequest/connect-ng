package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/SUSE/connect-ng/internal/connect"
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)

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

func (f Connectd) AnnounceSystem(clientParams, distroTarget string) string {
	loadConfig(clientParams)

	login, password, err := connect.AnnounceSystem(distroTarget, "")
	if err != nil {
		return errorToJSON(err)
	}
	var res struct {
		Credentials []string `json:"credentials"`
	}
	res.Credentials = []string{login, password, ""}
	jsn, _ := json.Marshal(&res)

	return string(jsn)
}

func (f Connectd) UpdateSystem(clientParams, distroTarget string) string {
	loadConfig(clientParams)

	if err := connect.UpdateSystem(distroTarget, ""); err != nil {
		return errorToJSON(err)
	}
	return "{}"
}

func (f Connectd) DeactivateSytem(clientParams string) string {
	loadConfig(clientParams)

	if err := connect.DeregisterSystem(); err != nil {
		return errorToJSON(err)
	}
	return "{}"
}

func Credentials(path string) string {
	creds, err := connect.ReadCredentials(path)
	if err != nil {
		return errorToJSON(err)
	}

	jsn, _ := json.Marshal(creds)

	return string(jsn)
}

func (f Connectd) CreateCredentialsFile(login, password, token, path string) string {
	if err:=connect.CreateCredentials(login, password, token, path); err!=nil {
		return errorToJSON(err)
	}
	return "{}"
}

func (f Connectd) CurlrcCredentials() string {
	creds, _ := connect.ReadCurlrcCredentials()
	jsn, _:= json.Marshal(creds)
	return string(jsn)
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

	if reply != dbus.RequestNameReplyPrimaryOwner {
		fmt.Fprintln(os.Stderr, "name already taken")
		os.Exit(1)
	}
	fmt.Println("Listening on com.suse.Connect /com/suse/Connect ...")
	select {}
}
