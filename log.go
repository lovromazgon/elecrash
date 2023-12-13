package elecrash

import (
	"fmt"
	"log/slog"
	"os"
)

var logger *slog.Logger

func InitLogger() func() {
	f, err := os.Create("elecrash.log")
	if err != nil {
		// file exists
		f, err = os.Open("elecrash.log")
	}
	if err != nil {
		panic(fmt.Errorf("could not open log file: %w", err))
	}
	logger = slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{}))
	return func() { _ = f.Close() }
}
