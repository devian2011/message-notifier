package swap

import (
	"encoding/gob"
	"os"
)

func Swap(fileName string, data interface{}) error {

	file, fileOpenErr := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	if fileOpenErr != nil {
		return fileOpenErr
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(data)
}

func Load(fileName string, data interface{}) error {
	file, fileOpenErr := os.OpenFile(fileName, os.O_RDWR, 0777)
	if os.IsNotExist(fileOpenErr) {
		return nil
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	return decoder.Decode(data)
}
