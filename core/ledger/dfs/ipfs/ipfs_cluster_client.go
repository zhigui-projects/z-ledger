package ipfs

import (
	"context"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	files "github.com/ipfs/go-ipfs-files"
	sh "github.com/pengisgood/go-ipfs-api"
)

type AddFileOutput struct {
	Name string
	Cid  map[string]string
	Size int64
}

type ClusterClient struct {
	url     string //example: http://1.2.3.4:1234
	httpcli *http.Client
}

func NewClusterClient(url string) *ClusterClient {
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}

	c := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	return &ClusterClient{url, c}
}

func (c *ClusterClient) Add(r io.Reader) (string, error) {
	fr := files.NewReaderFile(r)
	slf := files.NewSliceDirectory([]files.DirEntry{files.FileEntry("", fr)})
	fileReader := files.NewMultiFileReader(slf, true)

	opts := map[string]string{
		"encoding":        "json",
		"stream-channels": "true",
	}
	req := &sh.Request{
		Ctx:     context.Background(),
		ApiBase: c.url,
		Command: "add",
		Opts:    opts,
		Body:    fileReader,
	}

	res, err := req.Send(c.httpcli)
	if err != nil {
		return "", err
	}
	var out AddFileOutput
	if err := res.Decode(&out); err != nil {
		return "", err
	}

	return out.Cid["/"], nil
}
