package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/snappy"
	"github.com/klauspost/compress/zstd"
	"github.com/umbracle/minimal/rlp"
	"github.com/valyala/gozstd"
)

const eth1Folder = "./fixtures/eth1/"

func main() {
	// retrieveHeaders()
}

func retrieveHeaders() {
	c, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		panic(err)
	}

	start := 800000
	count := 100
	for i := start; i < start+count; i++ {
		fmt.Println(i)

		block, err := c.BlockByNumber(context.Background(), big.NewInt(int64(i)))
		if err != nil {
			panic(err)
		}
		dst, err := rlp.EncodeToBytes(block.Header())
		if err != nil {
			panic(err)
		}
		if err := ioutil.WriteFile(filepath.Join(eth1Folder, fmt.Sprintf("header_%d", i)), dst, 0755); err != nil {
			fmt.Printf("Unable to write file: %v", err)
		}
		time.Sleep(1 * time.Second)
	}
}

type CompressFormat int

const (
	SNAPPY CompressFormat = iota
	GOZSTD
	ZSTD
)

func (c CompressFormat) String() string {
	switch c {
	case SNAPPY:
		return "snappy"
	case GOZSTD:
		return "gozstd"
	case ZSTD:
		return "zstd"
	default:
		panic("err")
	}
}

type decompressFn func(dst, src []byte) ([]byte, error)

func getDecompressFn(format CompressFormat) decompressFn {
	switch format {
	case SNAPPY:
		return decompressSnappy
	case GOZSTD:
		return decompressGoZstd
	case ZSTD:
		return decompressZstd
	default:
		panic("err")
	}
}

type compressFn func(dst, src []byte) []byte

func getCompressFn(format CompressFormat) compressFn {
	switch format {
	case SNAPPY:
		return compressSnappy
	case GOZSTD:
		return compressGoZstd
	case ZSTD:
		return compressZstd
	default:
		panic("err")
	}
}

func decompressGoZstd(dst, src []byte) ([]byte, error) {
	return gozstd.Decompress(dst, src)
}

func compressGoZstd(dst, src []byte) []byte {
	return gozstd.Compress(dst, src)
}

var encoder = zstd.Encoder{}
var decoder = zstd.Decoder{}

func decompressZstd(dst, src []byte) ([]byte, error) {
	// This function gets stuck
	return decoder.DecodeAll(src, dst)
}

func compressZstd(dst, src []byte) []byte {
	return encoder.EncodeAll(src, dst)
}

func compressSnappy(dst, src []byte) []byte {
	dst = extendByteSlice(dst, snappy.MaxEncodedLen(len(src)))
	return snappy.Encode(dst, src)
}

func decompressSnappy(dst, src []byte) ([]byte, error) {
	size, err := snappy.DecodedLen(src)
	if err != nil {
		return nil, err
	}
	dst = extendByteSlice(dst, size)
	return snappy.Decode(dst, src)
}

func extendByteSlice(b []byte, needLen int) []byte {
	b = b[:cap(b)]
	if n := needLen - cap(b); n > 0 {
		b = append(b, make([]byte, n)...)
	}
	return b[:needLen]
}

var bufPool byteBufferPool

type byteBuffer struct {
	b []byte
}

func (b *byteBuffer) Reset() {
	b.b = b.b[:0]
}

type byteBufferPool struct {
	pool sync.Pool
}

func (b *byteBufferPool) Get() *byteBuffer {
	v := b.pool.Get()
	if v == nil {
		return &byteBuffer{}
	}
	return v.(*byteBuffer)
}

func (b *byteBufferPool) Put(v *byteBuffer) {
	v.Reset()
	b.pool.Put(v)
}
