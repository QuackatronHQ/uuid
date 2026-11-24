// Copyright 2017 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !js

package uuid

import "net"

var interfaces []net.Interface // cached list of interfaces, protected by nodeMu

// getHardwareInterface returns the name and hardware address of interface name.
// If name is "" then the name and hardware address of one of the system's
// interfaces is returned.  If no interfaces are found (name does not exist or
// there are no interfaces) then "", nil is returned.
//
// Only addresses of at least 6 bytes are returned. Interfaces that are down
// or are loopback devices are skipped, since their addresses are often
// shared or zero.
//
// The caller must hold nodeMu.
func getHardwareInterface(name string) (string, []byte) {
	if interfaces == nil {
		var err error
		interfaces, err = net.Interfaces()
		if err != nil {
			return "", nil
		}
	}
	for _, ifs := range interfaces {
		if name == "" && (ifs.Flags&net.FlagUp == 0 || ifs.Flags&net.FlagLoopback != 0) {
			continue
		}
		if len(ifs.HardwareAddr) >= 6 && (name == "" || name == ifs.Name) {
			return ifs.Name, ifs.HardwareAddr
		}
	}
	return "", nil
}
