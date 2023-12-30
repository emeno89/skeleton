package example

import "context"

type Service struct{}

func newService() *Service {
	return &Service{}
}

func (s *Service) GetManyId(_ context.Context, ids ...string) ([]Data, error) {
	result := make([]Data, len(ids))

	for i, val := range ids {
		result[i] = Data{
			Id: val,
		}
	}

	return result, nil
}
