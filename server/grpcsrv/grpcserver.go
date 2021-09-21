package grpcsrv

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
	"todo/api/v1/pb"
	"todo/config"
	"todo/logger"
	"todo/model"
	"todo/service"
	"todo/storage"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedUsersServer
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
	pb.RegisterUsersServer(server, &Server{
		config:  c,
		log:     l,
		service: s,
	})
	return server

}

func (s *Server) AddUser(ctx context.Context, in *pb.AddUserRequest) (*pb.AddUserReply, error) {
	user := model.User{
		UserName:  in.UserName,
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Password:  in.Password,
		Location:  getCustomLocation(in.Location),
	}

	id, err := s.service.AddUser(ctx, user)

	if err != nil {
		s.log.Errorf("Could not add user %v", err)
		return nil, err
	}

	return &pb.AddUserReply{Id: id}, nil
}

func (s *Server) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.UpdateUserReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}

	userid, ok := md["userid"]
	if !ok {
		s.log.Errorf("%q", "Could not get user.")
		return nil, fmt.Errorf("%q: %w", "userid is not provided.", model.ErrBadRequest)
	}

	user := model.User{
		UserName:  in.UserName,
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Password:  in.Password,
		Location:  getCustomLocation(in.Location),
	}

	err := s.service.UpdateUser(ctx, userid[0], user)

	if err != nil {
		s.log.Errorf("Could not update user %v", err)
		return nil, err
	}

	return &pb.UpdateUserReply{}, nil
}

func (s *Server) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.DeleteUserReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}

	userid, ok := md["userid"]
	if !ok {
		s.log.Errorf("%q", "Could not get user.")
		return nil, fmt.Errorf("%q: %w", "userid is not provided.", model.ErrBadRequest)
	}

	err := s.service.DeleteUser(ctx, userid[0])
	if err != nil {
		s.log.Errorf("%q: %w", "Could not delete user.", err)
		return nil, err
	}
	return &pb.DeleteUserReply{}, nil
}

func (s *Server) GetAllUsers(ctx context.Context, in *pb.GetAllUsersRequest) (*pb.GetAllUsersReply, error) {
	filter := storage.UserFilter{}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}

	username, ok := md["username"]
	if ok {
		filter.UserName = username[0]
	}

	users, err := s.service.GetUsers(ctx, filter)
	if err != nil {
		s.log.Errorf("%q: 	", "Could not get all users.")
		return nil, err
	}

	usersReply := &pb.GetAllUsersReply{}
	for _, u := range users {
		usersReply.Users = append(usersReply.Users, &pb.User{
			Id:        u.Id,
			UserName:  u.UserName,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Location:  u.Location.String(),
		})
	}
	return usersReply, nil
}

func (s *Server) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}

	userid, ok := md["userid"]
	if !ok {
		s.log.Errorf("%q", "Could not get user.")
		return nil, fmt.Errorf("%q: %w", "userid is not provided.", model.ErrBadRequest)
	}

	user, err := s.service.GetUser(ctx, userid[0])
	if err != nil {
		s.log.Errorf("%q: %w", "Could not get user.", err)
		return nil, err
	}

	userReply := &pb.GetUserReply{}
	userReply.User = &pb.User{
		Id:        user.Id,
		UserName:  user.UserName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Location:  user.Location.String(),
	}

	return userReply, nil
}

func (s *Server) LoginUser(ctx context.Context, in *pb.LoginRequest) (*pb.LoginReply, error) {
	credentials := model.Credentials{UserName: in.GetUserName(), Password: in.GetPassword()}
	token, err := s.service.LoginUser(ctx, credentials)

	if err != nil {
		s.log.Errorf("%q: %w", "Can't parse credentials.", err)
		return nil, err
	}

	return &pb.LoginReply{
		Token: token.TokenString,
	}, nil
}
func getCustomLocation(s string) model.CustomLocation {
	loc, err := time.LoadLocation(s)
	if err != nil {
		return model.CustomLocation{}
	}
	return model.CustomLocation{loc}
}

func (s *Server) AddTodo(ctx context.Context, in *pb.AddTodoRequest) (*pb.AddTodoReply, error) {
	todo := model.TodoItem{
		Name:   in.Name,
		Date:   in.Date.AsTime(),
		Status: in.Status,
	}

	id, err := s.service.AddTodo(ctx, todo)

	if err != nil {
		s.log.Errorf("Could not add todo %v", err)
		return nil, err
	}

	return &pb.AddTodoReply{Id: id}, nil
}

func (s *Server) GetAllTodos(ctx context.Context, in *pb.GetAllTodosRequest) (*pb.GetAllTodosReply, error) {
	filter := storage.TodoFilter{}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}

	status, ok := md["status"]
	if ok {
		filter.Status = status[0]
	}

	date, ok := md["fromdate"]
	if ok {
		fromDate, err := time.Parse(time.RFC3339, date[0])
		if err != nil {
			s.log.Errorf("Could not parse date in GetAllTodos %v", err)
			return nil, err
		}
		filter.FromDate = &fromDate
	}

	date, ok = md["todate"]
	if ok {
		toDate, err := time.Parse(time.RFC3339, date[0])
		if err != nil {
			s.log.Errorf("Could not parse date in GetAllTodos %v", err)
			return nil, err
		}
		filter.ToDate = &toDate
	}

	todos, err := s.service.GetTodos(ctx, filter)
	if err != nil {
		s.log.Errorf("%q: 	", "Could not get all todos.")
		return nil, err
	}

	todosReply := &pb.GetAllTodosReply{}
	for _, u := range todos {
		todosReply.Todos = append(todosReply.Todos, &pb.Todo{
			Id:     u.Id,
			Status: u.Status,
			Name:   u.Name,
			Date:   timestamppb.New(u.Date),
		})
	}
	return todosReply, nil
}

func (s *Server) GetTodo(ctx context.Context, in *pb.GetTodoRequest) (*pb.GetTodoReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}

	todoid, ok := md["todoid"]
	if !ok {
		s.log.Errorf("%q", "Could not get todo.")
		return nil, fmt.Errorf("%q: %w", "todoid is not provided.", model.ErrBadRequest)
	}

	todo, err := s.service.GetTodo(ctx, todoid[0])
	if err != nil {
		s.log.Errorf("%q: %w", "Could not get todo.", err)
		return nil, err
	}

	todoReply := &pb.GetTodoReply{}
	todoReply.Todo = &pb.Todo{
		Id:     todo.Id,
		Name:   todo.Name,
		Status: todo.Status,
		Date:   timestamppb.New(todo.Date),
	}

	return todoReply, nil
}

func (s *Server) UpdateTodo(ctx context.Context, in *pb.UpdateTodoRequest) (*pb.UpdateTodoReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}

	todoid, ok := md["todoid"]
	if !ok {
		s.log.Errorf("%q", "Could not update todo.")
		return nil, fmt.Errorf("%q: %w", "todoid is not provided.", model.ErrBadRequest)
	}

	todo := model.TodoItem{
		Name:   in.Name,
		Status: in.Status,
		Date:   in.Date.AsTime(),
	}

	err := s.service.UpdateTodo(ctx, todoid[0], todo)

	if err != nil {
		s.log.Errorf("Could not update todo %v", err)
		return nil, err
	}

	return &pb.UpdateTodoReply{}, nil
}

func (s *Server) DeleteTodo(ctx context.Context, in *pb.DeleteTodoRequest) (*pb.DeleteTodoReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}

	todoid, ok := md["todoid"]
	if !ok {
		s.log.Errorf("%q", "Could not delete todo.")
		return nil, fmt.Errorf("%q: %w", "todoid is not provided.", model.ErrBadRequest)
	}

	err := s.service.DeleteTodo(ctx, todoid[0])
	if err != nil {
		s.log.Errorf("%q: %w", "Could not delete todo.", err)
		return nil, err
	}
	return &pb.DeleteTodoReply{}, nil
}
