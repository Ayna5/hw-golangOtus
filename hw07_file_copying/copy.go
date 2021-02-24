package main

import (
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/pkg/errors"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromStat, err := os.Stat(fromPath)
	if err != nil {
		return errors.Wrapf(err, "cannot get file stat for path: %s", fromPath)
	}

	if fromStat.Size() == 0 {
		return ErrUnsupportedFile
	}

	if fromStat.Size() <= offset {
		return ErrOffsetExceedsFileSize
	}

	from, err := os.OpenFile(fromPath, os.O_RDONLY, os.ModeDir)
	if err != nil {
		return errors.Wrapf(err, "cannot open file for path: %s", fromPath)
	}
	defer from.Close()

	to, err := os.Create(toPath)
	if err != nil {
		return errors.Wrapf(err, "cannot create file for path: %s", toPath)
	}
	defer to.Close()

	_, err = from.Seek(offset, io.SeekStart)
	if err != nil {
		return errors.Wrapf(err, "cannot execute seek")
	}

	if limit == 0 {
		limit = fromStat.Size() - offset
	}

	var chunk int64 = 1024
	bar := pb.Full.Start64(limit)
	bar.Set(pb.Bytes, true)
	for {
		if chunk > limit {
			chunk = limit
		}
		if limit == 0 {
			break
		}
		written, err := io.CopyN(to, from, chunk)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return errors.Wrapf(err, "cannot execute io.CopyN")
		}
		bar.Add64(written)
		limit -= written
		time.Sleep(time.Millisecond * 100)
	}
	bar.Finish()
	return nil
}
