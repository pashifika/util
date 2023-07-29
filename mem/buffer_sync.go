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

// Simple byte buffer for marshaling data.

import (
	"bytes"
	"errors"
	"io"
	"sync"
	"unicode/utf8"

	"github.com/pashifika/util/conv"
)

// A SyncFakeIO is a variable-sized buffer of bytes with Read and Write methods.
// The zero value for SyncFakeIO is an empty buffer ready to use.
type SyncFakeIO struct {
	m        sync.RWMutex
	buf      []byte // contents are the bytes buf[off : len(buf)]
	off      int64  // read at &buf[off], write at &buf[len(buf)]
	lastRead readOp // last read operation, so that Unread* can work correctly.

	ManualReset bool // don't auto reset cache
}

// Bytes returns a slice of length b.Len() holding the unread portion of the buffer.
// The slice is valid for use only until the next buffer modification (that is,
// only until the next call to a method like Read, Write, Reset, or Truncate).
// The slice aliases the buffer content at least until the next buffer modification,
// so immediate changes to the slice will affect the result of future reads.
func (fio *SyncFakeIO) Bytes() []byte {
	fio.m.RLock()
	b := fio.buf[fio.off:]
	fio.m.RUnlock()
	return b
}

// String returns the contents of the unread portion of the buffer
// as a string. If the SyncFakeIO is a nil pointer, it returns "<nil>".
//
// To build strings more efficiently, see the strings.Builder type.
func (fio *SyncFakeIO) String() string {
	if fio == nil {
		// Special case, useful in debugging.
		return "<nil>"
	}
	fio.m.RLock()
	str := conv.BytesToString(fio.buf[fio.off:])
	fio.m.RUnlock()
	return str
}

// empty reports whether the unread portion of the buffer is empty.
func (fio *SyncFakeIO) empty() bool { return len(fio.buf) <= int(fio.off) }

func (fio *SyncFakeIO) len() int { return len(fio.buf) - int(fio.off) }

// Cap returns the capacity of the buffer's underlying byte slice, that is, the
// total space allocated for the buffer's data.
func (fio *SyncFakeIO) Cap() int {
	fio.m.RLock()
	n := cap(fio.buf)
	fio.m.RUnlock()
	return n
}

func (fio *SyncFakeIO) Size() int64 {
	fio.m.RLock()
	size := int64(len(fio.buf))
	fio.m.RUnlock()
	return size
}

// Len returns the number of bytes of the unread portion of the buffer;
// b.Len() == len(b.Bytes()).
func (fio *SyncFakeIO) Len() int {
	fio.m.RLock()
	n := fio.len()
	fio.m.RUnlock()
	return n
}

// Truncate discards all but the first n unread bytes from the buffer
// but continues to use the same allocated storage.
// It panics if n is negative or greater than the length of the buffer.
func (fio *SyncFakeIO) Truncate(n int) {
	fio.m.Lock()
	defer fio.m.Unlock()
	if n == 0 {
		fio.reset()
		return
	}

	fio.lastRead = opInvalid
	if n < 0 || n > fio.len() {
		panic("bytes.SyncFakeIO: truncation out of range")
	}
	fio.buf = fio.buf[:int(fio.off)+n]
}

func (fio *SyncFakeIO) reset() {
	fio.buf = fio.buf[:0]
	fio.off = 0
	fio.lastRead = opInvalid
}

// Reset resets the buffer to be empty,
// but it retains the underlying storage for use by future writes.
// Reset is the same as Truncate(0).
func (fio *SyncFakeIO) Reset() {
	fio.m.Lock()
	fio.reset()
	fio.m.Unlock()
}

// ResetTo resets the Reader to be reading from b.
func (fio *SyncFakeIO) ResetTo(b []byte) { *fio = SyncFakeIO{buf: b, off: 0, lastRead: opRead} }

// tryGrowByReslice is a inlineable version of grow for the fast-case where the
// internal buffer only needs to be resliced.
// It returns the index where bytes should be written and whether it succeeded.
func (fio *SyncFakeIO) tryGrowByReslice(n int) (int, bool) {
	if l := len(fio.buf); n <= cap(fio.buf)-l {
		fio.buf = fio.buf[:l+n]
		return l, true
	}
	return 0, false
}

// grow grows the buffer to guarantee space for n more bytes.
// It returns the index where bytes should be written.
// If the buffer can't grow it will panic with ErrTooLarge.
func (fio *SyncFakeIO) grow(n int) int {
	m := fio.len()
	// If buffer is empty, reset to recover space.
	if m == 0 && fio.off != 0 && !fio.ManualReset {
		fio.reset()
	}
	// Try to grow by means of a reslice.
	if i, ok := fio.tryGrowByReslice(n); ok {
		return i
	}
	if fio.buf == nil && n <= smallBufferSize {
		fio.buf = make([]byte, n, smallBufferSize)
		return 0
	}
	c := cap(fio.buf)
	if n <= c/2-m {
		// We can slide things down instead of allocating a new
		// slice. We only need m+n <= c to slide, but
		// we instead let capacity get twice as large so we
		// don't spend all our time copying.
		copy(fio.buf, fio.buf[fio.off:])
	} else if c > maxInt-c-n {
		panic(ErrTooLarge)
	} else {
		// Not enough space anywhere, we need to allocate.
		buf := makeSlice(2*c + n)
		copy(buf, fio.buf[fio.off:])
		fio.buf = buf
	}
	// Restore b.off and len(b.buf).
	fio.off = 0
	fio.buf = fio.buf[:m+n]
	return m
}

// Grow grows the buffer's capacity, if necessary, to guarantee space for
// another n bytes. After Grow(n), at least n bytes can be written to the
// buffer without another allocation.
// If n is negative, Grow will panic.
// If the buffer can't grow it will panic with ErrTooLarge.
func (fio *SyncFakeIO) Grow(n int) {
	if n < 0 {
		panic("bytes.SyncFakeIO.Grow: negative count")
	}
	fio.m.Lock()
	m := fio.grow(n)
	fio.buf = fio.buf[:m]
	fio.m.Unlock()
}

// Write appends the contents of p to the buffer, growing the buffer as
// needed. The return value n is the length of p; err is always nil. If the
// buffer becomes too large, Write will panic with ErrTooLarge.
func (fio *SyncFakeIO) Write(p []byte) (n int, err error) {
	fio.m.Lock()
	defer fio.m.Unlock()
	fio.lastRead = opInvalid
	m, ok := fio.tryGrowByReslice(len(p))
	if !ok {
		m = fio.grow(len(p))
	}
	return copy(fio.buf[m:], p), nil
}

// WriteString appends the contents of s to the buffer, growing the buffer as
// needed. The return value n is the length of s; err is always nil. If the
// buffer becomes too large, WriteString will panic with ErrTooLarge.
func (fio *SyncFakeIO) WriteString(s string) (n int, err error) {
	fio.m.Lock()
	defer fio.m.Unlock()
	fio.lastRead = opInvalid
	m, ok := fio.tryGrowByReslice(len(s))
	if !ok {
		m = fio.grow(len(s))
	}
	return copy(fio.buf[m:], conv.StringToBytes(s)), nil
}

// ReadFrom reads data from r until EOF and appends it to the buffer, growing
// the buffer as needed. The return value n is the number of bytes read. Any
// error except io.EOF encountered during the read is also returned. If the
// buffer becomes too large, ReadFrom will panic with ErrTooLarge.
func (fio *SyncFakeIO) ReadFrom(r io.Reader) (n int64, err error) {
	fio.m.Lock()
	defer fio.m.Unlock()
	fio.lastRead = opInvalid
	for {
		i := fio.grow(MinRead)
		fio.buf = fio.buf[:i]
		m, e := r.Read(fio.buf[i:cap(fio.buf)])
		if m < 0 {
			panic(errNegativeRead)
		}

		fio.buf = fio.buf[:i+m]
		n += int64(m)
		if e == io.EOF {
			return n, nil // e is EOF, so return nil explicitly
		}
		if e != nil {
			return n, e
		}
	}
}

// WriteTo writes data to w until the buffer is drained or an error occurs.
// The return value n is the number of bytes written; it always fits into an
// int, but it is int64 to match the io.WriterTo interface. Any error
// encountered during the write is also returned.
func (fio *SyncFakeIO) WriteTo(w io.Writer) (n int64, err error) {
	fio.m.Lock()
	defer fio.m.Unlock()
	fio.lastRead = opInvalid
	if nBytes := fio.len(); nBytes > 0 {
		m, e := w.Write(fio.buf[fio.off:])
		if m > nBytes {
			panic("bytes.SyncFakeIO.WriteTo: invalid Write count")
		}
		fio.off += int64(m)
		n = int64(m)
		if e != nil {
			return n, e
		}
		// all bytes should have been written, by definition of
		// Write method in io.Writer
		if m != nBytes {
			return n, io.ErrShortWrite
		}
	}
	// SyncFakeIO is now empty; reset.
	if !fio.ManualReset {
		fio.reset()
	}
	return n, nil
}

// WriteByte appends the byte c to the buffer, growing the buffer as needed.
// The returned error is always nil, but is included to match bufio.Writer's
// WriteByte. If the buffer becomes too large, WriteByte will panic with
// ErrTooLarge.
func (fio *SyncFakeIO) WriteByte(c byte) error {
	fio.m.Lock()
	err := fio.writeByte(c)
	fio.m.Unlock()
	return err
}

func (fio *SyncFakeIO) writeByte(c byte) error {
	fio.lastRead = opInvalid
	m, ok := fio.tryGrowByReslice(1)
	if !ok {
		m = fio.grow(1)
	}
	fio.buf[m] = c
	return nil
}

// WriteRune appends the UTF-8 encoding of Unicode code point r to the
// buffer, returning its length and an error, which is always nil but is
// included to match bufio.Writer's WriteRune. The buffer is grown as needed;
// if it becomes too large, WriteRune will panic with ErrTooLarge.
func (fio *SyncFakeIO) WriteRune(r rune) (n int, err error) {
	// Compare as uint32 to correctly handle negative runes.
	if uint32(r) < utf8.RuneSelf {
		fio.m.Lock()
		err = fio.writeByte(byte(r))
		fio.m.Unlock()
		return 1, err
	}
	fio.m.Lock()
	fio.lastRead = opInvalid
	m, ok := fio.tryGrowByReslice(utf8.UTFMax)
	if !ok {
		m = fio.grow(utf8.UTFMax)
	}
	n = utf8.EncodeRune(fio.buf[m:m+utf8.UTFMax], r)
	fio.buf = fio.buf[:m+n]
	fio.m.Unlock()
	return n, nil
}

// Read reads the next len(p) bytes from the buffer or until the buffer
// is drained. The return value n is the number of bytes read. If the
// buffer has no data to return, err is io.EOF (unless len(p) is zero);
// otherwise it is nil.
func (fio *SyncFakeIO) Read(p []byte) (n int, err error) {
	fio.m.Lock()
	defer fio.m.Unlock()
	if fio.off >= int64(len(fio.buf)) {
		return 0, io.EOF
	}
	fio.lastRead = opInvalid
	if fio.empty() {
		// SyncFakeIO is empty, reset to recover space.
		if !fio.ManualReset {
			fio.reset()
		}
		if len(p) == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}
	n = copy(p, fio.buf[fio.off:])
	fio.off += int64(n)
	if n > 0 {
		fio.lastRead = opRead
	}
	return n, nil
}

// Next returns a slice containing the next n bytes from the buffer,
// advancing the buffer as if the bytes had been returned by Read.
// If there are fewer than n bytes in the buffer, Next returns the entire buffer.
// The slice is only valid until the next call to a read or write method.
func (fio *SyncFakeIO) Next(n int) []byte {
	fio.m.Lock()
	defer fio.m.Unlock()
	fio.lastRead = opInvalid
	m := fio.len()
	if n > m {
		n = m
	}
	data := fio.buf[fio.off : fio.off+int64(n)]
	fio.off += int64(n)
	if n > 0 {
		fio.lastRead = opRead
	}
	return data
}

// ReadByte reads and returns the next byte from the buffer.
// If no byte is available, it returns error io.EOF.
func (fio *SyncFakeIO) ReadByte() (byte, error) {
	fio.m.Lock()
	if fio.empty() {
		// SyncFakeIO is empty, reset to recover space.
		if !fio.ManualReset {
			fio.reset()
		}
		fio.m.Unlock()
		return 0, io.EOF
	}
	c := fio.buf[fio.off]
	fio.off++
	fio.lastRead = opRead
	fio.m.Unlock()
	return c, nil
}

// ReadRune reads and returns the next UTF-8-encoded
// Unicode code point from the buffer.
// If no bytes are available, the error returned is io.EOF.
// If the bytes are an erroneous UTF-8 encoding, it
// consumes one byte and returns U+FFFD, 1.
func (fio *SyncFakeIO) ReadRune() (r rune, size int, err error) {
	fio.m.Lock()
	if fio.empty() {
		// SyncFakeIO is empty, reset to recover space.
		if !fio.ManualReset {
			fio.reset()
		}
		fio.m.Unlock()
		return 0, 0, io.EOF
	}
	c := fio.buf[fio.off]
	if c < utf8.RuneSelf {
		fio.off++
		fio.lastRead = opReadRune1
		fio.m.Unlock()
		return rune(c), 1, nil
	}
	r, n := utf8.DecodeRune(fio.buf[fio.off:])
	fio.off += int64(n)
	fio.lastRead = readOp(n)
	fio.m.Unlock()
	return r, n, nil
}

// UnreadRune unreads the last rune returned by ReadRune.
// If the most recent read or write operation on the buffer was
// not a successful ReadRune, UnreadRune returns an error.  (In this regard
// it is stricter than UnreadByte, which will unread the last byte
// from any read operation.)
func (fio *SyncFakeIO) UnreadRune() error {
	fio.m.Lock()
	if fio.lastRead <= opInvalid {
		fio.m.Unlock()
		return errors.New("bytes.SyncFakeIO: UnreadRune: previous operation was not a successful ReadRune")
	}
	if fio.off >= int64(fio.lastRead) {
		fio.off -= int64(fio.lastRead)
	}
	fio.lastRead = opInvalid
	fio.m.Unlock()
	return nil
}

// UnreadByte unreads the last byte returned by the most recent successful
// read operation that read at least one byte. If a write has happened since
// the last read, if the last read returned an error, or if the read read zero
// bytes, UnreadByte returns an error.
func (fio *SyncFakeIO) UnreadByte() error {
	fio.m.Lock()
	if fio.lastRead == opInvalid {
		fio.m.Unlock()
		return errUnreadByte
	}
	fio.lastRead = opInvalid
	if fio.off > 0 {
		fio.off--
	}
	fio.m.Unlock()
	return nil
}

// ReadBytes reads until the first occurrence of delim in the input,
// returning a slice containing the data up to and including the delimiter.
// If ReadBytes encounters an error before finding a delimiter,
// it returns the data read before the error and the error itself (often io.EOF).
// ReadBytes returns err != nil if and only if the returned data does not end in
// delim.
func (fio *SyncFakeIO) ReadBytes(delim byte) (line []byte, err error) {
	fio.m.Lock()
	slice, err := fio.readSlice(delim)
	// return a copy of slice. The buffer's backing array may
	// be overwritten by later calls.
	line = append(line, slice...)
	fio.m.Unlock()
	return line, err
}

// readSlice is like ReadBytes but returns a reference to internal buffer data.
func (fio *SyncFakeIO) readSlice(delim byte) (line []byte, err error) {
	i := bytes.IndexByte(fio.buf[fio.off:], delim)
	end := fio.off + int64(i) + 1
	if i < 0 {
		end = int64(len(fio.buf))
		err = io.EOF
	}
	line = fio.buf[fio.off:end]
	fio.off = end
	fio.lastRead = opRead
	return line, err
}

// ReadString reads until the first occurrence of delim in the input,
// returning a string containing the data up to and including the delimiter.
// If ReadString encounters an error before finding a delimiter,
// it returns the data read before the error and the error itself (often io.EOF).
// ReadString returns err != nil if and only if the returned data does not end
// in delim.
func (fio *SyncFakeIO) ReadString(delim byte) (line string, err error) {
	fio.m.Lock()
	slice, err := fio.readSlice(delim)
	fio.m.Unlock()
	return conv.BytesToString(slice), err
}

// ReadAt implements the io.ReaderAt interface.
func (fio *SyncFakeIO) ReadAt(b []byte, off int64) (n int, err error) {
	// cannot modify state - see io.ReaderAt
	if off < 0 {
		return 0, errors.New("SyncFakeIO.Multi.ReadAt: negative offset")
	}
	fio.m.RLock()
	if off >= int64(len(fio.buf)) {
		fio.m.RUnlock()
		return 0, io.EOF
	}
	n = copy(b, fio.buf[off:])
	if n < len(b) {
		err = io.EOF
	}
	fio.m.RUnlock()
	return
}

// Seek implements the io.Seeker interface.
func (fio *SyncFakeIO) Seek(offset int64, whence int) (int64, error) {
	fio.m.Lock()
	fio.lastRead = opRead
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = fio.off + offset
	case io.SeekEnd:
		abs = int64(len(fio.buf)) + offset
	default:
		fio.m.Unlock()
		return 0, errors.New("bytes.Reader.Seek: invalid whence")
	}
	if abs < 0 {
		fio.m.Unlock()
		return 0, errors.New("bytes.Reader.Seek: negative position")
	}
	fio.off = abs
	fio.m.Unlock()
	return abs, nil
}

// SeekStart if you want to read the complete data after writing, must be used it
func (fio *SyncFakeIO) SeekStart() {
	fio.m.Lock()
	fio.lastRead = opRead
	fio.off = 0
	fio.m.Unlock()
}

// SeekEnd if you want to continue writing data after reading, must be used it
func (fio *SyncFakeIO) SeekEnd() {
	fio.m.Lock()
	fio.lastRead = opRead
	fio.off = int64(len(fio.buf))
	fio.m.Unlock()
}

// Close implements the io.Closer interface.
func (fio *SyncFakeIO) Close() error {
	fio.Reset()
	return nil
}

// WriteAt writes a slice of bytes to a buffer starting at the position provided
// The number of bytes written will be returned, or error. Can overwrite previous
// written slices if the write ats overlap.
func (fio *SyncFakeIO) WriteAt(p []byte, pos int64) (n int, err error) {
	pLen := len(p)
	expLen := pos + int64(pLen)
	fio.m.Lock()
	if int64(len(fio.buf)) < expLen {
		if int64(cap(fio.buf)) < expLen {
			newBuf := make([]byte, expLen, expLen)
			copy(newBuf, fio.buf)
			fio.buf = newBuf
		}
		fio.buf = fio.buf[:expLen]
	}
	copy(fio.buf[pos:], p)
	fio.m.Unlock()
	return pLen, nil
}

// NewSyncFakeIO creates and initializes a new SyncFakeIO using buf as its
// initial contents. The new SyncFakeIO takes ownership of buf, and the
// caller should not use buf after this call. NewFakeIO is intended to
// prepare a SyncFakeIO to read existing data. It can also be used to set
// the initial size of the internal buffer for writing. To do that,
// buf should have the desired capacity but a length of zero.
//
// In most cases, new(SyncFakeIO) (or just declaring a SyncFakeIO variable) is
// sufficient to initialize a SyncFakeIO.
//
//goland:noinspection GoUnusedExportedFunction
func NewSyncFakeIO(buf []byte) *SyncFakeIO { return &SyncFakeIO{buf: buf} }

// NewSyncFakeIOString creates and initializes a new SyncFakeIO using string s as its
// initial contents. It is intended to prepare a buffer to read an existing
// string.
//
// In most cases, new(SyncFakeIO) (or just declaring a SyncFakeIO variable) is
// sufficient to initialize a SyncFakeIO.
//
//goland:noinspection GoUnusedExportedFunction
func NewSyncFakeIOString(s string) *SyncFakeIO {
	return &SyncFakeIO{buf: []byte(s)}
}
