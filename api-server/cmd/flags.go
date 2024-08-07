package main

import "flag"

func ParseFlags(p *params) {
	flag.StringVar(&p.ServerAddress, "a", p.ServerAddress, "address to run server")
	flag.StringVar(&p.BaseAddress, "b", p.BaseAddress, "base address of the resulting shortened URL")
	flag.StringVar(&p.FileStoragePath, "f", p.FileStoragePath, "file storage path")
	flag.StringVar(&p.DatabaseDSN, "d", p.DatabaseDSN, "data connection Database")
	flag.Parse()
}
