package main

import (
	"fmt"
	"log"
	"os"

	"git.sr.ht/~sircmpwn/getopt"

	"git.sr.ht/~sircmpwn/go-bare/schema"
)

func main() {
	log.SetFlags(0)
	opts, optind, err := getopt.Getopts(os.Args, "hs:p:")
	if err != nil {
		log.Fatalf("error: %e", err)
	}
	pkg := "gen"
	skip := make(map[string]interface{})
	for _, opt := range opts {
		switch opt.Option {
		case 'p':
			pkg = opt.Value
		case 's':
			skip[opt.Value] = nil
		case 'h':
			log.Println("Usage: gen [-p <package>] [-s <skip type>] <input.bare> <output.go>")
			os.Exit(0)
		}
	}

	args := os.Args[optind:]
	if len(args) != 2 {
		log.Fatal("Usage: gen [-p <package>] <input.bare> <output.go>")
	}
	in := args[0]
	out := args[1]

	inf, err := os.Open(in)
	if err != nil {
		log.Fatalf("error opening %s: %e", in, err)
	}
	defer inf.Close()

	types, err := schema.Parse(inf)
	if err != nil {
		log.Fatalf("error parsing %s: %e", in, err)
	}

	outf, err := os.Create(out)
	if err != nil {
		log.Fatalf("error opening %s for writing: %e", out, err)
	}
	defer outf.Close()
	fmt.Fprintf(outf, "package %s\n", pkg)

	if len(skip) != 0 {
		var typesp []schema.SchemaType
		for _, ty := range types {
			if _, ok := skip[ty.Name()]; !ok {
				typesp = append(typesp, ty)
			}
		}
		types = typesp
	}

	genTypes(outf, types)
}
