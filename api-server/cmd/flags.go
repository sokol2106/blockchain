package main

import "flag"

func ParseFlags(p *params) {
	flag.StringVar(&p.RunAddress, "a", p.RunAddress, "address to run server")
	flag.StringVar(&p.DatabaseDSN, "d", p.DatabaseDSN, "data connection Database")
	flag.Parse()
}
