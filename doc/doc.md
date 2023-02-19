# Difference with Glob for Windows and Linux
19/12/2023

`filepath.Match(file, pattern)` is not same between Windows and Linux.
On Linux `filepath.Match("test1.txt", "*.txt")=false` and on Windows
`filepath.Match("test1.txt", "*.txt")=true`
