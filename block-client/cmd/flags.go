package main

import "flag"

func ParseFlags(p *params) {
	flag.StringVar(&p.ServerAddress, "a", p.ServerAddress, "address to run server")
	flag.StringVar(&p.NoncePattern, "n", p.NoncePattern, "nonce pattern to use")
	flag.Parse()
}
