package queue

import(
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)


func (q *queue) saveToDisk(filename string) error {

	data, err := json.Marshal(&q)
	if err != nil {
		log.Println("json.Marshal() error: Failed to backup.")
		return err
	}
	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		log.Println("ioutil.WriteFile() error: Failed to backup.")
		return err
	}
	return nil
}


func (q *queue) loadFromDisk(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		log.Println("Backup file found, processing...")

		data, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Println("loadFromDisk() error: Failed to read file.")
		}
		if err := json.Unmarshal(data, q); err != nil {
			log.Println("loadFromDisk() error: Failed to Unmarshal.")
		}
	}
	return nil
}