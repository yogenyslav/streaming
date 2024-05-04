package database

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func TruncateTable(t *testing.T, ctx context.Context, pg *pgxpool.Pool, tables ...string) {
	t.Helper()
	_, err := pg.Exec(ctx, fmt.Sprintf(`
		truncate %s restart identity
	`, strings.Join(tables, ",")))
	require.NoError(t, err)
}
