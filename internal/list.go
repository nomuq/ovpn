// list.go -- list one or many user certs
//
// (c) 2018 Sudhi Herle; License GPLv2
//
// This software does not come with any express or implied
// warranty; it is provided "as is". No claim  is made to its
// suitability for any purpose.

package internal

import (
	"fmt"
	"os"
	"time"

	"github.com/opencoff/go-pki"
	flag "github.com/opencoff/pflag"
)

func ListCert(db string, args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	fs.Usage = func() {
		listUsage(fs)
	}

	var showCA bool
	var envpw string

	fs.BoolVarP(&showCA, "ca", "", false, "Display the CA certificate")
	fs.StringVarP(&envpw, "env-password", "E", "", "Use password from environment var `E`")

	err := fs.Parse(args)
	if err != nil {
		Die("%s", err)
	}

	ca := OpenCA(db, envpw)
	defer ca.Close()

	if showCA {
		fmt.Printf("CA Certificate:\n%s\n", Cert(*ca.Certificate))
	}

	args = fs.Args()

	if len(args) == 0 {
		// always print the abbreviated root-CA
		c := &pki.Cert{
			Certificate: ca.Certificate,
		}
		printcert(c, true)

		var certs []*pki.Cert

		certs, err := ca.GetServers()
		if err != nil {
			Die("can't fetch servers: %s", err)
		}

		users, err := ca.GetClients()
		if err != nil {
			Die("can't fetch users: %s", err)
		}
		certs = append(certs, users...)

		cas, err := ca.GetCAs()
		if err != nil {
			Die("can't fetch CAs: %s", err)
		}
		for i := range certs {
			printcert(certs[i], false)
		}

		for i := range cas {
			c := cas[i]
			if c.SerialNumber.Cmp(ca.SerialNumber) == 0 {
				continue
			}
			z := &pki.Cert{
				Certificate: c.Certificate,
				IsCA:        true,
			}
			printcert(z, false)
		}

		return
	}

	for _, cn := range args {
		c, err := ca.Find(cn)
		if err != nil {
			warn("Can't find Common Name %s", cn)
			continue
		}
		printcert(c, false)
	}
}

func printcert(z *pki.Cert, rootCA bool) {
	var pref, server string

	now := time.Now().UTC()
	if now.After(z.NotAfter) {
		pref = fmt.Sprintf("EXPIRED %s", z.NotAfter)
	} else {
		pref = fmt.Sprintf("valid until %s", z.NotAfter)
	}

	if z.IsServer {
		server = "server"
	} else if z.IsCA {
		server = "CA (I)"
	} else if rootCA {
		server = "Root-CA"
	}

	fmt.Printf("%-16s  %7.7s %#x (%s)\n", z.Subject.CommonName, server, z.SerialNumber, pref)
	Print("%s\n", Cert(*z.Certificate))
}

func listUsage(fs *flag.FlagSet) {
	fmt.Printf(`%s list: List one or more issued certificates

Usage: %s DB list [options] [CN...]

Where 'DB' is the CA Database file and 'CN' is common name in the certificate.

Options:
`, os.Args[0], os.Args[0])

	fs.PrintDefaults()
	os.Exit(0)
}
