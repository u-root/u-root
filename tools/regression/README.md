# regression

CI regression / build status tester for individual u-root cmdlets built with tinygo.
This tool is designed to be run in the CI, but can also be used to run locally.
The default output format is Markdown to stdout.

## Example

```shell
TINYGO=/path/to/tinygo/build/tinygo 
TINYGOVER=$(/path/to/tinygo/build/tinygo version | awk {'print $3'})
GOVER=$(go version | awk {'print $3'} )

go run tools/regression/main.go \
    -tinygo $TINYGO \
    -commit-tinygo $TINYGOVER \
    -commit-uroot $(git rev-parse HEAD) \
    -version-go $(GOVER) \
    cmds/core/ls* cmds/core/sshd
```