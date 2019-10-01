
# Tests

## Compression

One header:

```
$ go test -v ./... -run=XXX -bench=CompressOne
```

100 headers:

```
$ go test -v ./... -run=XXX -bench=CompressAll
```

## Decompression

One header:

```
$ go test -v ./... -run=XXX -bench=DecompressOne
```

100 headers:

```
$ go test -v ./... -run=XXX -bench=DecompressAll
```
