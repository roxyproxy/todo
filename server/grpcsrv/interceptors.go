package grpcsrv

import (
	"context"
	"todo/config"
	"todo/model"
	"todo/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthMD represents AuthMD structure.
type AuthMD struct {
	service service.Handlers
	config  *config.Config
}

// UnaryInterceptor middleware for grps server.
func (a *AuthMD) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if info.FullMethod != "/users.Users/AddUser" && info.FullMethod != "/users.Users/LoginUser" {
			ctx, err = a.authorize(ctx)
			if err != nil {
				return nil, err
			}
		}
		return handler(ctx, req)
	}
}

func (a *AuthMD) authorize(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}

	authHeader, ok := md["authorization"]
	if !ok {
		return ctx, status.Errorf(codes.Unauthenticated, "Authorization token is not supplied")
	}

	token := authHeader[0]

	claims, err := a.service.ValidateToken(ctx, token)
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, err.Error())
	}

	ctx = context.WithValue(ctx, model.KeyUserID("userid"), claims.UserID)

	return ctx, nil
}
