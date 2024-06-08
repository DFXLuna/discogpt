package discogpt

import (
	"context"

	chroma "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/collection"
	huggingface "github.com/amikos-tech/chroma-go/hf"
	"github.com/samber/lo"
)

// http://localhost:8000
// http://localhost:8080/embed
// "collection-1"

const (
	numResults = 4
)

func NewChromaMod(baseURL string, teiURL string, collectionName string) (GenerationRequestModifier, error) {
	client, err := chroma.NewClient(baseURL)
	if err != nil {
		return nil, err
	}

	ef, err := huggingface.NewHuggingFaceEmbeddingInferenceFunction(teiURL)
	if err != nil {
		return nil, err
	}

	myCollection, err := client.NewCollection(context.Background(),
		collection.WithName(collectionName),
		collection.WithEmbeddingFunction(ef))
	if err != nil {
		return nil, err
	}

	return func(ocr *oaiCompletionsReq) error {
		results, err := myCollection.Query(context.TODO(),
			lo.FilterMap(ocr.Messages, func(item oaiMessage, _ int) (string, bool) {
				if item.Role == oaiUser {
					return item.Content, true
				}
				return "", false
			}),
			numResults, nil, nil, nil)
		if err != nil {
			return err
		}
		memories := "[The following are memories that may inform this interaction: "
		for _, res := range results.Documents {
			for _, str := range res {
				memories += str + "\n"
			}
		}
		memories += "]"
		memoriesMessage := oaiMessage{
			Role:    oaiSystem,
			Content: memories,
		}
		ocr.Messages = append(ocr.Messages[:len(ocr.Messages)-1], memoriesMessage, ocr.Messages[len(ocr.Messages)-1])
		return nil
	}, nil

	// rs, err := types.NewRecordSet(types.WithEmbeddingFunction(ef),
	// 	types.WithIDGenerator(types.NewULIDGenerator()))
	// 	if err != nil {
	// 		return err
	// 	}

	// rs.WithRecord(types.WithDocument("My name is Emilia and I am a programmer."))

	// _, err = rs.BuildAndValidate(context.Background())
	// if err != nil {
	// 	return err
	// }

	// _, err = myCollection.AddRecords(context.Background(), rs)
	// if err != nil {
	// 	return err
	// }

	// sz, err := myCollection.Count(context.Background())

	// docs, err := myCollection.Query(context.TODO(), []string{"Who is Emilia"}, 5, nil, nil, nil)
	// if err != nil {
	// 	return err
	// }

}