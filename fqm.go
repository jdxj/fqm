package main

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	Ext               = ".qmcflac"
	DefaultBufferSize = 8 * (1 << 20) // x * 1MB
)

var (
	seedMap = [][]byte{
		{0x4a, 0xd6, 0xca, 0x90, 0x67, 0xf7, 0x52},
		{0x5e, 0x95, 0x23, 0x9f, 0x13, 0x11, 0x7e},
		{0x47, 0x74, 0x3d, 0x90, 0xaa, 0x3f, 0x51},
		{0xc6, 0x09, 0xd5, 0x9f, 0xfa, 0x66, 0xf9},
		{0xf3, 0xd6, 0xa1, 0x90, 0xa0, 0xf7, 0xf0},
		{0x1d, 0x95, 0xde, 0x9f, 0x84, 0x11, 0xf4},
		{0x0e, 0x74, 0xbb, 0x90, 0xbc, 0x3f, 0x92},
		{0x00, 0x09, 0x5b, 0x9f, 0x62, 0x66, 0xa1},
	}
)

func NewMask() *Mask {
	m := &Mask{
		x:     -1,
		y:     8,
		dx:    1,
		index: -1,
	}
	return m
}

type Mask struct {
	x     int64
	y     int64
	dx    int64
	index int64
}

func (m *Mask) NextMask() (ret byte) {
	if m.x < 0 {
		m.dx = 1
		m.y = (8 - m.y) % 8
		ret = 0xc3
	} else if m.x > 6 {
		m.dx = -1
		m.y = 7 - m.y
		ret = 0xd8
	} else {
		ret = seedMap[m.y][m.x]
	}
	m.x += m.dx

	m.index++
	if m.index == 0x8000 || (m.index > 0x8000 && (m.index+1)%0x8000 == 0) {
		ret = m.NextMask()
	}
	return
}

func NewFQm(input, output string) *FQm {
	fqm := &FQm{
		input:  input,
		output: output,
		reader: bufio.NewReaderSize(nil, DefaultBufferSize),
		writer: bufio.NewWriterSize(nil, DefaultBufferSize),
	}
	return fqm
}

type FQm struct {
	input  string
	output string

	reader *bufio.Reader
	writer *bufio.Writer
}

func (fqm *FQm) Decrypt() error {
	src, err := os.OpenFile(fqm.input, os.O_RDONLY, 0400)
	if err != nil {
		return err
	}
	defer src.Close()
	fqm.reader.Reset(src)

	ext := filepath.Ext(fqm.input)
	base := filepath.Base(fqm.input)
	base = strings.TrimSuffix(base, ext) + ".flac"
	outPath := filepath.Join(fqm.output, base)
	des, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer des.Close()
	defer des.Sync()
	fqm.writer.Reset(des)
	defer fqm.writer.Flush()

	m := NewMask()
	for {
		b, err := fqm.reader.ReadByte()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		err = fqm.writer.WriteByte(b ^ m.NextMask())
		if err != nil {
			return err
		}
	}
}
