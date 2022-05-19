/*
==================================================================================
  Copyright (c) 2019 AT&T Intellectual Property.
  Copyright (c) 2019 Nokia

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
==================================================================================
*/

package xapp

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Utils struct {
	baseDir string
	status  string
}

func NewUtils() *Utils {
	b := Config.GetString("controls.symptomdata.baseDir")
	if b == "" {
		b = "/tmp/symptomdata/"
	}

	return &Utils{
		baseDir: b,
	}
}

func (u *Utils) FileExists(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

func (u *Utils) CreateDir(path string) error {
	if u.FileExists(path) {
		os.RemoveAll(path)
	}
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	os.Chmod(path, os.ModePerm)
	return nil
}

func (u *Utils) DeleteFile(fileName string) {
	os.Remove(fileName)
}

func (u *Utils) AddFileToZip(zipWriter *zip.Writer, filePath string, filename string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = strings.TrimPrefix(filename, filePath)
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	if info.Size() > 0 {
		_, err = io.Copy(writer, fileToZip)
	}
	return err
}

func (u *Utils) ZipFiles(newZipFile *os.File, filePath string, files []string) error {
	defer newZipFile.Close()
	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()
	for _, file := range files {
		if err := u.AddFileToZip(zipWriter, filePath, file); err != nil {
			return err
		}
	}

	return nil
}

func (u *Utils) ZipFilesToTmpFile(baseDir string, tmpfilename string, fileList []string) (string, error) {
	//Generate zip file
	tmpFile, err := ioutil.TempFile("", tmpfilename)
	if err != nil {
		return "", fmt.Errorf("Failed to create a tmp file: %w", err)
	}
	err = u.ZipFiles(tmpFile, baseDir, fileList)
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("Failed to zip the files: %w", err)
	}
	return tmpFile.Name(), nil
}

func (u *Utils) GetFileFromZip(file *zip.File, filePath string) (string, error) {
	filename := filepath.Join(filePath, file.Name)

	if file.FileInfo().IsDir() {
		os.MkdirAll(filename, os.ModePerm)
		return "", nil
	}

	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return "", fmt.Errorf("mkdir failed %s", filepath.Dir(filename))
	}

	dstFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return "", fmt.Errorf("openfile failed %s", filename)
	}
	defer dstFile.Close()

	fileInArchive, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("zip file open failed %s", file.Name)
	}
	defer fileInArchive.Close()

	if _, err := io.Copy(dstFile, fileInArchive); err != nil {
		return "", fmt.Errorf("copy failed %s -> %s", file.Name, filename)
	}
	return filename, nil

}

func (u *Utils) UnZipFiles(zippedFile string, filePath string) ([]string, error) {
	retval := []string{}
	zipReader, err := zip.OpenReader(zippedFile)
	if err != nil {
		return retval, fmt.Errorf("Failed to open zip reader: %w", err)
	}
	defer zipReader.Close()

	//fmt.Printf("Reading zipfile: %s\n", zippedFile)
	for _, file := range zipReader.File {
		fname, err := u.GetFileFromZip(file, filePath)
		if err != nil {
			return retval, fmt.Errorf("Failed to unzip the files: %w", err)
		}
		if len(fname) > 0 {
			retval = append(retval, fname)
		}
	}
	return retval, nil
}

func (u *Utils) FetchFiles(filePath string, fileList []string) []string {
	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		fmt.Printf("ioutil.ReadDir failed: %+v\n", err)
		return nil
	}
	for _, file := range files {
		if !file.IsDir() {
			fileList = append(fileList, filepath.Join(filePath, file.Name()))
		} else {
			fileList = u.FetchFiles(filepath.Join(filePath, file.Name()), fileList)
		}
	}
	return fileList
}

func (u *Utils) WriteToFile(fileName string, data string) error {
	f, err := os.Create(fileName)
	defer f.Close()

	if err != nil {
		Logger.Error("Unable to create file %s': %+v", fileName, err)
	} else {
		_, err := io.WriteString(f, data)
		if err != nil {
			Logger.Error("Unable to write to file '%s': %+v", fileName, err)
			u.DeleteFile(fileName)
		}
	}
	return err
}
