package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/schollz/progressbar/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	from, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("trying to open source path: %w", err)
	}
	from.Seek(offset, 0)
	to, err := os.OpenFile(toPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open or create destination file: %w", err)
	}

	progressChannel := make(chan any)

	go retranslateData(from, to, limit, progressChannel)

	bar := progressbar.DefaultBytes(limit)
	for data := range progressChannel {
		switch data.(type) {
		case error:
			return fmt.Errorf("copy data: %w", data.(error))
		case int:
			bar.Add(data.(int))
		default:
			fmt.Errorf("unexpected data type %T", data)
		}
	}

	return nil
}

func retranslateData(from io.Reader, to io.Writer, limit int64, progress chan<- any) {
	defer close(progress)
	curLimit := limit
	buf := make([]byte, 1024)
	for {
		maxBufIndex := int64(cap(buf))
		if limit > 0 && curLimit < int64(cap(buf)) {
			maxBufIndex = curLimit
		}
		size, err := from.Read(buf[:maxBufIndex])
		curLimit -= int64(size)
		if size > 0 {
			progress <- size
			_, writeErr := to.Write(buf[0:size])
			if writeErr != nil {
				progress <- fmt.Errorf("write to destination: %w", writeErr)
				return
			}
		}

		if err != nil && errors.Is(err, io.EOF) {
			return
		} else if err != nil {
			progress <- fmt.Errorf("read from source: %w", err)
			return
		}

		if limit > 0 && curLimit == 0 {
			return
		}
	}
}
