package types

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type contextKey string

const (
	UserIDKey     contextKey = "user_id"
	BusinessIDKey contextKey = "business_id"
	UserRoleKey   contextKey = "user_role"
	SubdomainKey  contextKey = "subdomain"
	RequestIDKey  contextKey = "request_id"
)

func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return uuid.Nil, errors.New("user_id not found in context")
	}
	return userID, nil
}

func WithBusinessID(ctx context.Context, businessID uuid.UUID) context.Context {
	return context.WithValue(ctx, BusinessIDKey, businessID)
}

func GetBusinessIDFromContext(ctx context.Context) (uuid.UUID, error) {
	businessID, ok := ctx.Value(BusinessIDKey).(uuid.UUID)
	if !ok || businessID == uuid.Nil {
		return uuid.Nil, errors.New("business_id not found in context")
	}
	return businessID, nil
}

func WithUserRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, UserRoleKey, role)
}

func GetUserRoleFromContext(ctx context.Context) (string, error) {
	role, ok := ctx.Value(UserRoleKey).(string)
	if !ok || role == "" {
		return "", errors.New("user_role not found in context")
	}
	return role, nil
}

func WithSubdomain(ctx context.Context, subdomain string) context.Context {
	return context.WithValue(ctx, SubdomainKey, subdomain)
}

func GetSubdomainFromContext(ctx context.Context) (string, error) {
	subdomain, ok := ctx.Value(SubdomainKey).(string)
	if !ok {
		return "", errors.New("subdomain not found in context")
	}
	return subdomain, nil
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

func GetRequestIDFromContext(ctx context.Context) string {
	requestID, _ := ctx.Value(RequestIDKey).(string)
	if requestID == "" {
		return "unknown"
	}
	return requestID
}
