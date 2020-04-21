/*
Copyright Zhigui.com. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ledgerconfig

import "github.com/spf13/viper"

const confArchiveEnable = "ledger.archive.enable"

//IsArchiveEnabled exposes the archiveEnabled variable
func IsArchiveEnabled() bool {
	if viper.IsSet(confArchiveEnable) {
		return viper.GetBool(confArchiveEnable)
	}
	return false
}
