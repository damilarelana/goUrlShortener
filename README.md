### goUrlShortener

A simple Golang implementation of a Url Shortener that:

* examines path of incoming request
* determines if re-direction is required
* uses flags (`-yaml`, `-json`, `-sql`) to source the required content from files or database instead of inline strings

***

The code leverages the following packages:

* [yaml](gopkg.in/yaml.v3)
* [json](https://golang.org/pkg/encoding/json/)
* [sql](https://golang.org/pkg/database/sql/)
* [postgres driver](https://github.com/lib/pq)
* `flags`
* `net/http`
* `net/url`
* `fmt`
* `io`
* `log`
* `error`
* `reflect`
* `errors`
* `os`


***

### Example
To run the code nominally without any flags, start the local webserver by running:
```bash
    $ ./main/main
```
Then test the fallback `MapHandler`, by pointing browser to `127.0.0.1:8080/urlshort-godoc`. The browser should redirect you to `https://godoc.org/github.com/gophercises/urlshort`. 

To test the `YAMLHandler`, re-run the application with the YAML file flag i.e.
```bash
    $ ./main/main -yaml="pathsData.yaml"
```
This uses the default ymal file (`pathsData.ymal`). You can change this file, by providing (your custom yaml file's full path) to the the `-yaml` flag i.e. `-yaml="<pathToYAMLFile>"`. Test this scenario by pointing browser to `127.0.0.1:8080/urlshort-final-yaml`. The browser should redirect to `https://github.com/damilarelana/goUrlShortener/tree/master/main`.



To test the `JSONHandler`, re-run the application with the JSON file flag i.e.
```bash
    $ ./main/main -json="pathsData.json"
```
This uses the default json file (`pathsData.json`). You can change this file, by providing (your custom json file's full path) to the the `-json` flag i.e. `-json="<pathToJSONFile>"`. Test this scenario by pointing browser to `127.0.0.1:8080/urlshort-final-json`. The browser should redirect to `https://github.com/damilarelana/goUrlShortener/tree/master/main`.


To test `SQLHandler`, you supply the sql database path in the format `<protocol>://<host>:<port>/<database>?dbUser=<dbUserValue>&dbUserPassword=<dbUserPasswordValue>`.
```bash
    $ ./main/main -sql="http//127.0.0.1:5432/go_test_db?dbUser=postgres&dbUserPassword=brainiac"
```
 This specific local PostgresSQL database instance was tested by pointing browser to `127.0.0.1:8080/urlshort-final-sql`. The browser redirects to `https://github.com/damilarelana/goUrlShortener/tree/master/main`
***

### To Do

+ [x] YAML implementation - accepts a YAML file as a flag, loads the YAML content from file instead of from a string
+ [x] JSON implementation - accepts a JSON file as a flag, loads the JSON content from file instead of from a string
+ [x] SQL implementation  - reads the required content from a database instead of from a file
