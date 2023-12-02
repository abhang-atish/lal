package main

import (
	"fmt"
	"log"
	"os"
)

type ManifestFile struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	IsMaster   bool   `json:"is_master"`
	Key        string `json:"key"`
	Bytes      []byte `json:"bytes"`
	StreamName string `json:"stream_name"`
}

func (mf *ManifestFile) Read(force bool) {
	if mf.Bytes != nil && !force {
		return
	}

	file, err := os.Open(mf.Path)
	if err != nil {
		e := fmt.Errorf("event='FAILED_TO_READ_FILE' filename=%s err=%s", mf.Path, err.Error())
		log.Println(e)
		return
	}

	defer func(file *os.File) { _ = file.Close() }(file)

	fileInfo, err := file.Stat()
	if err != nil {
		e := fmt.Errorf("event='FAILED_TO_READ_FILE' filename=%s err=%s", mf.Path, err.Error())
		log.Println(e)
		return
	}

	var size = fileInfo.Size()
	buffer := make([]byte, size)
	_, _ = file.Read(buffer)
	mf.Bytes = buffer
}

func (mf *ManifestFile) SaveToDisk() error {
	var file *os.File
	var err error

	file, err = os.Create(mf.Path)
	if err != nil {
		return fmt.Errorf("event=%s file=%s error=%v", "FAILED_TO_CREATE_ON_DISK", mf.Path, err.Error())
	}
	defer file.Close()
	_, err = file.Write(mf.Bytes)
	if err != nil {
		return fmt.Errorf("event=%s file=%s error=%v", "FAILED_TO_WRITE_ON_DISK", mf.Path, err.Error())
	}
	file.Sync()
	return nil
}

func (mf *ManifestFile) GetOldCopy() ([]byte, error) {

	old_copy, err := GetValue([]byte(fmt.Sprintf("manifest_m3u8_%s", mf.Key)))
	if err == nil {
		return old_copy, nil
	}
	mf.Read(true)
	if mf.Bytes == nil || len(mf.Bytes) == 0 {
		return nil, fmt.Errorf("empty")
	}
	return mf.Bytes, nil
}
