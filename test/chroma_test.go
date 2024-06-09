package test

import (
	"context"
	"fmt"
	"testing"

	chroma "github.com/amikos-tech/chroma-go"
	huggingface "github.com/amikos-tech/chroma-go/hf"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/stretchr/testify/require"
)

func TestChroma(t *testing.T) {
	require := require.New(t)

	client, err := chroma.NewClient("http://localhost:8000")
	require.NoError(err)

	ef, err := huggingface.NewHuggingFaceEmbeddingInferenceFunction("http://localhost:8080/embed")
	require.NoError(err)

	myCollection, err := client.CreateCollection(context.Background(), "chadbot-test", map[string]interface{}{}, true,
		ef, types.L2)
	require.NoError(err)

	rs, err := types.NewRecordSet(types.WithEmbeddingFunction(ef),
		types.WithIDGenerator(types.NewULIDGenerator()))
	require.NoError(err)

	rs.WithRecord(types.WithDocument("My name is Emilia and I am a programmer."))

	_, err = rs.BuildAndValidate(context.Background())
	require.NoError(err)

	_, err = myCollection.AddRecords(context.Background(), rs)
	require.NoError(err)

	sz, err := myCollection.Count(context.Background())
	require.NoError(err)

	fmt.Println("Size of collection: ", sz)

	docs, err := myCollection.Query(context.TODO(), []string{"Who is Emilia"}, 5, nil, nil, nil)
	require.NoError(err)

	fmt.Printf("results: %+v\n", docs)
}
