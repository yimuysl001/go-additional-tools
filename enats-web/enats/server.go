package enats

import (
	"flag"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/nats-io/nats-server/v2/server"
	"go.uber.org/automaxprocs/maxprocs"
	"os"
)

var usageStr = `
Usage: nats-server [options]

Server Options:
    -a, --addr, --net <host>         Bind to host address (default: 0.0.0.0)
    -p, --port <port>                Use port for clients (default: 4222)
    -n, --name
        --server_name <server_name>  Server name (default: auto)
    -P, --pid <file>                 File to store PID
    -m, --http_port <port>           Use port for http monitoring
    -ms,--https_port <port>          Use port for https monitoring
    -c, --config <file>              Configuration file
    -t                               Test configuration and exit
    -sl,--signal <signal>[=<pid>]    Send signal to nats-server process (ldm, stop, quit, term, reopen, reload)
                                     <pid> can be either a PID (e.g. 1) or the path to a PID file (e.g. /var/run/nats-server.pid)
        --client_advertise <string>  Client URL to advertise to other servers
        --ports_file_dir <dir>       Creates a ports file in the specified directory (<executable_name>_<pid>.ports).

Logging Options:
    -l, --log <file>                 File to redirect log output
    -T, --logtime                    Timestamp log entries (default: true)
    -s, --syslog                     Log to syslog or windows event log
    -r, --remote_syslog <addr>       Syslog server addr (udp://localhost:514)
    -D, --debug                      Enable debugging output
    -V, --trace                      Trace the raw protocol
    -VV                              Verbose trace (traces system account as well)
    -DV                              Debug and trace
    -DVV                             Debug and verbose trace (traces system account as well)
        --log_size_limit <limit>     Logfile size limit (default: auto)
        --max_traced_msg_len <len>   Maximum printable length for traced messages (default: unlimited)

JetStream Options:
    -js, --jetstream                 Enable JetStream functionality
    -sd, --store_dir <dir>           Set the storage directory

Authorization Options:
        --user <user>                User required for connections
        --pass <password>            Password required for connections
        --auth <token>               Authorization token required for connections

TLS Options:
        --tls                        Enable TLS, do not verify clients (default: false)
        --tlscert <file>             Server certificate file
        --tlskey <file>              Private key for server certificate
        --tlsverify                  Enable TLS, verify client certificates
        --tlscacert <file>           Client certificate CA for verification

Cluster Options:
        --routes <rurl-1, rurl-2>    Routes to solicit and connect
        --cluster <cluster-url>      Cluster URL for solicited routes
        --cluster_name <string>      Cluster Name, if not set one will be dynamically generated
        --no_advertise <bool>        Do not advertise known cluster information to clients
        --cluster_advertise <string> Cluster URL to advertise to other servers
        --connect_retries <number>   For implicit routes, number of connect retries
        --cluster_listen <url>       Cluster url from which members can solicit routes

Profiling Options:
        --profile <port>             Profiling HTTP port

Common Options:
    -h, --help                       Show this message
    -v, --version                    Show version
        --help_tls                   TLS help
`

// usage will print out the flag options for the server.
func usage() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}

func NatsServer(name ...string) {

	exe := "nats-server"
	if len(name) > 0 && len(name[0]) > 0 {
		exe = name[0]
	}

	// Create a FlagSet and sets the usage
	fs := flag.NewFlagSet(exe, flag.ExitOnError)
	fs.Usage = usage

	flagboll := true
	for _, arg := range os.Args {
		if arg == "-c" || arg == "--config" {
			flagboll = false
			break
		}

	}
	if flagboll {
		var confPath = []string{"nats.conf", "config/nats.conf", "resources/nats.conf"}
		for _, s := range confPath {
			_, err := os.Stat(s)
			if err == nil {
				os.Args = append(os.Args, "-c", s)
				break
			}

		}
	}

	// Configure the options from the flags/config file
	opts, err := server.ConfigureOptions(fs, os.Args[1:],
		server.PrintServerAndExit,
		fs.Usage,
		server.PrintTLSHelpAndDie)
	if err != nil {
		g.Log().Error(gctx.GetInitCtx(), err.Error())
		os.Exit(0)
		//server.PrintAndDie(fmt.Sprintf("%s: %s", exe, err))
	} else if opts.CheckConfig {
		g.Log().Errorf(gctx.GetInitCtx(), "%s: configuration file %s is valid (%s)\n", exe, opts.ConfigFile, opts.ConfigDigest())
		//fmt.Fprintf(os.Stderr, "%s: configuration file %s is valid (%s)\n", exe, opts.ConfigFile, opts.ConfigDigest())
		os.Exit(0)
	}
	if opts.LogFile != "" && !gfile.Exists(gfile.Dir(opts.LogFile)) {
		gfile.Mkdir(gfile.Dir(opts.LogFile))
	}
	if opts.PidFile != "" && !gfile.Exists(gfile.Dir(opts.PidFile)) {
		gfile.Mkdir(gfile.Dir(opts.PidFile))
	}
	s, err := server.NewServer(opts)
	if err != nil {
		server.PrintAndDie(fmt.Sprintf("%s: %s", exe, err))
	}
	// Configure the logger based on the flags.
	s.ConfigureLogger()
	gctx.SetInitCtx(gctx.New())

	g.Log().Debug(gctx.GetInitCtx(), "服务启动")
	// Start things up. Block here until done.
	if err := server.Run(s); err != nil {
		server.PrintAndDie(err.Error())
	}
	g.Log().Debug(gctx.GetInitCtx(), "启动完成")
	// Adjust MAXPROCS if running under linux/cgroups quotas.
	undo, err := maxprocs.Set(maxprocs.Logger(s.Debugf))
	if err != nil {
		s.Warnf("Failed to set GOMAXPROCS: %v", err)
	} else {
		defer undo()
	}

	s.WaitForShutdown()
}
