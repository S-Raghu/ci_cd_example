package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var registry []string
var whiteList []string
var fileHash = make(map[string]string)

func main() {
	startup()
}

func startup() {
	action := "./build.sh build docker deploy"
	log.Println(fileHash)
	log.Println("creating filehash")
	createWhitelist()
	log.Println("building registry")
	buildRegistry()
	go doEvery(15*time.Second, verifyHashes, action)
	for {
	}
}

func runAction(action string) {
	log.Println("Taking action and running...: " + action)
	cmd := exec.Command("/bin/sh", "-c", action)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		log.Printf("error")
	}
	log.Println(outb.String())
	log.Println(errb.String())
	log.Println("--------------------------------------------------------------------------------")
}

func doEvery(d time.Duration, f func(time.Time, string), action string) {
	for x := range time.Tick(d) {
		f(x, action)
	}
}

func handleErr(err error) {
	if err != nil {
		log.Println("error")
	}
}

func stopwatch(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func verifyHashes(t time.Time, action string) {
	log.Println("verifying hashes...")
	for _, fn := range registry {
		oldHash := retrieveHash(fn)
		newHash := calculateHash(fn)
		if !(compareHash(oldHash, newHash) == 0) {
			insertRecord(fn, newHash)
			log.Println(fn + "old hash" + oldHash + "new hash" + newHash + "changes detected - updating hash, action required")
			runAction(action)
		}
	}
}

func createWhitelist() {
	file, err := os.Open("./ignore")
	if err != nil {
		log.Println("no .ignore file found, race condition will ensue if jobs edit files -- will not create whitelist")
	} else {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			log.Println(scanner.Text())
			whiteList = append(whiteList, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
}

func calculateHash(absoluteFilePath string) string {
	f, err := os.Open(absoluteFilePath)
	handleErr(err)
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(h.Sum(nil))
}

func insertRecord(absoluteFilePath string, hash string) {
	fileHash[absoluteFilePath] = hash
}

func retrieveHash(absoluteFilePath string) string {
	val, _ := fileHash[absoluteFilePath]
	return val
}

func recursiveDirectoryCrawl(dirName string) {
	files, err := ioutil.ReadDir(dirName)
	handleErr(err)
	for _, f := range files {
		fileOrDir, err := os.Stat(dirName + "/" + f.Name())
		handleErr(err)
		switch mode := fileOrDir.Mode(); {
		case mode.IsDir():
			if !(f.Name() == ".git") {
				recursiveDirectoryCrawl(dirName + "/" + f.Name())
			}
		case mode.IsRegular():
			toAdd := true
			for _, whiteListed := range whiteList {
				if f.Name() == whiteListed {
					toAdd = false
					log.Println(f.Name() + "is whitelisted, not adding to registry")
				}
			}
			if toAdd {
				absolutePath := dirName + "/" + f.Name()
				registry = append(registry, absolutePath)
			}
		}
	}
}

func compareHash(old string, new string) int {
	return strings.Compare(old, new)
}

func buildRegistry() {
	log.Println("starting directory scan")
	recursiveDirectoryCrawl(".")
	log.Println("computing hashes and creating map entries")
	for _, fn := range registry {
		hash := calculateHash(fn)
		insertRecord(fn, hash)
	}
}
