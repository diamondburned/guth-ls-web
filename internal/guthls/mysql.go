//go:build !nomysql

package guthls

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/keegancsmith/sqlf"
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
func (p *MySQLProvider) LeaderboardForUser(steamID SteamID, f ...LeaderboardQueryFlags) (*UserLeaderboard, error) {
	query, scan := leaderboardQuery(f, sqlf.Sprintf("guth_ls.SteamID = %s", steamID))

	row := p.DB.QueryRow(query.Query(sqlf.SimpleBindVar), query.Args()...)

	l, err := scan(row.Scan)
	if err != nil {
		return nil, err
	}

	return &l, nil
}

// Leaderboard implements Provider.
func (p *MySQLProvider) Leaderboard(f ...LeaderboardQueryFlags) (Leaderboard, error) {
	query, scan := leaderboardQuery(f)

	rows, err := p.DB.Query(query.Query(sqlf.SimpleBindVar), query.Args()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lb Leaderboard
	for rows.Next() {
		u, err := scan(rows.Scan)
		if err != nil {
			return nil, err
		}
		lb = append(lb, u)
	}

	return lb, rows.Err()
}

type scanFunc func(...interface{}) error

const userQuery = "" +
	"SELECT utime.steamid, utime.playername, utime.team, utime.totaltime, utime.lastvisit " +
	"FROM utime"

func scanUser(scan scanFunc) (User, error) {
	var u User
	err := scan(&u.SteamID, &u.PlayerName, &u.Team, &u.TotalTime, &u.LastVisit)
	return u, newScanError(err, "user")
}

type leaderboardScanner = func(scanFunc) (UserLeaderboard, error)

func leaderboardQuery(flags []LeaderboardQueryFlags, filters ...*sqlf.Query) (*sqlf.Query, leaderboardScanner) {
	head := sqlf.Sprintf("guth_ls.SteamID, guth_ls.XP, guth_ls.LVL")
	join := &sqlf.Query{}

	flag := NewLeaderboardQueryFlags(flags)
	if flag&LeaderboardQueryUser != 0 {
		head = sqlf.Sprintf("%s, utime.playername, utime.team, utime.totaltime, utime.lastvisit", head)
		join = sqlf.Sprintf("%s INNER JOIN utime ON utime.steamid = guth_ls.SteamID", join)
	}
	if flag&LeaderboardQueryRank != 0 {
		head = sqlf.Sprintf("%s, ranks.rank", head)
		join = sqlf.Sprintf("%s LEFT JOIN ranks ON ranks.steamid = guth_ls.SteamID", join)
	}

	head = sqlf.Sprintf("SELECT %s FROM guth_ls %s", head, join)
	if len(filters) > 0 {
		head = sqlf.Sprintf("%s WHERE %s", head, sqlf.Join(filters, "AND"))
	}
	head = sqlf.Sprintf("%s ORDER BY guth_ls.LVL DESC, guth_ls.XP DESC", head)

	return head, func(scan scanFunc) (l UserLeaderboard, err error) {
		v := []interface{}{&l.SteamID, &l.XP, &l.Level}

		if flag&LeaderboardQueryUser != 0 {
			l.User = &User{}
			v = append(v, &l.User.PlayerName, &l.User.Team, &l.User.TotalTime, &l.User.LastVisit)

			defer func() {
				l.User.SteamID = l.SteamID
			}()
		}

		if flag&LeaderboardQueryRank != 0 {
			var rank sql.NullString
			v = append(v, &rank)

			defer func() {
				l.Rank = rank.String
			}()
		}

		err = newScanError(scan(v...), "leaderboard")
		return
	}
}

func newScanError(err error, thing string) error {
	if err == nil {
		return nil
	}
	return errors.Wrapf(err, "cannot scan %s", thing)
}
