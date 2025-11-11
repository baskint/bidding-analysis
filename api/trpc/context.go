package trpc

import (
	"context"
)

// ContextKey is a type used for context keys to prevent collisions with other packages.
type ContextKey string

const (
	// Exported context keys (start with a capital letter) that are used by middleware and handlers.
	ContextKeyUserID   ContextKey = "user_id"
	ContextKeyUsername ContextKey = "username"
)

// ContextWithUser adds the UserID and Username to the provided context.
// Exported to be used by the authMiddleware.
func ContextWithUser(ctx context.Context, userID, username string) context.Context {
	ctx = context.WithValue(ctx, ContextKeyUserID, userID)
	ctx = context.WithValue(ctx, ContextKeyUsername, username)
	return ctx
}

// GetUserIDFromContext safely retrieves the UserID from the context.
// Exported to be used by analytics handlers.
func GetUserIDFromContext(ctx context.Context) string {
	// CRITICAL FIX: Use the exported constant ContextKeyUserID
	if userID, ok := ctx.Value(ContextKeyUserID).(string); ok {
		return userID
	}
	return ""
}

// GetUsernameFromContext safely retrieves the Username from the context.
// Exported to be used by handlers like getMe.
func GetUsernameFromContext(ctx context.Context) string {
	// CRITICAL FIX: Use the exported constant ContextKeyUsername
	if username, ok := ctx.Value(ContextKeyUsername).(string); ok {
		return username
	}
	return ""
}
