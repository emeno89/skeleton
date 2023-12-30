package user

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"net/http"
	"testing"
)

type userListerMock struct{ mock.Mock }

func (s *userListerMock) GetByAuthId(ctx context.Context, id string) (Data, bool, error) {
	args := s.Called(ctx, id)

	return args.Get(0).(Data), args.Bool(1), args.Error(2)
}

type jwtJenMock struct{ mock.Mock }

func (s *jwtJenMock) ParseUserId(token string) (string, error) {
	args := s.Called(token)

	return args.Get(0).(string), args.Error(1)
}

type userSuite struct {
	suite.Suite
	jwtJen     *jwtJenMock
	userLister *userListerMock
	srv        *Service
}

func (s *userSuite) SetupTest() {
	s.jwtJen = &jwtJenMock{}
	s.userLister = &userListerMock{}

	s.srv = newService(
		s.jwtJen,
		s.userLister,
		zap.NewNop(),
	)
}

func (s *userSuite) TestParseFromRequest() {
	type test struct {
		jwt    string
		userId string
		err    error
		want   RequestInfo
	}

	userIdOk := "6516e4302470c9c9368f17e2"

	tests := []test{
		{jwt: "Bearer 123456", err: errors.New("invalid")},                        //invalid
		{jwt: "Bearer 654321", userId: userIdOk, want: RequestInfo{Id: userIdOk}}, //ok
	}

	for _, val := range tests {
		s.jwtJen.On("ParseUserId", val.jwt).Return(val.userId, val.err)

		r, _ := http.NewRequest("", "", nil)
		r.Header.Add("Authorization", val.jwt)

		rInfo := s.srv.ParseFromRequest(r)

		s.Assert().Equal(val.want, rInfo)
	}
}

func (s *userSuite) TestFromIncomingContext() {
	type test struct {
		ctx  context.Context
		usr  Data
		want bool
	}

	usrOk := Data{Id: "6516e4302470c9c9368f17e2"}

	tests := []test{
		{ctx: context.WithValue(context.Background(), dataKey{}, usrOk), usr: usrOk, want: true},
		{ctx: context.WithValue(context.Background(), dataKey{}, Data{}), usr: Data{}},
		{ctx: context.Background()},
	}

	for _, val := range tests {
		usr, ok := s.srv.FromIncomingContext(val.ctx)

		s.Assert().Equal(val.usr, usr)
		s.Assert().Equal(val.want, ok)
	}
}

func (s *userSuite) TestCurrentUser() {
	type test struct {
		ctx context.Context
		usr Data
		err error
	}

	usrOk := Data{Id: "6516e4302470c9c9368f17e2"}

	tests := []test{
		{ctx: context.Background(), err: ErrNotFound},
		{ctx: context.WithValue(context.Background(), dataKey{}, usrOk), usr: usrOk},
		{ctx: context.WithValue(context.Background(), dataKey{}, Data{}), err: ErrNotFound},
	}

	for _, val := range tests {
		usr, err := s.srv.CurrentUser(val.ctx)

		s.Assert().Equal(val.usr, usr)
		s.Assert().Equal(val.err, err)
	}
}

func (s *userSuite) TestToOutgoingContextNoToken() {
	r := &http.Request{}

	s.jwtJen.On("ParseUserId", "").Return("", errors.New("invalid"))

	ctx, err := s.srv.ToOutgoingContext(r)

	s.Assert().Equal(r.Context(), ctx)
	s.Assert().Nil(err)
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(userSuite))
}
