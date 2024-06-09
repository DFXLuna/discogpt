package discogpt

import (
	"context"

	chroma "github.com/amikos-tech/chroma-go"
	huggingface "github.com/amikos-tech/chroma-go/hf"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/samber/lo"
)

const (
	numResults = 4
)

func NewChromaMod(baseURL string, teiURL string, collectionName string, log Logger) (GenerationRequestModifier, error) {
	client, err := chroma.NewClient(baseURL)
	if err != nil {
		return nil, err
	}

	ef, err := huggingface.NewHuggingFaceEmbeddingInferenceFunction(teiURL)
	if err != nil {
		return nil, err
	}

	myCollection, err := client.CreateCollection(context.TODO(), collectionName, map[string]interface{}{}, true, ef, types.L2)
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
		log.Debugf("Modified memories: %+v", ocr.Messages)
		err = makeRecord(context.TODO(), ef, myCollection, ocr.Messages[len(ocr.Messages)-1].Content, log)
		if err != nil {
			log.Errorf("Ignoring error sending message to Chroma")
			return nil
		}
		msgs, err := myCollection.Count(context.TODO())
		if err != nil {
			log.Errorf("Ignoring error counting messages in Chroma")
			return nil
		}
		log.Debugf("Found %d messages in Chroma", msgs)
		return nil
	}, nil
}

func makeRecord(ctx context.Context, ef types.EmbeddingFunction,
	myCollection *chroma.Collection, message string, log Logger) error {
	rs, err := types.NewRecordSet(types.WithEmbeddingFunction(ef),
		types.WithIDGenerator(types.NewULIDGenerator()))
	if err != nil {
		return err
	}

	rs.WithRecord(types.WithDocument(message))

	_, err = rs.BuildAndValidate(ctx)
	if err != nil {
		return err
	}

	_, err = myCollection.AddRecords(ctx, rs)
	if err != nil {
		return err
	}
	log.Debugf("Stored message in Chroma")
	return nil
}
