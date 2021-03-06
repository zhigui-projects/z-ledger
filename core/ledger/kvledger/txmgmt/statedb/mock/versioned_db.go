// Code generated by counterfeiter. DO NOT EDIT.
package mock

import (
	"sync"

	"github.com/hyperledger/fabric-chaincode-go/shim/entitydefinition"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/version"
)

type VersionedDB struct {
	ApplyUpdatesStub        func(*statedb.UpdateBatch, *version.Height) error
	applyUpdatesMutex       sync.RWMutex
	applyUpdatesArgsForCall []struct {
		arg1 *statedb.UpdateBatch
		arg2 *version.Height
	}
	applyUpdatesReturns struct {
		result1 error
	}
	applyUpdatesReturnsOnCall map[int]struct {
		result1 error
	}
	BytesKeySupportedStub        func() bool
	bytesKeySupportedMutex       sync.RWMutex
	bytesKeySupportedArgsForCall []struct {
	}
	bytesKeySupportedReturns struct {
		result1 bool
	}
	bytesKeySupportedReturnsOnCall map[int]struct {
		result1 bool
	}
	CloseStub        func()
	closeMutex       sync.RWMutex
	closeArgsForCall []struct {
	}
	ExecuteConditionQueryStub        func(string, entitydefinition.Search) (interface{}, error)
	executeConditionQueryMutex       sync.RWMutex
	executeConditionQueryArgsForCall []struct {
		arg1 string
		arg2 entitydefinition.Search
	}
	executeConditionQueryReturns struct {
		result1 interface{}
		result2 error
	}
	executeConditionQueryReturnsOnCall map[int]struct {
		result1 interface{}
		result2 error
	}
	ExecuteQueryStub        func(string, string) (statedb.ResultsIterator, error)
	executeQueryMutex       sync.RWMutex
	executeQueryArgsForCall []struct {
		arg1 string
		arg2 string
	}
	executeQueryReturns struct {
		result1 statedb.ResultsIterator
		result2 error
	}
	executeQueryReturnsOnCall map[int]struct {
		result1 statedb.ResultsIterator
		result2 error
	}
	ExecuteQueryWithMetadataStub        func(string, string, map[string]interface{}) (statedb.QueryResultsIterator, error)
	executeQueryWithMetadataMutex       sync.RWMutex
	executeQueryWithMetadataArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 map[string]interface{}
	}
	executeQueryWithMetadataReturns struct {
		result1 statedb.QueryResultsIterator
		result2 error
	}
	executeQueryWithMetadataReturnsOnCall map[int]struct {
		result1 statedb.QueryResultsIterator
		result2 error
	}
	GetLatestSavePointStub        func() (*version.Height, error)
	getLatestSavePointMutex       sync.RWMutex
	getLatestSavePointArgsForCall []struct {
	}
	getLatestSavePointReturns struct {
		result1 *version.Height
		result2 error
	}
	getLatestSavePointReturnsOnCall map[int]struct {
		result1 *version.Height
		result2 error
	}
	GetStateStub        func(string, string) (*statedb.VersionedValue, error)
	getStateMutex       sync.RWMutex
	getStateArgsForCall []struct {
		arg1 string
		arg2 string
	}
	getStateReturns struct {
		result1 *statedb.VersionedValue
		result2 error
	}
	getStateReturnsOnCall map[int]struct {
		result1 *statedb.VersionedValue
		result2 error
	}
	GetStateMultipleKeysStub        func(string, []string) ([]*statedb.VersionedValue, error)
	getStateMultipleKeysMutex       sync.RWMutex
	getStateMultipleKeysArgsForCall []struct {
		arg1 string
		arg2 []string
	}
	getStateMultipleKeysReturns struct {
		result1 []*statedb.VersionedValue
		result2 error
	}
	getStateMultipleKeysReturnsOnCall map[int]struct {
		result1 []*statedb.VersionedValue
		result2 error
	}
	GetStateRangeScanIteratorStub        func(string, string, string) (statedb.ResultsIterator, error)
	getStateRangeScanIteratorMutex       sync.RWMutex
	getStateRangeScanIteratorArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
	}
	getStateRangeScanIteratorReturns struct {
		result1 statedb.ResultsIterator
		result2 error
	}
	getStateRangeScanIteratorReturnsOnCall map[int]struct {
		result1 statedb.ResultsIterator
		result2 error
	}
	GetStateRangeScanIteratorWithMetadataStub        func(string, string, string, map[string]interface{}) (statedb.QueryResultsIterator, error)
	getStateRangeScanIteratorWithMetadataMutex       sync.RWMutex
	getStateRangeScanIteratorWithMetadataArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 map[string]interface{}
	}
	getStateRangeScanIteratorWithMetadataReturns struct {
		result1 statedb.QueryResultsIterator
		result2 error
	}
	getStateRangeScanIteratorWithMetadataReturnsOnCall map[int]struct {
		result1 statedb.QueryResultsIterator
		result2 error
	}
	GetVersionStub        func(string, string) (*version.Height, error)
	getVersionMutex       sync.RWMutex
	getVersionArgsForCall []struct {
		arg1 string
		arg2 string
	}
	getVersionReturns struct {
		result1 *version.Height
		result2 error
	}
	getVersionReturnsOnCall map[int]struct {
		result1 *version.Height
		result2 error
	}
	OpenStub        func() error
	openMutex       sync.RWMutex
	openArgsForCall []struct {
	}
	openReturns struct {
		result1 error
	}
	openReturnsOnCall map[int]struct {
		result1 error
	}
	ValidateKeyValueStub        func(string, []byte) error
	validateKeyValueMutex       sync.RWMutex
	validateKeyValueArgsForCall []struct {
		arg1 string
		arg2 []byte
	}
	validateKeyValueReturns struct {
		result1 error
	}
	validateKeyValueReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *VersionedDB) ApplyUpdates(arg1 *statedb.UpdateBatch, arg2 *version.Height) error {
	fake.applyUpdatesMutex.Lock()
	ret, specificReturn := fake.applyUpdatesReturnsOnCall[len(fake.applyUpdatesArgsForCall)]
	fake.applyUpdatesArgsForCall = append(fake.applyUpdatesArgsForCall, struct {
		arg1 *statedb.UpdateBatch
		arg2 *version.Height
	}{arg1, arg2})
	fake.recordInvocation("ApplyUpdates", []interface{}{arg1, arg2})
	fake.applyUpdatesMutex.Unlock()
	if fake.ApplyUpdatesStub != nil {
		return fake.ApplyUpdatesStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.applyUpdatesReturns
	return fakeReturns.result1
}

func (fake *VersionedDB) ApplyUpdatesCallCount() int {
	fake.applyUpdatesMutex.RLock()
	defer fake.applyUpdatesMutex.RUnlock()
	return len(fake.applyUpdatesArgsForCall)
}

func (fake *VersionedDB) ApplyUpdatesCalls(stub func(*statedb.UpdateBatch, *version.Height) error) {
	fake.applyUpdatesMutex.Lock()
	defer fake.applyUpdatesMutex.Unlock()
	fake.ApplyUpdatesStub = stub
}

func (fake *VersionedDB) ApplyUpdatesArgsForCall(i int) (*statedb.UpdateBatch, *version.Height) {
	fake.applyUpdatesMutex.RLock()
	defer fake.applyUpdatesMutex.RUnlock()
	argsForCall := fake.applyUpdatesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *VersionedDB) ApplyUpdatesReturns(result1 error) {
	fake.applyUpdatesMutex.Lock()
	defer fake.applyUpdatesMutex.Unlock()
	fake.ApplyUpdatesStub = nil
	fake.applyUpdatesReturns = struct {
		result1 error
	}{result1}
}

func (fake *VersionedDB) ApplyUpdatesReturnsOnCall(i int, result1 error) {
	fake.applyUpdatesMutex.Lock()
	defer fake.applyUpdatesMutex.Unlock()
	fake.ApplyUpdatesStub = nil
	if fake.applyUpdatesReturnsOnCall == nil {
		fake.applyUpdatesReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.applyUpdatesReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *VersionedDB) BytesKeySupported() bool {
	fake.bytesKeySupportedMutex.Lock()
	ret, specificReturn := fake.bytesKeySupportedReturnsOnCall[len(fake.bytesKeySupportedArgsForCall)]
	fake.bytesKeySupportedArgsForCall = append(fake.bytesKeySupportedArgsForCall, struct {
	}{})
	fake.recordInvocation("BytesKeySupported", []interface{}{})
	fake.bytesKeySupportedMutex.Unlock()
	if fake.BytesKeySupportedStub != nil {
		return fake.BytesKeySupportedStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.bytesKeySupportedReturns
	return fakeReturns.result1
}

func (fake *VersionedDB) BytesKeySupportedCallCount() int {
	fake.bytesKeySupportedMutex.RLock()
	defer fake.bytesKeySupportedMutex.RUnlock()
	return len(fake.bytesKeySupportedArgsForCall)
}

func (fake *VersionedDB) BytesKeySupportedCalls(stub func() bool) {
	fake.bytesKeySupportedMutex.Lock()
	defer fake.bytesKeySupportedMutex.Unlock()
	fake.BytesKeySupportedStub = stub
}

func (fake *VersionedDB) BytesKeySupportedReturns(result1 bool) {
	fake.bytesKeySupportedMutex.Lock()
	defer fake.bytesKeySupportedMutex.Unlock()
	fake.BytesKeySupportedStub = nil
	fake.bytesKeySupportedReturns = struct {
		result1 bool
	}{result1}
}

func (fake *VersionedDB) BytesKeySupportedReturnsOnCall(i int, result1 bool) {
	fake.bytesKeySupportedMutex.Lock()
	defer fake.bytesKeySupportedMutex.Unlock()
	fake.BytesKeySupportedStub = nil
	if fake.bytesKeySupportedReturnsOnCall == nil {
		fake.bytesKeySupportedReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.bytesKeySupportedReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *VersionedDB) Close() {
	fake.closeMutex.Lock()
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct {
	}{})
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if fake.CloseStub != nil {
		fake.CloseStub()
	}
}

func (fake *VersionedDB) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *VersionedDB) CloseCalls(stub func()) {
	fake.closeMutex.Lock()
	defer fake.closeMutex.Unlock()
	fake.CloseStub = stub
}

func (fake *VersionedDB) ExecuteConditionQuery(arg1 string, arg2 entitydefinition.Search) (interface{}, error) {
	fake.executeConditionQueryMutex.Lock()
	ret, specificReturn := fake.executeConditionQueryReturnsOnCall[len(fake.executeConditionQueryArgsForCall)]
	fake.executeConditionQueryArgsForCall = append(fake.executeConditionQueryArgsForCall, struct {
		arg1 string
		arg2 entitydefinition.Search
	}{arg1, arg2})
	fake.recordInvocation("ExecuteConditionQuery", []interface{}{arg1, arg2})
	fake.executeConditionQueryMutex.Unlock()
	if fake.ExecuteConditionQueryStub != nil {
		return fake.ExecuteConditionQueryStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.executeConditionQueryReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *VersionedDB) ExecuteConditionQueryCallCount() int {
	fake.executeConditionQueryMutex.RLock()
	defer fake.executeConditionQueryMutex.RUnlock()
	return len(fake.executeConditionQueryArgsForCall)
}

func (fake *VersionedDB) ExecuteConditionQueryCalls(stub func(string, entitydefinition.Search) (interface{}, error)) {
	fake.executeConditionQueryMutex.Lock()
	defer fake.executeConditionQueryMutex.Unlock()
	fake.ExecuteConditionQueryStub = stub
}

func (fake *VersionedDB) ExecuteConditionQueryArgsForCall(i int) (string, entitydefinition.Search) {
	fake.executeConditionQueryMutex.RLock()
	defer fake.executeConditionQueryMutex.RUnlock()
	argsForCall := fake.executeConditionQueryArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *VersionedDB) ExecuteConditionQueryReturns(result1 interface{}, result2 error) {
	fake.executeConditionQueryMutex.Lock()
	defer fake.executeConditionQueryMutex.Unlock()
	fake.ExecuteConditionQueryStub = nil
	fake.executeConditionQueryReturns = struct {
		result1 interface{}
		result2 error
	}{result1, result2}
}

func (fake *VersionedDB) ExecuteConditionQueryReturnsOnCall(i int, result1 interface{}, result2 error) {
	fake.executeConditionQueryMutex.Lock()
	defer fake.executeConditionQueryMutex.Unlock()
	fake.ExecuteConditionQueryStub = nil
	if fake.executeConditionQueryReturnsOnCall == nil {
		fake.executeConditionQueryReturnsOnCall = make(map[int]struct {
			result1 interface{}
			result2 error
		})
	}
	fake.executeConditionQueryReturnsOnCall[i] = struct {
		result1 interface{}
		result2 error
	}{result1, result2}
}

func (fake *VersionedDB) ExecuteQuery(arg1 string, arg2 string) (statedb.ResultsIterator, error) {
	fake.executeQueryMutex.Lock()
	ret, specificReturn := fake.executeQueryReturnsOnCall[len(fake.executeQueryArgsForCall)]
	fake.executeQueryArgsForCall = append(fake.executeQueryArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("ExecuteQuery", []interface{}{arg1, arg2})
	fake.executeQueryMutex.Unlock()
	if fake.ExecuteQueryStub != nil {
		return fake.ExecuteQueryStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.executeQueryReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *VersionedDB) ExecuteQueryCallCount() int {
	fake.executeQueryMutex.RLock()
	defer fake.executeQueryMutex.RUnlock()
	return len(fake.executeQueryArgsForCall)
}

func (fake *VersionedDB) ExecuteQueryCalls(stub func(string, string) (statedb.ResultsIterator, error)) {
	fake.executeQueryMutex.Lock()
	defer fake.executeQueryMutex.Unlock()
	fake.ExecuteQueryStub = stub
}

func (fake *VersionedDB) ExecuteQueryArgsForCall(i int) (string, string) {
	fake.executeQueryMutex.RLock()
	defer fake.executeQueryMutex.RUnlock()
	argsForCall := fake.executeQueryArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *VersionedDB) ExecuteQueryReturns(result1 statedb.ResultsIterator, result2 error) {
	fake.executeQueryMutex.Lock()
	defer fake.executeQueryMutex.Unlock()
	fake.ExecuteQueryStub = nil
	fake.executeQueryReturns = struct {
		result1 statedb.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *VersionedDB) ExecuteQueryReturnsOnCall(i int, result1 statedb.ResultsIterator, result2 error) {
	fake.executeQueryMutex.Lock()
	defer fake.executeQueryMutex.Unlock()
	fake.ExecuteQueryStub = nil
	if fake.executeQueryReturnsOnCall == nil {
		fake.executeQueryReturnsOnCall = make(map[int]struct {
			result1 statedb.ResultsIterator
			result2 error
		})
	}
	fake.executeQueryReturnsOnCall[i] = struct {
		result1 statedb.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *VersionedDB) ExecuteQueryWithMetadata(arg1 string, arg2 string, arg3 map[string]interface{}) (statedb.QueryResultsIterator, error) {
	fake.executeQueryWithMetadataMutex.Lock()
	ret, specificReturn := fake.executeQueryWithMetadataReturnsOnCall[len(fake.executeQueryWithMetadataArgsForCall)]
	fake.executeQueryWithMetadataArgsForCall = append(fake.executeQueryWithMetadataArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 map[string]interface{}
	}{arg1, arg2, arg3})
	fake.recordInvocation("ExecuteQueryWithMetadata", []interface{}{arg1, arg2, arg3})
	fake.executeQueryWithMetadataMutex.Unlock()
	if fake.ExecuteQueryWithMetadataStub != nil {
		return fake.ExecuteQueryWithMetadataStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.executeQueryWithMetadataReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *VersionedDB) ExecuteQueryWithMetadataCallCount() int {
	fake.executeQueryWithMetadataMutex.RLock()
	defer fake.executeQueryWithMetadataMutex.RUnlock()
	return len(fake.executeQueryWithMetadataArgsForCall)
}

func (fake *VersionedDB) ExecuteQueryWithMetadataCalls(stub func(string, string, map[string]interface{}) (statedb.QueryResultsIterator, error)) {
	fake.executeQueryWithMetadataMutex.Lock()
	defer fake.executeQueryWithMetadataMutex.Unlock()
	fake.ExecuteQueryWithMetadataStub = stub
}

func (fake *VersionedDB) ExecuteQueryWithMetadataArgsForCall(i int) (string, string, map[string]interface{}) {
	fake.executeQueryWithMetadataMutex.RLock()
	defer fake.executeQueryWithMetadataMutex.RUnlock()
	argsForCall := fake.executeQueryWithMetadataArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *VersionedDB) ExecuteQueryWithMetadataReturns(result1 statedb.QueryResultsIterator, result2 error) {
	fake.executeQueryWithMetadataMutex.Lock()
	defer fake.executeQueryWithMetadataMutex.Unlock()
	fake.ExecuteQueryWithMetadataStub = nil
	fake.executeQueryWithMetadataReturns = struct {
		result1 statedb.QueryResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *VersionedDB) ExecuteQueryWithMetadataReturnsOnCall(i int, result1 statedb.QueryResultsIterator, result2 error) {
	fake.executeQueryWithMetadataMutex.Lock()
	defer fake.executeQueryWithMetadataMutex.Unlock()
	fake.ExecuteQueryWithMetadataStub = nil
	if fake.executeQueryWithMetadataReturnsOnCall == nil {
		fake.executeQueryWithMetadataReturnsOnCall = make(map[int]struct {
			result1 statedb.QueryResultsIterator
			result2 error
		})
	}
	fake.executeQueryWithMetadataReturnsOnCall[i] = struct {
		result1 statedb.QueryResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *VersionedDB) GetLatestSavePoint() (*version.Height, error) {
	fake.getLatestSavePointMutex.Lock()
	ret, specificReturn := fake.getLatestSavePointReturnsOnCall[len(fake.getLatestSavePointArgsForCall)]
	fake.getLatestSavePointArgsForCall = append(fake.getLatestSavePointArgsForCall, struct {
	}{})
	fake.recordInvocation("GetLatestSavePoint", []interface{}{})
	fake.getLatestSavePointMutex.Unlock()
	if fake.GetLatestSavePointStub != nil {
		return fake.GetLatestSavePointStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getLatestSavePointReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *VersionedDB) GetLatestSavePointCallCount() int {
	fake.getLatestSavePointMutex.RLock()
	defer fake.getLatestSavePointMutex.RUnlock()
	return len(fake.getLatestSavePointArgsForCall)
}

func (fake *VersionedDB) GetLatestSavePointCalls(stub func() (*version.Height, error)) {
	fake.getLatestSavePointMutex.Lock()
	defer fake.getLatestSavePointMutex.Unlock()
	fake.GetLatestSavePointStub = stub
}

func (fake *VersionedDB) GetLatestSavePointReturns(result1 *version.Height, result2 error) {
	fake.getLatestSavePointMutex.Lock()
	defer fake.getLatestSavePointMutex.Unlock()
	fake.GetLatestSavePointStub = nil
	fake.getLatestSavePointReturns = struct {
		result1 *version.Height
		result2 error
	}{result1, result2}
}

func (fake *VersionedDB) GetLatestSavePointReturnsOnCall(i int, result1 *version.Height, result2 error) {
	fake.getLatestSavePointMutex.Lock()
	defer fake.getLatestSavePointMutex.Unlock()
	fake.GetLatestSavePointStub = nil
	if fake.getLatestSavePointReturnsOnCall == nil {
		fake.getLatestSavePointReturnsOnCall = make(map[int]struct {
			result1 *version.Height
			result2 error
		})
	}
	fake.getLatestSavePointReturnsOnCall[i] = struct {
		result1 *version.Height
		result2 error
	}{result1, result2}
}

func (fake *VersionedDB) GetState(arg1 string, arg2 string) (*statedb.VersionedValue, error) {
	fake.getStateMutex.Lock()
	ret, specificReturn := fake.getStateReturnsOnCall[len(fake.getStateArgsForCall)]
	fake.getStateArgsForCall = append(fake.getStateArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("GetState", []interface{}{arg1, arg2})
	fake.getStateMutex.Unlock()
	if fake.GetStateStub != nil {
		return fake.GetStateStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getStateReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *VersionedDB) GetStateCallCount() int {
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	return len(fake.getStateArgsForCall)
}

func (fake *VersionedDB) GetStateCalls(stub func(string, string) (*statedb.VersionedValue, error)) {
	fake.getStateMutex.Lock()
	defer fake.getStateMutex.Unlock()
	fake.GetStateStub = stub
}

func (fake *VersionedDB) GetStateArgsForCall(i int) (string, string) {
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	argsForCall := fake.getStateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *VersionedDB) GetStateReturns(result1 *statedb.VersionedValue, result2 error) {
	fake.getStateMutex.Lock()
	defer fake.getStateMutex.Unlock()
	fake.GetStateStub = nil
	fake.getStateReturns = struct {
		result1 *statedb.VersionedValue
		result2 error
	}{result1, result2}
}

func (fake *VersionedDB) GetStateReturnsOnCall(i int, result1 *statedb.VersionedValue, result2 error) {
	fake.getStateMutex.Lock()
	defer fake.getStateMutex.Unlock()
	fake.GetStateStub = nil
	if fake.getStateReturnsOnCall == nil {
		fake.getStateReturnsOnCall = make(map[int]struct {
			result1 *statedb.VersionedValue
			result2 error
		})
	}
	fake.getStateReturnsOnCall[i] = struct {
		result1 *statedb.VersionedValue
		result2 error
	}{result1, result2}
}

func (fake *VersionedDB) GetStateMultipleKeys(arg1 string, arg2 []string) ([]*statedb.VersionedValue, error) {
	var arg2Copy []string
	if arg2 != nil {
		arg2Copy = make([]string, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.getStateMultipleKeysMutex.Lock()
	ret, specificReturn := fake.getStateMultipleKeysReturnsOnCall[len(fake.getStateMultipleKeysArgsForCall)]
	fake.getStateMultipleKeysArgsForCall = append(fake.getStateMultipleKeysArgsForCall, struct {
		arg1 string
		arg2 []string
	}{arg1, arg2Copy})
	fake.recordInvocation("GetStateMultipleKeys", []interface{}{arg1, arg2Copy})
	fake.getStateMultipleKeysMutex.Unlock()
	if fake.GetStateMultipleKeysStub != nil {
		return fake.GetStateMultipleKeysStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getStateMultipleKeysReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *VersionedDB) GetStateMultipleKeysCallCount() int {
	fake.getStateMultipleKeysMutex.RLock()
	defer fake.getStateMultipleKeysMutex.RUnlock()
	return len(fake.getStateMultipleKeysArgsForCall)
}

func (fake *VersionedDB) GetStateMultipleKeysCalls(stub func(string, []string) ([]*statedb.VersionedValue, error)) {
	fake.getStateMultipleKeysMutex.Lock()
	defer fake.getStateMultipleKeysMutex.Unlock()
	fake.GetStateMultipleKeysStub = stub
}

func (fake *VersionedDB) GetStateMultipleKeysArgsForCall(i int) (string, []string) {
	fake.getStateMultipleKeysMutex.RLock()
	defer fake.getStateMultipleKeysMutex.RUnlock()
	argsForCall := fake.getStateMultipleKeysArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *VersionedDB) GetStateMultipleKeysReturns(result1 []*statedb.VersionedValue, result2 error) {
	fake.getStateMultipleKeysMutex.Lock()
	defer fake.getStateMultipleKeysMutex.Unlock()
	fake.GetStateMultipleKeysStub = nil
	fake.getStateMultipleKeysReturns = struct {
		result1 []*statedb.VersionedValue
		result2 error
	}{result1, result2}
}

func (fake *VersionedDB) GetStateMultipleKeysReturnsOnCall(i int, result1 []*statedb.VersionedValue, result2 error) {
	fake.getStateMultipleKeysMutex.Lock()
	defer fake.getStateMultipleKeysMutex.Unlock()
	fake.GetStateMultipleKeysStub = nil
	if fake.getStateMultipleKeysReturnsOnCall == nil {
		fake.getStateMultipleKeysReturnsOnCall = make(map[int]struct {
			result1 []*statedb.VersionedValue
			result2 error
		})
	}
	fake.getStateMultipleKeysReturnsOnCall[i] = struct {
		result1 []*statedb.VersionedValue
		result2 error
	}{result1, result2}
}

func (fake *VersionedDB) GetStateRangeScanIterator(arg1 string, arg2 string, arg3 string) (statedb.ResultsIterator, error) {
	fake.getStateRangeScanIteratorMutex.Lock()
	ret, specificReturn := fake.getStateRangeScanIteratorReturnsOnCall[len(fake.getStateRangeScanIteratorArgsForCall)]
	fake.getStateRangeScanIteratorArgsForCall = append(fake.getStateRangeScanIteratorArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	fake.recordInvocation("GetStateRangeScanIterator", []interface{}{arg1, arg2, arg3})
	fake.getStateRangeScanIteratorMutex.Unlock()
	if fake.GetStateRangeScanIteratorStub != nil {
		return fake.GetStateRangeScanIteratorStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getStateRangeScanIteratorReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *VersionedDB) GetStateRangeScanIteratorCallCount() int {
	fake.getStateRangeScanIteratorMutex.RLock()
	defer fake.getStateRangeScanIteratorMutex.RUnlock()
	return len(fake.getStateRangeScanIteratorArgsForCall)
}

func (fake *VersionedDB) GetStateRangeScanIteratorCalls(stub func(string, string, string) (statedb.ResultsIterator, error)) {
	fake.getStateRangeScanIteratorMutex.Lock()
	defer fake.getStateRangeScanIteratorMutex.Unlock()
	fake.GetStateRangeScanIteratorStub = stub
}

func (fake *VersionedDB) GetStateRangeScanIteratorArgsForCall(i int) (string, string, string) {
	fake.getStateRangeScanIteratorMutex.RLock()
	defer fake.getStateRangeScanIteratorMutex.RUnlock()
	argsForCall := fake.getStateRangeScanIteratorArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *VersionedDB) GetStateRangeScanIteratorReturns(result1 statedb.ResultsIterator, result2 error) {
	fake.getStateRangeScanIteratorMutex.Lock()
	defer fake.getStateRangeScanIteratorMutex.Unlock()
	fake.GetStateRangeScanIteratorStub = nil
	fake.getStateRangeScanIteratorReturns = struct {
		result1 statedb.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *VersionedDB) GetStateRangeScanIteratorReturnsOnCall(i int, result1 statedb.ResultsIterator, result2 error) {
	fake.getStateRangeScanIteratorMutex.Lock()
	defer fake.getStateRangeScanIteratorMutex.Unlock()
	fake.GetStateRangeScanIteratorStub = nil
	if fake.getStateRangeScanIteratorReturnsOnCall == nil {
		fake.getStateRangeScanIteratorReturnsOnCall = make(map[int]struct {
			result1 statedb.ResultsIterator
			result2 error
		})
	}
	fake.getStateRangeScanIteratorReturnsOnCall[i] = struct {
		result1 statedb.ResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *VersionedDB) GetStateRangeScanIteratorWithMetadata(arg1 string, arg2 string, arg3 string, arg4 map[string]interface{}) (statedb.QueryResultsIterator, error) {
	fake.getStateRangeScanIteratorWithMetadataMutex.Lock()
	ret, specificReturn := fake.getStateRangeScanIteratorWithMetadataReturnsOnCall[len(fake.getStateRangeScanIteratorWithMetadataArgsForCall)]
	fake.getStateRangeScanIteratorWithMetadataArgsForCall = append(fake.getStateRangeScanIteratorWithMetadataArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 map[string]interface{}
	}{arg1, arg2, arg3, arg4})
	fake.recordInvocation("GetStateRangeScanIteratorWithMetadata", []interface{}{arg1, arg2, arg3, arg4})
	fake.getStateRangeScanIteratorWithMetadataMutex.Unlock()
	if fake.GetStateRangeScanIteratorWithMetadataStub != nil {
		return fake.GetStateRangeScanIteratorWithMetadataStub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getStateRangeScanIteratorWithMetadataReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *VersionedDB) GetStateRangeScanIteratorWithMetadataCallCount() int {
	fake.getStateRangeScanIteratorWithMetadataMutex.RLock()
	defer fake.getStateRangeScanIteratorWithMetadataMutex.RUnlock()
	return len(fake.getStateRangeScanIteratorWithMetadataArgsForCall)
}

func (fake *VersionedDB) GetStateRangeScanIteratorWithMetadataCalls(stub func(string, string, string, map[string]interface{}) (statedb.QueryResultsIterator, error)) {
	fake.getStateRangeScanIteratorWithMetadataMutex.Lock()
	defer fake.getStateRangeScanIteratorWithMetadataMutex.Unlock()
	fake.GetStateRangeScanIteratorWithMetadataStub = stub
}

func (fake *VersionedDB) GetStateRangeScanIteratorWithMetadataArgsForCall(i int) (string, string, string, map[string]interface{}) {
	fake.getStateRangeScanIteratorWithMetadataMutex.RLock()
	defer fake.getStateRangeScanIteratorWithMetadataMutex.RUnlock()
	argsForCall := fake.getStateRangeScanIteratorWithMetadataArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *VersionedDB) GetStateRangeScanIteratorWithMetadataReturns(result1 statedb.QueryResultsIterator, result2 error) {
	fake.getStateRangeScanIteratorWithMetadataMutex.Lock()
	defer fake.getStateRangeScanIteratorWithMetadataMutex.Unlock()
	fake.GetStateRangeScanIteratorWithMetadataStub = nil
	fake.getStateRangeScanIteratorWithMetadataReturns = struct {
		result1 statedb.QueryResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *VersionedDB) GetStateRangeScanIteratorWithMetadataReturnsOnCall(i int, result1 statedb.QueryResultsIterator, result2 error) {
	fake.getStateRangeScanIteratorWithMetadataMutex.Lock()
	defer fake.getStateRangeScanIteratorWithMetadataMutex.Unlock()
	fake.GetStateRangeScanIteratorWithMetadataStub = nil
	if fake.getStateRangeScanIteratorWithMetadataReturnsOnCall == nil {
		fake.getStateRangeScanIteratorWithMetadataReturnsOnCall = make(map[int]struct {
			result1 statedb.QueryResultsIterator
			result2 error
		})
	}
	fake.getStateRangeScanIteratorWithMetadataReturnsOnCall[i] = struct {
		result1 statedb.QueryResultsIterator
		result2 error
	}{result1, result2}
}

func (fake *VersionedDB) GetVersion(arg1 string, arg2 string) (*version.Height, error) {
	fake.getVersionMutex.Lock()
	ret, specificReturn := fake.getVersionReturnsOnCall[len(fake.getVersionArgsForCall)]
	fake.getVersionArgsForCall = append(fake.getVersionArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("GetVersion", []interface{}{arg1, arg2})
	fake.getVersionMutex.Unlock()
	if fake.GetVersionStub != nil {
		return fake.GetVersionStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getVersionReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *VersionedDB) GetVersionCallCount() int {
	fake.getVersionMutex.RLock()
	defer fake.getVersionMutex.RUnlock()
	return len(fake.getVersionArgsForCall)
}

func (fake *VersionedDB) GetVersionCalls(stub func(string, string) (*version.Height, error)) {
	fake.getVersionMutex.Lock()
	defer fake.getVersionMutex.Unlock()
	fake.GetVersionStub = stub
}

func (fake *VersionedDB) GetVersionArgsForCall(i int) (string, string) {
	fake.getVersionMutex.RLock()
	defer fake.getVersionMutex.RUnlock()
	argsForCall := fake.getVersionArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *VersionedDB) GetVersionReturns(result1 *version.Height, result2 error) {
	fake.getVersionMutex.Lock()
	defer fake.getVersionMutex.Unlock()
	fake.GetVersionStub = nil
	fake.getVersionReturns = struct {
		result1 *version.Height
		result2 error
	}{result1, result2}
}

func (fake *VersionedDB) GetVersionReturnsOnCall(i int, result1 *version.Height, result2 error) {
	fake.getVersionMutex.Lock()
	defer fake.getVersionMutex.Unlock()
	fake.GetVersionStub = nil
	if fake.getVersionReturnsOnCall == nil {
		fake.getVersionReturnsOnCall = make(map[int]struct {
			result1 *version.Height
			result2 error
		})
	}
	fake.getVersionReturnsOnCall[i] = struct {
		result1 *version.Height
		result2 error
	}{result1, result2}
}

func (fake *VersionedDB) Open() error {
	fake.openMutex.Lock()
	ret, specificReturn := fake.openReturnsOnCall[len(fake.openArgsForCall)]
	fake.openArgsForCall = append(fake.openArgsForCall, struct {
	}{})
	fake.recordInvocation("Open", []interface{}{})
	fake.openMutex.Unlock()
	if fake.OpenStub != nil {
		return fake.OpenStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.openReturns
	return fakeReturns.result1
}

func (fake *VersionedDB) OpenCallCount() int {
	fake.openMutex.RLock()
	defer fake.openMutex.RUnlock()
	return len(fake.openArgsForCall)
}

func (fake *VersionedDB) OpenCalls(stub func() error) {
	fake.openMutex.Lock()
	defer fake.openMutex.Unlock()
	fake.OpenStub = stub
}

func (fake *VersionedDB) OpenReturns(result1 error) {
	fake.openMutex.Lock()
	defer fake.openMutex.Unlock()
	fake.OpenStub = nil
	fake.openReturns = struct {
		result1 error
	}{result1}
}

func (fake *VersionedDB) OpenReturnsOnCall(i int, result1 error) {
	fake.openMutex.Lock()
	defer fake.openMutex.Unlock()
	fake.OpenStub = nil
	if fake.openReturnsOnCall == nil {
		fake.openReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.openReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *VersionedDB) ValidateKeyValue(arg1 string, arg2 []byte) error {
	var arg2Copy []byte
	if arg2 != nil {
		arg2Copy = make([]byte, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.validateKeyValueMutex.Lock()
	ret, specificReturn := fake.validateKeyValueReturnsOnCall[len(fake.validateKeyValueArgsForCall)]
	fake.validateKeyValueArgsForCall = append(fake.validateKeyValueArgsForCall, struct {
		arg1 string
		arg2 []byte
	}{arg1, arg2Copy})
	fake.recordInvocation("ValidateKeyValue", []interface{}{arg1, arg2Copy})
	fake.validateKeyValueMutex.Unlock()
	if fake.ValidateKeyValueStub != nil {
		return fake.ValidateKeyValueStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.validateKeyValueReturns
	return fakeReturns.result1
}

func (fake *VersionedDB) ValidateKeyValueCallCount() int {
	fake.validateKeyValueMutex.RLock()
	defer fake.validateKeyValueMutex.RUnlock()
	return len(fake.validateKeyValueArgsForCall)
}

func (fake *VersionedDB) ValidateKeyValueCalls(stub func(string, []byte) error) {
	fake.validateKeyValueMutex.Lock()
	defer fake.validateKeyValueMutex.Unlock()
	fake.ValidateKeyValueStub = stub
}

func (fake *VersionedDB) ValidateKeyValueArgsForCall(i int) (string, []byte) {
	fake.validateKeyValueMutex.RLock()
	defer fake.validateKeyValueMutex.RUnlock()
	argsForCall := fake.validateKeyValueArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *VersionedDB) ValidateKeyValueReturns(result1 error) {
	fake.validateKeyValueMutex.Lock()
	defer fake.validateKeyValueMutex.Unlock()
	fake.ValidateKeyValueStub = nil
	fake.validateKeyValueReturns = struct {
		result1 error
	}{result1}
}

func (fake *VersionedDB) ValidateKeyValueReturnsOnCall(i int, result1 error) {
	fake.validateKeyValueMutex.Lock()
	defer fake.validateKeyValueMutex.Unlock()
	fake.ValidateKeyValueStub = nil
	if fake.validateKeyValueReturnsOnCall == nil {
		fake.validateKeyValueReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.validateKeyValueReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *VersionedDB) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.applyUpdatesMutex.RLock()
	defer fake.applyUpdatesMutex.RUnlock()
	fake.bytesKeySupportedMutex.RLock()
	defer fake.bytesKeySupportedMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	fake.executeConditionQueryMutex.RLock()
	defer fake.executeConditionQueryMutex.RUnlock()
	fake.executeQueryMutex.RLock()
	defer fake.executeQueryMutex.RUnlock()
	fake.executeQueryWithMetadataMutex.RLock()
	defer fake.executeQueryWithMetadataMutex.RUnlock()
	fake.getLatestSavePointMutex.RLock()
	defer fake.getLatestSavePointMutex.RUnlock()
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	fake.getStateMultipleKeysMutex.RLock()
	defer fake.getStateMultipleKeysMutex.RUnlock()
	fake.getStateRangeScanIteratorMutex.RLock()
	defer fake.getStateRangeScanIteratorMutex.RUnlock()
	fake.getStateRangeScanIteratorWithMetadataMutex.RLock()
	defer fake.getStateRangeScanIteratorWithMetadataMutex.RUnlock()
	fake.getVersionMutex.RLock()
	defer fake.getVersionMutex.RUnlock()
	fake.openMutex.RLock()
	defer fake.openMutex.RUnlock()
	fake.validateKeyValueMutex.RLock()
	defer fake.validateKeyValueMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *VersionedDB) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ statedb.VersionedDB = new(VersionedDB)
