package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBrackets(t *testing.T) {
	t.Run("entireString", func(t *testing.T) {
		result, err := parseBrackets("{}", "string")
		require.NoError(t, err)
		assert.Equal(t, "string", result)
	})
	t.Run("literalInsertion", func(t *testing.T) {
		result, err := parseBrackets("{{string}}", "")
		require.NoError(t, err)
		assert.Equal(t, "{string}", result)
	})
	t.Run("emptyLiteral", func(t *testing.T) {
		result, err := parseBrackets("{{}}", "")
		require.NoError(t, err)
		assert.Equal(t, "{}", result)
	})
	t.Run("nthChar", func(t *testing.T) {
		result, err := parseBrackets("{:2}", "input")
		require.NoError(t, err)
		assert.Equal(t, "p", result)
	})
	t.Run("indexRange", func(t *testing.T) {
		result, err := parseBrackets("{:2:4}", "index")
		require.NoError(t, err)
		assert.Equal(t, "de", result)
	})
	t.Run("charDelimitedIndex", func(t *testing.T) {
		result, err := parseBrackets("{/:1}", "some/path")
		require.NoError(t, err)
		assert.Equal(t, "path", result)
	})
	t.Run("stringDelimitedIndex", func(t *testing.T) {
		result, err := parseBrackets("{/p:1}", "some/path")
		require.NoError(t, err)
		assert.Equal(t, "ath", result)
	})
	t.Run("charDelimitedNegativeIndex", func(t *testing.T) {
		result, err := parseBrackets("{/:-1}", "some/path")
		require.NoError(t, err)
		assert.Equal(t, "path", result)
	})
}

func TestInsertFilename(t *testing.T) {
	t.Run("entireString", func(t *testing.T) {
		result, err := InsertFilename("{}", "some string")
		require.NoError(t, err)
		assert.Equal(t, "some string", result)
	})
	t.Run("literalInsertion", func(t *testing.T) {
		result, err := InsertFilename("{{some string}}", "")
		require.NoError(t, err)
		assert.Equal(t, "{some string}", result)
	})
	t.Run("nthChar", func(t *testing.T) {
		result, err := InsertFilename("{:2}", "input")
		require.NoError(t, err)
		assert.Equal(t, "p", result)
	})
	t.Run("indexRange", func(t *testing.T) {
		result, err := InsertFilename("{:2:4}", "index")
		require.NoError(t, err)
		assert.Equal(t, "de", result)
	})
}

func TestGeneratePatterns(t *testing.T) {
	t.Run("entireString", func(t *testing.T) {
		result, err := GeneratePatterns("{}", []string{"a", "b", "c"})
		require.NoError(t, err)
		assert.Equal(t, []string{"a", "b", "c"}, result)
	})
	t.Run("literalInsertion", func(t *testing.T) {
		result, err := GeneratePatterns("{{string}}", []string{""})
		require.NoError(t, err)
		assert.Equal(t, []string{"{string}"}, result)
	})
	t.Run("nthChar", func(t *testing.T) {
		result, err := GeneratePatterns("{:2}", []string{"input", "output"})
		require.NoError(t, err)
		assert.Equal(t, []string{"p", "t"}, result)
	})
	t.Run("indexRange", func(t *testing.T) {
		result, err := GeneratePatterns("{:2:4}", []string{"index", "test"})
		require.NoError(t, err)
		assert.Equal(t, []string{"de", "st"}, result)
	})
}
