package loader

import (
	"context"
	"github.com/graph-gophers/dataloader"
	"skeleton/srv-gql/internal/api/graphql/graph/model"
	"skeleton/srv-gql/internal/pkg/client/example"
)

type loadersKey struct{}

func ToOutgoingContext(
	ctx context.Context,
	exampleClient example.Client,
) context.Context {
	loader := newLoader(
		newReader(exampleClient),
	)

	return context.WithValue(ctx, loadersKey{}, loader)
}

func FromIncomingContext(ctx context.Context) *Loader {
	return ctx.Value(loadersKey{}).(*Loader)
}

type Loader struct {
	batch
}

func newLoader(reader *reader) *Loader {
	return &Loader{
		batch: makeBatch(reader),
	}
}

func (s *Loader) GetExampleItem(ctx context.Context, id string) (*model.ExampleItem, error) {
	thunk := s.getExampleItemLoader.Load(ctx, dataloader.StringKey(id))

	result, err := thunk()
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	return result.(*model.ExampleItem), nil
}
