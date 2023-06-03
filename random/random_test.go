package randomutils_test

import (
	"math"
	"testing"

	randomutils "github.com/foxxorcat/library-go/random"
)

func TestRandom(t *testing.T) {
	t.Run("RandomUTF8Str", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			str := randomutils.RandomUTF8Str(i)
			if len(str) != i {
				t.Fatal(len(str), i, str)
			}
		}
	})

	t.Run("RandomBytes", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			str := randomutils.RandomBytes(i)
			if len(str) != i {
				t.Fatal(len(str), i, str)
			}
		}
	})

	t.Run("RandomASCII", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			str := randomutils.RandomASCII(i)
			if len(str) != i {
				t.Fatal(len(str), i, str)
			}
		}
	})

}

func BenchmarkFastRandn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = randomutils.FastRandn(math.MaxUint32)
	}
}

func Benchmark(b *testing.B) {
	b.Run("RandomUTF8Str", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = randomutils.RandomUTF8Str(4096)
		}
	})

	b.Run("RandomBytes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = randomutils.RandomBytes(4096)
		}
	})

	b.Run("RandomASCII", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = randomutils.RandomASCII(4096)
		}
	})

}
