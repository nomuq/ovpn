// main.go - simple cert manager
//
// (c) 2018 Sudhi Herle; License GPLv2
//
// This software does not come with any express or implied
// warranty; it is provided "as is". No claim  is made to its
// suitability for any purpose.

package main

import (
	"fmt"
	"os"
	"ovpn/internal"
	"path"
	"strings"

	"github.com/opencoff/go-utils"
	flag "github.com/opencoff/pflag"
)


func main() {
	flag.SetInterspersed(false)

	verFlag := flag.BoolP("version", "", false, "Show version info and quit")
	flag.BoolVarP(&internal.Verbose, "verbose", "v", false, "Show verbose output")

	flag.Usage = func() {
		fmt.Printf(
			`%s - Opinionated OpenVPN cert tool

Usage: %s [options] DB CMD [args..]

Where 'DB' points to the certificate database, and 'CMD' is one of:

    init             Initialize a new CA and cert store
    intermediate-ca  Create a new intermediate CA
    server           Create a new server certificate
    list, show       List one or all certificates in the DB
    export           Export a OpenVPN server or client configuration
    delete           Delete a user and revoke their certificate
    user, client     Create a new user/client certificate
    crl              List revoked certificates or generate CRL
    passwd           Change the DB encryption password

Options:
`, path.Base(os.Args[0]), os.Args[0])
		flag.PrintDefaults()
		os.Stdout.Sync()
		os.Exit(0)
	}

	flag.Parse()

	if *verFlag {
		fmt.Printf("%s - %s [%s; %s]\n", path.Base(os.Args[0]), internal.ProductVersion, internal.RepoVersion, internal.Buildtime)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) < 1 {
		internal.Die("Insufficient arguments!\nTry '%s -h'\n", os.Args[0])
	}

	db := args[0]

	// handle the common case of people forgetting the DB
	switch strings.ToLower(db) {
	case "help", "hel", "he":
		flag.Usage()
		os.Exit(1)
	default:
	}

	args = args[1:]
	if len(args) < 1 {
		internal.Die("Insufficient arguments; missing command!\nTry '%s -h'\n", os.Args[0])
	}

	var cmds = map[string]func(string, []string){
		"init":            internal.InitCmd,
		"server":          internal.ServerCert,
		"user":            internal.UserCert,
		"delete":          internal.Delete,
		"client":          internal.UserCert,
		"export":          internal.ExportCert,
		"show":            internal.ListCert,
		"list":            internal.ListCert,
		"crl":             internal.ListCRL,
		"passwd":          internal.ChangePasswd,
		"intermediate-ca": internal.IntermediateCA,
	}

	words := make([]string, len(cmds))
	for k := range cmds {
		words = append(words, k)
	}
	ab := utils.Abbrev(words)

	cmd := strings.ToLower(args[0])
	canon, ok := ab[cmd]
	if !ok {
		internal.Die("unknown command '%s'; Try '%s --help'", cmd, os.Args[0])
	}

	fp, ok := cmds[canon]
	if !ok {
		internal.Die("can't map command '%s'", canon)
	}

	fp(db, args[1:])
}

