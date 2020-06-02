package common

type Config struct {
	Type     string
	HdfsConf *HdfsConfig
	IpfsConf *IpfsConfig
}

type HdfsConfig struct {
	User                string
	NameNodes           []string
	UseDatanodeHostname bool
}

type IpfsConfig struct {
}
