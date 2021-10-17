package main

var listen = "127.0.0.1:8000"

var certdir = "/home/martin"

var domainRedirects = map[string]string{
	"arp242.net": "www.arp242.net",
}

// Redirect pages on a domain.
//
// Redirects are always 301.
//
// Files will *overide* the redirects; they're only checked on a 404.
//
// Wildcards (*) can only be at the end, and {} is replaced with the contents.
var redirects = map[string]map[string]string{
	"www.arp242.net": map[string]string{
		// Redirect old /weblog/ paths.
		"/weblog*": "/{}",

		// # Code pages, just redirect to GitHub
		"/code":            "/#code",
		"/code/index.html": "/#code",
		"/code/*":          "https://github.com/arp242/{}",

		// Old weblog titles; use matching for some because Jekyll generated
		// weird filenames over the years :-/
		"/weblog/run_multiple_services_*":    "/pf-switch.html",
		"/weblog/minimal_apache*":            "/apache-svn.html",
		"/weblog/manage_unreal*":             "/ut-cache.html",
		"/weblog/online_unreal*()":           "/ut-browser.html",
		"/weblog/tunnelling_ssh*-l":          "/ssh-tunnel.html",
		"/weblog/creating_temporary*":        "/php-mktemp.html",
		"/weblog/digging_for_hosts*":         "/freebsd-dig.html",
		"/weblog/making_flag_shih*.html":     "/flagshihtzu-formtastic.html",
		"/weblog/intercept_outgoing*":        "/rails-devmail.html",
		"/weblog/uninstalling_emacs*":        "/apt-get.html",
		"/weblog/generate_passwords*":        "/cli-passwords.html",
		"/weblog/security-of-python*":        "/pickle-marshal-security.html",
		"/weblog/i-dont-like-git*":           "/git-hg.html",
		"/weblog/some-thoughts-on-cdns.html": "/cdn.html",
		"/weblog/a-primer*":                  "/python-str-bytes.html",
		"/weblog/yaml*":                      "/yaml-config.html",
		"/weblog/json_as_configuration*":     "/json-config.html",
	},
}
