package datasource

import (
	"fmt"
	"os"
)

type GitDataSource struct {}

func (ds *GitDataSource) Fetch(from, to string) ([]string, error) {
	fmt.Printf("*Fetching data from %s into %s...\n", from, to)
	
}

func createFolderIfNotExist(path string) error  {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(path, os.ModePerm); err != nil {

			}
		}
	}
}