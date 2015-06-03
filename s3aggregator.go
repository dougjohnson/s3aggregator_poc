package main

import (
  "gopkg.in/amz.v1/aws"
  "gopkg.in/amz.v1/s3"
  "fmt"
  "time"
  "runtime"
  "os"
)

func main() {
  runtime.GOMAXPROCS(runtime.NumCPU())

  cnt := 500

  start := time.Now()

  auth, err := aws.EnvAuth()
  if err != nil {
      panic(err.Error())
  }

  s := s3.New(auth, aws.USEast)
  bucket := s.Bucket("dj-sample-data")

  c := make(chan int)

  for i := 1; i <= cnt; i++ {
    go put(bucket, i, c)
  }

  written := 0

  for written < cnt {
    written += <-c
  }

  fmt.Fprintf(os.Stderr, "wrote %d files in %v\n", cnt, time.Since(start))

  start = time.Now()

  r,err := bucket.List("tmp/", "/", "", 1000)
  if err != nil {
      panic(err.Error())
  }

  fmt.Fprintf(os.Stderr, "read list of %d in %v\n", len(r.Contents), time.Since(start))

  start = time.Now()
  c2 := make(chan []byte)
  for i := 0; i < len(r.Contents); i++ {
    path := r.Delimiter + r.Contents[i].Key
    go get(bucket, path, c2)
  }

  var aggregate []byte
  for i := 0; i < len(r.Contents); i++ {
    aggregate = append(aggregate, <-c2...)
  }
  //fmt.Printf("aggregated files:\n%s", string(aggregate))
  fmt.Fprintf(os.Stderr, "aggregated %d files in %v\n", len(r.Contents), time.Since(start))

}

func put(bucket *s3.Bucket, i int, c chan int) {
  bytes := []byte(fmt.Sprintf("file number %d\n", i))

  err := bucket.Put(fmt.Sprintf("tmp/%d",i), bytes, "binary/octet-stream", s3.Private)
  if err != nil {
    panic(err.Error())
  }
  c <- 1
}

func get(bucket *s3.Bucket, path string, c chan []byte) {
  bytes,err := bucket.Get(path)
  if err != nil {
    fmt.Println(path)
    panic(err.Error())
  }
  c <- bytes
}
