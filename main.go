package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func main() {
	sep := flag.String("sep", "\n", "separator")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: strmatch regexp [regexp ...]")
		os.Exit(2)
	}

	errCh := make(chan error)

	output := readLines(os.Stdin, errCh)
	output = filterLines(output, errCh, args)
	printLines(output, errCh, []byte(*sep))

	err := <-errCh
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func printLines(input chan string, errCh chan error, sep []byte) {
	go func() {
		defer func() { errCh <- nil }()
		count := 0
		for s := range input {
			if count > 0 {
				os.Stdout.Write(sep)
			}
			count++
			os.Stdout.Write([]byte(s))
		}
	}()
}

func filterLines(input chan string, errCh chan error, args []string) chan string {
	output := make(chan string)

	go func() {
		reList := make([]*regexp.Regexp, len(args))

		for i, a := range args {
			re, err := regexp.Compile(a)
			if err != nil {
				errCh <- err
				return
			}
			reList[i] = re
		}

		defer close(output)

		for s := range input {
			s = strings.TrimRight(s, "\r\n")
			for _, re := range reList {
				if groups := re.FindStringSubmatch(s); len(groups) > 0 {
					if len(groups) > 1 {
						groups = groups[1:]
					}
					for _, g := range groups {
						output <- g
					}
					break
				}
			}
		}
	}()

	return output
}

func readLines(input io.Reader, errCh chan error) chan string {
	output := make(chan string)

	go func() {
		defer close(output)
		reader := bufio.NewReader(os.Stdin)
		for {
			s, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					errCh <- err
				}
				return
			}

			output <- s
		}
	}()

	return output
}
