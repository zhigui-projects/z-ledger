package common

type Config struct {
	// dfs type, supported values: hdfs, ipfs
	Type     string
	HdfsConf *HdfsConfig
	IpfsConf *IpfsConfig
}

type HdfsConfig struct {
	// HDFS user
	User string

	// HDFS namenode addresses
	NameNodes []string

	// Docker DNS config is needed when this value is true
	UseDatanodeHostname bool
}

type IpfsConfig struct {
	// IPFS api url, example: 127.0.0.1:5001
	Url string

	// IPFS cluster api url, example: 127.0.0.1:9094
	ClusterUrl string

	// Store dir for file name and cid mapping as index
	CidIndexDir string
}
