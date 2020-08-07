package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	minimock "github.com/gojuno/minimock/v3"
	"github.com/hexdigest/gowrap/generator"
	"github.com/hexdigest/gowrap/pkg"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

var (
	//do not modify the following vars
	//the values are being injected at the compile time by goreleaser
	version   = "dev"
	commit    = "dev"
	buildDate = time.Now().Format(time.RFC3339)
)

var helpers = template.FuncMap{
	"title": strings.Title,
	"in": func(s string, in ...string) bool {
		s = strings.Trim(s, " ")
		for _, i := range in {
			if s != "" && strings.Contains(i, s) {
				return true
			}
		}
		return false
	},
}

type (
	options struct {
		interfaces []interfaceInfo
		noGenerate bool
		suffix     string
	}

	interfaceInfo struct {
		Type       string
		ImportPath string
		WriteTo    string
	}
)

func main() {
	opts, err := processArgs(os.Args[1:], os.Stdout, os.Stderr)
	if err != nil {
		if err == errInvalidArguments {
			os.Exit(2)
		}

		die("%v", err)
	}

	if opts == nil { //help requested
		os.Exit(0)
	}

	if err = run(opts); err != nil {
		die("%v", err)
	}
}

func run(opts *options) (err error) {
	var (
		sourcePackage *packages.Package
		astPackage    *ast.Package
		fs            *token.FileSet
	)

	for i, in := range opts.interfaces {
		if i == 0 || in.ImportPath != opts.interfaces[i-1].ImportPath {
			if sourcePackage, err = pkg.Load(in.ImportPath); err != nil {
				return err
			}

			fs = token.NewFileSet()
			if astPackage, err = pkg.AST(fs, sourcePackage); err != nil {
				return errors.Wrap(err, "failed to load package sources")
			}
		}

		interfaces, err := findInterfaces(astPackage, in.Type)
		if err != nil {
			return err
		}

		gopts := generator.Options{
			SourcePackage:      sourcePackage.PkgPath,
			SourcePackageAlias: "mm_" + sourcePackage.Name,
			HeaderTemplate:     minimock.HeaderTemplate,
			BodyTemplate:       minimock.BodyTemplate,
			HeaderVars: map[string]interface{}{
				"GenerateInstruction": !opts.noGenerate,
				"Version":             version,
			},
			Funcs: helpers,
		}

		if err := processPackage(gopts, interfaces, in.WriteTo, opts.suffix); err != nil {
			return err
		}
	}

	return nil
}

func processPackage(opts generator.Options, interfaces []string, writeTo, suffix string) (err error) {
	for _, name := range interfaces {
		opts.InterfaceName = name
		opts.OutputFile, err = destinationFile(name, writeTo, suffix)
		if err != nil {
			return errors.Wrapf(err, "failed to generate mock for %s", name)
		}

		if err := generate(opts); err != nil {
			return err
		}

		fmt.Printf("minimock: %s\n", opts.OutputFile)
	}

	return nil
}

func isGoFile(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return false, err
		}

		dir := filepath.Dir(path)
		if stat, err = os.Stat(dir); err != nil {
			return false, err
		}

		if !stat.IsDir() {
			return false, errors.Errorf("not a directory: %s", dir)
		}

		return strings.HasSuffix(path, ".go"), nil
	}

	return strings.HasSuffix(path, ".go") && !stat.IsDir(), nil
}

func destinationFile(interfaceName, writeTo, suffix string) (string, error) {
	ok, err := isGoFile(writeTo)
	if err != nil {
		return "", err
	}

	var path string

	if ok {
		path = writeTo
	} else {
		path = filepath.Join(writeTo, minimock.CamelToSnake(interfaceName)+suffix)
	}

	if !strings.HasPrefix(path, "/") && !strings.HasPrefix(path, ".") {
		path = "./" + path
	}

	return path, nil
}

func generate(o generator.Options) (err error) {
	g, err := generator.NewGenerator(o)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer([]byte{})

	if err = g.Generate(buf); err != nil {
		return errors.Wrap(err, "failed to generate mock")
	}

	return ioutil.WriteFile(o.OutputFile, buf.Bytes(), 0644)
}

func findInterfaces(p *ast.Package, pattern string) ([]string, error) {
	var names []string
	for _, f := range p.Files {
		for _, d := range f.Decls {
			if gd, ok := d.(*ast.GenDecl); ok && gd.Tok == token.TYPE {
				for _, spec := range gd.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						if _, ok := ts.Type.(*ast.InterfaceType); ok && match(ts.Name.Name, pattern) {
							names = append(names, ts.Name.Name)
						}
					}
				}
			}
		}
	}

	if len(names) == 0 {
		return nil, errors.Errorf("failed to find any interfaces matching %s in %s", pattern, p.Name)
	}

	return names, nil
}

func match(s, pattern string) bool {
	return pattern == "*" || s == pattern
}

func usage(fs *flag.FlagSet, w io.Writer) {
	const usageTemplate = `Usage: {{bold "minimock"}} [{{bold "-i"}} source.interface] [{{bold "-o"}} output/dir/or/file.go] [{{bold "-g"}}]
{{.}}
Examples:

  Generate mocks for all interfaces that can be found in the current directory:
    {{bold "minimock"}}

  Generate mock for the io.Writer interface and put it into the "./buffer" package:
    {{bold "minimock"}} {{bold "-i"}} io.Writer {{bold "-o"}} ./buffer

  Generate mocks for the fmt.Stringer and all interfaces from the "io" package and put them into the "./buffer" package:
    {{bold "minimock"}} {{bold "-i"}} fmt.Stringer,io.* {{bold "-o"}} ./buffer

For more information please visit https://github.com/gojuno/minimock
`

	t := template.Must(template.New("usage").Funcs(template.FuncMap{
		"bold": func(s string) string { return "\033[1m" + s + "\033[0m" },
	}).Parse(usageTemplate))

	buf := bytes.NewBuffer([]byte{})

	fs.SetOutput(buf)
	fs.PrintDefaults()

	if err := t.Execute(w, buf.String()); err != nil {
		panic(err) //something went completely wrong, i.e. OOM, closed pipe, etc
	}
}

func showVersion(w io.Writer) {
	const versionTemplate = `MiniMock version {{bold .Version}}
Git commit: {{bold .Commit}}
Build date: {{bold .BuildDate}}
`

	t := template.Must(template.New("version").Funcs(template.FuncMap{
		"bold": func(s string) string { return "\033[1m" + s + "\033[0m" },
	}).Parse(versionTemplate))

	versionInfo := struct {
		Version   string
		Commit    string
		BuildDate string
	}{version, commit, buildDate}

	if err := t.Execute(w, versionInfo); err != nil {
		panic(err) //something went completely wrong, i.e. OOM, closed pipe, etc
	}
}

var errInvalidArguments = errors.New("invalid arguments")

func processArgs(args []string, stdout, stderr io.Writer) (*options, error) {
	var opts options

	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.BoolVar(&opts.noGenerate, "g", false, "don't put go:generate instruction into the generated code")
	fs.StringVar(&opts.suffix, "s", "_mock_test.go", "mock file suffix")

	input := fs.String("i", "*", "comma-separated names of the interfaces to mock, i.e fmt.Stringer,io.Reader\nuse io.* notation to generate mocks for all interfaces in the \"io\" package")
	output := fs.String("o", "", "comma-separated destination file names or packages to put the generated mocks in,\nby default the generated mock is placed in the source package directory")
	help := fs.Bool("h", false, "show this help message")
	version := fs.Bool("version", false, "display version information and exit")

	fs.Usage = func() { usage(fs, stderr) }

	if err := fs.Parse(args); err != nil {
		return nil, errInvalidArguments
	}

	if *version {
		showVersion(stdout)
		return nil, nil
	}

	if *help {
		usage(fs, stdout)
		return nil, nil
	}

	interfaces := strings.Split(*input, ",")

	var writeTo = make([]string, len(interfaces))
	if *output != "" {
		//if only one output package specified
		if to := strings.Split(*output, ","); len(to) == 1 {
			for i := range writeTo {
				writeTo[i] = to[0]
			}
		} else {
			writeTo = to
		}
	}

	if len(writeTo) != len(interfaces) {
		return nil, errors.Errorf("count of the source interfaces doesn't match the output files count")
	}

	if err := checkDuplicateOutputFiles(writeTo); err != nil {
		return nil, err
	}

	for i, in := range interfaces {
		info, err := makeInterfaceInfo(in, writeTo[i])
		if err != nil {
			return nil, err
		}

		opts.interfaces = append(opts.interfaces, *info)
	}

	return &opts, nil
}

// checkDuplicateOutputFiles finds first non-unique Go file
func checkDuplicateOutputFiles(fileNames []string) error {
	for i := range fileNames {
		ok, err := isGoFile(fileNames[i])
		if err != nil {
			return err
		}

		if !ok {
			continue
		}

		ipath, err := filepath.Abs(fileNames[i])
		if err != nil {
			return err
		}

		for j := range fileNames {
			jpath, err := filepath.Abs(fileNames[j])
			if err != nil {
				return err
			}

			if j != i && ipath == jpath {
				return errors.Errorf("duplicate output file name: %s", ipath)
			}
		}
	}

	return nil
}

func makeInterfaceInfo(typ, writeTo string) (*interfaceInfo, error) {
	info := interfaceInfo{WriteTo: writeTo}

	dot := strings.LastIndex(typ, ".")
	slash := strings.LastIndex(typ, "/")

	if slash > dot {
		return nil, errors.Errorf("invalid interface type: %s", typ)
	}

	if dot >= 0 {
		info.Type = typ[dot+1:]
		info.ImportPath = typ[:dot]
	} else {
		info.Type = typ
		info.ImportPath = "./"
	}

	return &info, nil
}

func die(format string, args ...interface{}) {
	if _, err := fmt.Fprintf(os.Stderr, "minimock: "+format+"\n", args...); err != nil {
		panic(err)
	}

	os.Exit(1)
}
