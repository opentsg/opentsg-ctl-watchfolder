# app semver - use Semantic Versioning in your app

Package semver manages data for reporting semantic versions of your app.

The goal of the package is to provide a structured way to include semantic
version, hash and note information in a CI/CD pipeline:

```sh
# things you can do with the library
app -v          # → 1.2.3-rc
app --version   # → Full App Title 1.2.3-rc+f3c6
app --note      # → Some commit message maybe
app --gitHash   # → d57e36f5f84302c1          - first 16 chars
app --arch      # → arm64                     - GOARCH - expected architecture
```

This should allow you to manage docker images & documentation by using the app
to generate its own version.

Note that idiomatic golang required a prefix of **v** when adding versions to
packages. The underlying [semver][sv] spec does not require a **v** prefix.

This package does not use a **v** prefix and assumes that you'll insert the
letter **v** when needed in your workflow.

## information flow

* `/releases.yaml` is used to track the high level releases of the app
* `semver/semver.go` mixes the yaml and the linker data to display a version
* Linker data is provided with the go build command:
  ```sh
  go build -ldflags "-X gitlab.com/mrmxf/clog/cloglib/semver.SemVerInfo=\"'x_x_x_x_HASH_x_x_x_x_x_VALUE_x_x_x_x_AS_x_x_x_x_x_TEXT_x_x_x_x_x_x|2024-07-09||myclog|Command_Line_Of_Go'" .
  ```
  The format of the linker string is `commithash|date|suffix|appname|apptitle`:
  * `commithash` - the 40 digit hash from the repo used to build clog
  * `date` - an ISO 8601 date string representing the release/build date
  * `appName` - the command name usually used to run the program e.g clog
  * `app title` - a text string to describe the title. Due to a limitation in
    the construction of linker strings, **all spaces must be represented by
    underscores**. The build will replace underscores with spaces at runtime.
    See the sample above

## usage

Create an embedded yaml (or json) file to track the releases.

```golang
// Typically this is in the root folder making it easy to find
//go:embed releases.yaml
var ClogFs embed.FS
```

```yaml
# Dates must be in YYYY-MM-DD international ISO 8601 format
- version:   "0.4.3"
  date:      2024-06-25
  codename:  drainage
  note:      public release candidate
- version:   "0.4.2"
  date:      2024-02-06
  codename:  drainage
  note:      Probably the commit message
```

The package assumes that the first entry is the one that is to be used and that
the rest are some kind of history - even if you're regressing the semantic version
because you're going backwards to a previous branch.

You can add extra fields, however, the ones shown are required.


## linking

A release build script can inject LDLinkerData{} variables:

```shell
  GOOS=$OS GOARCH=$CPU go build -ldflags \
 "-X main.Los='mac' -X main.Lcpu='arm64'  -X main.Lhash=$(git rev-list -1 HEAD) -X main.Ldate=$DT -X main.Lappname=$APP -X Lsuffix="rc" \
 -o tmp/executable
````

to use in your code:

```
    package main
    dummy data to be overridden by linker injection for production
    var Los = "default"
    var Lcpu = "default"
    var Lcommit = "default"
    var Ldate = "default"
    var Lsuffix = "default"
    var Lappname = "default"
    //...
    semver.Initialise(semver.LinkerData{
        BuildOs: Los,
        BuildCpu: Lcpu,
    })
    //...
    fmt.Printf("current version=%s", semver.Info.Short)
```

[sv]: https://semver.org/
