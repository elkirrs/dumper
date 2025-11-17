package main

import (
	"context"
	"dumper/internal/app"
	conf "dumper/internal/config/local"
	"dumper/internal/crypt"
	appDomain "dumper/internal/domain/app"
	"dumper/pkg/logging"
	"dumper/pkg/utils"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var (
	version = "dev"
	date    = "unknown"
	appKey  = "app_key"
	appName = "Dumper"
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
	input := flag.String("input", "", "Decrypt path file")
	cryptType := flag.String("crypt", "", "Crypt file: backup | config")
	pass := flag.String("password", "", "Password to crypt file (optional)")
	mode := flag.String("mode", "", "Mode: encrypt | decrypt | recovery")
	recoveryKey := flag.String("token", "", "Recovery token for recovery")

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
		Recovery:   *recoveryKey,
		AppSecret:  appKey,
	}

	if flags.Crypt != "" {
		cryptApp := crypt.NewApp(ctx, &flags)
		if err := utils.RunWithCtx(ctx, func() error { cryptApp.Run(); return nil }); err != nil {
			fmt.Printf("Crypt error: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if showVersion {
		fmt.Println(appName)
		fmt.Printf("Version: %s \nDate: %s\n", version, date)
		return
	}

	config, err := conf.Load(*configPath, appKey)
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
		switch {
		case errors.Is(err, context.Canceled):
			fmt.Println("Closed dumper...")
			logging.L(ctx).Error("Closed dumper", logging.ErrAttr(err))
			os.Exit(0)
		default:
			fmt.Printf("\napplication run error : %v \n", err)
			logging.L(ctx).Error("Failed to run app", logging.ErrAttr(err))
			os.Exit(1)
		}
	}

	fmt.Println("Finished dumper...")
	logging.L(ctx).Info("Finished dumper")
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
