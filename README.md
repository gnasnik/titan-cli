# titan-cli

`titan-cli` is a command-line interface tool for interacting with the titan network.

## Usage
```
 » titan-cli --help
NAME:
   titan cli - titan's toolset

USAGE:
   titan cli [global options] command [command options] [arguments...]

COMMANDS:
   download, d  get file from titan network
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help

```

## Example:

Export the titan locator address. 
```shell
export LOCATOR_API_INFO=https://120.78.83.177:5000

```
Use the download command to fetch the desired file. Specify the CID with the -c flag, and the output file location with the -o flag.

```
 » ./titan-cli download -c QmbmPTRvxq8W9DKa9CAHk3GmuEdYBkCXScVFko1D4CiJRn -o /tmp/download.car
```

