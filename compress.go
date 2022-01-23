package gobuild

import (
	"context"
	"io"
	"io/fs"
	"os"

	"github.com/mholt/archiver/v4"
)

// CompressType describe compress strategy.
type CompressType string

const (
	// CompressRaw do not compress binaries, just copy them to output path.
	CompressRaw CompressType = "raw"

	// CompressTarGz compress all binaries into tar.gz format
	CompressTarGz CompressType = "tar.gz"

	// CompressZip compress all binaries into zip format
	CompressZip CompressType = "zip"
)

func moveWithoutCompress(outputPath, inputPath string) error {
	input, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	output, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	input.Close()
	if err != nil {
		return err
	}

	os.Chmod(outputPath, fs.FileMode(OutputMode))
	return os.Remove(inputPath)
}

// Compress files to outputPath with format
func Compress(outputPath string, files map[string]string, format archiver.Archiver) error {
	targets, err := archiver.FilesFromDisk(nil, files)
	if err != nil {
		return err
	}

	output, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer output.Close()

	return format.Archive(context.Background(), output, targets)
}
