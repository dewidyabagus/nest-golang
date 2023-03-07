package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

var databases = []string{"products", "accounting", "transaction"}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	var result = make([]rune, length)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}

func sqlDump(db string) ([]byte, error) {
	var dump = randomString(32)

	time.Sleep(7 * time.Second)

	return []byte(fmt.Sprintf("%s_%s", db, dump)), nil
}

func archiveFile(db string, dump []byte) (string, error) {
	fl, err := os.CreateTemp("./archives", fmt.Sprintf("%s_%s_*.zip", db, time.Now().Format("20060102")))
	if err != nil {
		return "", err
	}
	defer fl.Close()

	fl.Write(dump)

	time.Sleep(2 * time.Second)

	return fl.Name(), nil
}

func uploadFile(db, fullPath string) error {
	var dst = fmt.Sprintf("./storages/%s_%s.zip", db, time.Now().Format("20060102"))

	if err := os.Rename(fullPath, dst); err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	return nil
}

type Backup struct {
	DB       string
	DumpFile []byte
}

type Archive struct {
	DB       string
	FullPath string
}

func main() {
	var (
		chanBackup  = make(chan Backup)
		chanArchive = make(chan Archive)
	)

	go func() {
		defer close(chanBackup)

		for _, db := range databases {
			log.Println("Dump database:", db)

			dump, err := sqlDump(db)
			if err != nil {
				log.Printf("backup database %s failed, error: %v \n", db, err)
			}

			chanBackup <- Backup{DB: db, DumpFile: dump}
		}
	}()

	go func() {
		defer close(chanArchive)

		for backup := range chanBackup {
			log.Println("Archive dump:", backup.DB)

			fullPath, err := archiveFile(backup.DB, backup.DumpFile)
			if err != nil {
				log.Printf("archive database dump %s failed, error: %v \n", backup.DB, err)
			}

			chanArchive <- Archive{DB: backup.DB, FullPath: fullPath}
		}
	}()

	for archive := range chanArchive {
		log.Println("Upload file backup:", archive.DB)

		if err := uploadFile(archive.DB, archive.FullPath); err != nil {
			log.Printf("upload file backup %s failed. error: %v \n", archive.DB, err)
		}
	}

	log.Println("Proses backup database berhasil dijalankan")
}
