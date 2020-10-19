package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/uber/jaeger-client-go"
	"regexp"
)

const (
	forbiddenMessageFormat = "forbidden (reason: %s)"
	notFoundMessageFormat = "not found (reason: %s)"
	proxyAuthRequiredMessageFormat = "proxy auth required (reason: %s)"
	requestTimeoutMessageFormat = "request time out (reason: %s)"
	conflictMessageFormat = "conflict (reason: %s)"
	internalServerMessageFormat = "internal server error (reason: %s)"
	serviceUnavailableMessageFormat = "service unavailable (reason: %s)"
)

var (
	adminUUIDRegex = regexp.MustCompile("^admin-\\d{12}")
	studentUUIDRegex = regexp.MustCompile("^student-\\d{12}")
	clubUUIDRegex = regexp.MustCompile("^club-\\d{12}")
)

func (_ _default) getContextFromMetadata(ctx context.Context) (parsedCtx context.Context, proxyAuthenticated bool, reason string) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		proxyAuthenticated = false
		reason = "metadata not exists"
		return
	}

	reqID, ok := md.Get("X-Request-Id")
	if !ok {
		proxyAuthenticated = false
		reason = "X-Request-Id not exists"
		return
	}

	_, err := uuid.Parse(reqID)
	if err != nil {
		proxyAuthenticated = false
		reason = "X-Request-ID invalid, err: " + err.Error()
		return
	}

	spanCtx, ok := md.Get("Span-Context")
	if !ok {
		proxyAuthenticated = false
		reason = "Span-Context not exists"
		return
	}

	parentSpan, err := jaeger.ContextFromString(spanCtx)
	if err != nil {
		proxyAuthenticated = false
		reason = "Span-Context invalid, err: " + err.Error()
		return
	}

	proxyAuthenticated = true
	reason = ""

	parsedCtx = context.Background()
	parsedCtx = context.WithValue(parsedCtx, "X-Request-Id", reqID)
	parsedCtx = context.WithValue(parsedCtx, "Span-Context", parentSpan)

	if cUUID, ok := md.Get("ClubUUID"); ok { parsedCtx = context.WithValue(parsedCtx, "ClubUUID", cUUID) }
	if cUUID, ok := md.Get("RecruitmentUUID"); ok { parsedCtx = context.WithValue(parsedCtx, "RecruitmentUUID", cUUID) }

	return
}

func contains(slice []string, item string) bool {
	for _, element := range slice {
		if item == element {
			return true
		}
	}
	return false
}
