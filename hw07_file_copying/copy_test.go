package main

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

var outPath string
var tmpDir string

func setup(t *testing.T, offset, limit int64) {
	tmpDir, err := ioutil.TempDir("", "copy")
	if err != nil {
		t.Fatal("can't create temp dir: ", err)
	}

	var builder strings.Builder
	builder.WriteString(tmpDir)
	builder.WriteString("out_offset")
	builder.WriteString(strconv.Itoa(int(offset)))
	builder.WriteString("_limit")
	builder.WriteString(strconv.Itoa(int(limit)))
	builder.WriteString(".txt")
	outPath = builder.String()
}

func TestCopy_Valid(t *testing.T) {
	tests := []struct {
		name     string
		fromPath string
		toPath   string
		offset   int64
		limit    int64
		err      error
	}{
		{
			name:     "test case when offset=0 and limit=0",
			fromPath: "testdata/input.txt",
			offset:   0,
			limit:    0,
			toPath:   "out_offset0_limit0.txt",
		},
		{
			name:     "test case when offset=0 and limit=10",
			fromPath: "testdata/input.txt",
			offset:   0,
			limit:    10,
			toPath:   "out_offset0_limit10.txt",
		},
		{
			name:     "test case when offset=0 and limit=1000",
			fromPath: "testdata/input.txt",
			offset:   0,
			limit:    1000,
			toPath:   "out_offset0_limit1000.txt",
		},
		{
			name:     "test case when offset=0 and limit=10000",
			fromPath: "testdata/input.txt",
			offset:   0,
			limit:    10000,
			toPath:   "out_offset0_limit10000.txt",
		},
		{
			name:     "test case when offset=100 and limit=1000",
			fromPath: "testdata/input.txt",
			offset:   100,
			limit:    1000,
			toPath:   "out_offset100_limit1000.txt",
		},
		{
			name:     "test case when offset=6000 and limit=1000",
			fromPath: "testdata/input.txt",
			offset:   6000,
			limit:    1000,
			toPath:   "out_offset6000_limit1000.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.RemoveAll(tmpDir)
			setup(t, tt.offset, tt.limit)
			err := Copy(tt.fromPath, outPath, tt.offset, tt.limit)

			expected, err := ioutil.ReadFile("testdata/" + tt.toPath)
			require.NoError(t, err)

			actual, err := ioutil.ReadFile(outPath)
			require.NoError(t, err)

			require.Equal(t, expected, actual)
			require.Nil(t, err)
		})
	}
}

func TestCopy_InValid(t *testing.T) {
	tests := []struct {
		name     string
		fromPath string
		toPath   string
		offset   int64
		limit    int64
		err      error
	}{
		{
			name:     "test case when returns error unsupported file",
			fromPath: "testdata/input1.txt",
			offset:   1,
			limit:    0,
			err:      ErrUnsupportedFile,
		},
		{
			name:     "test case when returns error offset exceeds file size",
			fromPath: "testdata/input.txt",
			offset:   8000,
			limit:    1000,
			err:      ErrOffsetExceedsFileSize,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.RemoveAll(tmpDir)
			setup(t, tt.offset, tt.limit)
			err := Copy(tt.fromPath, outPath, tt.offset, tt.limit)

			require.Equal(t, errors.Unwrap(tt.err), errors.Unwrap(err))
		})
	}
}
