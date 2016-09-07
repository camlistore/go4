// +build ignore

/*
Copyright 2016 The Go4 Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Tests that both cloudlaunch and wkfs/gcs work.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

	compute "google.golang.org/api/compute/v1"

	"go4.org/cloud/cloudlaunch"
	"go4.org/wkfs"
	_ "go4.org/wkfs/gcs"
	storageapi "google.golang.org/api/storage/v1"
)

var launchConfig = &cloudlaunch.Config{
	Name:         "serveoncloud",
	BinaryBucket: "camlitests", // REPLACE WITH YOUR OWN, CAN'T DO IT WITH FLAGS
	GCEProjectID: "camlitests", // REPLACE WITH YOUR OWN, CAN'T DO IT WITH FLAGS
	Scopes: []string{
		storageapi.DevstorageFullControlScope,
		compute.ComputeScope,
	},
}

var httpAddr = flag.String("http", ":80", "HTTP address")

func serveHTTP(w http.ResponseWriter, r *http.Request) {
	rc, err := wkfs.Open(path.Join("/gcs", launchConfig.BinaryBucket, r.URL.Path))
	if err != nil {
		http.Error(w, fmt.Sprintf("could not open %v: %v", r.URL.Path, err), 500)
		return
	}
	defer rc.Close()
	http.ServeContent(w, r, r.URL.Path, time.Now(), rc)
}

func main() {
	launchConfig.MaybeDeploy()
	flag.Parse()

	http.HandleFunc("/", serveHTTP)

	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
