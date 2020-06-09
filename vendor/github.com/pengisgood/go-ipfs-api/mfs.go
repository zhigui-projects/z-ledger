package shell

import (
	"context"
	files "github.com/ipfs/go-ipfs-files"
	"io"
)

type MfsLsEntry struct {
	Name string
	Type uint8
	Size uint64
	Hash string
}

type filesLsOutput struct {
	Entries []*MfsLsEntry
}

type filesStatOutput struct {
	Hash string
	Size uint64
	Type string
	//Blocks         int
	//CumulativeSize uint64
	//Local          bool
	//SizeLocal      uint64
	//WithLocality   bool
}

// FilesLs entries at the given path using the MFS commands
func (s *Shell) FilesLs(path string) ([]*MfsLsEntry, error) {
	var out filesLsOutput
	err := s.Request("files/ls", path).
		Option("long", true).
		Exec(context.Background(), &out)
	if err != nil {
		return nil, err
	}

	return out.Entries, nil
}

// FilesStat stat file at the given path using the MFS commands
func (s *Shell) FilesStat(path string) (*filesStatOutput, error) {
	out := &filesStatOutput{}
	err := s.Request("files/stat", path).
		Exec(context.Background(), out)

	if err != nil {
		return nil, err
	}

	return out, err
}

// FilesWrite write file at the given path using the MFS commands
func (s *Shell) FilesWrite(path string, data io.Reader) error {
	fr := files.NewReaderFile(data)
	slf := files.NewSliceDirectory([]files.DirEntry{files.FileEntry("", fr)})
	fileReader := files.NewMultiFileReader(slf, true)

	_, err := s.Request("files/write", path).
		Option("create", true).
		Option("parents", true).
		Body(fileReader).
		Send(context.Background())
	return err
}

// FilesRead read file at the given path using the MFS commands
func (s *Shell) FilesRead(path string) (io.ReadCloser, error) {
	resp, err := s.Request("files/read", path).
		Send(context.Background())
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}

	return resp.Output, nil
}
