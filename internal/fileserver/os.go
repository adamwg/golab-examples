package fileserver

import (
	"io"
	"os"
)

type OS interface {
	Stat(path string) (os.FileInfo, error)
	Open(path string) (io.ReadCloser, error)
}

type realOS struct{}

func (*realOS) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

func (*realOS) Open(path string) (io.ReadCloser, error) {
	return os.Open(path)
}
