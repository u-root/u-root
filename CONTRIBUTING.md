# Contributing to u-root

We need help with this project, so contributions are welcome.

## Code of Conduct

Conduct collaboration around u-root in accordance to the [Code of
Conduct](https://github.com/u-root/u-root/wiki/Code-of-Conduct).

## Communication

- [Join slack](https://u-root.slack.com) (Get an invite [here](http://slack.u-root.com).)
- [Join the mailing list](https://groups.google.com/forum/#!forum/u-root)

## Coding Style

The ``u-root`` project aims to follow the standard formatting recommendations
and language idioms set out in the [Effective Go](https://golang.org/doc/effective_go.html)
guide, for example [formatting](https://golang.org/doc/effective_go.html#formatting)
and [names](https://golang.org/doc/effective_go.html#names).

`gofmt` and `golint` are law, although this is not automatically enforced
yet and some housecleaning needs done to achieve that.

We have a few rules not covered by these tools:

- Standard imports are separated from other imports. Example:
    ```
    import (
      "regexp"
      "time"

      dhcp "github.com/krolaw/dhcp4"
    )
    ```

- ``u-root`` uses [govendor](https://github.com/kardianos/govendor)
for its dependency management.  Re-vendoring will generally be
handled by the [maintainers](MAINTAINERS.md).

## Developer Sign-Off

For purposes of tracking code-origination, we follow a simple sign-off
process.  If you can attest to the [Developer Certificate of
Origin](https://developercertificate.org/) then you append in each git
commit text a line such as:
```
Signed-off-by: Your Name <username@youremail.com>
```
## Patch Format

Well formatted patches aide code review pre-merge and code archaeology in
the future.  The abstract form should be:
```
<component>: Change summary

More detailed explanation of your changes: Why and how.
Wrap it to 72 characters.
See [here] (http://chris.beams.io/posts/git-commit/)
for some more good advices.

Signed-off-by: <contributor@foo.com>
```

An example from this repo:
```
tcz: quiet it down

It had a spurious print that was both annoying and making
boot just a tad slower.

Signed-off-by: Ronald G. Minnich <rminnich@gmail.com>
```

## Pull Requests

We accept github pull requests.

Fork the project on github, work in your fork and in branches, push
these to your github fork, and when ready, do a github pull requests
against https://github.com/u-root/u-root.

## Code Reviews

Look at the area of code you're modifying, its history, and consider
tagging some of the [maintainers](MAINTAINERS.md).  when doing a
pull request in order to instigate some code review.

## Quality Controls

This needs enhancing.  ``scripts/`` and ``travis.sh`` include some
initial plumbing.  [Travis CI](https://travis-ci.org/) will run on
your github fork and its branches and also on your PR's to ``u-root``.

## Discussion

``u-root`` is on Slack as "u-root.slack.com".  Please sign up
[here](http://slack.u-root.com/).

Issues can be reported via Github.  A PR which addresses a Github issue can
reference that with a simple "Fixes #NNN" on a line in the commit message,
where "NNN" is the issue number.
