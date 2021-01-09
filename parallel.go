package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
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
parallel [flags] [command+flags] {} [list of files]

Example:
parallel -j 4 -t 6000 sed -i '$ a\New line' {} ./folder/*.md

Flags:`

func main() {
	deadline := flag.Uint64("t", 0, "Set the maximum runtime in `milliseconds` for each execution. Defaults to no deadline (zero).")
	threads := flag.Uint64("j", 1, "Set the number of `routines` to use.")
	flag.Usage = func() {
		fmt.Println(ABOUT)
		flag.PrintDefaults()
	}
	flag.Parse()

	argStart := len(os.Args) - flag.NArg()

	fileOffset := -1
	for i, arg := range os.Args[argStart:] {
		if arg == "{}" {
			fileOffset = argStart + i
		}
	}

	// Check args
	switch {
	case fileOffset == -1:
		log.Fatalln("Missing '{}' from input.")
	case *threads == 0:
		log.Fatalln("Incorrect thread count!")
	}

	program := os.Args[argStart:fileOffset]
	files := os.Args[fileOffset+1:]

	if len(program) == 0 {
		log.Fatalln("No program supplied.")
	}

	startTime := time.Now()
	log.Printf("Started execution at %v", startTime)

	// Spawn workers
	work, done := make(chan string), make(chan bool, len(files))
	for id := uint64(0); id < *threads; id++ {
		go func() {
			for f := range work {
				// Create timeout context
				ctx, cancel := context.Background(), func() {}
				if *deadline != 0 {
					ctx, cancel = context.WithTimeout(context.Background(), time.Duration(*deadline)*time.Millisecond)
				}
				defer cancel()

				// Setup command, redirect pipes
				cmd := exec.CommandContext(ctx, program[0], f)
				if len(program) > 1 {
					cmd = exec.CommandContext(ctx, program[0], append(program[1:], f)...)
				}
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				if err := cmd.Run(); err != nil {
					log.Printf("Encountered error on thread ID %d: %v", id, err)
				} else {
					log.Printf("Thread ID %d ran on file %s", id, f)
				}
				done <- true
			}
		}()
	}

	// Create work
	for _, f := range files {
		work <- f
	}
	close(work)

	// Wait for jobs to complete, exit.
	for range files {
		<-done
	}
	endTime := time.Now()
	log.Printf("Operation terminated at %v. Time taken: %v", endTime, endTime.Sub(startTime))
	os.Exit(0)
}
