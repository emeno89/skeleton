package locale

import (
	"context"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/metadata"
	"testing"
)

type LocaleSuite struct {
	suite.Suite
}

func (s *LocaleSuite) TestAddAcceptLanguage() {
	locales := []string{"ru-RU", "uz-UZ"}

	ctx := context.TODO()

	for _, val := range locales {
		ctx = AddAcceptLanguage(ctx, val)
	}

	md, ok := metadata.FromOutgoingContext(ctx)

	getRes := md.Get(acceptLanguageHeader)

	s.Assert().Equal(true, ok)
	s.Assert().Len(getRes, 2)

	s.Assert().Equal(locales, getRes)
}

func TestLocaleSuite(t *testing.T) {
	suite.Run(t, new(LocaleSuite))
}
