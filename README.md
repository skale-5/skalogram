# SKALOGRAM

Education purpose Nerds Instagram Clone

## Level 0 - Command Line App

![rick.jpg](docs/rick.png)


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

Dont forget to make binary executable: `chmod +x skalogram-cli`

### Compile

```
$ cd cli/
mkdir build/
export GOOS=linux   && export GOARCH=amd64 && go build -o "build/skalogram-cli_${GOOS}_${GOARCH}"
export GOOS=linux   && export GOARCH=386   && go build -o "build/skalogram-cli_${GOOS}_${GOARCH}"
export GOOS=darwin  && export GOARCH=amd64 && go build -o "build/skalogram-cli_${GOOS}_${GOARCH}"
export GOOS=darwin  && export GOARCH=arm64 && go build -o "build/skalogram-cli_${GOOS}_${GOARCH}"
export GOOS=windows && export GOARCH=amd64 && go build -o "build/skalogram-cli_${GOOS}_${GOARCH}.exe"
export GOOS=windows && export GOARCH=386   && go build -o "build/skalogram-cli_${GOOS}_${GOARCH}.exe"

$ ./skalogram-cli -help
Usage of ./skalogram-cli:
  -db-path string
        path to the skalogram database (default "./.db.json")
  -image string
        path to the image you want to print (required)
```