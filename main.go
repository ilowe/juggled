/*
The juggled daemon maps docker containers to HTTP hostnames.

Hosts File

You can provide a hosts file to force certain types of mappings pre-hoc by providing
a JSON file and using the --hostmap/-H flag.

The JSON should look like:

	{
		"hosts": [
			{
				"id": "68394fda8458",
				"ssl": true,
				"auth": false
			}
		]
	}

Usage

You can obtain "online" help by running "juggled -h"; here is the output for convenience:

	Usage of juggled:
	  -H, --hostmap="": (optional) a file containing host mappings
	  -w, --http=":80": The address to listen on for HTTP connections
	  -s, --https=":443": The address to listen on for HTTPS connections
	  -q, --quiet=false: be quiet
	  -c, --sslCert="": The certificate file to use for TLS/SSL
	  -k, --sslKey="": The key file to use for TLS/SSL
	  -v, --verbose=false: be verbose
	  -V, --version=false: output version and exit
*/
package main

import (
	"fmt"
	"os"

	flag "github.com/ogier/pflag"

	"github.com/ilowe/juggled/juggler"
	"github.com/ilowe/log"
)

func main() {
	var verbose, quiet bool
	var httpPort, httpsPort string
	var sslCert, sslKey string

	var hostmapFile string
	var outputVersion bool

	flag.BoolVarP(&verbose, "verbose", "v", false, "be verbose")
	flag.BoolVarP(&quiet, "quiet", "q", false, "be quiet")
	flag.BoolVarP(&outputVersion, "version", "V", false, "output version and exit")

	flag.StringVarP(&hostmapFile, "hostmap", "H", "", "(optional) a file containing host mappings")

	flag.StringVar(&httpPort, "http", "", "The address to listen on for HTTP connections")
	flag.StringVar(&httpsPort, "https", "", "The address to listen on for HTTPS connections")

	flag.StringVarP(&sslCert, "sslCert", "c", "", "The certificate file to use for TLS/SSL")
	flag.StringVarP(&sslKey, "sslKey", "k", "", "The key file to use for TLS/SSL")

	flag.Parse()

	if outputVersion {
		fmt.Println(juggler.Version)
		os.Exit(0)
	}

	switch {
	case verbose:
		log.Verbose()
	case quiet:
		log.Quiet()
	default:
		log.Normal()
	}

	j := juggler.NewJuggler(sslCert, sslKey)

	if hostmapFile != "" {
		j.LoadHostmapFile(hostmapFile)
	}

	if httpPort == "" || httpsPort == "" && os.Geteuid() != 0 {
		log.Errorln("Non-root user cannot use ports under 1024!")
		os.Exit(-1)
	}

	j.Juggle(httpPort, httpsPort)
}