package db

import (
	"encoding/binary"
	"encoding/gob"
	"io"
	"os"
)

type TileDbWriter struct {
	offsetMap map[string]int64
	file      interface {
		io.WriteCloser
		io.Seeker
	}
	saveOffset int64
}

func NewWriter(fileName string) (*TileDbWriter, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	offsetMap := make(map[string]int64)
	return &TileDbWriter{offsetMap, file, 8}, nil
}

func (w *TileDbWriter) Write(marker Marker, data []byte) error {
	_, err := w.file.Seek(w.saveOffset, io.SeekStart)
	if err != nil {
		return err
	}
	_, err = w.file.Write(data)
	if err != nil {
		return err
	}
	w.offsetMap[marker.String()] = w.saveOffset
	w.saveOffset += int64(len(data))
	return nil
}

func (w *TileDbWriter) Close() error {
	// 写文件表偏移
	_, err := w.file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	err = binary.Write(w.file, binary.LittleEndian, w.saveOffset)
	if err != nil {
		return err
	}
	// 写文件表
	_, err = w.file.Seek(w.saveOffset, io.SeekStart)
	if err != nil {
		return err
	}
	err = gob.NewEncoder(w.file).Encode(w.offsetMap)
	if err != nil {
		return err
	}

	return w.file.Close()
}
