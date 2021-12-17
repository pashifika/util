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
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// FileOpen os full path.
//
// w  open the file write-only. (support create a new file)
//
// r  open the file read-only.
//
// a  append data to the file when writing. (aw is support create a new file)
//
// rw open the file read-write. (support create a new file)
func FileOpen(path, mode string) (*os.File, error) {
	var wmode int
	switch mode {
	case "w":
		wmode = os.O_WRONLY | os.O_CREATE
	case "r":
		wmode = os.O_RDONLY
	case "a":
		wmode = os.O_WRONLY | os.O_APPEND
	case "wa", "aw":
		wmode = os.O_WRONLY | os.O_APPEND | os.O_CREATE
	case "rw", "wr":
		wmode = os.O_RDWR | os.O_CREATE
	default:
		wmode = os.O_RDONLY
	}
	fp, err := os.OpenFile(path, wmode, 0664)
	if err != nil {
		return nil, err
	}
	return fp, err
}

// ByteToFile write bytes to file.
// (best to the small file)
func ByteToFile(path string, buf []byte) (err error) {
	if Exists(path) {
		err = os.Remove(path)
		if err != nil {
			return
		}
	}
	return ioutil.WriteFile(path, buf, 0664)
}

// BufferToFile write buffer to file.
// (best to the big file)
func BufferToFile(path string, r io.Reader) (err error) {
	if Exists(path) {
		err = os.Remove(path)
		if err != nil {
			return
		}
	}
	f, err := FileOpen(path, "w")
	if err != nil {
		return
	}
	//noinspection ALL
	defer f.Close()

	// make a write buffer
	w := bufio.NewWriter(f)
	var n int
	// make a buffer to keep chunks that are read
	buf := make([]byte, 4096)
	for {
		// read a chunk
		n, err = r.Read(buf)
		if err != nil && err != io.EOF {
			return
		}
		if n == 0 {
			break
		}

		// write a chunk
		if n, err = w.Write(buf[:n]); err != nil {
			return
		}
		if n != 4096 {
			err = fmt.Errorf("write chunk size error, file:%s", path)
			return
		}
	}

	return w.Flush()
}
