package user

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"skeleton/srv-gql/internal/pkg/security"
)

type jwtGen interface {
	ParseUserId(token string) (string, error)
}

func newJwtJen(securityGen *security.JwtGen) jwtGen {
	return securityGen
}

type userLister interface {
	GetByAuthId(ctx context.Context, id string) (Data, bool, error)
}

func newUserLister(lister *exampleLister) userLister {
	return lister
}

type exampleLister struct{}

func newExampleLister() *exampleLister {
	return &exampleLister{}
}

func (s *exampleLister) GetByAuthId(_ context.Context, id string) (Data, bool, error) {
	return Data{
		Id:   id,
		Name: "Example name",
	}, true, nil
}

type Service struct {
	jwtGen jwtGen
	lister userLister
	logger *zap.Logger
}

func newService(
	jwtGen jwtGen,
	lister userLister,
	logger *zap.Logger,
) *Service {
	return &Service{
		jwtGen: jwtGen,
		lister: lister,
		logger: logger,
	}
}

func (s *Service) ToOutgoingContext(r *http.Request) (context.Context, error) {
	rInfo := s.ParseFromRequest(r)

	ctx := r.Context()

	if rInfo.Id == "" {
		return ctx, nil
	}

	data, ok, err := s.lister.GetByAuthId(ctx, rInfo.Id)
	if err != nil {
		s.logger.Error("ToOutgoingContext err", zap.Any("rInfo", rInfo), zap.Error(err))
		return ctx, err
	}

	if !ok {
		return ctx, nil
	}

	return context.WithValue(r.Context(), dataKey{}, data), nil
}

func (s *Service) FromIncomingContext(ctx context.Context) (Data, bool) {
	usr, ok := ctx.Value(dataKey{}).(Data)
	if !ok {
		return Data{}, false
	}

	if usr.Id == "" {
		return Data{}, false
	}

	return usr, true
}

func (s *Service) CurrentUser(ctx context.Context) (Data, error) {
	usr, ok := s.FromIncomingContext(ctx)
	if !ok {
		return Data{}, ErrNotFound
	}

	return usr, nil
}

func (s *Service) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx, err := s.ToOutgoingContext(r)
			if err != nil {
				http.Error(w, "Auth check err", http.StatusInternalServerError)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}

func (s *Service) ParseFromRequest(r *http.Request) RequestInfo {
	userId, err := s.jwtGen.ParseUserId(r.Header.Get("Authorization"))
	if err != nil {
		return RequestInfo{}
	}

	return RequestInfo{
		Id: userId,
	}
}
