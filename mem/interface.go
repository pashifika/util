// Package mem
/*
 * Version: 1.0.0
 * Copyright (c) 2022. Pashifika
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
package mem

import (
	"io"
)

type FakeReader interface {
	// Len returns the number of bytes of the unread portion of the slice.
	Len() int

	// Size returns the original length of the underlying byte slice.
	// Size is the number of bytes available for reading via ReadAt.
	// The returned value is always the same and is not affected by calls
	// to any other method.
	Size() int64

	// Read implements the io.Reader interface.
	Read(b []byte) (n int, err error)

	// ReadAt implements the io.ReaderAt interface.
	ReadAt(b []byte, off int64) (n int, err error)

	// ReadByte implements the io.ByteReader interface.
	ReadByte() (byte, error)

	// UnreadByte complements ReadByte in implementing the io.ByteScanner interface.
	UnreadByte() error

	// ReadRune implements the io.RuneReader interface.
	ReadRune() (ch rune, size int, err error)

	// UnreadRune complements ReadRune in implementing the io.RuneScanner interface.
	UnreadRune() error

	// Seek implements the io.Seeker interface.
	Seek(offset int64, whence int) (int64, error)

	// WriteTo implements the io.WriterTo interface.
	WriteTo(w io.Writer) (n int64, err error)

	// ResetTo resets the Reader to be reading from b.
	ResetTo(b []byte)
}

type FakeWriter interface {
	// Write implements the io.Writer interface.
	Write(p []byte) (n int, err error)

	// WriteAt implements the io.WriterAt interface.
	WriteAt(p []byte, off int64) (n int, err error)

	// WriteRune writes a single Unicode code point, returning
	// the number of bytes written and any error.
	WriteRune(r rune) (n int, err error)

	// WriteString implements the io.StringWriter interface.
	WriteString(s string) (n int, err error)

	// WriteByte implements the io.ByteWriter interface.
	WriteByte(c byte) error

	// Seek implements the io.Seeker interface.
	Seek(offset int64, whence int) (int64, error)

	// ReadFrom implements the io.ReaderFrom interface.
	ReadFrom(r io.Reader) (n int64, err error)
}
