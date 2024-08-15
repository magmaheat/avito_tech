package slg

import "log/slog"

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func SetupLogger(fn, reqID string) *slog.Logger {
	return slog.With(
		slog.String("fn", fn),
		slog.String("id_request", reqID),
	)
}
