package stateormdb

//func (vdb *VersionedDB) buildCommittersForNs(ns string, nsUpdates map[string]*statedb.VersionedValue) ([]*committer, error) {
//	//db, err := vdb.getNamespaceDBHandle(ns)
//	//if err != nil {
//	//	return nil, err
//	//}
//	//// for each namespace, build mutiple committers based on the maxBatchSize
//	//maxBatchSize := db.CouchInstance.MaxBatchUpdateSize()
//	//numCommitters := 1
//	//if maxBatchSize > 0 {
//	//	numCommitters = int(math.Ceil(float64(len(nsUpdates)) / float64(maxBatchSize)))
//	//}
//	//committers := make([]*committer, numCommitters)
//	//
//	//cacheEnabled := vdb.cache.Enabled(ns)
//	//
//	//for i := 0; i < numCommitters; i++ {
//	//	committers[i] = &committer{
//	//		db:             db,
//	//		batchUpdateMap: make(map[string]*batchableDocument),
//	//		namespace:      ns,
//	//		cacheKVs:       make(statedb.CacheKVs),
//	//		cacheEnabled:   cacheEnabled,
//	//	}
//	//}
//	//
//	//// for each committer, create a couchDoc per key-value pair present in the update batch
//	//// which are associated with the committer's namespace.
//	//revisions, err := vdb.getRevisions(ns, nsUpdates)
//	//if err != nil {
//	//	return nil, err
//	//}
//	//
//	//i := 0
//	//for key, vv := range nsUpdates {
//	//	kv := &keyValue{key: key, revision: revisions[key], VersionedValue: vv}
//	//	couchDoc, err := keyValToCouchDoc(kv)
//	//	if err != nil {
//	//		return nil, err
//	//	}
//	//	committers[i].batchUpdateMap[key] = &batchableDocument{CouchDoc: *couchDoc, Deleted: vv.Value == nil}
//	//	committers[i].addToCacheUpdate(kv)
//	//	if maxBatchSize > 0 && len(committers[i].batchUpdateMap) == maxBatchSize {
//	//		i++
//	//	}
//	//}
//	//return committers, nil
//}
