package customerimporter

import (
	"bytes"
	"testing"

	"example.com/interview/test_util"
)

func importFromBuffer(buf *bytes.Buffer, concurrency int, chunkSize int) {
	c := WithOptions(&options{concurrency: concurrency, chunkSize: chunkSize})
	c.readAndImport(buf)
}

func runBench(b *testing.B, lines int32, concurrency int, chunkSize int) {
	b.StopTimer()
	buf := test_util.BuildBufferFile(lines)

	b.StartTimer()
	importFromBuffer(buf, concurrency, chunkSize)
}

func BenchmarkBigFile1KConcurrency40(b *testing.B) {
	runBench(b, 1000, 40, 64000)

}

func BenchmarkBigFile10KConcurrency40(b *testing.B) {
	runBench(b, 10000, 140, 64000)
}

func BenchmarkBigFile100KConcurrency10(b *testing.B) {
	runBench(b, 100000, 10, 64000)
}
func BenchmarkBigFile100KConcurrency40(b *testing.B) {
	runBench(b, 100000, 40, 64000)
}

func BenchmarkBigFile100KConcurrency80(b *testing.B) {
	runBench(b, 100000, 80, 64000)
}

func BenchmarkBigFile100KConcurrency100(b *testing.B) {
	runBench(b, 100000, 100, 64000)
}

func BenchmarkBigFile100KConcurrency120(b *testing.B) {
	runBench(b, 100000, 120, 64000)
}

func BenchmarkBigFile100KConcurrency140(b *testing.B) {
	runBench(b, 100000, 140, 64000)
}

func BenchmarkBigFile1MConcurrency20(b *testing.B) {
	runBench(b, 1000000, 20, 64000)
}

func BenchmarkBigFile1MConcurrency40(b *testing.B) {
	runBench(b, 1000000, 40, 64000)
}

func BenchmarkBigFile1MConcurrency60(b *testing.B) {
	runBench(b, 1000000, 60, 64000)
}

func BenchmarkBigFile1MConcurrency80(b *testing.B) {
	runBench(b, 1000000, 60, 64000)
}

func BenchmarkBigFile1MConcurrency80Chunk128K(b *testing.B) {
	runBench(b, 1000000, 60, 128000)
}

func BenchmarkBigFile1MConcurrency80Chunk256K(b *testing.B) {
	runBench(b, 1000000, 60, 256000)
}
func BenchmarkBigFile1MConcurrency80Chunk512K(b *testing.B) {
	runBench(b, 1000000, 60, 512000)
}
