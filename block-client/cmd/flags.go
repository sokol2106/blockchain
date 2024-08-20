package main

import "flag"

func ParseFlags(p *params) {
	flag.StringVar(&p.ServerURL, "a", p.ServerURL, "URL server")
	flag.StringVar(&p.NoncePattern, "n", p.NoncePattern, "nonce pattern to use")
	flag.Parse()
}
