package main

import (
	"flag"
	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/options"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/labi-le/server/internal/server/basic"
	"github.com/labi-le/server/internal/server/storage"
	"github.com/labi-le/server/pkg/badgerdb"
	"github.com/labi-le/server/pkg/config"
	"github.com/labi-le/server/pkg/filesystem"
	"github.com/labi-le/server/pkg/log"
	"github.com/labi-le/server/pkg/response"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "debug mode")
	flag.Parse()

	cfg := MustConfig()

	logger := MustLogger(debugMode, cfg.GetLogLevel())

	server := MustServer(cfg, logger)

	reply := response.New(logger)
	MustBasic(server, reply)

	store := MustStorage(server, logger, cfg, reply)
	defer store.Close()

	if cfg.GetEnableHTTPS() {
		go UpTLSServer(logger, server, cfg)
	}

	if httpServerErr := server.Listen(cfg.GetServerConn()); httpServerErr != nil {
		logger.Warn(httpServerErr)
	}

}

func MustConfig() config.Config {
	cfg, err := config.NewFromENV()
	if err != nil {
		panic(err)
	}
	return cfg
}

func MustServer(cfg config.Config, logger log.Logger) *fiber.App {
	r := fiber.New(fiber.Config{
		DisableStartupMessage: false,
		BodyLimit:             cfg.GetMaxUploadSize(),
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

	return r
}

func MustBasic(r *fiber.App, reply *response.Reply) {
	basic.RegisterHandlers(r, reply)
}

func MustStorage(r *fiber.App, log log.Logger, cfg config.Config, reply *response.Reply) badgerdb.Store {
	defaultOpt := badger.DefaultOptions("db")
	defaultOpt.SyncWrites = false
	defaultOpt.TableLoadingMode = options.LoadToRAM
	// defaultOpt.ValueLogLoadingMode = options.MemoryMap

	db, dbErr := badger.Open(defaultOpt)
	client := badgerdb.NewWithBadger(db)
	if dbErr != nil {
		log.Error(dbErr)
	}

	storage.RegisterHandlers(
		r,
		storage.NewService(
			storage.NewStore(client,
				filesystem.New(cfg.GetVirtualFSPath()),
			),
		),
		cfg.GetOwnerKey(),
		reply,
	)

	return client
}

func UpTLSServer(logger log.Logger, r *fiber.App, cfg config.Config) {
	logger.Info("Starting server in production mode")
	go func() {
		if httpsServerErr := r.Listener(autocert.NewListener(cfg.GetWhiteListDomains()...)); httpsServerErr != nil {
			logger.Error(httpsServerErr)
		}

	}()
}

// MustLogger returns a logger based on the given parameters.
// see zapcore.level
func MustLogger(debug bool, level string) log.Logger {
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
	l, prodErr := zap.NewProduction()
	if prodErr != nil {
		panic(prodErr)
	}
	l.Core().Enabled(atomicLevel.Level())

	return log.NewWithZap(l)
}
