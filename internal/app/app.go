package app

//log.Info().Msgf(")
//log.Warn().Msg("")
//log.Error().Msgf("%s:  > %s", op, err)
//log.Fatal().Msgf("%s:  > %s", op, err)

//tests
//t.Errorf("\n%s: \n\twant:\n\t\"%s\" \n\tget: \n\t\"%s\"", tC.desc, tC.want, res)

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alekslesik/neuro-news/internal/app/handler"
	"github.com/alekslesik/neuro-news/internal/app/repository"
	"github.com/alekslesik/neuro-news/internal/app/service"
	"github.com/alekslesik/neuro-news/internal/pkg/grabber"
	"github.com/alekslesik/neuro-news/internal/pkg/router"
	"github.com/alekslesik/neuro-news/internal/pkg/server"
	"github.com/alekslesik/neuro-news/internal/pkg/template"

	"github.com/alekslesik/neuro-news/pkg/config"
	"github.com/alekslesik/neuro-news/pkg/logger"
)

type Application struct {
	ctx context.Context
	ccl context.CancelFunc
	cfg *config.Config
	log *logger.Logger
	db  *sql.DB
	tp  *template.Template
	svs *service.Services
	grb *grabber.Grabber
	rtr *router.Router
	srv *server.Server
	// middleware *middleware.Middleware
	// session    *session.Session
	// template *template.Template
	// mailer   *mailer.Mailer
}

func New(context context.Context, cancel context.CancelFunc) (*Application, error) {

	// config init
	config := configInit()

	// flag init
	flagInit(config)

	// logger init
	logger := loggerInit(config)

	// db init
	db := dbInit(config, logger)

	// template init
	templates := templateInit(logger)

	// repository init
	repositories := repository.New(db, logger)

	// grabber init
	grabber := grabberInit(logger, config)

	// services init
	services := service.New(repositories, logger, templates, grabber)

	// handlers init
	handler := handler.New(services, logger, templates)

	// TODO add todo init
	// appMiddleware := middleware.New()

	// router init
	router := routerInit(handler)

	server := serverInit(config, logger, router)

	// mail service init
	// appMailer := mailer.New(appConfig.SMTPConfig)

	// session init
	// appSession := session.New()

	// template init
	// appTemplate := template.New()

	return &Application{
		ctx: context,
		ccl: cancel,
		cfg: config,
		log: logger,
		rtr: router,
		db:  db,
		tp:  templates,
		svs: services,
		grb: grabber,
		srv: server,
		// middleware: appMiddleware,
		// session:    appSession,
		// model:      model,
		// template:   appTemplate,
		// mailer:     appMailer,
	}, nil
}

func (a *Application) Run() error {
	log := a.log.With().Str("app", "Run()").Logger()

	var err error
	errChan := make(chan error)
	done := make(chan bool)
	defer close(done)

	// delta := int64(1)

	go func() {
		// Kandinsky
		for {
			<-done
			time.Sleep(time.Minute * time.Duration(a.cfg.App.Delta))
			article, err := a.svs.GetArticleService().GrabNewArticle()
			if err != nil {
				log.Warn().Msgf("get new article error > %s", err)
				done <- true
				continue
			}

			image, err := a.svs.GetImageService().GenerateImageKand(article)
			if err != nil {
				log.Warn().Msgf("generate new image error > %s", err)
				done <- true
				continue
			}

			err = a.svs.GetImageService().InsertImage(image)
			if err != nil {
				log.Warn().Msgf("insert generated image to DB error > %s", err)
				done <- true
				continue
			}

			err = a.svs.GetArticleService().InsertArticleImage(image, article)
			if err != nil {
				log.Warn().Msgf("insert article to DB error > %s", err)
				done <- true
				continue
			}

			log.Info().Msgf("article insert to DB through kandinsky package")

			done <- true
		}
	}()

	go func() {
		// Fruity
		for {
			<-done
			time.Sleep(time.Minute * time.Duration(a.cfg.App.Delta))
			article, err := a.svs.GetArticleService().GrabNewArticle()
			if err != nil {
				log.Warn().Msgf("get new article error > %s", err)
				done <- true
				continue
			}

			image, err := a.svs.GetImageService().GenerateImageFruity(article)
			if err != nil {
				log.Warn().Msgf("generate new image error > %s",  err)
				done <- true
				continue
			}

			err = a.svs.GetImageService().InsertImage(image)
			if err != nil {
				log.Warn().Msgf("insert generated image to DB error > %s",  err)
				done <- true
				continue
			}

			err = a.svs.GetArticleService().InsertArticleImage(image, article)
			if err != nil {
				log.Warn().Msgf("insert article to DB error > %s",  err)
				done <- true
				continue
			}

			log.Info().Msgf("article insert to DB through Fruity API",)
			done <- true
		}
	}()

	done <- true

	// db close
	defer a.closeDB()
	// logfile close
	defer a.log.File.Close()

	// Set signals handler
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	switch a.srv.Addr {
	case "localhost:80", "localhost:8080", ":80", ":8080":
		go func() {
			err = a.srv.ListenAndServe()
		}()
		log.Info().Msgf("server started on %s/", a.srv.Addr)

	case "localhost:443", "localhost:8443", ":443", ":8443":
		go func() {
			err = a.srv.ListenAndServeTLS(a.cfg.TLS.CertPath, a.cfg.TLS.KeyPath)
		}()
		log.Info().Msgf("Server started on https://localhost:%s/", a.srv.Addr)

	default:
		log.Error().Msgf("address or port are not exists > %s", a.srv.Addr)
	}

	if err != nil && err != http.ErrServerClosed {
		errChan <- err
	}

	select {
	case <-a.ctx.Done():
		log.Warn().Msg("Context signal received, initiating shutdown")
		a.srv.Shutdown(a.ctx)
		// TODO temprorary
		// time.Sleep(2 * time.Second)

	case err := <-errChan:
		log.Error().Msgf("server failure > %s",  err)
		return err

	case <-signals:
		log.Warn().Msg("Signal received, initiating shutdown")
		a.srv.Shutdown(a.ctx)
		// time.Sleep(2 * time.Second)
	}

	return nil
}

func (a *Application) closeDB() {
	const op = "app.Close()"

	if err := a.db.Close(); err != nil {
		a.log.Error().Msgf("%s: failed to close data base > %s", op, err)
	}
}
