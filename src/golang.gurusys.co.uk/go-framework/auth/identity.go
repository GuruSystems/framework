package auth

import (
	pb "golang.gurusys.co.uk/apis/auth"
)

type User struct {
	FirstName string
	LastName  string
	Email     string
	ID        string
	Groups    []*Group
}

type Group struct {
	ID   string
	Name string
}

// injected into a context
type AuthInfo struct {
	UserID string
}

type Authenticator interface {
	// give token -> return userid
	Authenticate(token string) (string, error)
	GetUserDetail(userid string) (*User, error)
	// given a previous challenge and an email, will return a token if challenge and password stuff matches
	CreateVerifiedToken(email string, password string) string
	CreateUser(*pb.CreateUserRequest) (string, error)
	GetUserByEmail(*pb.UserByEmailRequest) ([]*User, error)
	AddUserToGroup(req *pb.AddToGroupRequest) ([]*User, error)
	RemoveUserFromGroup(req *pb.RemoveFromGroupRequest) ([]*User, error)
	ListUsersInGroup(req *pb.ListGroupRequest) ([]*User, error)
}
