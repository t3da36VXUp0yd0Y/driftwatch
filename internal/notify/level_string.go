package notify

import (
	"fmt"
	"strings"
)

var levelNames = map[Level]string{
	LevelAll:       "all",
	LevelDriftOnly: "drift-only",
}

// String returns the human-readable name of the Level.
func (l Level) String() string {
	if name, ok := levelNames[l]; ok {
		return name
	}
	return fmt.Sprintf("Level(%d)", int(l))
}

// ParseLevel converts a string to a Level.
// It returns an error if the string does not match a known level.
func ParseLevel(s string) (Level, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "all":
		return LevelAll, nil
	case "drift-only", "driftonly":
		return LevelDriftOnly, nil
	default:
		return LevelAll, fmt.Errorf("unknown notify level %q: must be one of [all, drift-only]", s)
	}
}
