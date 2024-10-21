package main

import (
	"fmt"
	"os"
	"strings"
	"tileServer/db"
)

func makeDb(src, dist string) error {
	writer, err := db.NewWriter(dist)
	if err != nil {
		return err
	}
	defer writer.Close()

	zls, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, zdir := range zls {
		if !zdir.IsDir() {
			continue
		}
		z := zdir.Name()
		xls, err := os.ReadDir(src + "/" + zdir.Name())
		if err != nil {
			return err
		}
		for _, xdir := range xls {
			if !xdir.IsDir() {
				continue
			}
			x := xdir.Name()
			yls, err := os.ReadDir(src + "/" + zdir.Name() + "/" + xdir.Name())
			if err != nil {
				return err
			}
			for _, yf := range yls {
				if yf.IsDir() {
					continue
				}
				split := strings.Split(yf.Name(), ".")

				y := strings.Join(split[:len(split)-1], ".")
				data, err := os.ReadFile(src + "/" + zdir.Name() + "/" + xdir.Name() + "/" + yf.Name())
				if err != nil {
					return err
				}
				err = writer.Write(db.NewMarker(x, y, z), data)
				fmt.Printf("Write %s %s %s\n", x, y, z)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
