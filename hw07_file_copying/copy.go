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
	fileStat, statErr := os.Stat(fromPath)
	if statErr != nil {
		return fmt.Errorf("get input file stats: %w", statErr)
	}

	if fileStat.Size() < offset {
		return errors.New("worng offset")
	}
	if !fileStat.Mode().IsRegular() {
		return errors.New("source is not a regular file")
	}
	if fileStat.Size() == 0 {
		return errors.New("worng source file size")
	}
	from, openErr := os.Open(fromPath)
	if openErr != nil {
		return fmt.Errorf("trying to open source path: %w", openErr)
	}
	defer func() {
		err := from.Close()
		if err != nil {
			fmt.Println("error on file close", err)
		}
	}()

	_, seekErr := from.Seek(offset, 0)
	if seekErr != nil {
		return fmt.Errorf("input file seek failure: %w", seekErr)
	}

	to, openErr2 := os.OpenFile(toPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if openErr2 != nil {
		return fmt.Errorf("open or create destination file: %w", openErr2)
	}
	defer func() {
		err := to.Close()
		if err != nil {
			fmt.Println("error on file close", err)
		}
	}()

	if limit == 0 {
		limit = fileStat.Size() - offset
	}

	progressChannel := make(chan any)

	go retransmitData(from, to, limit, progressChannel)

	bar := progressbar.DefaultBytes(limit, fmt.Sprintf("Transmit data from %s to %s", fromPath, toPath))
	for data := range progressChannel {
		switch typedData := data.(type) {
		case error:
			return fmt.Errorf("copy data: %w", typedData)
		case int:
			bar.Add(typedData)
		default:
			return fmt.Errorf("unexpected data type %T", data)
		}
	}

	return nil
}

func retransmitData(from io.Reader, to io.Writer, limit int64, progress chan<- any) {
	defer close(progress)
	if limit <= 0 {
		progress <- errors.New("wrong limit value")
		return
	}
	curLimit := limit
	buf := make([]byte, 1024)
	for {
		maxBufIndex := int64(cap(buf))
		if curLimit < int64(cap(buf)) {
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

		if curLimit == 0 {
			return
		}
	}
}
