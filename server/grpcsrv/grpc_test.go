package grpcsrv

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	//"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"io/ioutil"
	"net"
	"testing"
	//"time"
	"todo/api/v1/pb"
	conf "todo/config"
	"todo/logger"
	//"todo/model"
	mockstore "todo/server/mocks"
	//"todo/service"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func TestServerWithMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mockstore.NewMockStorage(ctrl)
	log := logger.New(ioutil.Discard)
	lis = bufconn.Listen(bufSize)

	server := NewGrpcServer(m, conf.New(), log)
	fmt.Printf("grpc server listening at %v", lis.Addr())
	go func() {
		if err := server.Serve(lis); err != nil {
			fmt.Printf("grpc failed to serve: %v", err)
		}
	}()

	//config := conf.New()
	//service := service.NewService(m, config)

	/*
		username := "Roxy"
		password := "SecretPassword12!"
		l, _ := time.LoadLocation("America/New_York")
		location := model.CustomLocation{Location: l}
		credentials := model.Credentials{UserName: username, Password: password}

		hash, err := service.HashPassword(password)
		assert.NoError(t, err)
		user := model.User{
			Id:        "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			UserName:  username,
			FirstName: "Roxy",
			LastName:  "Proxy",
			Password:  hash,
			Location:  location,
		}


		token, _ := service.GenerateToken(user.Id, config.SecretKey)
	*/

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, lis.Addr().String(), grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	m.EXPECT().AddUser(gomock.Any()).Return("123", nil)

	client := pb.NewUsersClient(conn)
	resp, err := client.AddUser(ctx, &pb.AddUserRequest{
		UserName:  "RoxyP",
		FirstName: "Roxy",
		LastName:  "Proxy",
		Password:  "Test",
		Location:  "America/New_York",
	})
	if err != nil {
		t.Fatalf("AddUser failed: %v", err)
	}
	fmt.Printf("Response: %+v", resp)
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}
