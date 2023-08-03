package handlers

import (
	"context"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"

	"github.com/rarimo/issuer/internal/service/core/issuer"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	issuerCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxIssuer(issuer issuer.Issuer) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, issuerCtxKey, issuer)
	}
}

func Issuer(r *http.Request) issuer.Issuer {
	return r.Context().Value(issuerCtxKey).(issuer.Issuer)
}
