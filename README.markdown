Simple static file server in Go.

Features:

- Directory-based vhosts (convention over configuration).
- Compile all static resources in to the binary.
- Redirects for pages (/old-url → new-url) and domains (foo.com → www.foo.com).

Usage
-----
Sites are in `www/<hostname>/files`:

    % mkdir -p www/example.com
    % echo '<h1>Hello, world!</h1>' > www/example.com/index.html

If you want to use something outside of `www` you need to set up a bind mount:

    % mkdir -p www/mysite.com
    % doas mount -obind,ro ~/arp242.net www/www.arp242.net

You can use wildcards if you want (`STAR.example.com`, or `STAR` for a
catch-all). This needs to be in all-caps (we can't use `*` as Go embed doesn't
like it).

Now just build and start it:

    $ go build
    $ ./zsrv

There is a little `build.sh` script which ensures there's no cgo dependencies.

The file size is about 7.7M (stripped) or 3M (compressed with upx), plus the
website data you have.

This is not intended for very large websites. I use it to host my own website
(https://www.arp242.net).

Configuration
-------------
Edit `config.go` and recompile.
