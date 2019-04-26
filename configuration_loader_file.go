package http

import (
	"io/ioutil"
	"os"
	pathlib "path"
)

type FileConfigurationLoader struct {
	Directory string
}

func NewFileConfigurationLoader(dir string) *FileConfigurationLoader {
	return &FileConfigurationLoader{
		Directory: dir,
	}
}

func (loader *FileConfigurationLoader) Load(id string) ([]byte, error) {
	file, err := os.Open(pathlib.Join(loader.Directory, id))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}
