package ioutils_test

import (
	"errors"
	"io"
	"strings"
	"sync"
	"testing"

	ioutils "github.com/foxxorcat/library-go/io"
	randomutils "github.com/foxxorcat/library-go/random"
)

func TestBufferingReader(t *testing.T) {
	str := randomutils.RandomASCII(256)

	t.Run("ReaderAt", func(t *testing.T) {
		ra := ioutils.NewBufferingReader(strings.NewReader(str))
		testReaderAt(t, str, ra)
	})

	t.Run("ReadSeeker", func(t *testing.T) {
		ra := ioutils.NewBufferingReader(strings.NewReader(str))
		testReaderSeek(t, str, ra)
	})
}

func testReaderSeek(t *testing.T, raw string, rs io.ReadSeeker) {
	// 正确范围读取测试
	end, err := rs.Seek(0, io.SeekEnd)
	if err != nil {
		t.Fatal(end, err)
	}
	data, err := io.ReadAll(rs)
	if err != nil || len(data) != 0 {
		t.Fatal(end, io.SeekEnd, err)
	}

	ce, err := rs.Seek(-end/2, io.SeekCurrent)
	if err != nil || ce != end/2 {
		t.Fatal(ce, io.SeekStart, err)
	}
	data, err = io.ReadAll(rs)
	if err != nil || string(data) != raw[ce:] {
		t.Fatal(ce, io.SeekEnd, err)
	}

	start, err := rs.Seek(0, io.SeekStart)
	if err != nil || start != 0 {
		t.Fatal(start, err)
	}
	data, err = io.ReadAll(rs)
	if err != nil || string(data) != raw[start:] {
		t.Fatal(start, io.SeekEnd, err)
	}

	// 超出范围报错测试
	end, err = rs.Seek(1, io.SeekEnd)
	if err == nil {
		t.Fatal(end, io.SeekEnd)
	}

	start, err = rs.Seek(-1, io.SeekStart)
	if err == nil {
		t.Fatal(start, io.SeekEnd)
	}
}

func testReaderAt(t *testing.T, raw string, ra io.ReaderAt) {
	var wg sync.WaitGroup
	for x := 0; x < len(raw); x++ {
		for y := x; y < len(raw); y++ {
			wg.Add(1)
			go func(x, y int) {
				defer wg.Done()
				buf := make([]byte, len(raw))

				z, err := ra.ReadAt(buf[:y-x], int64(x))
				if err != nil && err != io.EOF {
					t.Error(err)
				}

				str1 := string(buf[:z])
				str2 := raw[x:y]
				if str1 != str2 {
					t.Errorf("read fail %s != %s", str1, str2)
				}
			}(x, y)
		}
	}
	wg.Wait()
}

func TestZoreReader(t *testing.T) {
	r := ioutils.Zero
	var buf [4096]byte
	n, err := r.Read(buf[:])
	if n != 4096 || err != nil {
		t.Fatal(n, err)
	}
	for _, v := range buf {
		if v != 0 {
			t.Fatal(v, errors.New("is not zero"))
		}
	}
}

func TestCrossReader(t *testing.T) {
	r1 := strings.NewReader("1470")
	r2 := strings.NewReader("235689")
	r := ioutils.NewCrossReader(r1, r2, 1, 2)
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "1234567890" {
		t.Fatal(errors.New("数据校验错误"))
	}
}

func TestRepeatReader(t *testing.T) {
	for i := 0; i < 255; i++ {
		var data [4096]byte
		r := ioutils.NewRepeatReader(byte(i))
		io.ReadFull(r, data[:])
		for _, v := range data {
			if v != byte(i) {
				t.Fatal(errors.New("数据校验失败"), v, i)
			}
		}
	}
}
