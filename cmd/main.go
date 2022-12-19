package main

import (
	"flag"
	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/options"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
	"labile-me-serv/internal/file"
	"labile-me-serv/pkg/badgerdb"
	"labile-me-serv/pkg/config"
	"labile-me-serv/pkg/filesystem"
	"labile-me-serv/pkg/log"
	"time"
)

func main() {
	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "debug mode")
	flag.Parse()

	cfg, err := config.NewFromENV()
	if err != nil {
		panic(err)
	}

	logger := GetLogger(debugMode, cfg.GetLogLevel())

	r := fiber.New(fiber.Config{
		DisableStartupMessage: false,
		BodyLimit:             cfg.GetMaxUploadSize(),
		ReadTimeout:           time.Minute * 5,
		WriteTimeout:          time.Minute * 30,
	})

	r.Use(log.LoggerMiddleware(logger))
	//r.Use(cache.New(cache.Config{
	//	Next: func(c *fiber.Ctx) bool {
	//		return c.Query("refresh") == "true"
	//	},
	//	Expiration:   time.Hour,
	//	CacheControl: true,
	//	CacheHeader:  "X-Cache",
	//}))

	r.Use(favicon.New(favicon.Config{
		File: "favicon.ico",
	}))

	defaultOpt := badger.DefaultOptions("db")
	defaultOpt.SyncWrites = false
	defaultOpt.TableLoadingMode = options.LoadToRAM
	// defaultOpt.ValueLogLoadingMode = options.MemoryMap

	db, dbErr := badger.Open(defaultOpt)
	client := badgerdb.NewWithBadger(db)
	if dbErr != nil {
		logger.Error(dbErr)
	}

	defer db.Close()

	fs := filesystem.New(cfg.GetVirtualFSPath())

	store := file.NewStore(client, fs)

	file.RegisterHandlers(r, logger, file.NewService(store), cfg)

	if cfg.GetEnableHTTPS() {
		go UpTLSServer(logger, r, cfg)
	}

	if httpServerErr := r.Listen(cfg.GetServerConn()); httpServerErr != nil {
		logger.Warn(httpServerErr)
	}

}

func UpTLSServer(logger log.Logger, r *fiber.App, cfg config.Config) {
	logger.Info("Starting server in production mode")
	go func() {
		if httpsServerErr := r.Listener(autocert.NewListener(cfg.GetWhiteListDomains()...)); httpsServerErr != nil {
			logger.Error(httpsServerErr)
		}

	}()
}

// GetLogger returns a logger based on the given parameters.
// see zapcore.level
func GetLogger(debug bool, level string) log.Logger {
	if level == "disable" {
		return log.NilLogger{}
	}

	if debug {
		l := zap.Must(zap.NewDevelopment())
		// set level to debug
		l.Core().Enabled(zap.DebugLevel)
		return log.NewWithZap(l)
	}

	// parse level
	atomicLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		panic(err)
	}

	// create logger
	l, _ := zap.NewProduction()
	l.Core().Enabled(atomicLevel.Level())

	return log.NewWithZap(l)
}
