package pkg_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"krashanoff.com/parallel/pkg"
)

func TestValidPattern(t *testing.T) {
	t.Run("entireString", func(t *testing.T) {
		assert.True(t, pkg.ValidPattern("{}"))
	})
}

func TestInsertFilename(t *testing.T) {
	t.Run("entireString", func(t *testing.T) {
		result, err := pkg.InsertFilename("{}", "some string")
		require.NoError(t, err)
		assert.Equal(t, "some string", result)
	})
	t.Run("literalInsertion", func(t *testing.T) {
		result, err := pkg.InsertFilename("{{some string}}", "")
		require.NoError(t, err)
		assert.Equal(t, "{some string}", result)
	})
	t.Run("nthChar", func(t *testing.T) {
		result, err := pkg.InsertFilename("{:2}", "input")
		require.NoError(t, err)
		assert.Equal(t, "p", result)
	})
}
