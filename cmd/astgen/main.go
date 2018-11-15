package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"text/template"

	"github.com/flier/astq/pkg/query"
)

const (
	tplExt = ".tpl"
	goExt  = ".go"
)

var (
	generator   = &Generator{filepath.Base(os.Args[0]), "1.0"}
	outFile     string
	pkgPath     string
	filters     string
	tplName     string
	userVars    string
	showVersion bool
	parseMode   = parser.AllErrors | parser.ParseComments
)

type Generator struct {
	Name, Version string
}

func (g *Generator) String() string {
	return fmt.Sprintf("%s v%s", g.Name, g.Version)
}

func init() {
	flag.StringVar(&userVars, "D", "", "define a key-value pair to parametrize the template (example \"-D key=value\" or \"-D key\")")
	flag.StringVar(&filters, "f", "", "filter names from top-level declarations (example \"-f Foo,Bar\")")
	flag.StringVar(&outFile, "o", "-", "the output file name")
	flag.StringVar(&pkgPath, "p", "-", "the source package import path")
	flag.StringVar(&tplName, "t", "./template/", "the template to use")
	flag.BoolVar(&showVersion, "v", false, "show the version")
}

func openOutput() (w io.Writer, err error) {
	if outFile == "-" {
		w = os.Stdout
	} else if f, err := os.Create(outFile); err == nil {
		w = f
	}

	return
}

func openTemplate() (tpl *template.Template, err error) {
	var u *url.URL

	if u, err = url.Parse(tplName); err == nil && (u.Scheme == "http" || u.Scheme == "https") {
		var res *http.Response

		if res, err = http.DefaultClient.Get(u.RequestURI()); err != nil {
			return
		}

		if res.StatusCode < 200 || res.StatusCode >= 300 {
			return nil, fmt.Errorf("fetch template failed, status: %d (%s)", res.StatusCode, res.Status)
		}

		var body []byte

		if body, err = ioutil.ReadAll(res.Body); err != nil {
			return
		}

		return template.New(filepath.Base(u.Path)).Parse(string(body))
	}

	tplLocations := []string{
		tplName,
		tplName + tplExt,
		"template/" + tplName,
		filepath.Join("template", tplName+tplExt),
	}

	for _, path := range tplLocations {
		if stat, err := os.Stat(path); err == nil {
			var tplFiles []string

			if stat.IsDir() {
				filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if info.IsDir() {
						return filepath.SkipDir
					}
					if strings.HasSuffix(path, tplExt) {
						tplFiles = append(tplFiles, path)
					}
					return nil
				})
			} else {
				tplFiles = append(tplFiles, path)
			}

			if tpl, err = template.ParseFiles(tplFiles...); err != nil {
				return tpl, nil
			}
		} else if matches, err := filepath.Glob(path); err == nil && len(matches) > 0 {
			if tpl, err = template.ParseFiles(matches...); err != nil {
				return tpl, nil
			}
		}
	}

	if tpl == nil {
		err = errors.New("template not found")
	}

	return
}

func parseGoPackage() (pkgs map[string]*ast.Package, err error) {
	fset := token.NewFileSet()

	if pkgPath == "-" {
		return parseGoFile(fset, "-", os.Stdin)
	}

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	pkgLocations := []string{
		pkgPath,
		filepath.Join(".", pkgPath),
		filepath.Join("vendor", pkgPath),
		filepath.Join(gopath, "src", pkgPath),
	}

	for _, path := range pkgLocations {
		var stat os.FileInfo

		if stat, err = os.Stat(path); err == nil {
			if stat.IsDir() {
				pkgs, err = parser.ParseDir(fset, path, nil, parseMode)
			} else {
				pkgs, err = parseGoFile(fset, path, nil)
			}

			break
		}
	}

	return
}

func parseGoFile(fset *token.FileSet, filename string, src io.Reader) (pkgs map[string]*ast.Package, err error) {
	var file *ast.File

	file, err = parser.ParseFile(fset, filename, src, parseMode)
	if err != nil {
		return
	}

	files := make(map[string]*ast.File)
	files[filename] = file

	pkg := &ast.Package{
		Name:  file.Name.Name,
		Files: files,
	}

	pkgs = make(map[string]*ast.Package)
	pkgs[file.Name.Name] = pkg

	return
}

func injectEnvVars(data map[string]interface{}) error {
	inject := func(name, defval string) {
		if env, found := os.LookupEnv(name); found {
			data[name] = env
		} else {
			data[name] = defval
		}
	}

	inject("GOARCH", build.Default.GOARCH)
	inject("GOOS", build.Default.GOOS)
	inject("GOROOT", build.Default.GOROOT)
	inject("GOPATH", build.Default.GOPATH)

	envVars := make(map[string]string)

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)

		if len(parts) > 1 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			envVars[key] = value
		} else {
			key := strings.TrimSpace(env)

			envVars[key] = ""
		}
	}

	data["ENV"] = envVars

	return nil
}

func parseUserDefinedVars(data map[string]interface{}) error {
	if len(userVars) > 0 {
		for _, v := range strings.Split(userVars, ",") {
			parts := strings.SplitN(v, "=", 2)

			if len(parts) == 2 {
				key := parts[0]
				value := parts[1]

				if strings.Contains(value, ".") {
					if n, err := strconv.ParseFloat(v, 64); err == nil {
						data[key] = n
						continue
					}
				} else if n, err := strconv.ParseInt(value, 10, 64); err == nil {
					data[key] = n
					continue
				} else if b, err := strconv.ParseBool(value); err == nil {
					data[key] = b
					continue
				}

				data[key] = value
			} else {
				data[v] = true
			}
		}
	}

	return nil
}

func main() {
	flag.Parse()

	if showVersion {
		fmt.Println(generator)
		return
	}

	data := make(map[string]interface{})

	data["GoVersion"] = runtime.Version()
	data["Generator"] = generator

	if err := injectEnvVars(data); err != nil {
		log.Fatalf("fail to inject environment variables, %v", err)
	}

	if err := parseUserDefinedVars(data); err != nil {
		log.Fatalf("fail to parse user defined variables, %v", err)
	}

	pkgs, err := parseGoPackage()
	if err != nil {
		log.Fatalf("fail to parse GO package `%s`, %v", pkgPath, err)
	}

	data["Packages"] = query.FromPackages(pkgs)

	for _, pkg := range pkgs {
		data["Package"] = query.FromPackage(pkg)

		for _, file := range pkg.Files {
			data["File"] = query.FromFile(file)

			break
		}

		break
	}

	tpl, err := openTemplate()
	if err != nil {
		log.Fatalf("fail to parse template `%s`, %v", tplName, err)
	}

	out, err := openOutput()
	if err != nil {
		log.Fatalf("fail to open output file %s, %v", outFile, err)
	}

	var buf bytes.Buffer

	if err := tpl.Execute(&buf, data); err != nil {
		log.Fatalf("fail to generate template, %v", err)
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatalf("%s\nfail to format generated code, %v", string(buf.Bytes()), err)
	}

	if _, err = out.Write(src); err != nil {
		log.Fatalf("fail to write output file, %v", err)
	}
}
