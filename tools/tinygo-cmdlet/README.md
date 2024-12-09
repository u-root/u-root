# tinygo-cmdlet

This tool builds all u-root cmdlets passed to it and generates a build report.
The data format is defined for `CSV`, `JSON` and as the output for `Github`.

## CSV
```
subcommand,error
```

## JSON
TODO

## Github
The github output has the following format
```
# tinygo status

## version
* tinygo: <version>
* golang: <version>
* u-root: <commit>

## build status

### completed
* <basename> 

### failed

#### <failed basename>
\`\`\`
<error output>
\`\`\`
```