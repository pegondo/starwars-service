package errors_test

import (
	"testing"

	"github.com/pegondo/starwars/service/internal/errors"
	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	testCases := []struct {
		name           string
		errCode        string
		errMsg         string
		expectedErrStr string
	}{
		{
			name:           "empty_code_and_empty_msg",
			errCode:        "",
			errMsg:         "",
			expectedErrStr: " - ",
		},
		{
			name:           "code_and_empty_msg",
			errCode:        "<err-code>",
			errMsg:         "",
			expectedErrStr: "<err-code> - ",
		},
		{
			name:           "empty_code_and_msg",
			errCode:        "",
			errMsg:         "<err-msg>",
			expectedErrStr: " - <err-msg>",
		},
		{
			name:           "code_and_msg",
			errCode:        "<err-code>",
			errMsg:         "<err-msg>",
			expectedErrStr: "<err-code> - <err-msg>",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := errors.New(tc.errCode, tc.errMsg)
			require.NotNil(t, err)

			errStr := err.Error()
			require.Equal(t, tc.expectedErrStr, errStr)
		})
	}
}
