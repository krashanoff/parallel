# parallel

A small CLI application that looks to take GNU's
[parallel](https://www.gnu.org/software/parallel/)
and run with its ideas in Go.

## Example

```sh
# first, let's create some files.
for i in {1..24}; do
  head -c `expr $i \* 500` /dev/urandom > "test/$i"
done

# cat all files with four threads in
# operation.
parallel -j 4 cat {} \; ./test/*

# cat all files with twelve threads in
# operation, but cull jobs that take
# over 50ms of execution time.
parallel -j 12 -t 50 cat {} \; ./test/*

# taking input from stdin.
find test | grep '1.*' | parallel -j 6 -t 100 cat {} \; -

# a slightly more realistic application:
# threading pandoc on each markdown file
# in a folder of notes.
parallel -j 6 pandoc -f markdown -t latex -o {}.pdf {} \; ./notes/*.md
```

## Running

You can get the binary built for your system from the
[releases page](https://github.com/krashanoff/parallel/releases).

## Building From Source

```sh
go install krashanoff.com/parallel
```

## Stuff to add
* Web-server for hosting and/or distributing jobs.
