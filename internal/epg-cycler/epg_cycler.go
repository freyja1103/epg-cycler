package epgcycler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	edcbapiclient "github.com/freyja1103/epg-cycler/edcb-api-client"
	syobocalapiclient "github.com/freyja1103/epg-cycler/syobocal-api-client"
)

type EPGCycler interface {
	SimpleTidy(ctx context.Context) error
}
type epgCycler struct {
	Config      *Config
	EDCBAPI     edcbapiclient.API
	SyobocalAPI syobocalapiclient.API
}

type Config struct {
	ReserveCutoffHour int
	File              *ProgramFile
}

type ProgramFile struct {
	BaseName        string
	OriginPath      string
	DestinationPath string
}

func NewEPGCycler(edcbapi edcbapiclient.API, syobocalapi syobocalapiclient.API, config *Config) EPGCycler {
	return &epgCycler{
		Config:      config,
		EDCBAPI:     edcbapi,
		SyobocalAPI: syobocalapi,
	}
}

func (e *epgCycler) SimpleTidy(ctx context.Context) error {
	rec, err := e.EDCBAPI.GetEnumRecInfo(ctx, nil)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get enum rec info", slog.Any("error", err))
		return err
	}
	if len(rec) == 0 {
		return errors.New("no recording info found")
	}

	endTime := rec[0].StartTime.Add(time.Duration(rec[0].Duration))
	prog, err := e.SyobocalAPI.ProgLookup(ctx, &syobocalapiclient.ProgLookupParams{
		ChIDs: []string{SIDToChID[fmt.Sprintf("%d", rec[0].ServiceID)]},
		Range: &syobocalapiclient.Range{
			From: rec[0].StartTime,
			To:   &endTime,
		}})

	title, err := e.SyobocalAPI.TitleLookup(ctx, &syobocalapiclient.TitleLookupParams{
		TIDs: []string{fmt.Sprintf("%d", prog[0].TID)},
	})

	if err = e.tidyDirectory(ctx, title[0].Title); err != nil {
		return err
	}

	entry, err := e.EDCBAPI.GetEnumReserveInfo(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get enum reserve info", slog.Any("error", err))
		return err
	}
	hasReserve, _, err := e.HasReserve(entry)
	if err != nil {
		slog.ErrorContext(ctx, "failed to check remaining reserve", slog.Any("error", err))
		return err
	}
	if hasReserve {
		return nil
	}

	if err = execShutdown(); err != nil {
		slog.ErrorContext(ctx, "failed to execute shutdown", slog.Any("error", err))
		return err
	}
	return nil
}

func (e *epgCycler) tidyDirectory(ctx context.Context, title string) error {
	programFolderPath := filepath.Join(e.Config.File.DestinationPath, title)
	err := os.MkdirAll(programFolderPath, 0755)
	if err != nil && os.IsNotExist(err) {
		slog.ErrorContext(ctx, "failed to mkdir: ", slog.Any("error", err))
		return err
	}
	if errors.Is(err, os.ErrNotExist) {
		slog.InfoContext(ctx, "created directory", slog.String("path", programFolderPath))
	}

	for ext := range extensions {
		if err = os.Rename(e.Config.File.OriginPath, filepath.Join(programFolderPath, e.Config.File.BaseName+extensions[ext])); err != nil {
			slog.ErrorContext(ctx, "failed to move file: ", slog.Any("error", err))
			return err
		}
	}

	return nil
}

func execShutdown() error {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("C:\\Windows\\System32\\shutdown.exe", "/s", "/t", "60", "/f", "/c", "shutdown by epg-cycler after 60s")
		if cmd.Err != nil {
			return cmd.Err
		}
		slog.Info("execute shutdown")
		// Runだとシャットダウンし終えるまで処理が進まなくなるのでStartを使う
		if err := cmd.Start(); err != nil {
			return err
		}
		return nil
	}
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("shutdown", "-h now")
		if cmd.Err != nil {
			return cmd.Err
		}
		slog.Info("execute shutdown")
		if err := cmd.Start(); err != nil {
			return err
		}
		return nil
	}

	return errors.New("unsupported OS for shutdown command")
}

var extensions = map[string]string{
	"ts":      ".ts",
	"error":   ".err",
	"program": ".ts.program.txt",
}

var SIDToChID = map[string]string{
	"1024":  "1",  // NHK総合
	"1032":  "2",  // NHK Eテレ
	"1056":  "3",  // フジテレビ
	"1040":  "4",  // 日本テレビ
	"1048":  "5",  // TBS
	"1064":  "6",  // テレビ朝日
	"1072":  "7",  // テレビ東京
	"24632": "8",  // tvk1
	"23608": "19", // TOKYO MX1
	"29752": "14", // テレ玉
}
