### goUrlShortener

A simple Golang implementation of a Url Shortener that:

* examines path of incoming request
* determines if re-direction is required
* uses flags (`-yaml`, `-json`, `-sql`) to sources the required content from files instead of inline strings

***

The code leverages the following packages:

* `flags` package
* `net` package
* `fmt` package
* `gorilla mux` package

### To Do

+ [ ] YAML implementation - accepts a YAML file as a flag, loads the YAML content from file instead of from a string
+ [ ] JSON implementation - accepts a JSON file as a flag, loads the JSON content from file instead of from a string
+ [ ] SQL implementation  - reads the required content from a database instead of from a file