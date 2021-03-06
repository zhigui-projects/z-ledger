/*
Copyright Zhigui.com. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ledgerconfig

import "github.com/spf13/viper"

const confArchiveEnable = "ledger.archive.enable"

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
