`snl`
===

Sample N number or percentage of Lines from a file or STDIN.

Usage:

        $ snl [[sample size]%] [file path]
        $ snl [[sample size]%] -

I thought it would be fun to write a BASH script to return random lines from a file. Turns out that people have already written one-liners to return one line from a file as written here:

[http://www.hilarymason.com/blog/how-to-get-a-random-line-from-a-file-in-bash/](http://www.hilarymason.com/blog/how-to-get-a-random-line-from-a-file-in-bash/)

Once I started looking at their solutions it became clear it would be about as fun to write a program in C that would return N random lines from a file; but that feeling didn't last because I wanted more instant gratification. So this is my attempt at a Go implementation.


Install:

        $ go get github.com/stuntgoat/snl
        $ go install github.com/stuntgoat/snl

TODO:

      - tests
      - manpage
