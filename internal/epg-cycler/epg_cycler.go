package epgcycler

import (
	"bufio"
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	edcbapiclient "github.com/freyja1103/epg-cycler/edcb-api-client"
)

type EPGCycler interface {
	SimpleTidy(ctx context.Context, originPath string) error
}
type epgCycler struct {
	Config  *Config
	EDCBAPI edcbapiclient.API
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

func NewEPGCycler(edcbapi edcbapiclient.API, config *Config) EPGCycler {
	return &epgCycler{
		Config:  config,
		EDCBAPI: edcbapi,
	}
}

func (e *epgCycler) SimpleTidy(ctx context.Context, originPath string) error {
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
	// めんどくさいから全角のままでいいかも・・・
	return nil
}

func (e *epgCycler) GetProgramName() (fullname string, err error) {
	file, err := os.Open(e.Config.File.OriginPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	var thirdLine string
	for scanner.Scan() {
		lineNum++
		if lineNum == 3 {
			thirdLine = scanner.Text()
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	if thirdLine == "" {
		return "", ErrInvalidProgramFile
	}
	// parts := strings.Split(thirdLine, "　") // 全角スペース
	// if len(parts) > 1 {
	// 	return parts[1], nil
	// }
	return thirdLine, nil
}

func (e *epgCycler) tidyDirectory(ctx context.Context, p *Program) error {
	programFolderPath := filepath.Join(e.Config.File.DestinationPath, p.Title)
	err := os.MkdirAll(programFolderPath, 0755)
	if err != nil && os.IsNotExist(err) {
		slog.ErrorContext(ctx, "failed to mkdir: ", slog.Any("error", err))
		return err
	}
	if errors.Is(err, os.ErrNotExist) {
		slog.InfoContext(ctx, "created directory", slog.String("path", programFolderPath))
	}
	return nil
}

var extensions = map[string]string{
	"ts":      ".ts",
	"error":   ".err",
	"program": ".ts.program.txt",
}
