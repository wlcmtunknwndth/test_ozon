package slogAttr

import "log/slog"

func SlogErr(op string, err error) slog.Attr {
	return slog.Any(op, err)
}

func SlogInfo(msg string, data any) slog.Attr {
	return slog.Any(msg, data)
}
