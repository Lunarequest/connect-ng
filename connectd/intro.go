package main

import "github.com/godbus/dbus/v5/introspect"

const intro = `
<node>
	<interface name="com.suse.Connect">
		<method name="AnnounceSystem">
			<arg direction="in" type="s" />
			<arg direction="in" type="s" />
			<arg direction="out" type="s" />
		</method>
		<method name="UpdateSystem">
			<arg direction="in" type="s" />
			<arg direction="in" type="s" />
			<arg direction="out" type="s" />
		</method>
		<method name="DeactivateSytem">
			<arg direction="in" type="s" />
			<arg direction="out" type="s" />
		</method>
		<method name="Credentials">
			<arg direction="in" type="s"/>
			<arg direction="out" type="s"/>
		</method>
		<method name="CreateCredentialsFile">
			<arg direction="in" type="s"/>
			<arg direction="in" type="s"/>
			<arg direction="in" type="s"/>
			<arg direction="in" type="s"/>
			<arg direction="out" type="s"/>
		</method>
		<method name="CurlrcCredentials">
			<arg direction="out" type="s"/>
		</method>
		<method name="ShowProduct">
			<arg direction="in" type="s"/>
			<arg direction="in" type="s"/>
			<arg direction="out" type="s"/>
		</method>
		<method name="ActivateProduct">
			<arg direction="in" type="s"/>
			<arg direction="in" type="s"/>
			<arg direction="in" type="s"/>
			<arg direction="out" type="s"/>
		</method>
		<method name="ActivatedProducts">
			<arg direction="in" type="s"/>
			<arg direction="out" type="s"/>
		</method>
		<method name="DeactivateProduct">
			<arg direction="in" type="s"/>
			<arg direction="in" type="s"/>
			<arg direction="out" type="s"/>
		</method>
		<method name="GetConfig">
			<arg direction="in" type="s"/>
			<arg direction="out" type="s"/>
		</method>
		<method name="WriteConfig">
			<arg direction="in" type="s"/>
			<arg direction="out" type="s"/>
		</method>
		<method name="Status">
			<arg direction="in" type="s" />
			<arg direction="out" type="s" />
		</method>
		<method name="UpdateCertificates">
			<arg direction="out" type="s" />
		</method>
		<method name="ReloadCertificates">
			<arg direction="out" type="s" />
		</method>
		<method name="ListInstallerUpdates">
			<arg direction="in" type="s" />
			<arg direction="in" type="s" />
			<arg direction="out" type="s" />
		</method>
		<method name="SystemMigrations">
			<arg direction="in" type="s" />
			<arg direction="in" type="s" />
			<arg direction="out" type="s" />
		</method>
		<method name="OFFlineSystemMigration">
			<arg direction="in" type="s" />
			<arg direction="in" type="s" />
			<arg direction="in" type="s" />
			<arg direction="out" type="s" />
		</method>
		<method name="UpdateProduct">
			<arg direction="in" type="s" />
			<arg direction="in" type="s" />
			<arg direction="out" type="s" />
		</method>
		<method name="Synchronize">
			<arg direction="in" type="s" />
			<arg direction="in" type="s" />
			<arg direction="out" type="s" />
		</method>
		<method name="SearchPackage">
			<arg direction="in" type="s" />
			<arg direction="out" type="s" />
		</method>
		<method name="OFFlineSystemMigration">
			<arg direction="in" type="s" />
			<arg direction="in" type="s" />
			<arg direction="in" type="s" />
			<arg direction="out" type="s" />
		</method>
		<method name="Version">
			<arg direction="in" type="b"/>
			<arg direction="out" type="s"/>
		</method>
	</interface>` + introspect.IntrospectDataString + `</node>`
