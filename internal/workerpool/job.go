package workerpool

import (
	"fmt"
	"io"
	"io/fs"
	"os"
)

type Job struct {
	DirPath  string
	FileInfo fs.FileInfo
}

func (j *Job) execute() (map[string]int, error) {
	result := map[string]int{}

	file, err := os.Open(fmt.Sprintf("%s/%s", j.DirPath, j.FileInfo.Name()))
	if err != nil {
		return result, err
	}
	defer file.Close()

	data := make([]byte, 1)

	for {
		_, err := file.Read(data)
		if err == io.EOF {
			break
		}

		ch := string(data[0])
		if num, ok := result[ch]; !ok {
			result[ch] = 1
		} else {
			result[ch] = num + 1
		}
	}

	return result, nil
}
