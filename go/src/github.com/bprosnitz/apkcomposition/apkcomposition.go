package main

import (
  "flag"
  "fmt"
  "log"
  "strings"
  "archive/zip"
)

type categoryPath []string

type categoryNode struct {
  compressed, uncompressed uint64
  children map[string]*categoryNode
}

func newCategoryNode() *categoryNode {
  return &categoryNode{
    children: map[string]*categoryNode{},
  }
}

func main() {
  flag.Parse()

  if flag.NArg() != 1 {
    log.Fatalf("Please specify apk")
  }

  r, err := zip.OpenReader(flag.Arg(0))
  if err != nil {
    log.Fatal(err)
  }
  root := newCategoryNode()
  processApk(r, root)
  printCategories(root, "[root]", "")
  r.Close()
}

func processApk(r *zip.ReadCloser, root *categoryNode) {
  for _, f := range r.File {
    addToCategories(f.CompressedSize64, f.UncompressedSize64, categorize(f), root)
  }
}

func addToCategories(compressed, uncompressed uint64, path categoryPath, root *categoryNode) {
  root.compressed += compressed
  root.uncompressed += uncompressed

  if len(path) > 0 {
    next, ok := root.children[path[0]]
    if !ok {
      next = newCategoryNode()
      root.children[path[0]] = next
    }
    addToCategories(compressed, uncompressed, path[1:], next)
  }
}

func categorize(f *zip.File) categoryPath {
  parts := strings.Split(f.Name, ".")
  if len(parts) >= 2 {
    return categoryPath{parts[len(parts)-1]}
  }
  return categoryPath{"[blank]"}
}

func printCategories(root *categoryNode, name string, prefix string) {
  fmt.Printf("%s%24s  %s (uncompressed: %s)\n", prefix, name, sizeStr(root.compressed), sizeStr(root.uncompressed))
  for childName, child := range root.children {
    printCategories(child, childName, prefix+"\t")
  }
}

func sizeStr(size uint64) string {
  if size > 1024*1024*1024 {
    return fmt.Sprintf("%.1f GB", float64(size)/float64(1024*1024*1024))
  }
  if size > 1024*1024 {
    return fmt.Sprintf("%.1f MB", float64(size)/float64(1024*1024))
  }
  if size > 1024 {
    return fmt.Sprintf("%.1f KB", float64(size)/float64(1024))
  }
  return fmt.Sprintf("%.1f B", float64(size))
}
