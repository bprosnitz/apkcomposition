package main

import (
  "flag"
  "fmt"
  "log"
  "archive/zip"
)

func main() {
  flag.Parse()

  if flag.NArg() != 1 {
    log.Fatalf("Please specify apk")
  }

  r, err := zip.OpenReader(flag.Arg(0))
  if err != nil {
    log.Fatal(err)
  }
  processApk(r)
  r.Close()
}

func processApk(r *zip.ReadCloser) {

  for _, f := range r.File {
    fmt.Printf("Contents of %s:\n", f.Name)
  }
}
