package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: strmatch regexp [regexp ...]")
		os.Exit(2)
	}

	errCh := make(chan error)

	output := readLines(os.Stdin, errCh)
	output = filterLines(output, errCh, args)
	printLines(output, errCh)

	err := <-errCh
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func printLines(input chan string, errCh chan error) {
	go func() {
		defer func() { errCh <- nil }()
		for s := range input {
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
			for _, re := range reList {
				if groups := re.FindStringSubmatch(s); len(groups) > 0 {
					if len(groups) > 1 {
						groups = groups[1:]
					}
					for _, g := range groups {
						output <- g
						output <- "\n"
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
