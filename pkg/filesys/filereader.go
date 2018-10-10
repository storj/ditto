package filesys

import (
	"os"
	"github.com/minio/minio/pkg/hash"
	"io"
	"errors"
)

//Interface---------------------------
type FReader interface {
	io.Reader
	FileInfo() os.FileInfo
	Reader() io.Reader
}

type HFReader interface {
	FReader
	HashReader() *hash.Reader
}

//-------------------------------------------

//Implementation-----------------------------
type fReader struct {
	r     io.Reader
	finfo os.FileInfo
}

func (f *fReader) Read(b []byte) (int, error) {
	if f.r == nil {
		return 0, errors.New("reader is nil")
	}

	return f.r.Read(b)
}

func (f *fReader) FileInfo() os.FileInfo {
	return f.finfo
}

func (f *fReader) Reader() io.Reader {
	return f.r
}

//------------------------------------------
type hfReader struct {
	hr    *hash.Reader
	finfo os.FileInfo
}

func (f *hfReader) Read(b []byte) (int, error) {
	if f.hr == nil {
		return 0, errors.New("reader is nil")
	}

	return f.hr.Read(b)
}

func (f *hfReader) FileInfo() os.FileInfo {
	return f.finfo
}

func (f *hfReader) Reader() io.Reader {
	return f.HashReader()
}

func (f *hfReader) HashReader() *hash.Reader {
	return f.hr
}

//------------------------------------------

//Interface---------------------------------
type FileReader interface {
	ReadFile(string) (FReader, error)
}

type HashFileReader interface {
	ReadFileH(lpath string) (HFReader, error)
}

//-------------------------------------------

//Implementation-----------------------------
type baseFileReader struct {
	FsOpen
}

func (f *baseFileReader) ReadFile(lpath string) (FReader, error) {
	file, err := f.Open(lpath)
	if err != nil {
		return nil, err
	}

	finfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	return &fReader{file, finfo}, nil
}

//---------------------------------------------------------
type hashFileReader struct {
	fr FileReader
}

func NewHashFileReader() HashFileReader {
	return &hashFileReader{&baseFileReader{osOpen(os.Open)}}
}

func (f *hashFileReader) ReadFile(lpath string) (FReader, error) {
	return f.ReadFileH(lpath)
}

func (f *hashFileReader) ReadFileH(lpath string) (HFReader, error) {
	freader, err := f.fr.ReadFile(lpath)
	if err != nil {
		return nil, err
	}

	//TODO: calculate hash
	hreader, err := hash.NewReader(freader, freader.FileInfo().Size(), "", "")
	if err != nil {
		return nil, err
	}

	return &hfReader{hreader, freader.FileInfo()}, nil
}

//--------------------------------------------------------------
