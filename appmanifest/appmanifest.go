package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/groob/plist"
)

// MaxChunkSize is the max size of each file chunk that needs to be hashed
const MaxChunkSize = 10 << 20 // 10MB

// http://help.apple.com/deployment/osx/#/ior5df10f73a
type manifest struct {
	ManifestItems []manifestItem `plist:"items"`
}

type manifestItem struct {
	Assets []asset `plist:"assets"`
	// Apple claims the metadata struct is required,
	// but testing shows otherwise.
	Metadata *metadata `plist:"metadata,omitempty"`
}

type asset struct {
	Kind    string   `plist:"kind"`
	MD5Size int64    `plist:"md5-size"`
	MD5s    []string `plist:"md5s"`
	URL     string   `plist:"url"`
}

type metadata struct {
	bundleInfo
	Items    []bundleInfo `plist:"items,omitempty"`
	Kind     string       `plist:"kind"`
	Subtitle string       `plist:"subtitle"`
	Title    string       `plist:"title"`
}

type bundleInfo struct {
	BundleIdentifier string `plist:"bundle-identifier"`
	BundleVersion    string `plist:"bundle-version"`
}

var (
	version = "unreleased"
	gitHash = "unknown"
)

const usage = `appmanifest [options] /path/to/some.pkg`

func main() {
	flVersion := flag.Bool("version", false, "prints the version")
	flURL := flag.String("url", "", "url of the pkg as it will be on the server")

	// set usage
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", usage)
		flag.PrintDefaults()
	}
	flag.Parse()

	if *flVersion {
		fmt.Printf("appmanifest - %v\n", version)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) == 0 {
		log.Println("must specify a path to a pkg")
		fmt.Println(usage)
		os.Exit(1)
	}

	path := args[0]
	if err := createAppManifest(path, *flURL, os.Stdout); err != nil {
		log.Fatal(err)
	}

}

// create manifest and return back a writer
func createAppManifest(path, url string, writer io.Writer) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// get file info
	info, err := file.Stat()
	if err != nil {
		return err
	}
	size := info.Size()

	// create a list of md5s
	md5s, err := calculateMD5s(file)
	if err != nil {
		return err
	}

	// create an asset
	ast := asset{
		Kind:    "software-package",
		MD5Size: size,
		MD5s:    md5s,
		URL:     url,
	}

	// make a manifest
	m := manifest{
		ManifestItems: []manifestItem{
			manifestItem{
				Assets: []asset{ast},
			},
		},
	}

	// write a plist
	enc := plist.NewEncoder(writer)
	enc.Indent("  ")
	return enc.Encode(&m)
}

// reads a file and returns a slice of hashes, one for each
// 10mb chunk
func calculateMD5s(file io.Reader) ([]string, error) {
	var md5s []string
	buf := make([]byte, MaxChunkSize)
	for {
		n, err := file.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			return nil, err
		}
		md5s = append(md5s, fmt.Sprintf("%x", md5.Sum(buf)))
	}
	return md5s, nil
}
