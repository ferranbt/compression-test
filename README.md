
# Benchmarks

## Compression

One header:

```
$ go test -v ./... -run=XXX -bench=CompressOne
```

```
goos: linux
goarch: amd64
pkg: github.com/ferranbt/compression-test
BenchmarkCompressOne/Snappy-8         	 5000000	       312 ns/op	       0 B/op	       0 allocs/op
BenchmarkCompressOne/GoZstd-8         	  200000	      8506 ns/op	       0 B/op	       0 allocs/op
BenchmarkCompressOne/Zstd-8           	  100000	     17868 ns/op	       0 B/op	       0 allocs/op
```

100 headers:

```
$ go test -v ./... -run=XXX -bench=CompressAll
```

```
goos: linux
goarch: amd64
pkg: github.com/ferranbt/compression-test
BenchmarkCompressAll/Snappy-8         	   50000	     61289 ns/op	       0 B/op	       0 allocs/op
BenchmarkCompressAll/GoZstd-8         	    2000	    692607 ns/op	       1 B/op	       0 allocs/op
BenchmarkCompressAll/Zstd-8           	     500	   2341084 ns/op	       3 B/op	
```

## Decompression

One header:

```
$ go test -v ./... -run=XXX -bench=DecompressOne
```

```
goos: linux
goarch: amd64
pkg: github.com/ferranbt/compression-test
BenchmarkDecompressOne/Snappy-8         	10000000	       188 ns/op	       0 B/op       0 allocs/op
BenchmarkDecompressOne/GoZstd-8         	 5000000	       355 ns/op	       0 B/op       0 allocs/op
```

100 headers:

```
$ go test -v ./... -run=XXX -bench=DecompressAll
```

```
goos: linux
goarch: amd64
pkg: github.com/ferranbt/compression-test
BenchmarkDecompressAll/Snappy-8         	  100000	     23061 ns/op	       0 B/op       0 allocs/op
BenchmarkDecompressAll/GoZstd-8         	   50000	     39547 ns/op	       0 B/op       0 allocs/op
```
