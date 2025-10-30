package main

import (
	"context"
	"dumper/internal/app"
	conf "dumper/internal/config/local"
	"dumper/internal/crypt"
	appDomain "dumper/internal/domain/app"
	"dumper/pkg/logging"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var (
	version   = "dev"
	date      = "unknown"
	appSecret = "wTke8p;yGRM#$9Fh1kkYf$o_S@qnEt0Y"
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
	input := flag.String("input", "", "Decrypt path backup file")
	cryptType := flag.String("crypt", "", "Crypt file: dump | config")
	pass := flag.String("password", "", "Password to decrypt backup file")
	mode := flag.String("mode", "", "Mode: encrypt | decrypt | recover")
	recovery := flag.String("recovery", "", "Recovery token for recovery")

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		_, _ = fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	flags := appDomain.Flags{
		ConfigFile: *configPath,
		DbNameList: *dbName,
		All:        *all,
		FileLog:    *fileLog,
		Input:      *input,
		Crypt:      *cryptType,
		Password:   *pass,
		Mode:       *mode,
		Recovery:   *recovery,
		AppSecret:  appSecret,
	}

	if flags.Crypt != "" {
		cryptApp := crypt.NewApp(ctx, &flags)
		err := cryptApp.Run()
		if err != nil {
			fmt.Printf("Crypt error: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if showVersion {
		fmt.Printf("Version: %s \nDate: %s\n", version, date)
		return
	}

	config, err := conf.Load(*configPath, appSecret)
	if err != nil {
		fmt.Printf("configuration loading error : %v \n", err)
		os.Exit(1)
	}

	logger := runLog(&flags, *config.Settings.Logging)

	defer func(logger *logging.Logs) {
		_ = logger.Close()
	}(logger)

	ctx = logging.ContextWithLogger(ctx, logger.Logger)

	logging.L(ctx).Info("Configuration loaded")

	a := app.NewApp(ctx, config, &flags)

	logging.L(ctx).Info("Starting application...")

	if err := a.MustRun(); err != nil {
		logging.L(ctx).Error("Failed to run app", logging.ErrAttr(err))
		fmt.Printf("\napplication run error : %v \n", err)
		os.Exit(1)
	}

	logging.L(ctx).Info("Finished dumper...")
	os.Exit(0)
}

func runLog(
	e *appDomain.Flags,
	isLogging bool,
) *logging.Logs {
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
