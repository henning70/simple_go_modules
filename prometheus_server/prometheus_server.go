package prometheus_server

import (
    "fmt"
    "os"
    "net/http"
    "flag"
    "sync"
    "runtime"

    "github.com/rs/zerolog"

    "github.com/prometheus/common/promlog"
    "github.com/prometheus/exporter-toolkit/web"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"

    "github.com/henning70/simple_go_modules"
)


var (
   downdetector_metrics = promauto.NewGaugeVec(prometheus.GaugeOpts{
       Name: "downdetector_metrics",
       Help: "Downdetector metrics",
   },
       []string{"name", "country", "company"},
   )

   downdetector_error = promauto.NewGaugeVec(prometheus.GaugeOpts{
       Name: "downdetector_error",
       Help: "Downdetector error information",
   },
       []string{"name", "country", "company", "error"},
   )

   module_main = "server.go"
   mux = http.NewServeMux()
   mu  sync.Mutex

   logger      = promlog.New(&promlog.Config{})

   defaultExporter      string
   defaultListenAddress string
   defaultTLSConfig     = "tls_config.yml"
)

func Init() {
   webConfig, listenAddress, logLevel := init_webconfig()

   // Initialise logging
   init_logger(logLevel)

   go init_exporter_listener(listenAddress, webConfig)
}

func init_logger(logLevel *string) {
    if *logLevel == "debug" {
        logging.Debug = true
        logging.DebugLogFile, _ = os.Create(logging.DebugLog)
        logging.Debugging = zerolog.New(logging.DebugLogFile).Level(zerolog.DebugLevel)

        logmsg := fmt.Sprintln("Debugging started")
        logging.DebugLogging(module_main, logmsg)
    }
}

func init_webconfig() (web.FlagConfig, *string, *string) {
    // Initialise command line flags
    commandLine := flag.NewFlagSet("downdetector_exporter", flag.ExitOnError)

    // Initialise default values
    var (
        listenAddress = commandLine.String("web.listen-address", defaultListenAddress, "Address to listen on for web interface and telemetry.")
        tlsConfigFile = commandLine.String("web.config.file", defaultTLSConfig, "Path to config yaml file that can enable TLS or authentication.")
        logLevel      = commandLine.String("log.level", "info", "Only log messages with the given severity or above. Valid levels: [debug, info, warn, error].")
    )

    // Initialise web config for ListenAndServe
    webConfig := web.FlagConfig{
        WebListenAddresses: func() *[]string { a := make([]string, 1); return &a }(),
        WebSystemdSocket:   func() *bool { b := false; return &b }(),
        WebConfigFile:      tlsConfigFile,
    }

    // Set value for webConfig.WebSystemdSocket if OS is Linux
    if runtime.GOOS == "linux" {
        webConfig.WebSystemdSocket = commandLine.Bool("web.systemd-socket", false, "Use systemd socket activation listeners instead of port listeners (Linux only).")
    }

    // Parse command line flags
    _ = commandLine.Parse(os.Args[1:])

    return webConfig, listenAddress, logLevel
}

func init_exporter_listener(listenAddress *string, webConfig web.FlagConfig) {
    mux.Handle("/metrics", promhttp.Handler())

    server := &http.Server{}
    server.Handler = mux

    fmt.Printf("Starting exporter at %s", *listenAddress)

    (*webConfig.WebListenAddresses)[0] = *listenAddress

    if err := web.ListenAndServe(server, &webConfig, logger); err != nil {
        fmt.Printf("err: %v\n", err)
        fmt.Printf("web.ListenAndServe terminated with, %s", err)
        os.Exit(1)
    }
}
