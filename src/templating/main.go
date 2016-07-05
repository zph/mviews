package main

/*
Credit for this code and go generate concept comes from https://github.com/kelcecil/go-gamesdb/blob/master/generategamesdb/generate_gamesdb_call.go and http://kelcecil.com/golang/2015/01/09/using-go-generate-in-go-1-dot-4.html by Kel Cecil.
*/
import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"unicode"
	"unicode/utf8"
)

type Options struct {
	Output          string
	Kind            string
	LowerKind       string
	Name            string
	LookId          string
	RefreshInterval int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var entryTemplate string
var pwd string

func init() {
	pwd, _ := os.LookupEnv("PWD")
	e, err := ioutil.ReadFile(pwd + "/src/templating/entry.go.tmpl")
	check(err)
	entryTemplate = string(e)
}

func main() {
	options := parseFlags()
	log.Printf("Options %+v", options)
	tmpl, err := template.New("entry").Parse(entryTemplate)
	if err != nil {
		panic(err)
	}
	createDirectoryIfNotExist(options.Output)
	f, err := os.Create(options.Output)
	if err != nil {
		panic(err)
	}
	writer := bufio.NewWriter(f)
	defer func() {
		writer.Flush()
		f.Close()
	}()

	err = tmpl.Execute(writer, options)
	if err != nil {
		panic(err)
	}
}

func createDirectoryIfNotExist(file string) {
	directory := filepath.Dir(file)
	_, err := os.Stat(directory)
	if os.IsNotExist(err) {
		os.MkdirAll(directory, 0755)
	}
}

func lowerFirst(s string) string {
	// Credit: DisposaBoy https://groups.google.com/forum/#!topic/golang-nuts/WfpmVDQFecU
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}

func parseFlags() Options {
	k := flag.String("kind", "", "The type of item being rebuilt")
	n := flag.String("name", "", "name")
	id := flag.String("looker", "", "Looker Endpoint Id")
	interval := flag.Int("interval", 300, "Refresh Interval")
	flag.Parse()
	pwd, _ := os.LookupEnv("PWD")
	var out = pwd + "/src/mviews/generated_" + *n + ".go"

	return Options{
		Output:          out,
		Kind:            *k,
		LowerKind:       lowerFirst(*k),
		Name:            *n,
		LookId:          *id,
		RefreshInterval: *interval,
	}
}
