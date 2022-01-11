package guthls

import (
	"fmt"
	"time"

	"github.com/leighmacdonald/steamid/steamid"
)

// Leaderboard describes a leaderboard.
type Leaderboard []UserLeaderboard

// LeaderboardQueryFlags contains flags for querying additional information from
// the database along with the leaderboard.
type LeaderboardQueryFlags uint8

const (
	LeaderboardQueryUser LeaderboardQueryFlags = 1 << iota
	LeaderboardQueryRank
)

// NewLeaderboardQueryFlags ORs the flags together.
func NewLeaderboardQueryFlags(flags []LeaderboardQueryFlags) LeaderboardQueryFlags {
	var flag LeaderboardQueryFlags
	for _, f := range flags {
		flag |= f
	}
	return flag
}

// UserLeaderboard describes a player in the leaderboard (or a row in the
// guth-ls table).
type UserLeaderboard struct {
	SteamID SteamID
	XP      int
	Level   int

	User *User  // only if LeaderboardQueryUser
	Rank string // only if LeaderboardQueryRank
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

// ProfileURL returns the profile URL for the Steam ID. If the SteamID is
// invalid, then an empty string is returned.
func (id SteamID) ProfileURL() string {
	s := steamid.SIDToSID64(steamid.SID(id))
	if s == 0 {
		return ""
	}
	return fmt.Sprintf("https://steamcommunity.com/profiles/%d", s)
}

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

// GlobalStats is the global statistics of the server.
type GlobalStats struct {
	Players   int
	TotalTime Seconds
}

// Provider describes the database getters.
type Provider interface {
	// User gets a user by Steam ID.
	User(SteamID) (*User, error)
	// Users gets all users.
	Users() ([]User, error)
	// LeaderboardForUser gets the user's leaderboard entry by Steam ID.
	LeaderboardForUser(SteamID, ...LeaderboardQueryFlags) (*UserLeaderboard, error)
	// Leaderboard gets the entire leaderboard. The returned leaderboard is
	// guaranteed to be sorted by Level then XP. If fetchUser is true, then the
	// returned UserLeaderboard instances should have the User field filled.
	Leaderboard(...LeaderboardQueryFlags) (Leaderboard, error)
	// GlobalStats gets the global server statistics.
	GlobalStats() (*GlobalStats, error)
}
