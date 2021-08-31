package internal

import (
	"crypto/x509"
	"fmt"
	"github.com/opencoff/go-pki"
	"os"
)


// This will be filled in by "build"
var RepoVersion string = "UNDEFINED"
var Buildtime string = "UNDEFINED"
var ProductVersion string = "UNDEFINED"

var Verbose bool

type Cert x509.Certificate

func (z Cert) String() string {
	c := x509.Certificate(z)
	s, err := pki.CertificateText(&c)
	if err != nil {
		s = fmt.Sprintf("can't stringify %x (%s)", c.SerialNumber, err)
	}
	return s
}

// Only show output if needed
func Print(format string, v ...interface{}) {
	if Verbose {
		s := fmt.Sprintf(format, v...)
		if n := len(s); s[n-1] != '\n' {
			s += "\n"
		}
		os.Stdout.WriteString(s)
		os.Stdout.Sync()
	}
}
