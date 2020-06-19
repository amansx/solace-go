package main

import "os"
import "path"
import "path/filepath"

func pluginPath(pluginName string) string {
	ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
    exPath := filepath.Dir(ex)
    return path.Join(exPath, pluginName)
}