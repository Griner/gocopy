package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

var offset int
var limit int
var from string
var to string

func init() {
	flag.IntVar(&offset, "offset", 0, "offset in input file")
	flag.IntVar(&limit, "limit", 0, "size")
	flag.StringVar(&from, "from", "", "source path")
	flag.StringVar(&to, "to", "", "destination path")
}

func main() {

	flag.Parse()

	_, err := CopyFile(to, from, limit, offset)

	if err != nil {
		log.Fatalf("Copy error: %s\n", err)
	}

}

func CopyFile(dstPath, srcPath string, limit, offset int) (written int64, err error) {

	src, err := os.Open(srcPath)
	if err != nil {
		log.Fatalf("Source error: %s\n", err)
	}
	defer src.Close()

	if offset > 0 {
		_, err := src.Seek(int64(offset), io.SeekStart)
		if err != nil {
			return 0, fmt.Errorf("Source seek error %s\n", err)
		}
	}

	dst, err := os.Create(dstPath)
	if err != nil {
		return 0, fmt.Errorf("Destination error: %s\n", err)
	}
	defer dst.Close()

	return Copy(dst, src, int64(limit))

}

func Copy(dst io.Writer, src io.Reader, n int64) (written int64, err error) {

	progress := make(chan int, 1)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(p chan int) {
		defer wg.Done()
		for i := range p {
			fmt.Printf("Bytes: %d\r", i)
		}
		fmt.Println("")
	}(progress)

	if n > 0 {
		src = io.LimitReader(src, n)
	}

	written, err = io.Copy(dst, NewProgressReader(src, progress))
	if written == n {
		return n, nil
	}
	if written < n && err == nil {
		// src stopped early; must have been EOF.
		err = io.EOF
	}

	wg.Wait()

	return
}

func NewProgressReader(reader io.Reader, progress chan int) io.Reader {
	return &ProgressReader{reader, progress, 0}
}

type ProgressReader struct {
	R         io.Reader
	progress  chan int
	bytesRead int
}

func (pr *ProgressReader) Read(p []byte) (n int, err error) {

	n, err = pr.R.Read(p)
	pr.bytesRead += n

	if n > 0 {
		pr.progress <- pr.bytesRead
	}

	if err == io.EOF {
		close(pr.progress)
	}

	return
}
