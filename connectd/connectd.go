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
	if err := connect.CreateCredentials(login, password, token, path); err != nil {
		return errorToJSON(err)
	}
	return "{}"
}

func (f Connectd) CurlrcCredentials() string {
	creds, _ := connect.ReadCurlrcCredentials()
	jsn, _ := json.Marshal(creds)
	return string(jsn)
}

func (f Connectd) ShowProduct(clientParams, product string) string {
	loadConfig(clientParams)

	var productQuery connect.Product
	if err := json.Unmarshal([]byte(product), &productQuery); err != nil {
		return errorToJSON(connect.JSONError{Err: err})
	}
	productData, err := connect.ShowProduct(productQuery)

	if err != nil {
		return errorToJSON(err)
	}
	jsn, err := json.Marshal(productData)
	if err != nil {
		return errorToJSON(err)
	}
	return string(jsn)
}

func (f Connectd) ActivateProduct(clientParams, product, email string) string {
	loadConfig(clientParams)

	var p connect.Product

	if err := json.Unmarshal([]byte(product), &p); err != nil {
		return errorToJSON(connect.JSONError{Err: err})
	}

	service, err := connect.ActivateProduct(p, email)
	if err != nil {
		return errorToJSON(err)
	}
	jsn, err := json.Marshal(service)
	if err != nil {
		return errorToJSON(err)
	}
	return string(jsn)
}

func ActivatedProducts(clientParams string) string {
	loadConfig(clientParams)
	products, err := connect.ActivatedProducts()
	if err != nil {
		return errorToJSON(err)
	}
	jsn, err := json.Marshal(products)
	if err != nil {
		return errorToJSON(err)
	}
	return string(jsn)
}

func (f Connectd) DeactivateProduct(clientParams, product string) string {
	loadConfig(clientParams)
	var p connect.Product

	if err := json.Unmarshal([]byte(product), &p); err != nil {
		return errorToJSON(connect.JSONError{Err: err})
	}

	service, err := connect.DeactivateProduct(p)
	if err != nil {
		return errorToJSON(err)
	}
	jsn, err := json.Marshal(service)
	if err != nil {
		return errorToJSON(err)
	}
	return string(jsn)
}

func (f Connectd) GetConfig(path string) string {
	c := connect.NewConfig()
	c.Path = path
	c.Load()
	jsn, err := json.Marshal(c)
	if err != nil {
		errorToJSON(err)
	}
	return string(jsn)
}

func (f Connectd) WriteConfig(clientParams string) string {
	loadConfig(clientParams)
	if err := connect.CFG.Save(); err != nil {
		return errorToJSON(err)
	}
	return "{}"
}

func (f Connectd) Status(format string) (string, *dbus.Error) {
	connect.CFG.Load()
	output, err := connect.GetProductStatuses(format)

	if err != nil {
		return output, dbus.MakeFailedError(err)
	}
	return output, nil
}

func (f Connectd) UpdateCertificates() string {
	if err := connect.UpdateCertificates(); err != nil {
		return errorToJSON(err)
	}
	return "{}"
}

func (f Connectd) ReloadCertificates() string {
	if err := connect.ReloadCertPool(); err != nil {
		return errorToJSON(err)
	}
	return ""
}

func (f Connectd) ListInstallerUpdates(clientParams, product string) string {
	loadConfig(clientParams)
	var productQuery connect.Product
	if err := json.Unmarshal([]byte(product), &productQuery); err != nil {
		return errorToJSON(connect.JSONError{Err: err})
	}
	repos, err := connect.InstallerUpdates(productQuery)
	if err != nil {
		return errorToJSON(err)
	}
	jsn, err := json.Marshal(repos)
	if err != nil {
		return errorToJSON(err)
	}
	return string(jsn)
}

func (f Connectd) SystemMigrations(clientParams, products string) string {
	loadConfig(clientParams)
	installed := make([]connect.Product, 0)
	if err := json.Unmarshal([]byte(products), &installed); err != nil {
		return errorToJSON(connect.JSONError{Err: err})
	}
	migrations, err := connect.ProductMigrations(installed)
	if err != nil {
		return errorToJSON(err)
	}
	jsn, err := json.Marshal(migrations)
	if err != nil {
		return errorToJSON(err)
	}
	return string(jsn)
}

func (f Connectd) OFFlineSystemMigration(clientParams, products, targetBaseProduct string) string {
	loadConfig(clientParams)
	installed := make([]connect.Product, 0)
	err := json.Unmarshal([]byte(products), &installed)
	if err != nil {
		return errorToJSON(connect.JSONError{Err: err})
	}
	var target connect.Product
	if err := json.Unmarshal([]byte(targetBaseProduct), &target); err != nil {
		return errorToJSON(connect.JSONError{Err: err})
	}
	migrations, err := connect.OfflineProductMigrations(installed, target)
	if err != nil {
		return errorToJSON(err)
	}
	jsn, err := json.Marshal(migrations)
	if err != nil {
		return errorToJSON(err)
	}
	return string(jsn)
}

func (f Connectd) UpdateProduct(clientParams, product string) string {
	loadConfig(clientParams)

	var prod connect.Product
	err := json.Unmarshal([]byte(product), &prod)
	if err != nil {
		return errorToJSON(connect.JSONError{Err: err})
	}
	service, err := connect.UpgradeProduct(prod)
	if err != nil {
		return errorToJSON(err)
	}
	jsn, err := json.Marshal(service)
	if err != nil {
		return errorToJSON(err)
	}
	return string(jsn)
}

func (f Connectd) Synchronize(clientParams, products string) string {
	loadConfig(clientParams)

	prods := make([]connect.Product, 0)
	err := json.Unmarshal([]byte(products), &prods)
	if err != nil {
		return errorToJSON(connect.JSONError{Err: err})
	}
	activated, err := connect.SyncProducts(prods)
	if err != nil {
		return errorToJSON(err)
	}
	jsn, err := json.Marshal(activated)
	if err != nil {
		return errorToJSON(err)
	}
	return string(jsn)
}

func (f Connectd) SystemActivations(clientParams string) string {
	loadConfig(clientParams)

	// converting from map to list as expected by Ruby clients
	actList := make([]connect.Activation, 0)
	actMap, err := connect.SystemActivations()
	if err != nil {
		return errorToJSON(err)
	}
	for _, a := range actMap {
		actList = append(actList, a)
	}
	jsn, err := json.Marshal(actList)
	if err != nil {
		return errorToJSON(err)
	}
	return string(jsn)
}

func (f Connectd) SearchPackage(clientParams, product, query string) string {
	loadConfig(clientParams)

	var p connect.Product
	err := json.Unmarshal([]byte(product), &p)
	if err != nil {
		return errorToJSON(connect.JSONError{Err: err})
	}
	results, err := connect.SearchPackage(query, p)
	if err != nil {
		return errorToJSON(err)
	}
	jsn, err := json.Marshal(results)
	if err != nil {
		return errorToJSON(err)
	}
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
