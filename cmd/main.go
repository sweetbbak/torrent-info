package main

import (
	"fmt"
	"log"
	"os"
	"time"
	torrent "torrent-info"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	Files   bool `short:"f" long:"files" description:"only list files"`
	Verbose bool `short:"v" long:"verbose" description:"print debugging information and verbose output"`
}

var Debug = func(string, ...interface{}) {}

func PTorrent(args []string) error {
	for _, file := range args {
		fi, err := os.Open(file)
		if err != nil {
			return err
		}
		defer fi.Close()

		tor, err := torrent.Parse(fi)
		if err != nil {
			return err
		}

		printout(tor, file)
	}

	return nil
}

func printout(t *torrent.Torrent, file string) {
	fmt.Printf("\x1b[4mFile:\x1b[0m %s\n", file)
	fmt.Printf("\x1b[4mHash:\x1b[0m %s\n", t.InfoHash)

	d := t.CreatedAt
	if !d.Equal(time.Unix(0, 0)) {
		fmt.Printf("\x1b[4mDate:\x1b[0m %s\n", t.CreatedAt.String())
	}

	if t.CreatedBy != "" {
		fmt.Printf("\x1b[4mCreated\x1b[0m by: %s\n", t.CreatedBy)
	}
	if t.Comment != "" {
		fmt.Printf("\x1b[4mComment:\x1b[0m %s\n", t.Comment)
	}

	fmt.Println("\n\x1b[4m-----Files------\x1b[0m")

	for _, f := range t.Files {
		for _, ff := range f.Path {
			fmt.Printf("%s\n", ff)
		}
	}
	fmt.Println("\x1b[4m----------------\x1b[0m")
	fmt.Println()
}

func main() {
	args, err := flags.Parse(&opts)
	if flags.WroteHelp(err) {
		os.Exit(0)
	}
	if err != nil {
		log.Fatal(err)
	}

	if opts.Verbose {
		Debug = log.Printf
	}

	if err := PTorrent(args); err != nil {
		log.Fatal(err)
	}
}
