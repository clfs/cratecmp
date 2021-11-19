package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	APIRoot   = "https://crates.io/api/v1/crates"
	UserAgent = "cratecmp (https://github.com/clfs/cratecmp)"
)

var (
	nameFlag    = flag.String("name", "", "crate name")
	versionFlag = flag.String("version", "", "crate version")
)

func main() {
	flag.Parse()

	if *nameFlag == "" {
		log.Print("name is required")
		flag.Usage()
		return
	}

	if *versionFlag == "" {
		log.Print("version is required")
		flag.Usage()
		return
	}

	client := &http.Client{}

	url := fmt.Sprintf("%s/%s/%s/download", APIRoot, *nameFlag, *versionFlag)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("%s", resp.Status)
		return
	}

	f, err := os.Create(fmt.Sprintf("%s-%s.crate", *nameFlag, *versionFlag))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
}
