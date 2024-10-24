# kbcli

this is the cli app for kbkitt

## How to run?

* show help by default.

```sh
make run
```

* or using go command.

```sh
go run cmd/kbcli/main.go
```

## How to test?

```sh
make test
```

## How to setup?

Run `configure` command to setup your environment

```sh
kbkitt configure
```

you will be asked to configure kbkitt in your machine typically in `{HOME_DIR}/.kbkitt`.

```sh
.
├── config.yaml
├── media
│   └── btc.jpeg
└── sync.yaml
```

config.yaml will contain basic information about where the remote server is, in case you want to centralize your kbs.

```yaml
➜  .kbkitt cat config.yaml
version: 0.1.0
fileForSyncPath: {HOME_DIR}/.kbkitt/sync.yaml
dirForMediaPath: {HOME_DIR}/.kbkitt/media
server:
    url: http://localhost:3030
```

* `fileForSyncPath` file that will keep kbs you cannot sent to the central server.
* `dirForMediaPath` directory that will store the media resources (images, docs, videos, etc) that you save as kbs.
* `server.url` kbkitt remote server.

## How to Add a KB?

Once you have setup your kbkitt, this is how you could add new kbs.

```sh
➜  kbkitt add --help

add a new knowledge base such as: concepts, commands, prompts, etc.

Usage:
  kb add [flags]

Flags:
  -c, --class string       category of knowledge base
  -h, --help               help for add
  -k, --key string         knowledge base key
  -n, --notes string       knowledge base notes
  -r, --reference string   author or refence of this kb
  -t, --tags strings       comma separated tags for this kb
  -u, --ux                 add KB in interactive mode
  -v, --value string       knowledge base value
```

The application will ask you for all the necessary parameters if you did not enter them from the beginning.

```sh
➜  kbkitt add
```
...
```yaml
key:
value:
notes:
class:
reference:
tags:
```

in case you want to have more freedom to organize your kb data, you can use an improved cli gui.

```sh
➜  kbkitt add -u
# OR
➜  kbkitt add --ux
```

you will see this prompt

```sh
 Adding a new KB:

Key

Category

Value
...

Notes

Reference

Tags
keyword1 keyword2 keyword3 keywordN

Continue ->

• tab fields • shift+tab fields • ctrl+c: quit
```

If you want to save all at once, you can use all the parameters that the add command has.

```sh
➜  kbkitt add \
-k btc -v crypto -n currencies -c crypto \
-t btc,crypto,currencies,blockchain \
-r dementor

...KB to save...

Key: btc
Value: crypto
Notes: currencies
Category: crypto
Reference: dementor
Tags: [btc crypto currencies blockchain]


> do you want to save it? [y/n]:
```

## References

* bubbletea examples
https://github.com/charmbracelet/bubbletea/blob/master/examples