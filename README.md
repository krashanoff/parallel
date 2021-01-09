# parallel

Software that looks to take GNU's [parallel](https://www.gnu.org/software/parallel/)
and run with it in Go.

## Example

```sh
# first, let's create some files.
for i in {1..24}; do
  head -c `expr $i \* 500` /dev/urandom > "$i"
done

# cat all files with four threads in
# operation.
parallel -j 4 cat {} ./*

# cat all files with twelve thread in
# operation, but limit each job to
# 500ms of execution time.
parallel -j 12 -t 500 cat {} ./*
```

## Stuff on the horizon
* Pass stdin as `-` in place of files list.
* Use `{}` wherever the current file's name should be substituted.
