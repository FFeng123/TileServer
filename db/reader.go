package db

import (
	"encoding/binary"
	"encoding/gob"
	"errors"
	"io"
	"os"
	"sort"
)

type TileDbReader struct {
	file      *os.File
	offsetMap map[string][2]int64
}

var ErrMarkerNotFound = errors.New("marker not found")

func NewReader(path string) (*TileDbReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	// 读取偏移表
	var tableOffset int64
	err = binary.Read(file, binary.LittleEndian, &tableOffset)
	if err != nil {
		return nil, err
	}
	_, err = file.Seek(tableOffset, io.SeekStart)
	if err != nil {
		return nil, err
	}
	var offsetMapSrc map[string]int64
	err = gob.NewDecoder(file).Decode(&offsetMapSrc)
	if err != nil {
		return nil, errors.New("decode offset map failed")
	}
	// 计算起止位置
	type FileItem struct {
		key   string
		start int64
	}
	items := make([]FileItem, 0, len(offsetMapSrc))
	for key, start := range offsetMapSrc {
		items = append(items, FileItem{key, start})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].start < items[j].start
	})
	offsetMap := make(map[string][2]int64, len(items))
	for i := 0; i < len(items)-1; i++ {
		if i == len(items)-1 {
			offsetMap[items[i].key] = [2]int64{items[i].start, tableOffset}
		} else {
			offsetMap[items[i].key] = [2]int64{items[i].start, items[i+1].start}
		}
	}

	return &TileDbReader{file, offsetMap}, nil
}
func (r *TileDbReader) Read(marker Marker) ([]byte, error) {
	offset, ok := r.offsetMap[marker.String()]
	if !ok {
		return nil, ErrMarkerNotFound
	}
	data := make([]byte, offset[1]-offset[0])
	_, err := r.file.ReadAt(data, offset[0])
	return data, err
}
func (r *TileDbReader) Close() error {
	return r.file.Close()
}
