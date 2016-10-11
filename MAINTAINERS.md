# For maintainers only

# We use the github hub tool for code review. We use govendor for maintaining dependencies.

### Setup your u-root Github Repository

Follow the instructions for using hub. At some point, you'll need to fork github.com/u-root/u-root, then
get it via go get or gitclone. In any event, you u-root repo should end up in
$GOPATH/src/github.com/u-root/u-root

``u-root`` uses [govendor](https://github.com/kardianos/govendor) for its dependency management.

### To manage dependencies

#### Add new dependencies

  - Edit your code to import foo/bar
  - Run `govendor add +external` from the top level

#### Remove dependencies 

  - Run `govendor remove foo/bar`

#### Update dependencies

  - Run `govendor remove +vendor`
  - Run `govendor add +external`
