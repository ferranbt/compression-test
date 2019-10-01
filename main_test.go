package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
)

//// Compress

func benchmarkCompressOne(b *testing.B, format CompressFormat) {
	src, err := ioutil.ReadFile(filepath.Join(eth1Folder, "header_800000"))
	if err != nil {
		panic(err)
	}

	b.ReportAllocs()

	f := getCompressFn(format)
	for i := 0; i < b.N; i++ {
		bb := bufPool.Get()
		bb.b = f(bb.b[:0], src)
		bufPool.Put(bb)
	}
}

func BenchmarkCompressOne(b *testing.B) {
	b.Run("Snappy", func(b *testing.B) {
		benchmarkCompressOne(b, SNAPPY)
	})
	b.Run("GoZstd", func(b *testing.B) {
		benchmarkCompressOne(b, GOZSTD)
	})
	b.Run("Zstd", func(b *testing.B) {
		benchmarkCompressOne(b, ZSTD)
	})
}

func benchmarkCompressAll(b *testing.B, format CompressFormat) {
	buf := make([][]byte, 0)

	files, err := ioutil.ReadDir(eth1Folder)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		src, err := ioutil.ReadFile(filepath.Join(eth1Folder, file.Name()))
		if err != nil {
			panic(err)
		}
		buf = append(buf, src)
	}

	b.ResetTimer()
	b.ReportAllocs()

	f := getCompressFn(format)
	for i := 0; i < b.N; i++ {
		for _, src := range buf {
			bb := bufPool.Get()
			bb.b = f(bb.b[:0], src)
			bufPool.Put(bb)
		}
	}
}

func BenchmarkCompressAll(b *testing.B) {
	b.Run("Snappy", func(b *testing.B) {
		benchmarkCompressAll(b, SNAPPY)
	})
	b.Run("GoZstd", func(b *testing.B) {
		benchmarkCompressAll(b, GOZSTD)
	})
	b.Run("Zstd", func(b *testing.B) {
		benchmarkCompressAll(b, ZSTD)
	})
}

//// Decompress

func benchmarkDecompressOne(b *testing.B, format CompressFormat) {
	aux, err := ioutil.ReadFile(filepath.Join(eth1Folder, "header_800000"))
	if err != nil {
		panic(err)
	}

	b.ReportAllocs()

	enc := getCompressFn(format)
	src := enc(nil, aux)

	f := getDecompressFn(format)
	for i := 0; i < b.N; i++ {
		bb := bufPool.Get()
		bb.b, _ = f(bb.b[:0], src)
		bufPool.Put(bb)
	}
}

func BenchmarkDecompressOne(b *testing.B) {
	b.Run("Snappy", func(b *testing.B) {
		benchmarkDecompressOne(b, SNAPPY)
	})
	b.Run("GoZstd", func(b *testing.B) {
		benchmarkDecompressOne(b, GOZSTD)
	})
}

func benchmarkDecompressAll(b *testing.B, format CompressFormat) {
	buf := make([][]byte, 0)
	enc := getCompressFn(format)

	files, err := ioutil.ReadDir(eth1Folder)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		src, err := ioutil.ReadFile(filepath.Join(eth1Folder, file.Name()))
		if err != nil {
			panic(err)
		}
		buf = append(buf, enc(nil, src))
	}

	b.ResetTimer()
	b.ReportAllocs()

	f := getDecompressFn(format)
	for i := 0; i < b.N; i++ {
		for _, src := range buf {
			bb := bufPool.Get()
			bb.b, _ = f(bb.b[:0], src)
			bufPool.Put(bb)
		}
	}
}

func BenchmarkDecompressAll(b *testing.B) {
	b.Run("Snappy", func(b *testing.B) {
		benchmarkDecompressAll(b, SNAPPY)
	})
	b.Run("GoZstd", func(b *testing.B) {
		benchmarkDecompressAll(b, GOZSTD)
	})
}

func testCompression(format CompressFormat) {
	f := getCompressFn(format)

	files, err := ioutil.ReadDir(eth1Folder)
	if err != nil {
		panic(err)
	}

	totalSize := 0
	compSize := 0

	for _, file := range files {
		src, err := ioutil.ReadFile(filepath.Join(eth1Folder, file.Name()))
		if err != nil {
			panic(err)
		}

		bb := bufPool.Get()
		bb.b = f(bb.b[:0], src)

		totalSize += len(src)
		compSize += len(bb.b)
		bufPool.Put(bb)
	}

	fmt.Printf("%s: %f\n", format.String(), float64(compSize*100)/float64(totalSize))
}

func TestCompression(t *testing.T) {
	testCompression(ZSTD)
	testCompression(GOZSTD)
	testCompression(SNAPPY)
}

func testDecompression(c CompressFormat) {
	enc := getCompressFn(c)
	dec := getDecompressFn(c)

	files, err := ioutil.ReadDir(eth1Folder)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		src, err := ioutil.ReadFile(filepath.Join(eth1Folder, file.Name()))
		if err != nil {
			panic(err)
		}

		res, err := dec(nil, enc(nil, src))
		if err != nil {
			panic(err)
		}

		if !bytes.Equal(src, res) {
			panic("bad")
		}
	}
}

func TestDecompression(t *testing.T) {
	t.Run("Snappy", func(t *testing.T) {
		testDecompression(SNAPPY)
	})
	t.Run("Gozstd", func(t *testing.T) {
		testDecompression(GOZSTD)
	})
	/*
		t.Run("Zstd", func(t *testing.T) {
			testDecompression(ZSTD)
		})
	*/
}
