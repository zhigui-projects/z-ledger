/*
Copyright Zhigui.com. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package stateormdb

import (
	"encoding/base64"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb/stateormdb/msgs"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/version"
)

func encodeVersionAndMetadata(version *version.Height, metadata []byte) (string, error) {
	msg := &msgs.ORMVersionFieldProto{
		VersionBytes: version.ToBytes(),
		Metadata:     metadata,
	}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(msgBytes), nil
}

func DecodeVersionAndMetadata(encodedstr string) (*version.Height, []byte, error) {
	versionFieldBytes, err := base64.StdEncoding.DecodeString(encodedstr)
	if err != nil {
		return nil, nil, err
	}
	versionFieldMsg := &msgs.ORMVersionFieldProto{}
	if err = proto.Unmarshal(versionFieldBytes, versionFieldMsg); err != nil {
		return nil, nil, err
	}
	ver, _, err := version.NewHeightFromBytes(versionFieldMsg.VersionBytes)
	if err != nil {
		return nil, nil, err
	}
	return ver, versionFieldMsg.Metadata, nil
}
