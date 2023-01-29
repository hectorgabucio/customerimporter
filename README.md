# Customer Importer - Héctor

This library reads the provided input file and returns a sorted array with the domains with more occurrences along with the occurence number.

It does it by splitting the file in chunks and letting independent goroutines read those chunks and extract the domains.

## Features
- You can provide your own store backend, being responsible for the sorting.
- You can provide a logger if you want to log errors when trying to read lines from the CSV.
- Flexible for big files and small machines: you can customize the chunk size (in bytes) and the concurrency (size of the goroutines workers pool) to fit your needs.

## Example usage

```go
import "example.com/interview/customerimporter"

package main

func main() {
    c := customerimporter.New()
    // or customerimporter.WithOptions() to customize and extend the functionality

    data, err := c.Import("customers.csv")
}
```

## Tests
- You can lint and run tests with ```make```
- You can run benchmarks using ```make bench```
- You can check the Makefile for more.