# SKALOGRAM

Education purpose Nerds Instagram Clone

## Level 0 - Command Line App
### Usage
```
$ ./skalogram-cli -help
Usage of ./skalogram-cli:
  -db-path string
        path to the skalogram database (default "./.db.json")
  -image string
        path to the image you want to print (required)
```

### Download

Compiled binaries are available in the [releases page](https://github.com/skale-5/skalogram/releases)

### Compile

```
$ cd cli/
$ go build -o skalogram-cli
$ ./skalogram-cli -help
Usage of ./skalogram-cli:
  -db-path string
        path to the skalogram database (default "./.db.json")
  -image string
        path to the image you want to print (required)
```