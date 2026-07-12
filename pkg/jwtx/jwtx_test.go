package jwtx

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBuildToken(t *testing.T) {
	secret := "test-secret"
	now := time.Now().Unix()
	token, err := BuildToken(secret, now, 86400, 1001)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestUserIDFromCtx_Int64(t *testing.T) {
	ctx := context.WithValue(context.Background(), UserIDKey, int64(1001))
	id := UserIDFromCtx(ctx)
	assert.Equal(t, int64(1001), id)
}

func TestUserIDFromCtx_Float64(t *testing.T) {
	ctx := context.WithValue(context.Background(), UserIDKey, float64(1001))
	id := UserIDFromCtx(ctx)
	assert.Equal(t, int64(1001), id)
}

func TestUserIDFromCtx_JSONNumber(t *testing.T) {
	ctx := context.WithValue(context.Background(), UserIDKey, json.Number("1001"))
	id := UserIDFromCtx(ctx)
	assert.Equal(t, int64(1001), id)
}

func TestUserIDFromCtx_String(t *testing.T) {
	ctx := context.WithValue(context.Background(), UserIDKey, "1001")
	id := UserIDFromCtx(ctx)
	assert.Equal(t, int64(1001), id)
}

func TestUserIDFromCtx_Empty(t *testing.T) {
	ctx := context.Background()
	id := UserIDFromCtx(ctx)
	assert.Equal(t, int64(0), id)
}
