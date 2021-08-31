// crl.go -- list one or many revoked certs and export them
//
// (c) 2018 Sudhi Herle; License GPLv2
//
// This software does not come with any express or implied
// warranty; it is provided "as is". No claim  is made to its
// suitability for any purpose.

package internal

import (
	"fmt"
	"io"
	"os"

	flag "github.com/opencoff/pflag"
)

func ListCRL(db string, args []string) {
	fs := flag.NewFlagSet("crl", flag.ExitOnError)
	fs.Usage = func() {
		crlUsage(fs)
	}

	var list bool
	var outfile, envpw string
	var crlvalid int

	fs.BoolVarP(&list, "list", "l", false, "List revoked certificates")
	fs.StringVarP(&outfile, "outfile", "o", "", "Write the CRL  to `F`")
	fs.IntVarP(&crlvalid, "validity", "V", 1, "Make the CRL valid for `N` days")
	fs.StringVarP(&envpw, "env-password", "E", "", "Use password from environment var `E`")

	err := fs.Parse(args)
	if err != nil {
		Die("%s", err)
	}

	ca := OpenCA(db, envpw)
	defer ca.Close()

	var out io.Writer = os.Stdout
	if len(outfile) > 0 && outfile != "-" {
		fd := mustOpen(outfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
		defer fd.Close()

		out = fd
	}

	if !list {
		pem, err := ca.CRL(crlvalid)
		if err != nil {
			Die("%s", err)
		}

		out.Write(pem)
	} else {
		rv, err := ca.ListRevoked()
		if err != nil {
			Die("can't list revoked certs: %s", err)
		}

		for _, z := range rv {
			fmt.Fprintf(out, "%-16s  %#x revoked on %s\n", z.Subject.CommonName, z.SerialNumber, z.When)
		}

		if err != nil {
			Die("%s", err)
		}
	}
}

func crlUsage(fs *flag.FlagSet) {
	fmt.Printf(`%s crl: Generate a new CRL or list revoked users

Usage: %s DB crl [options]

Where 'DB' is the CA Database file.

Options:
`, os.Args[0], os.Args[0])

	fs.PrintDefaults()
	os.Exit(0)
}
