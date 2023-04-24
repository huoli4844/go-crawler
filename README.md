# go-crawler

[go-crawler](https://github.com/lizongying/go-crawler)

## Usage

* log.filename: Log file path. You can replace {name} with ldflags.
* log.long_file: If set to true, the full file path is logged.
* log.level: DEBUG/INFO/WARN/ERROR

example

```shell
git clone github.com/lizongying/go-crawler-example
```

build

```shell
make
```

run

```shell
./releases/youtubeSpider -c example.yml
```
