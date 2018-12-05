package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/kr/pretty"
)

var api *http.ServeMux

type Product map[string]interface{}

var cacheProducts []Product

func loadProducts(filename string) ([]Product, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	var data = make(map[string]interface{})
	dec := json.NewDecoder(f)
	err = dec.Decode(&data)
	if err != nil {
		return nil, err
	}
	var products []Product
	log.Printf("data: %# v", pretty.Formatter(data))

	if data, ok := data["data"].(map[string]interface{}); ok {
		if data, ok = data["products"].(map[string]interface{}); ok {
			if edges, ok := data["edges"].([]interface{}); ok {
				for _, node := range edges {
					if node, ok := node.(map[string]interface{}); ok {
						if product, ok := node["node"].(map[string]interface{}); ok {
							products = append(products, product)
						}
					}
				}
				return products, nil
			}
		}
	}

	return nil, errors.New("empty products")
}

func Products() []Product {
	if cacheProducts == nil {
		cacheProducts, _ = loadProducts("./products.json")
	}

	return cacheProducts
}

func getProductList(w http.ResponseWriter, r *http.Request) {
	var (
		products = Products()
	)

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(products); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	var (
		vars     = mux.Vars(r)
		products = Products()
		id       = vars["id"]
	)

	log.Printf("ID: %s", id)

	w.Header().Set("Content-Type", "application/json")

	for _, prd := range products {
		log.Printf("prd: %#v", prd)
		if prdId, ok := prd["id"].(string); ok {
			log.Printf("prdId %s", prdId)
			if prdId == id {
				enc := json.NewEncoder(w)
				if err := enc.Encode(prd); err != nil {
					writeErrorJson(w, err, http.StatusInternalServerError)
					return
				}

				return
			}
		}
	}

	writeErrorJson(w, errors.New("没有找到这个产品"), http.StatusNotFound)
}

func ajsPath() string {
	actpath := os.Getenv("ACTIVITYJS_PATH")
	if len(actpath) == 0 {
		curDir, _ := os.Getwd()
		return curDir
	} else {
		return actpath
	}
}

func workdir() string {
	actpath := os.Getenv("WORKDIR")
	if len(actpath) == 0 {
		curDir, _ := os.Getwd()
		return filepath.Join(curDir, "workspace")
	} else {
		return actpath
	}
}

func trackAndDel(dir string, timeout time.Duration) {
	defer func() {
		log.Printf("Cleaning dir: %s", dir)
		os.RemoveAll(dir)
	}()
	<-time.After(timeout)

}

func worker(w http.ResponseWriter, r *http.Request, fn func(gomain string) error) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErrorJson(w, err, http.StatusInternalServerError)
		return
	}

	dir, err := ioutil.TempDir(workdir(), "/tmp/activityjs")
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		go trackAndDel(dir, 10*time.Second)
	}()

	tmpfn := filepath.Join(dir, "main.go")
	if err := ioutil.WriteFile(tmpfn, buf, 0666); err != nil {
		writeErrorJson(w, err, http.StatusInternalServerError)
		return
	}

	if err := fn(tmpfn); err != nil {
		writeErrorJson(w, err, http.StatusInternalServerError)
		return
	}
}

func compile(w http.ResponseWriter, r *http.Request) {
	t := time.Now()

	worker(w, r, func(gomain string) error {
		log.Printf("compile %s", gomain)
		dir := filepath.Dir(gomain)
		outfile := filepath.Join(dir, "main.js")
		debugs, ok := r.URL.Query()["debug"]
		var debug bool
		if ok && len(debugs) > 0 {
			fmt.Sscanf(debugs[0], "%v", &debug)
		}
		log.Printf("%v", debugs)
		log.Printf("debug mode: %v", debug)
		output, err := gopherjs(gomain, outfile, workdir(), debug)
		log.Printf("COMPILE: %s", output)
		if err != nil {
			log.Printf("Compile Error %s", err)
			return err
		}

		fi, e := os.Stat(outfile)
		if e != nil {
			log.Printf("Stat Error %s", e)
			return e
		}
		// get the size
		size := fi.Size()
		ms := time.Since(t) / time.Millisecond

		fmt.Printf("Compiled file :%s is %d bytes long", outfile, fi.Size())
		renderJson(w, map[string]interface{}{
			"success":  true,
			"size":     size,
			"js":       filepath.Join("/compiled", strings.Replace(outfile, workdir(), "", -1)),
			"duration": ms,
		})
		return nil
	})

}

func downloadjs(w http.ResponseWriter, r *http.Request) {
	var (
		vars     = mux.Vars(r)
		filename = vars["filename"]
	)
	log.Printf("path: %s", r.URL.Path)
	log.Printf("filename %s", filename)

	fullname := filepath.Join(workdir(), filename)
	f, err := os.Open(fullname)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()
	defer os.RemoveAll(filepath.Base(fullname)) // clean up

	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	io.Copy(w, f)
}

func vet(w http.ResponseWriter, r *http.Request) {
	worker(w, r, func(gomain string) error {
		log.Printf("Vet %s", gomain)

		output, _ := govet(gomain, workdir())
		log.Printf("`VET`: %s", output)
		// if err != nil {
		// 	return err
		// }

		lineErrs := extractErrors(bytes.NewBuffer(output))

		renderJson(w, map[string]interface{}{
			"success":    true,
			"lineErrors": lineErrs,
		})
		return nil
	})
}

func gopherjs(pkgname string, output string, workdir string, debug bool) ([]byte, error) {
	var args = []string{
		"build",
		pkgname,
	}

	if debug {
		args = append(args, "--tags=debug")
	}

	args = append(args, "-o", output)
	os.RemoveAll(filepath.Join(workdir, "pkg/darwin_js"))
	os.RemoveAll(filepath.Join(workdir, "pkg/linux_js"))
	os.RemoveAll(filepath.Join(workdir, "pkg/windows_js"))

	build := exec.Command("gopherjs", args...)
	log.Printf("build args %v", args)
	build.Env = []string{"GOPATH=" + workdir, "GO111MODULE=off"}
	// build.Env = os.Environ()
	build.Dir = workdir
	return build.CombinedOutput()

}

func govet(pkgname string, workdir string) ([]byte, error) {
	build := exec.Command("go", "vet", pkgname)
	// "GOCACHE=" + "/Users/hysios/Library/Caches/go-build"
	build.Env = []string{"GOPATH=" + workdir, "GO111MODULE=off"}
	build.Dir = workdir
	return build.CombinedOutput()
}

var re = regexp.MustCompile(`^([\w\/.]+):(\d+):(\d+):\s+(.*?)$`)

type LineError struct {
	Filename string `json:"filename"`
	Line     int    `json:"line"`
	Pos      int    `json:"pos"`
	Msg      string `json:"msg"`
}

func lastLine(errs *[]LineError) *LineError {
	if len(*errs) == 0 {
		*errs = append(*errs, LineError{})
	}

	return &(*errs)[len(*errs)-1]
}

func atoi(in string) int {
	n, _ := strconv.Atoi(in)
	return n
}

func extractErrors(in io.Reader) []LineError {
	var (
		scanner = bufio.NewScanner(in)
		start   bool
		errs    = make([]LineError, 0)
	)

	for scanner.Scan() {
		txt := scanner.Text()

		matchs := re.FindAllStringSubmatch(txt, -1)
		if len(matchs) > 0 {
			match := matchs[0]
			start = true
			lerr := LineError{
				Filename: match[1],
				Line:     atoi(match[2]),
				Pos:      atoi(match[3]),
				Msg:      match[4],
			}
			errs = append(errs, lerr)
		} else if start {
			lerr := lastLine(&errs)
			lerr.Msg += txt
		}
	}

	return errs
}

func examples(w http.ResponseWriter, r *http.Request) {
	var (
		vars     = mux.Vars(r)
		filename = vars["filename"]
	)
	log.Printf("path: %s", r.URL.Path)
	log.Printf("filename %s", filename)
	dir, _ := os.Getwd()
	fullname := filepath.Join(dir, "examples", filename)
	f, err := os.Open(fullname)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	// w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// io.Copy(w, f)
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		writeErrorJson(w, err, http.StatusInternalServerError)
		return
	}
	fi, _ := f.Stat()
	renderJson(w, map[string]interface{}{
		"file": string(buf),
		"size": fi.Size(),
	})
}

func renderJson(w http.ResponseWriter, val map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	if err := enc.Encode(val); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func errorJson(err error, status int) []byte {
	var errs = map[string]interface{}{
		"error": err.Error(),
	}

	buf, _ := json.Marshal(errs)
	return buf
}

func writeErrorJson(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	fmt.Fprintf(w, "%s", errorJson(err, status))
}

func main() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/products", getProductList)
	api.HandleFunc("/products/{id}", getProduct)
	api.HandleFunc("/compile", compile)
	api.HandleFunc("/vet", vet)

	api.HandleFunc("/examples/{filename:.*}", examples)
	r.HandleFunc("/compiled/{filename:.*}", downloadjs)

	log.Printf("Listen on port http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
