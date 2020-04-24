/*
Copyright Zhigui.com. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ledgerconfig

import "github.com/spf13/viper"

const confArchiveEnable = "ledger.archive.enable"
const confHDFSNameNodes = "ledger.archive.hdfs.nameNodes"
const confHDFSUser = "ledger.archive.hdfs.user"

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
