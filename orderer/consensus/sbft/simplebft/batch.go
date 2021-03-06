/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package simplebft

import (
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
	sb "github.com/hyperledger/fabric/protos/orderer/sbft"
	"github.com/hyperledger/fabric/protoutil"
)

func (s *SBFT) makeBatch(seq uint64, prevHash []byte, data [][]byte) *sb.Batch {
	//datahash := merkleHashData(data)

	block := protoutil.UnmarshalBlockOrPanic(data[0])
	batchhead := &sb.BatchHeader{
		Seq:      seq,
		PrevHash: prevHash,
		DataHash: block.GetHeader().DataHash,
	}
	rawHeader := protoutil.MarshalOrPanic(batchhead)
	return &sb.Batch{
		Header:   rawHeader,
		Payloads: data,
	}
}

func (s *SBFT) checkBatch(b *sb.Batch, checkData bool, needSigs bool) (*sb.BatchHeader, error) {
	batchheader := &sb.BatchHeader{}
	err := proto.Unmarshal(b.Header, batchheader)
	if err != nil {
		return nil, err
	}

	if checkData {
		//datahash := merkleHashData(b.Payloads)
		block := protoutil.UnmarshalBlockOrPanic(b.Payloads[0])
		datahash := block.GetHeader().DataHash
		if !reflect.DeepEqual(datahash, batchheader.DataHash) {
			return nil, fmt.Errorf("malformed blocks: invalid hash")
		}
	}

	if batchheader.PrevHash == nil {
		// TODO check against root hash, which should be part of constructor
	} else if needSigs {
		if len(b.Signatures) < s.oneCorrectQuorum() {
			return nil, fmt.Errorf("insufficient number of signatures on blocks: need %d, got %d", s.oneCorrectQuorum(), len(b.Signatures))
		}
	}

	bh := b.Hash()
	for r, sig := range b.Signatures {
		err = s.sys.CheckSig(bh, r, sig)
		if err != nil {
			return nil, err
		}
	}

	return batchheader, nil
}
