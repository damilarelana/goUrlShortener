### goUrlShortener

A simple Golang implementation of a Url Shortener that:

* examines path of incoming request
* determines if re-direction is required
* uses flags (`-yaml`, `-json`, `-sql`) to sources the required content from files instead of inline strings

***

The code leverages the following packages:

*	[yaml](gopkg.in/yaml.v2)
* `flags` package
* `net` package
* `fmt` package
* `io` package
* `log` package
* `error` package
* `reflect` package


***

### Example
To run the code, we need to start the local webserver by running:
```bash
    $ ./main/main
```

To test the `MapHandler`, point browser to `127.0.0.1:8080/urlshort-godoc`. The browser should redirect you to `https://godoc.org/github.com/gophercises/urlshort` . In order to test the `YAMLHandler`, point browser to `127.0.0.1:8080/urlshort`. The browser should redirect you to `https://github.com/damilarelana/goUrlShortener`.


***

### To Do

+ [ ] YAML implementation - accepts a YAML file as a flag, loads the YAML content from file instead of from a string
+ [ ] JSON implementation - accepts a JSON file as a flag, loads the JSON content from file instead of from a string
+ [ ] SQL implementation  - reads the required content from a database instead of from a file