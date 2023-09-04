package main

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
)

func CreateZip(filename string, paths []string) error {
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer outFile.Close()

	archive := zip.NewWriter(outFile)
	defer archive.Close()

	for _, pattern := range paths {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return err
		}

		for _, match := range matches {
			i := 0

			err := filepath.Walk(match, func(path string, info os.FileInfo, err error) error {
				if i%5000 == 0 {
					log.Printf("Processed %d files", i)
				}

				i += 1

				relativePath, err := filepath.Rel(match, path)
				if err != nil {
					return err
				}

				if relativePath == "." {
					return nil
				}

				header, err := zip.FileInfoHeader(info)
				if err != nil {
					return err
				}

				header.Name = relativePath
				header.Method = zip.Deflate

				writer, err := archive.CreateHeader(header)
				if err != nil {
					return err
				}

				if info.Mode()&os.ModeSymlink == os.ModeSymlink {
					linkTarget, err := os.Readlink(path)
					if err != nil {
						return err
					}

					if filepath.IsAbs(linkTarget) {
						linkTarget, err = filepath.Rel(match, linkTarget)
						if err != nil {
							return err
						}
					}

					_, err = writer.Write([]byte(linkTarget))

					return err
				}

				if info.IsDir() {
					return nil
				}

				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()

				_, err = io.Copy(writer, file)

				return err
			})

			log.Printf("Processed %d files", i)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func UnpackZip(filename, path string) error {
	reader, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer reader.Close()

	i := 0

	for _, file := range reader.File {
		if i%5000 == 0 {
			log.Printf("Processed %d files", i)
		}

		i += 1

		filePath := filepath.Join(path, file.Name)

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		if file.Mode()&os.ModeSymlink == os.ModeSymlink {
			reader, err := file.Open()
			if err != nil {
				return err
			}
			defer reader.Close()

			buffer := make([]byte, file.FileInfo().Size())
			size, err := reader.Read(buffer)
			if err != nil && err != io.EOF {
				return err
			}

			target := string(buffer[:size])

			err = os.Symlink(target, filePath)
			if err != nil {
				return err
			}

			continue
		} else if file.FileInfo().IsDir() {
			continue
		}

		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return nil
		}
		defer outFile.Close()

		currentFile, err := file.Open()
		if err != nil {
			return err
		}
		defer currentFile.Close()

		if _, err = io.Copy(outFile, currentFile); err != nil {
			return err
		}
	}

	log.Printf("Processed %d files", i)

	return nil
}
