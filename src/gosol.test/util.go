package main

import "os"
import "path/filepath"

func execPath() string {
	ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
    exPath := filepath.Dir(ex)
    return exPath
}