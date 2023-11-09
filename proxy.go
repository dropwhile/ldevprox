package main

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type verboseFlag bool

func (v verboseFlag) BeforeApply() error {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Debug().Msg("debug logging enabled")
	return nil
}

type CLI struct {
	Listen   string           `name:"listen" short:"l" default:"127.0.0.1:8080" help:"listen address:port"`
	TLSCert  string           `name:"tls-cert" default:"server.crt" help:"tls cert file path"`
	TLSKey   string           `name:"tls-key" default:"server.key" help:"tls key file path"`
	Upstream string           `name:"upstream" default:"http://127.0.0.1:8000" help:"upstream url"`
	Verbose  verboseFlag      `name:"verbose" short:"v" help:"enable verbose logging"`
	Version  kong.VersionFlag `name:"version" short:"V" help:"Print version information and quit"`
}

func (cmd *CLI) Run() error {
	localProxyUrl, _ := url.Parse(cmd.Upstream)
	localProxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(localProxyUrl)
			r.SetXForwarded()
		},
	}
	// localProxy := httputil.NewSingleHostReverseProxy(localProxyUrl)
	http.Handle("/", localProxy)
	log.Info().Msgf("Serving on %s", cmd.Listen)
	if err := http.ListenAndServeTLS(cmd.Listen, cmd.TLSCert, cmd.TLSKey, nil); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msg("Server error")
		return err
	}
	return nil
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:          os.Stderr,
		PartsExclude: []string{zerolog.TimestampFieldName},
	})

	cli := CLI{}
	ctx := kong.Parse(&cli,
		kong.Name("ldevprox"),
		kong.Description("A local ssl proxy for development"),
		kong.UsageOnError(),
		kong.Vars{
			"version": "0.0.1",
		},
	)
	err := ctx.Run(&cli)
	ctx.FatalIfErrorf(err)
}
