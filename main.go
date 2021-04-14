package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"krashanoff.com/parallel/pkg/pattern"
)

//
// parallel
//
// clean-room adaptation of gnu's parallel
// with no external deps.
//
// krashanoff
//

const ABOUT = `parallel - run a command on many threads

Usage:
parallel [flags] [command+flags] \; [list of files]

Example:
parallel -j 4 -t 6000 sed -i '$ a\New line' {} \; ./folder/*.md

Use "parallel ... \; -" to read newline-delimited file paths from stdin.

Flags:`

type JobStatus struct {
	id  uint64
	err error
}

func main() {
	deadline := flag.Uint64("t", 0, "Set the maximum runtime in `milliseconds` for each execution. Defaults to no deadline (zero).")
	threads := flag.Uint64("j", 1, "Set the number of `routines` to use.")
	quiet := flag.Bool("q", false, "Suppress output from subprocesses' stdout.")
	flag.Usage = func() {
		fmt.Println(ABOUT)
		flag.PrintDefaults()
	}
	flag.Parse()

	argStart := len(os.Args) - flag.NArg()
	fileOffset := -1
	for i, arg := range os.Args[argStart:] {
		if arg == ";" {
			fileOffset = argStart + i
		}
	}

	// Check args
	switch {
	case fileOffset == -1:
		log.Fatalln("Missing ';' from input.")
	case *threads == 0:
		log.Fatalln("Incorrect thread count!")
	}

	program := strings.Join(os.Args[argStart:fileOffset], " ")
	files := os.Args[fileOffset+1:]
	numFiles := len(files)
	if numFiles == 1 && files[0] == "-" {
		log.Println("Using stdin for input.")
		if buf, err := ioutil.ReadAll(os.Stdin); err != nil {
			log.Println("Failed to read stdin.")
		} else {
			files = strings.Split(string(buf), "\n")
		}
	}

	if len(program) == 0 {
		log.Fatalln("No program supplied.")
	}

	startTime := time.Now()
	log.Printf("Started execution at %v", startTime)

	// Spawn workers
	work, done := make(chan []string, numFiles), make(chan JobStatus, numFiles)
	for i := uint64(0); i < *threads; i++ {
		id := i
		go func() {
			for cmdArgs := range work {
				// Create timeout context
				ctx, cancel := context.Background(), func() {}
				if *deadline != 0 {
					ctx, cancel = context.WithTimeout(context.Background(), time.Duration(*deadline)*time.Millisecond)
				}
				defer cancel()

				// Setup command, redirect pipes
				cmd := exec.CommandContext(ctx, cmdArgs[0])
				if len(program) > 1 {
					cmd = exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
				}
				if !*quiet {
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
				}

				log.Printf("Thread ID %d assigned job \"%v\"", id, strings.Join(cmdArgs, " "))
				done <- JobStatus{
					id,
					cmd.Run(),
				}
			}
		}()
	}

	// Create work, transforming input pattern.
	cmdArgs, err := pattern.GeneratePatterns(program, files)
	if err != nil {
		log.Fatalf("Failed to parse input string: %v", err)
	}
	for _, c := range cmdArgs {
		work <- strings.Split(c, " ")
	}
	close(work)

	// Wait for jobs to complete, exit.
	successful, culled := 0, 0
	for range files {
		job := <-done
		if job.err != nil {
			log.Printf("ERROR: Thread %d reported error: %v", job.id, job.err)
			culled++
		} else {
			successful++
		}
	}

	endTime := time.Now()
	log.Printf("Operation terminated at %v.", endTime)
	log.Printf("Total Time: %v", endTime.Sub(startTime))
	log.Printf("Successful Jobs: %d", numFiles-culled)
	log.Printf("Culled Jobs: %d", culled)

	os.Exit(0)
}
