package file_generator

import (
	"fmt"
	"math/rand"
	"os"
)

type Generator struct {
	path string
}

func NewGenerator(path string) *Generator {
	return &Generator{
		path: path,
	}
}

func (g *Generator) CreateFiles(numberOfFiles int) error {
	for i := 0; i < numberOfFiles; i++ {
		file, err := os.Create(fmt.Sprintf("%s/file%d.txt", g.path, i))
		if err != nil {
			return fmt.Errorf("error occurred while creating file: %v", err)
		}

		str := randomBytesArray()
		_, err = file.Write(str)
		if err != nil {
			return fmt.Errorf("error occurred while writing in file: %v", err)
		}

		file.Close()
	}

	return nil
}

func (g *Generator) CleanFiles() {
	readDirectory, _ := os.Open(g.path)
	allFiles, _ := readDirectory.Readdir(0)

	for f := range allFiles {
		file := allFiles[f]

		fileName := file.Name()
		filePath := fmt.Sprintf("%s/%s", g.path, fileName)

		os.Remove(filePath)
	}
}

func randomBytesArray() []byte {
	numberOfChars := rand.Intn(10000)
	res := ""
	for i := 0; i < numberOfChars; i++ {
		randASCII := rand.Intn(128)
		res += string(randASCII)
	}

	return []byte(res)
}
