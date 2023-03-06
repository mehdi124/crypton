package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type PublicInfo struct {
	Name    string
	Address string
}

func NewPublicInfo(name, address string) *PublicInfo {
	return &PublicInfo{
		Name:    name,
		Address: address,
	}
}

func (d *Driver) WritePublicInfo(collection, resource string, pb *PublicInfo) error {

	if collection == "" {
		return fmt.Errorf("missing collection - no place to save records")
	}

	if resource == "" {
		return fmt.Errorf("missing resource - unable to save records (no name)")
	}

	mutex := d.getOrCreateMutext(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)

	log.Println(dir, resource, "ssss")

	fnlPath := filepath.Join(dir, resource+".json")
	log.Println(fnlPath)
	tempPath := fnlPath + ".tmp"

	if err := os.MkdirAll(dir, 0775); err != nil {
		return err
	}

	b, err := json.MarshalIndent(pb, "", "\t")
	fmt.Println(b)
	if err != nil {
		return err
	}

	b = append(b, byte('\n'))

	if err := ioutil.WriteFile(tempPath, b, 0644); err != nil {
		return nil
	}

	return os.Rename(tempPath, fnlPath)
}

func (d *Driver) ReadAllPublicInfo(collection string) ([]string, error) {

	if collection == "" {
		return nil, fmt.Errorf("missing collection name - unable to read")
	}

	dir := filepath.Join(d.dir, collection)

	if _, err := stat(dir); err != nil {
		return nil, nil
	}

	files, _ := ioutil.ReadDir(dir)
	var records []string

	for _, file := range files {
		b, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))

		if err != nil {
			return nil, err
		}

		records = append(records, string(b))

	}

	return records, nil
}

func (d *Driver) DeletePublicInfo(collection, resource string) error {

	path := filepath.Join(collection, resource)

	mutex := d.getOrCreateMutext(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, path)

	switch fi, err := stat(dir); {

	case fi == nil, err != nil:
		return fmt.Errorf("unable to find file or directory name %v\n", path)

	case fi.Mode().IsDir():
		return os.RemoveAll(dir)

	case fi.Mode().IsRegular():
		return os.RemoveAll(dir + ".json")
	}

	return nil
}

/*func (d *Driver) getOrCreateMutext(collection string) *sync.Mutex {

	d.mutex.Lock()
	defer d.mutex.Unlock()

	m, ok := d.mutexes[collection]

	if !ok {
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}

	return m
}

func stat(path string) (fi os.FileInfo, err error) {

	if fi, err = os.Stat(path); os.IsNotExist(err) {

		fi, err = os.Stat(path + ".json")

	}

	return fi, err

}*/
