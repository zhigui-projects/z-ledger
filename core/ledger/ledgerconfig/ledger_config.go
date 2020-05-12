/*
Copyright Zhigui.com. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ledgerconfig

import "github.com/spf13/viper"

const confArchiveEnable = "ledger.archive.enable"
const confHDFSNameNodes = "ledger.archive.hdfs.nameNodes"
const confHDFSUser = "ledger.archive.hdfs.user"
const confUseDatanodeHostname = "ledger.archive.hdfs.useDatanodeHostname"

//GetDefaultMaxBlockfileSize returns default max block file size
func GetDefaultMaxBlockfileSize() int {
	//TODO: [maxpeng] change back after all tasks has been finished
	//return 64 * 1024 * 1024 // 64MB
	return 40 * 1024 // 40KB for testing
}

//IsArchiveEnabled exposes the archiveEnabled variable
func IsArchiveEnabled() bool {
	if viper.IsSet(confArchiveEnable) {
		return viper.GetBool(confArchiveEnable)
	}
	return false
}

func GetHDFSNameNodes() []string {
	if viper.IsSet(confHDFSNameNodes) {
		return viper.GetStringSlice(confHDFSNameNodes)
	}
	return []string{}
}

func GetHDFSUser() string {
	if viper.IsSet(confHDFSUser) {
		return viper.GetString(confHDFSUser)
	}
	return ""
}

func UseDatanodeHostname() bool {
	if viper.IsSet(confUseDatanodeHostname) {
		return viper.GetBool(confUseDatanodeHostname)
	}
	return false
}
