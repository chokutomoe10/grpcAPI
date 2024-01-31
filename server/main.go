package main

import (
	"context"
	"grpcApi/protobuf"
	"log"
	"math/rand"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

var users []*protobuf.UserInfo

type userServer struct {
	protobuf.UnimplementedCrudServer
}

func main() {
	initUsers()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	protobuf.RegisterCrudServer(s, &userServer{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func initUsers() {
	movie1 := &protobuf.UserInfo{Name: "Rachman", RoleId: "3",
		RoleName: "Admin", Email: "rachman@gmail.com"}
	movie2 := &protobuf.UserInfo{Name: "Kurniawan", RoleId: "3",
		RoleName: "Admin", Email: "kurniawan@gmail.com"}

	users = append(users, movie1)
	users = append(users, movie2)
}

func (s *userServer) GetAllUsers(in *protobuf.Empty,
	stream protobuf.Crud_GetAllUsersServer) error {
	log.Printf("Received: %v", in)
	for _, user := range users {
		if err := stream.Send(user); err != nil {
			return err
		}
	}
	return nil
}

func (s *userServer) CreateUser(ctx context.Context,
	in *protobuf.UserInfo) (*protobuf.Id, error) {
	log.Printf("Received: %v", in)
	res := protobuf.Id{}
	res.Value = strconv.Itoa(rand.Intn(100000000))
	in.Name = res.GetValue()
	users = append(users, in)
	return &res, nil
}

func (s *userServer) UpdateMovie(ctx context.Context,
	in *protobuf.UserInfo) (*protobuf.Status, error) {
	log.Printf("Received: %v", in)

	res := protobuf.Status{}
	for index, user := range users {
		if user.GetName() == in.GetName() {
			users = append(users[:index], users[index+1:]...)
			in.Name = user.GetName()
			users = append(users, in)
			res.Value = 1
			break
		}
	}

	return &res, nil
}

func (s *userServer) DeleteMovie(ctx context.Context,
	in *protobuf.Id) (*protobuf.Status, error) {
	log.Printf("Received: %v", in)

	res := protobuf.Status{}
	for index, user := range users {
		if user.GetName() == in.GetValue() {
			users = append(users[:index], users[index+1:]...)
			res.Value = 1
			break
		}
	}

	return &res, nil
}
