package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	action := Action{
		Action:          os.Args[1],
		Bucket:          os.Getenv("AWS_BUCKET"),
		Key:             fmt.Sprintf("%s.zip", os.Getenv("KEY")),
		Paths:           strings.Split(strings.TrimSpace(os.Getenv("PATHS")), "\n"),
		FailOnCacheMiss: strings.TrimSpace(os.Getenv("FAIL_ON_CACHE_MISS")) == "true",
	}

	if action.Action == SaveAction {
		if len(action.Paths) == 0 {
			log.Fatal("No paths provided")
		}

		if err := CreateZip(action.Key, action.Paths); err != nil {
			log.Fatal(err)
		}

		if err := UploadFile(action.Key, action.Bucket); err != nil {
			log.Fatal(err)
		}
	} else if action.Action == RestoreAction {
		exists, err := ObjectExists(action.Key, action.Bucket)
		if err != nil {
			log.Fatal(err)
		}

		if exists {
			if err := DownloadFile(action.Key, action.Bucket); err != nil {
				log.Fatal(err)
			}

			if err := UnpackZip(action.Key); err != nil {
				log.Fatal(err)
			}
		} else {
			if action.FailOnCacheMiss {
				log.Fatalf("No cache found for key \"%s\". Failing because \"fail-on-cache-miss\" was set to \"true\".", action.Key)
			} else {
				log.Printf("No cache found for key \"%s\". Continuing.", action.Key)
			}
		}
	} else {
		log.Fatalf("Provided action \"%s\" not recognized. Only save or restore are available.", action.Action)
	}
}
