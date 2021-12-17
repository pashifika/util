// Package conv
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
package conv

import (
	"bytes"
	"unsafe"
)

// BytesToString convert bytes to string
func BytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes convert string to bytes
func StringToBytes(s string) []byte {
	if len(s) == 0 {
		return nil
	}
	b := make([]byte, 0)
	b = append(b, s...)
	return b
}

// StringToBytesV2 convert string to bytes (buffered I/O)
func StringToBytesV2(s string, buffer int) []byte {
	if len(s) == 0 {
		return nil
	}
	b := bytes.NewBuffer(make([]byte, buffer))
	b.WriteString(s)
	return b.Bytes()
}
