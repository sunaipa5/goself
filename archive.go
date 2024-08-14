package goself

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func Download_Update_File(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error downloading file: %v", err)
	}

	return resp.Body, nil
}

func (options Options) Targz_extractor(file io.Reader) error {
	tmpFolderName := options.TmpFolderName
	if tmpFolderName == ""{
		options.TmpFolderName = ".update-tmp"
	}

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("error creating gzip reader: %v", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)


	if err := os.MkdirAll(options.TmpFolderName, os.ModePerm); err != nil {
		return fmt.Errorf("error creating .tmp directory: %v", err)
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading tar: %v", err)
		}

		target := options.TmpFolderName+"/"+options.AppName

		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(target, os.ModePerm); err != nil {
				return fmt.Errorf("error creating directory: %v", err)
			}
			continue
		}

		// Extract file
		outFile, err := os.Create(target)
		if err != nil {
			return fmt.Errorf("error creating file: %v", err)
		}
		defer outFile.Close()

		// Write file
		if _, err := io.Copy(outFile, tarReader); err != nil {
			return fmt.Errorf("error writing file: %v", err)
		}
	}

	return nil
}

func (options Options) ZipExtractor(file io.Reader) error {
	tmpFolderName := options.TmpFolderName
	if tmpFolderName == "" {
		tmpFolderName = ".update-tmp"
	}


	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	readerAt := bytes.NewReader(data)

	zipReader, err := zip.NewReader(readerAt, int64(readerAt.Len()))
	if err != nil {
		return fmt.Errorf("error creating zip reader: %v", err)
	}

	// Geçici dizini oluşturun
	if err := os.MkdirAll(tmpFolderName, os.ModePerm); err != nil {
		return fmt.Errorf("error creating .tmp directory: %v", err)
	}

	// Zip dosyasındaki her bir dosyayı/dosyayı ayıklayın
	for _, f := range zipReader.File {
		target := filepath.Join(tmpFolderName, f.Name)

		// Eğer dosya bir dizinse, dizini oluşturun
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, os.ModePerm); err != nil {
				return fmt.Errorf("error creating directory: %v", err)
			}
			continue
		}

		// Dosyayı açın
		zippedFile, err := f.Open()
		if err != nil {
			return fmt.Errorf("error opening file in zip: %v", err)
		}

		// Hedef dosyayı oluşturun
		outFile, err := os.Create(target)
		if err != nil {
			zippedFile.Close() // Dosyayı kapatın
			return fmt.Errorf("error creating file: %v", err)
		}

		// Dosyayı yazın
		if _, err := io.Copy(outFile, zippedFile); err != nil {
			zippedFile.Close()
			outFile.Close()
			return fmt.Errorf("error writing file: %v", err)
		}

		zippedFile.Close()
		outFile.Close()
	}

	return nil
}
