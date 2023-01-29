package test_util

import (
	"bufio"
	"bytes"
)

// BuildBufferFile generates a buffer file used as input for testing.
// lines is provided to indicate the number of lines of the generated file.
func BuildBufferFile(lines int32) *bytes.Buffer {

	csvSamples := []string{
		"Mildred,Hernandez,mhernandez0@github.io,Female,38.194.51.128\n",
		"Bonnie,Ortiz,bortiz1@cyberchimps.com,Female,197.54.209.129\n",
		"Dennis,Henry,dhenry2@hubpages.com,Male,155.75.186.217\n",
		"Justin,Hansen,jhansen3@360.cn,Male,251.166.224.119\n",
	}

	var buf bytes.Buffer

	writer := bufio.NewWriter(&buf)
	var i int32
	for i = 0; i < lines; i++ {
		idx := int(i) % len(csvSamples)
		_, _ = writer.WriteString(csvSamples[idx])
	}
	writer.Flush()
	return &buf
}
