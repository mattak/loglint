# loglint

lint build log by your custom rules. 

## Install

from homebrew 

```
$ brew tap mattak/loglint
$ brew install loglint
```

from source

```
$ git clone git@github.com:mattak/loglint.git
$ cd loglint
$ make
```

## Usage 

Prepare log file to inspect.

```
$ echo "Error: 1" >> build.log
$ echo "2" >> build.log
$ echo "Error: 3" >> build.log
```

Prepare lint rules to inspect.

```
$ cat << __LINT__ > '.loglint.json'
[
  {
    "name": "sample",
    "type": "error",
    "detections": ["^Error: "],
    "help": "you can fix error by this link https://..."
  }
]
__LINT__
```


detect from stdin

```
$ cat build.log | loglint | jq .
```

detect from file

```
$ loglint build.log | jq .
{
  "passed": false,
  "errors": [
    {
      "matches": [
        {
          "start": 0,
          "end": 1,
          "message": "Error: 1"
        },
        {
          "start": 2,
          "end": 3,
          "message": "Error: 3"
        }
      ],
      "help": "you can fix error by this link https://...",
      "name": "sample"
    }
  ],
  "warnings": null
}
```

detect by argument rule

```
$ cat build.log | loglint -e '[{"detections":["^Error: "]}]' | jq .
{
  "passed": true,
  "errors": null,
  "warnings": [
    {
      "matches": [
        {
          "start": 0,
          "end": 1,
          "message": "Error: 1"
        },
        {
          "start": 2,
          "end": 3,
          "message": "Error: 3"
        }
      ],
      "help": "",
      "name": ""
    }
  ]
}
```

detect by file rule. (default: .loglint.json)

```
$ cat build.log | loglint -f '.loglint.json.sample' | jq .
```

detect multiple lines

```
$ cat build.log | loglint -e '[{"detections":["^Error: 1", "^Error: 3"]}]' | jq .
{
  "passed": true,
  "errors": null,
  "warnings": [
    {
      "matches": [
        {
          "start": 0,
          "end": 3,
          "message": "Error: 1\n2\nError: 3"
        }
      ],
      "help": "",
      "name": ""
    }
  ]
}
```

## LICENSE

[MIT](./LICENSE.md)
