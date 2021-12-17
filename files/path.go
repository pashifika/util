// Package files
/*
 * Version: 1.0.0
 * Copyright (c) 2021. Pashifika
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package files

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

// PathBaseAddPrefix add prefix to the last element of path.
func PathBaseAddPrefix(path, prefix string) string {
	if path == "" {
		return ""
	}
	return filepath.Join(filepath.Dir(path), prefix+filepath.Base(path))
}

// PathBaseAddSuffix add suffix to the last element of path.
func PathBaseAddSuffix(path, suffix string) string {
	if path == "" {
		return ""
	}
	name := filepath.Base(path)
	ext := filepath.Ext(name)
	return filepath.Join(filepath.Dir(path), name[:len(name)-len(ext)]+suffix+ext)
}

// RemoveNameExt remove file name's extension used of path.
func RemoveNameExt(name string) string {
	return name[:len(name)-len(filepath.Ext(name))]
}

// Exists returns whether the given file or directory exists or not
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// MkdirIfNotExist used os.MkdirAll to make path's all dir
//goland:noinspection GoUnusedExportedFunction
func MkdirIfNotExist(path string) error {
	folder := filepath.Dir(path)
	if _, err := os.Stat(folder); err != nil {
		if err = os.MkdirAll(folder, 0700); err != nil {
			return err
		}
	}
	return nil
}

// GetFileList used regexp filtering files
func GetFileList(path, filter string, fullPath bool) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	re, err := regexp.Compile(filter)
	if err != nil {
		return nil, err
	}

	var res []string
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if re.MatchString(f.Name()) {
			if fullPath {
				res = append(res, filepath.Join(path, f.Name()))
			} else {
				res = append(res, f.Name())
			}
		}
	}
	return res, nil
}
