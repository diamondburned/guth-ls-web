//go:build !nomysql

package guthls

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

// MySQLProvider provides a provider implementation to get the leaderboard from
// the MySQL server.
type MySQLProvider struct {
	*sql.DB
}

var _ Provider = (*MySQLProvider)(nil)

// NewMySQLProvider creates a new MySQLProvider.
func NewMySQLProvider(dsn string) (*MySQLProvider, error) {
	c, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &MySQLProvider{c}, nil
}

// User implements Provider.
func (p *MySQLProvider) User(steamID SteamID) (*User, error) {
	row := p.DB.QueryRow(userQuery+" WHERE steamid = ?", steamID)

	u, err := scanUser(row.Scan)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// Users implements Provider.
func (p *MySQLProvider) Users() ([]User, error) {
	rows, err := p.DB.Query(userQuery)
	if err != nil {
		return nil, errors.Wrap(err, "query error")
	}

	defer rows.Close()

	var users []User
	for rows.Next() {
		u, err := scanUser(rows.Scan)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, rows.Err()
}

// LeaderboardForUser implements Provider.
func (p *MySQLProvider) LeaderboardForUser(steamID SteamID, fetchUser bool) (*UserLeaderboard, error) {
	f, q := lbQueryChoose(fetchUser)

	row := p.DB.QueryRow(
		q+" WHERE guth_ls.SteamID = ? ORDER BY guth_ls.LVL DESC, guth_ls.XP DESC", steamID)

	l, err := f(row.Scan)
	if err != nil {
		return nil, err
	}

	return &l, nil
}

// Leaderboard implements Provider.
func (p *MySQLProvider) Leaderboard(fetchUser bool) (Leaderboard, error) {
	f, q := lbQueryChoose(fetchUser)

	rows, err := p.DB.Query(
		q + " ORDER BY guth_ls.LVL DESC, guth_ls.XP DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lb Leaderboard
	for rows.Next() {
		u, err := f(rows.Scan)
		if err != nil {
			return nil, err
		}
		lb = append(lb, u)
	}

	return lb, rows.Err()
}

func lbQueryChoose(fetchUser bool) (func(scanFunc) (UserLeaderboard, error), string) {
	if fetchUser {
		return scanUserLeaderboardFetch, lbFetchQuery
	}
	return scanUserLeaderboard, lbQuery
}

type scanFunc func(...interface{}) error

const (
	userQuery = "" +
		"SELECT utime.steamid, utime.playername, utime.team, utime.totaltime, utime.lastvisit " +
		"FROM utime"

	lbQuery = "" +
		"SELECT guth_ls.SteamID, guth_ls.XP, guth_ls.LVL " +
		"FROM guth_ls"

	lbFetchQuery = "" +
		"SELECT guth_ls.SteamID, guth_ls.XP, guth_ls.LVL, utime.playername, utime.team, utime.totaltime, utime.lastvisit " +
		"FROM guth_ls " +
		"INNER JOIN utime ON utime.steamid = guth_ls.SteamID"
)

func scanUser(scan scanFunc) (User, error) {
	var u User
	err := scan(&u.SteamID, &u.PlayerName, &u.Team, &u.TotalTime, &u.LastVisit)
	return u, newScanError(err, "user")
}

func scanUserLeaderboard(scan scanFunc) (UserLeaderboard, error) {
	var l UserLeaderboard
	err := scan(&l.SteamID, &l.XP, &l.Level)
	return l, newScanError(err, "leaderboard")
}

func scanUserLeaderboardFetch(scan scanFunc) (UserLeaderboard, error) {
	var l UserLeaderboard
	l.User = &User{}
	err := scan(
		&l.SteamID, &l.XP, &l.Level,
		&l.User.PlayerName, &l.User.Team, &l.User.TotalTime, &l.User.LastVisit)
	l.User.SteamID = l.SteamID
	return l, err
}

func newScanError(err error, thing string) error {
	if err == nil {
		return nil
	}
	return errors.Wrapf(err, "cannot scan %s", thing)
}
