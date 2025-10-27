package main

import (
	"context"
	"dumper/internal/app"
	conf "dumper/internal/config/local"
	"dumper/internal/decrypt"
	"dumper/pkg/logging"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var (
	version = "dev"
	date    = "unknown"
)
var showVersion bool

func init() {
	flag.BoolVar(&showVersion, "version", false, "Print version info")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sign
		cancel()
	}()

	configPath := flag.String("config", "./config.yaml", "The path to the configuration file")
	dbName := flag.String("db", "", "Name of the backup database")
	all := flag.Bool("all", false, "Backup of all databases from the configuration")
	fileLog := flag.String("file-log", "dumper.log", "Log files from the configuration")
	dec := flag.Bool("dec", false, "Decrypt")
	decFile := flag.String("input", "", "Decrypt path backup file")
	crypt := flag.String("crypt", "", "Type encrypt [aes]")
	pass := flag.String("pass", "", "Password to decrypt backup file [required if crypt aes]")

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		_, _ = fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *dec {
		decApp := decrypt.NewApp(*decFile, *pass, *crypt)
		err := decApp.Decrypt()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Decrypt success")
		return
	}

	if showVersion {
		fmt.Printf("Version: %s \nDate: %s\n", version, date)
		return
	}

	env := app.Env{
		ConfigFile: *configPath,
		DbName:     *dbName,
		All:        *all,
		FileLog:    *fileLog,
	}

	config, err := conf.Load(*configPath)
	if err != nil {
		fmt.Printf("configuration loading error : %v \n", err)
		os.Exit(1)
	}

	logger := runLog(&env, *config.Settings.Logging)

	defer func(logger *logging.Logs) {
		_ = logger.Close()
	}(logger)

	ctx = logging.ContextWithLogger(ctx, logger.Logger)

	logging.L(ctx).Info("Configuration loaded")

	a := app.NewApp(ctx, config, &env)

	logging.L(ctx).Info("Starting application...")

	if err := a.MustRun(); err != nil {
		logging.L(ctx).Error("Failed to run app", logging.ErrAttr(err))
		fmt.Printf("\napplication run error : %v \n", err)
		os.Exit(1)
	}

	logging.L(ctx).Info("Finished dumper...")
	os.Exit(0)
}

func runLog(e *app.Env, isLogging bool) *logging.Logs {
	var opts []logging.LoggerOption

	opts = []logging.LoggerOption{
		logging.WithFile(e.FileLog),
		logging.WithEnabled(isLogging),
		logging.WithAddSource(false),
		logging.WithLevel("debug"),
		logging.WithIsJSON(true),
	}

	return logging.NewLogger(opts...)
}
