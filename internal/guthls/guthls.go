//go:build !nomysql

package guthls

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLProvider provides a provider implementation to get the leaderboard from
// the MySQL server.
type MySQLProvider struct {
	*sql.Conn
}
