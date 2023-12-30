package loader

import (
	"context"
	"github.com/graph-gophers/dataloader"
	"skeleton/shared/gqlutils"
	"skeleton/srv-gql/internal/api/graphql/graph/model"
	"skeleton/srv-gql/internal/pkg/client/example"
	"time"
)

const (
	waitDuration = 2 * time.Millisecond
)

type batch struct {
	getExampleItemLoader *dataloader.Loader
}

func makeBatch(reader *reader) batch {
	return batch{
		getExampleItemLoader: dataloader.NewBatchedLoader(reader.getExampleItem, dataloader.WithWait(waitDuration)),
	}
}

type reader struct {
	exampleClient example.Client
}

func newReader(exampleClient example.Client) *reader {
	return &reader{
		exampleClient: exampleClient,
	}
}

func (s *reader) getExampleItem(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	list, err := s.exampleClient.GetManyById(ctx, keys.Keys()...)

	mp := make(map[string]*model.ExampleItem)

	if err == nil {
		for _, val := range list {
			mp[val.Id] = &model.ExampleItem{
				Id: val.Id,
			}
		}
	}

	output := make([]*dataloader.Result, len(keys))

	for index, key := range keys {
		if err != nil {
			output[index] = &dataloader.Result{Error: gqlutils.ConvertGPRCError(ctx, err)}
		} else {
			output[index] = &dataloader.Result{Data: mp[key.String()]}
		}
	}

	return output
}
