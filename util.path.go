package elio

import "path/filepath"

// GetBasepath get base path
func GetBasepath(path string) string {
	appName := filepath.Base(path)
	basePath := path[:len(path)-len(appName)]
	return basePath
}

// GetBasename get base name
func GetBasename(path string) string {
	appName := filepath.Base(path)
	//log.Printf("app base: %s\n", appName)     // server.upmatcher.exe
	appExt := filepath.Ext(appName)
	//log.Printf("app ext: %s\n", appExt)       // .exe
	baseName := appName[:len(appName)-len(appExt)]
	//log.Printf("base name: %s\n", baseName)   // server.upmatcher

	return baseName
}
