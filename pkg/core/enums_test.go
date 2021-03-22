package core_test

import (
	"testing"

	"github.com/gustavooferreira/pgw-payment-processor-service/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCCFailReasonLoad(t *testing.T) {
	tests := map[string]struct {
		input          string
		expectedErr    bool
		expectedOutput core.CCFailReason
	}{
		"authorise fail": {
			input:          "authorise fail",
			expectedErr:    false,
			expectedOutput: core.CCFailReason_Authorise,
		},
		"capture fail": {
			input:          "capture fail",
			expectedErr:    false,
			expectedOutput: core.CCFailReason_Capture,
		},
		"refund fail": {
			input:          "refund fail",
			expectedErr:    false,
			expectedOutput: core.CCFailReason_Refund,
		},
		"void fail": {
			input:          "void fail",
			expectedErr:    false,
			expectedOutput: core.CCFailReason_Void,
		},
		"error returned for unknown reason to fail": {
			input:       "unknown fail",
			expectedErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var reason core.CCFailReason
			err := reason.Load(test.input)
			if test.expectedErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.expectedOutput, reason)
		})
	}
}
