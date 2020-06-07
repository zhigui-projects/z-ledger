/*
Copyright Zhigui.com. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package common

import (
	"io"
	"os"
)

type FsReader interface {
	io.ReaderAt
	io.Seeker
	io.Reader
	io.Closer
	Stat() os.FileInfo
}

type FsClient interface {
	ReadDir(dirname string) ([]os.FileInfo, error)
	Stat(name string) (os.FileInfo, error)
	CopyToRemote(src string, dst string) error
	Open(name string) (FsReader, error)
	Close() error
}
