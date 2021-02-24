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

	if strings.HasPrefix(filename, filePath) {
		filename = strings.TrimPrefix(filename, filePath)
	}
	header.Name = filename
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
			Logger.Error("AddFileToZip() failed: %+v", err.Error())
			return err
		}
	}

	return nil
}

func (u *Utils) FetchFiles(filePath string, fileList []string) []string {
	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		Logger.Error("ioutil.ReadDir failed: %+v", err)
		return nil
	}
	for _, file := range files {
		if !file.IsDir() {
			fileList = append(fileList, filepath.Join(filePath, file.Name()))
		} else {
			subPath := filepath.Join(filePath, file.Name())
			subFiles, _ := ioutil.ReadDir(subPath)
			for _, subFile := range subFiles {
				if !subFile.IsDir() {
					fileList = append(fileList, filepath.Join(subPath, subFile.Name()))
				} else {
					fileList = u.FetchFiles(filepath.Join(subPath, subFile.Name()), fileList)
				}
			}
		}
	}
	return fileList
}

func (u *Utils) WriteToFile(fileName string, data string) error {
	f, err := os.Create(fileName)
	defer f.Close()

	if err != nil {
		Logger.Error("Unable to create file %s': ", fileName, err.Error())
	} else {
		_, err := io.WriteString(f, data)
		if err != nil {
			Logger.Error("Unable to write to file '%s'", fileName)
			u.DeleteFile(fileName)
		}
	}
	return err
}
