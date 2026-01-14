package cli

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	edcbapiclient "github.com/freyja1103/epg-cycler/edcb-api-client"
	epgcycler "github.com/freyja1103/epg-cycler/internal/epg-cycler"
	syobocalapiclient "github.com/freyja1103/epg-cycler/syobocal-api-client"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Options struct {
	DestinationPath string
	OriginPath      string
	Title           string
	BaseName        string
	Episode         string //deprecated
	Process         targetProcesses
	Address         string
	IsTidyMode      bool
}

func Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	rotatingWriter := &lumberjack.Logger{
		Filename:   "epg-cycler-log/epg-cycler.log",
		MaxSize:    50,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   true,
	}
	logger := slog.New(slog.NewJSONHandler(rotatingWriter, nil))
	slog.SetDefault(logger)

	o := new(Options)

	flag.StringVar(&o.DestinationPath, "srcpath", "", "save video path")
	flag.StringVar(&o.OriginPath, "originpath", "", "origin video path")
	flag.StringVar(&o.Title, "title", "", "a program's name")
	flag.StringVar(&o.BaseName, "basename", "", "filename without ext")
	flag.StringVar(&o.Episode, "episode", "", "episode number (deprecated)")
	flag.Var(&o.Process, "process", "process that prevent shutdown")
	flag.StringVar(&o.Address, "address", "localhost:5510", "the server address in host:port format (e.g., localhost:5510)")
	flag.BoolVar(&o.IsTidyMode, "all", false, "The mode for sorting all the recording files in the directory at once")
	flag.Parse()

	err := o.run(ctx)
	if err != nil {
		return err
	}
	return nil
}

type targetProcesses []string

func (tp *targetProcesses) String() string {
	return fmt.Sprintf("%v", *tp)
}

func (tp *targetProcesses) Set(value string) error {
	*tp = append(*tp, value)
	return nil
}

func (o *Options) run(ctx context.Context) error {
	edcbapi := edcbapiclient.NewEDCBAPIClient(o.Address, nil)
	syobocalapi := syobocalapiclient.NewSyobocalAPIClient(nil)

	ec := epgcycler.NewEPGCycler(edcbapi, syobocalapi, &epgcycler.Config{
		ReserveCutoffHour: 4,
		File: &epgcycler.ProgramFile{
			BaseName:        o.BaseName,
			OriginPath:      o.OriginPath,
			DestinationPath: o.DestinationPath,
		},
	})

	if o.IsTidyMode {
		return nil
	}

	if !o.IsTidyMode {
		if err := ec.SimpleTidy(ctx); err != nil {
			slog.ErrorContext(ctx, "failed to tidy", slog.Any("error", err))
			return err
		}

	}
	return nil
}
