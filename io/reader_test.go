package ioutils_test

import (
	"bytes"
	"hash/crc32"
	"io"
	"testing"

	ioutils "github.com/foxxorcat/library-go/io"
	math_utils "github.com/foxxorcat/library-go/math"
	randomutils "github.com/foxxorcat/library-go/random"
	"github.com/pkg/errors"
)

func TestMultiReaderAt(t *testing.T) {
	data1 := randomutils.RandomBytes(1024 * 1024)
	r1 := bytes.NewReader(data1)
	data2 := randomutils.RandomBytes(1024 * 1024)
	r2 := bytes.NewReader(data2)
	data3 := randomutils.RandomBytes(1024 * 1024)
	r3 := bytes.NewReader(data3)

	hash := crc32.NewIEEE()
	hash.Write(data1)
	hash.Write(data2)
	hash.Write(data3)

	mr := ioutils.MultiReaderAt(r1, r2, r3)
	if err := testSizeReadAt(mr, int64(len(data1)+len(data2)+len(data3)), hash.Sum32()); err != nil {
		t.Error(err)
	}
}

func TestLimtReadSeeker(t *testing.T) {
	data := randomutils.RandomBytes(2 * 1024 * 1024)
	r := ioutils.LimitReadSeeker(bytes.NewReader(data), 1024*1024, 1024*1024)
	if err := testReadSeek(r, 1024*1024, crc32.ChecksumIEEE(data[1024*1024:])); err != nil {
		t.Error(err)
	}
}

func TestCrossReader(t *testing.T) {
	data1 := randomutils.RandomBytes(2 * 1024 * 1024)
	data2 := randomutils.RandomBytes(2 * 1024 * 1024)
	r := ioutils.CrossReader(bytes.NewReader(data1), bytes.NewReader(data2), 1024*1024, 1024*1024)

	hash := crc32.NewIEEE()
	hash.Write(data1[:1024*1024])
	hash.Write(data2[:1024*1024])
	hash.Write(data1[1024*1024:])
	hash.Write(data2[1024*1024:])

	if err := testReader(r, hash.Sum32()); err != nil {
		t.Error(err)
	}
}

func TestBufferingReader(t *testing.T) {
	data1 := randomutils.RandomBytes(2 * 1024 * 1024)
	r := ioutils.NewBufferReader(bytes.NewReader(data1))

	if err := testReadAt(r, int64(len(data1)), crc32.ChecksumIEEE(data1)); err != nil {
		t.Error(err)
	}

	if err := testReadSeek(r, int64(len(data1)), crc32.ChecksumIEEE(data1)); err != nil {
		t.Error(err)
	}
}

func TestReaderAtBuffer(t *testing.T) {
	data1 := randomutils.RandomBytes(2 * 1024 * 1024)
	r := ioutils.NewReaderAtBuffer(bytes.NewReader(data1), 4096, 12)
	if err := testReadAt(r, int64(len(data1)), crc32.ChecksumIEEE(data1)); err != nil {
		t.Error(err)
	}
}

func TestBufferReadSeeker(t *testing.T) {
	data1 := randomutils.RandomBytes(2 * 1024 * 1024)
	r := ioutils.NewBufferReadSeeker(bytes.NewReader(data1), 4096, 12)
	if err := testReadAt(r, int64(len(data1)), crc32.ChecksumIEEE(data1)); err != nil {
		t.Error(err)
	}

	if err := testReadSeek(r, int64(len(data1)), crc32.ChecksumIEEE(data1)); err != nil {
		t.Error(err)
	}
}

func testSizeReadAt(rs ioutils.SizeReaderAt, size int64, crc32_ uint32) error {
	if rs.Size() != size {
		return errors.Errorf("大小错误,实际大小:%d != %d", rs.Size(), size)
	}
	return testReadAt(rs, rs.Size(), crc32_)
}

func testReadAt(r ioutils.ReaderAt, size int64, crc32_ uint32) error {
	var (
		n   int
		err error
		buf [4096]byte
	)

	/* 正常功能测试 */
	for i := 0; i < int(math_utils.Log(size)); i++ {
		off := randomutils.FastRandn(uint32(size))
		n, err = r.ReadAt(buf[:], int64(off))

		if n < 0 || int64(n) > size || n > len(buf) {
			return errors.Errorf("n 超过正常范围")
		}

		lb := len(buf)
		// n != len(p) 必然返回错误
		if n != lb && err == nil {
			return errors.Errorf("ReadAt实现有误, (n=%d) != (len(p)=%d) 但 err=nil", n, lb)
		}

		// n == len(buf) 必须返回io.EOF或nil
		if n == lb && err != nil && err != io.EOF {
			return errors.Errorf("buf 已读满，但返回非法错误%s", err)
		}
	}

	/* 非法值测试 */
	n, err = r.ReadAt(buf[:], size+1)
	if n != 0 || err == nil {
		return errors.Errorf("超过范围的读取应该返回错误, 并且保证 n==0")
	}

	n, err = r.ReadAt(buf[:], -1)
	if n != 0 || err == nil {
		return errors.Errorf("负数偏移的读取应该返回错误, 并且保证 n==0")
	}

	/* 内容测试 */
	hash := crc32.NewIEEE()
	for off := int64(0); off < size; off += int64(n) {
		n, err = r.ReadAt(buf[:randomutils.FastRandn(4096)], off)
		hash.Write(buf[:n])
		if err != nil {
			if err != io.EOF {
				return errors.WithMessage(err, "内容测试读取错误")
			}
			break
		}
	}

	crc32_2 := hash.Sum32()
	if crc32_2 != crc32_ {
		return errors.Errorf("读取内容错误, 内容crc32不匹配 %d != %d", crc32_2, crc32_)
	}

	return nil
}

func testReadSeek(r ioutils.ReadSeeker, size int64, crc32_ uint32) error {
	n, err := ioutils.StreamSizeBySeeking(r, true)
	if err != nil {
		return errors.WithMessage(err, "StreamSizeBySeeking")
	}
	if n != size {
		return errors.New("seek 获取大小与数据大小不符合")
	}

	/* 正常使用测试 */
	var buf [4096]byte

	// Seek可用性测试
	for i := 0; i < int(math_utils.Log(size)); i++ {
		off := int64(randomutils.FastRandn(uint32(size)))
		noff, err := r.Seek(off, io.SeekStart)
		if err != nil || noff != off {
			return errors.Errorf("Seek错误 noff:%d != off:%d, err=%s", noff, off, err)
		}
	}
	for i := 0; i < int(math_utils.Log(size)); i++ {
		off := -int64(randomutils.FastRandn(uint32(size)))
		noff, err := r.Seek(off, io.SeekEnd)
		off = size + off
		if err != nil || noff != off {
			return errors.Errorf("Seek错误 noff:%d != off:%d, err=%s", noff, off, err)
		}
	}

	// 读文件头
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return err
	}
	if _, err := r.Read(buf[:]); err != nil && err != io.EOF {
		return err
	}

	// 读文件末尾
	if _, err := r.Seek(0, io.SeekEnd); err != nil {
		return err
	}
	if _, err := r.Read(buf[:]); err != io.EOF {
		return errors.WithMessage(err, "读取文件末尾应该返回io.EOF")
	}

	/* 非法值测试 */
	if _, err := r.Seek(-1, io.SeekStart); err == nil {
		return errors.Errorf("SeekStart 错误的范围，应该返回错误")
	}

	if _, err := r.Seek(-size-1, io.SeekCurrent); err == nil {
		return errors.Errorf("SeekCurrent 错误的范围，应该返回错误")
	}

	if _, err := r.Seek(-size-1, io.SeekEnd); err == nil {
		return errors.Errorf("SeekEnd 错误的范围，应该返回错误")
	}

	r.Seek(0, io.SeekStart)
	return testReader(r, crc32_)
}

func testReader(r ioutils.Reader, crc32_ uint32) error {
	hash := crc32.NewIEEE()
	_, err := io.Copy(hash, r)
	if err != nil && err != io.EOF {
		return err
	}

	if hash.Sum32() != crc32_ {
		return errors.Errorf("读取内容错误, 内容crc32不匹配")
	}
	return nil
}
