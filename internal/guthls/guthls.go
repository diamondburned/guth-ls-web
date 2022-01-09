package guthls

import "time"

// Leaderboard describes a leaderboard.
type Leaderboard []UserLeaderboard

// UserLeaderboard describes a player in the leaderboard (or a row in the
// guth-ls table).
type UserLeaderboard struct {
	SteamID SteamID
	XP      int
	Level   int
	// User is the optional field to be fetched if Provider.Leaderboard is given
	// fetchUser = true.
	User *User
}

// UserTimes contains multiple users.
type UserTimes []User

// User is a player, or a row in the utime (or utime_server) table.
type User struct {
	SteamID    SteamID
	PlayerName string
	Team       string
	TotalTime  Seconds
	LastVisit  UnixTime
}

// SteamID is a player's Steam ID.
type SteamID string

// Seconds is a number type that represents a second duration.
type Seconds int

// Duration converts seconds to time.Duration.
func (s Seconds) Duration() time.Duration {
	return time.Duration(s) * time.Second
}

// UnixTime is a number type that represents the Unix timestamp.
type UnixTime int64

// Time converts the Unix timestamp to time.Time.
func (t UnixTime) Time() time.Time {
	return time.Unix(int64(t), 0)
}

// Provider describes the database getters.
type Provider interface {
	// User gets a user by Steam ID.
	User(SteamID) (*User, error)
	// Users gets all users.
	Users() ([]User, error)
	// LeaderboardForUser gets the user's leaderboard entry by Steam ID.
	LeaderboardForUser(id SteamID, fetchUser bool) (*UserLeaderboard, error)
	// Leaderboard gets the entire leaderboard. The returned leaderboard is
	// guaranteed to be sorted by Level then XP. If fetchUser is true, then the
	// returned UserLeaderboard instances should have the User field filled.
	Leaderboard(fetchUser bool) (Leaderboard, error)
}
