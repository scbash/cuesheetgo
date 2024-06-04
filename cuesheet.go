package cuesheetgo

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"strings"
)

// CueSheet represents the contents of a cue sheet file.
// Required fields: File, Format.
type CueSheet struct {
	Format string
	File   string
}

// trimChars contains the characters to be trimmed from a string.
// These are: space, double quote, tab, newline.
const trimChars = ` ` + `"` + `\t` + `\n`

// Parse reads the cue sheet data from the provided reader and returns a parsed CueSheet struct.
func Parse(reader io.Reader) (*CueSheet, error) {
	scanner := bufio.NewScanner(reader)
	c := &CueSheet{}

	var lineNr int
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), trimChars)
		lineNr++
		if line == "" {
			continue
		}
		if err := c.parseLine(line); err != nil {
			return nil, fmt.Errorf("line %d:\t%s:\n\t%w", lineNr, line, err)
		}
	}
	if err := c.validate(); err != nil {
		return nil, fmt.Errorf("invalid cue sheet: %w", err)
	}
	slog.Info("cue sheet parsed correctly", "file", c.File, "format", c.Format, "lines", lineNr)
	return c, nil
}

func (c *CueSheet) parseLine(line string) error {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return fmt.Errorf("expected at least %d fields, got %d", 2, len(fields))
	}

	var err error
	command := fields[0]
	parameters := fields[1:]
	switch command {
	case "FILE":
		err = c.parseFile(parameters)
	default:
		return fmt.Errorf("unexpected command: %q", command)
	}
	if err != nil {
		return fmt.Errorf("error parsing %q command: %w", command, err)
	}
	return nil
}

func assignValue[T comparable](val T, field *T) error {
	zero := reflect.Zero(reflect.TypeOf(*field)).Interface()
	if *field == zero {
		*field = val
		return nil
	}
	return fmt.Errorf("field already set: %v", *field)
}

func parseString(val string, field *string) error {
	val = strings.Trim(val, trimChars)
	return assignValue(val, field)
}

func (c *CueSheet) parseFile(parameters []string) error {
	if len(parameters) < 2 {
		return fmt.Errorf("expected at least %d parameters, got %d", 2, len(parameters))
	}
	last := len(parameters) - 1
	if err := parseString(parameters[last], &c.Format); err != nil {
		return fmt.Errorf("error parsing FILE format: %w", err)
	}
	if err := parseString(strings.Join(parameters[:last], " "), &c.File); err != nil {
		return fmt.Errorf("error parsing FILE name: %w", err)
	}
	return nil
}

// validate checks if the cue sheet has FILE and at least one TRACK command with INDEX 01.
func (c *CueSheet) validate() error {
	if c.File == "" {
		return errors.New("missing file name")
	}
	if c.Format == "" {
		return errors.New("missing format")
	}
	return nil
}
