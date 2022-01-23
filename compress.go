package gobuild

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"io/fs"
	"os"
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

func compressTarGz(outputPath, inputPath, binName string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	return addFileToTarWriter(tarWriter, inputPath, binName)
}

func compressZip(outputPath, inputPath, binName string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	return addFileToZipWritter(zipWriter, inputPath, binName)
}

func compressRaw(outputPath, inputPath string) error {
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

func addFileToTarWriter(tarWriter *tar.Writer, filePath, fileName string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	header := &tar.Header{
		Name:    fileName,
		Size:    stat.Size(),
		Mode:    int64(stat.Mode()),
		ModTime: stat.ModTime(),
	}

	err = tarWriter.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(tarWriter, file)
	return err
}

func addFileToZipWritter(zipWriter *zip.Writer, filePath, fileName string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(stat)
	if err != nil {
		return err
	}

	header.Name = fileName
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}
