package alert

import (
	"fmt"
	"strings"
)

var levelNames = map[Level]string{
	LevelNone:     "NONE",
	LevelWarning:  "WARNING",
	LevelCritical: "CRITICAL",
}

// String returns the string representation of a Level.
func (l Level) String() string {
	if name, ok := levelNames[l]; ok {
		return name
	}
	return fmt.Sprintf("Level(%d)", int(l))
}

// ParseLevel converts a string to a Level value.
func ParseLevel(s string) (Level, error) {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "NONE":
		return LevelNone, nil
	case "WARNING", "WARN":
		return LevelWarning, nil
	case "CRITICAL", "CRIT":
		return LevelCritical, nil
	default:
		return LevelNone, fmt.Errorf("unknown alert level: %q", s)
	}
}
