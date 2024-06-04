package cuesheetgo

import (
	"errors"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

type parseTestCase struct {
	name         string
	path         string
	cueSheetData string
	expected     CueSheet
	expectedErr  error
}

var minimalCueSheet = CueSheet{
	File:   "sample.flac",
	Format: "WAVE",
}

func TestParse(t *testing.T) {
	testCases := []parseTestCase{
		{
			name:     "MinimalCueSheet",
			path:     "minimal.cue",
			expected: minimalCueSheet,
		},
		{
			name:         "EmptyCueSheet",
			cueSheetData: "empty.cue",
			expectedErr:  errors.New("missing file name"),
		},
		{
			name:        "RepeatedFileCommand",
			path:        path.Join("repeated", "commands", "file.cue"),
			expectedErr: errors.New("field already set: WAVE"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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
		})
	}
}
