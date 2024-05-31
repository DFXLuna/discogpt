package test

import (
	"context"
	discogpt "egrant/disco-gpt/pkg"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageGenerator(t *testing.T) {
	assert := assert.New(t)
	assert.True(true)

	require := require.New(t)

	gen, err := discogpt.NewOpenAIGenerator("http://192.168.1.6:5001", "", "")
	require.NoError(err, "shouldn't error on NewOpenAIGen")

	out, err := gen.Generate(context.Background(), "Tell me a joke, please.", "Emilia")
	require.NoError(err, "shouldn't error on generate")
	fmt.Printf("Response:\n %s\n", out)
}
