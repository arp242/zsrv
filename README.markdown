Simple static file server in Go.

Features:

- Directory-based vhosts (convention over configuration).
- Compile all static resources in to the binary.
- Redirects for pages (/old-url → new-url) and domains (foo.com → www.foo.com).

Usage
-----

Create directories or links for your files:

    $ mkdir example.com
    $ echo '<h1>Hello, world!</h1>' > example.com/index.html
    $ ln -s ~/my-site mysite.com

You can use wildcards if you want (`*.example.com`, or `*` for a catch-all).

Generate the `pack.go` to store all of the data in the binary:

    $ go generate

And now start it:

    $ go build
    $ ./zsrv

There is a little `build.sh` script to run all the steps.

The file size is about 7.7M (stripped) or 3M (compressed with upx), plus the
website data you have.

This is not intended for very large websites. I use it to host my own website
(https://arp242.net).

Configuration
-------------

Edit `config.go` and recompile.
