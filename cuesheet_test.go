package cuesheetgo

import (
	"errors"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	name         string
	path         string
	cueSheetData string
	expected     CueSheet
	expectedErr  error
}

var minimalCueSheet = CueSheet{
	FileName: "sample.flac",
	Format:   "WAVE",
	Tracks: []Track{
		{
			Type: "AUDIO",
		},
	},
}

func TestParseCueSheets(t *testing.T) {
	tcs := []testCase{
		{
			name:     "MinimalCueSheet",
			path:     "minimal.cue",
			expected: minimalCueSheet,
		},
		{
			name:        "EmptyCueSheet",
			path:        "empty.cue",
			expectedErr: errors.New("missing file name"),
		},
		{
			name:        "UnexpectedCommand",
			path:        path.Join("command", "unexpected.cue"),
			expectedErr: errors.New("unexpected command: UNSUPPORTED"),
		},
		{
			name:        "InsufficientLineFields",
			path:        path.Join("command", "insufficient.cue"),
			expectedErr: errors.New("expected at least 2 fields, got 1"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, runTest(tc))
	}
}

func TestParseFileCommand(t *testing.T) {
	tcs := []testCase{
		{
			name:        "RepeatedFileCommand",
			path:        path.Join("file", "repeated.cue"),
			expectedErr: errors.New("field already set: WAVE"),
		},
		{
			name:        "InsufficientFileParams",
			path:        path.Join("file", "insufficient.cue"),
			expectedErr: errors.New("expected 2 parameters, got 1"),
		},
		{
			name:        "ExcessiveFileParams",
			path:        path.Join("file", "excessive.cue"),
			expectedErr: errors.New("expected 2 parameters, got 3"),
		},
		{
			name:        "EmptyFileName",
			path:        path.Join("file", "empty_name.cue"),
			expectedErr: errors.New("missing file name"),
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, runTest(tc))
	}
}

func TestParseTrackCommand(t *testing.T) {
	tcs := []testCase{
		{
			name:        "InsufficientTrackParams",
			path:        path.Join("track", "insufficient.cue"),
			expectedErr: errors.New("expected 2 parameters, got 1"),
		},
		{
			name:        "ExcessiveTrackParams",
			path:        path.Join("track", "excessive.cue"),
			expectedErr: errors.New("expected 2 parameters, got 3"),
		},
		{
			name:        "MissingTracks",
			path:        path.Join("track", "missing.cue"),
			expectedErr: errors.New("missing tracks"),
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, runTest(tc))
	}
}

func runTest(tc testCase) func(t *testing.T) {
	return func(t *testing.T) {
		cwd, err := os.Getwd()
		require.NoError(t, err)
		p := path.Join(cwd, "testdata", tc.path)
		file, err := os.Open(p)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file.Close())
		}()
		cueSheet, err := Parse(file)
		if tc.expectedErr != nil {
			require.Contains(t, err.Error(), tc.expectedErr.Error())
			fmt.Println(err)
			return
		}
		require.NoError(t, err)
		require.Equal(t, tc.expected, *cueSheet)
	}
}
