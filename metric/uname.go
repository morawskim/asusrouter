package metric

import (
	"os"
	"os/exec"
	"strings"
)

type UnameInfo struct {
	Domainname string
	Nodename   string
	Release    string
	Sysname    string
	Version    string
	Machine    string
}

type unameInfoResult struct {
	error error
	UnameInfo
}

func (u *unameInfoResult) process(bytes []byte, err error) string {
	//if error already exists return
	if u.error != nil {
		return ""
	}

	//if error is not nil
	if err != nil {
		u.error = err
		return ""
	}

	return strings.TrimSpace(string(bytes))
}

// Based on https://github.com/openwrt/packages/blob/11d3b495147ca02eab62e6507378347e5d824af8/utils/prometheus-node-exporter-lua/files/usr/lib/lua/prometheus-collectors/uname.lua
func (u *unameInfoResult) parseFiles() {
	var bytes []byte
	var err error

	bytes, err = os.ReadFile("/proc/sys/kernel/osrelease")
	u.Release = u.process(bytes, err)

	bytes, err = os.ReadFile("/proc/sys/kernel/ostype")
	u.Sysname = u.process(bytes, err)

	bytes, err = os.ReadFile("/proc/sys/kernel/version")
	u.Version = u.process(bytes, err)

	bytes, err = os.ReadFile("/proc/sys/kernel/domainname")
	u.Domainname = u.process(bytes, err)

	bytes, err = os.ReadFile("/proc/sys/kernel/hostname")
	u.Nodename = u.process(bytes, err)

	cmd := exec.Command("uname", "-m")
	bytes, err = cmd.Output()
	u.Machine = u.process(bytes, err)
}

func UnameMetrics() (*UnameInfo, error) {
	result := unameInfoResult{}
	result.parseFiles()

	if result.error != nil {
		return nil, result.error
	}

	return &result.UnameInfo, nil
}
