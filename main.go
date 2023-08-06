package wstunnel

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	log "github.com/fangdingjun/go-log/v5"
	"gopkg.in/yaml.v2"
)

func main() {
	var cfgfile string
	var logfile string
	var loglevel string
	var payload string
	flag.StringVar(&cfgfile, "c", "config.yaml", "config file")
	flag.StringVar(&logfile, "log_file", "", "log file")
	flag.StringVar(&loglevel, "log_level", "INFO", "log level")
	flag.StringVar(&payload, "payload", "free.facebook.com", "payload")
	flag.Parse()

	data, err := os.ReadFile(cfgfile)
	if err != nil {
		log.Fatal(err)
	}

	var cfg conf
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatal(err)
	}

	if logfile != "" {
		log.Default.Out = &log.FixedSizeFileWriter{
			MaxCount: 4,
			Name:     logfile,
			MaxSize:  10 * 1024 * 1024,
		}
	}

	if lv, err := log.ParseLevel(loglevel); err == nil {
		log.Default.Level = lv
	}

	MakeServers(cfg, payload)

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	select {
	case s := <-ch:
		log.Printf("received signal %s, exit.", s)
	}
}
