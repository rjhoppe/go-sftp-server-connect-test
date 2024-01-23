package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func envVariable(key string) string {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	host := envVariable("host")
	username := envVariable("username")
	password := envVariable("password")
	port := envVariable("port")

	addr := fmt.Sprintf("%s:%s", host, port)

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		fmt.Printf("Failed to dial: %v", err)
		os.Exit(1)
	}

	defer conn.Close()

	client, err := sftp.NewClient(conn)
	if err != nil {
		fmt.Printf("Failed to create SFTP client: %v ", err)
		os.Exit(1)
	}

	defer client.Close()

	theFiles, err := listFiles(*client, "/")
	if err != nil {
		log.Fatalf("Failed to list files in /: %v", err)
	}

	log.Printf("Found Files in / Files")
	// Output each file name and size in bytes
	log.Printf("%19s %12s %s", "MOD TIME", "SIZE", "NAME")
	for _, theFile := range theFiles {
		log.Printf("%19s %12s %s", theFile.ModTime, theFile.Size, theFile.Name)
	}
}

type remoteFiles struct {
	Name    string
	Size    string
	ModTime string
}

func listFiles(client sftp.Client, remoteDir string) (theFiles []remoteFiles, err error) {

	files, err := client.ReadDir(remoteDir)
	if err != nil {
		return theFiles, fmt.Errorf("Unable to list remote dir: %v", err)
	}

	for _, f := range files {
		var name, modTime, size string

		name = f.Name()
		modTime = f.ModTime().Format("2006-01-02 15:04:05")
		size = fmt.Sprintf("%12d", f.Size())

		if f.IsDir() {
			name = name + "/"
			modTime = ""
			size = "PRE"
		}

		theFiles = append(theFiles, remoteFiles{
			Name:    name,
			Size:    size,
			ModTime: modTime,
		})
	}

	return theFiles, nil
}
