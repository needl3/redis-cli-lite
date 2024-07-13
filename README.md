## Redis client written in go 
[![Go](https://img.shields.io/badge/go-1.22.4-blue.svg)](https://golang.org/)

You can use this on cli but that has limitations. For example, there is not really a way to tell if the value is string, number or array from cli. So, everything is stored as string. We could implement something intelligent like:
1. Parse if input is number or string
2. Support for RESP from cli but this is unnecessarily complex and unusable


## TODOS
- [ ] Implement pretty printer 
- [ ] Change string type value to []byte - tokenizer.go:70
- [ ] TODO: Use _ instead of ignoring - tokenizer.go:95
- [ ] Use native integer instead of string - tokenizer.go:104
- [ ] Expand as golang library
- [ ] Support for more datatypes
