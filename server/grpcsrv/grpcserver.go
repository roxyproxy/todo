package grpcsrv

import (
	"context"
	"todo/api/v1/pb"
	"todo/config"
	"todo/logger"
	"todo/model"
	"todo/service"
	"todo/storage"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedUserServer
	storage storage.Storage
	config  *config.Config
	log     logger.Logger
	service service.Handlers
}

func NewGrpcServer(storage storage.Storage, c *config.Config, l logger.Logger) *grpc.Server {
	s := service.NewService(storage, c)
	authMD := AuthMD{service: s, config: c}
	opts := make([]grpc.ServerOption, 0)
	opts = append(opts, grpc.ChainUnaryInterceptor(authMD.UnaryInterceptor()))

	server := grpc.NewServer(opts...)
	pb.RegisterUserServer(server, &Server{
		config:  c,
		log:     l,
		service: s,
	})
	return server

}

func (s *Server) AddUser(ctx context.Context, in *pb.AddUserRequest) (*pb.AddUserReply, error) {
	user := model.User{
		UserName:  in.GetUserName(),
		FirstName: in.GetFirstName(),
		LastName:  in.GetLastName(),
		Password:  in.GetPassword(),
		//Location:  in.GetLocation(),
	}

	id, err := s.service.AddUser(user)

	if err != nil {
		s.log.Errorf("Could not add user %v", err)
		return nil, err
	}

	return &pb.AddUserReply{Id: id}, nil
}

func (s *Server) GetAllUsers(ctx context.Context, in *pb.GetAllUsersRequest) (*pb.GetAllUsersReply, error) {
	// userid := getUserFromContext(r)
	filter := storage.UserFilter{}
	if username, ok := ctx.Value("username").(string); ok {
		filter.UserName = username
	}

	users, err := s.service.GetUsers(filter, ctx)
	if err != nil {
		s.log.Errorf("%q: %w", "Could not get all users.")
		return nil, err
	}

	usersReply := &pb.GetAllUsersReply{}
	for _, u := range users {
		usersReply.Result = append(usersReply.Result, &pb.GetAllUsersReply_User{
			Id:        u.Id,
			UserName:  u.UserName,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			//Location:  u.Location,
		})
	}

	return usersReply, nil
}

func (s *Server) LoginUser(ctx context.Context, in *pb.LoginRequest) (*pb.LoginReply, error) {
	credentials := model.Credentials{UserName: in.GetUserName(), Password: in.GetPassword()}
	token, err := s.service.LoginUser(credentials)

	if err != nil {
		s.log.Errorf("%q: %w", "Can't parse credentials.", err)
		return nil, err
	}

	return &pb.LoginReply{
		Token: token.TokenString,
	}, nil
}
