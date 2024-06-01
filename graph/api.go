package graph

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"io"
	"strconv"
	"time"
)

func MarshallID(id uint64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, err := io.WriteString(w, strconv.Quote(fmt.Sprintf("%d", id)))
		if err != nil {
			return
		}
	})
}

func UnmarshalID(v any) (uint64, error) {
	id, ok := v.(string)
	if !ok {
		return 0, fmt.Errorf("ids must be strings")
	}
	i, e := strconv.ParseUint(id, 10, 64)
	return i, e
}

func MarshalTimestamp(t time.Time) graphql.Marshaler {
	timestamp := t.Format(time.RFC3339)

	return graphql.WriterFunc(func(w io.Writer) {
		_, err := io.WriteString(w, timestamp)
		if err != nil {
			return
		}
	})
}

func UnmarshalTimestamp(v any) (time.Time, error) {
	if timeStr, ok := v.(string); ok {
		return time.Parse(time.RFC3339, timeStr)
	}
	return time.Time{}, fmt.Errorf("wrong timestamp: need RFC3339")
}
