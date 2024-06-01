package test

import (
	"context"
	discogpt "egrant/disco-gpt/pkg"
	mock_discogpt "egrant/disco-gpt/pkg/mock"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestMessageGenerator(t *testing.T) {
	assert := assert.New(t)
	assert.True(true)

	require := require.New(t)
	ctrl := gomock.NewController(t)
	mock_log := mock_discogpt.NewMockLogger(ctrl)

	gen, err := discogpt.NewOpenAIGenerator("http://192.168.1.6:5001", "", "", mock_log)
	require.NoError(err, "shouldn't error on NewOpenAIGen")

	mock_log.EXPECT().Debugf("Generating for %v", "Emilia").Times(1)
	out, err := gen.Generate(context.Background(), "Tell me a joke, please.", "Emilia")
	require.NoError(err, "shouldn't error on generate")
	fmt.Printf("Response:\n %s\n", out)
}
