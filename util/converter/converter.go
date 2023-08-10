package converter

import (
	"time"

	"github.com/egiferdians/micro-auth/auth"
	"github.com/egiferdians/micro-auth/models"
	"github.com/egiferdians/micro-auth/protobuf/pb"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func timeGolangTimeToPbTimestamp(golangTime time.Time) *timestamp.Timestamp {
	timestamps, _ := ptypes.TimestampProto(golangTime)
	return timestamps
}

func ConvertToPBAuth(Auth *auth.Authenticated) *pb.DataAuth {
	return &pb.DataAuth{
		User:           ConvertToPBUser(Auth.User),
		RefreshToken:        Auth.RefreshToken,
		AccessToken: Auth.AccessToken,
	}
}

func ConvertToPBUser(User *models.User) *pb.DataUser {
	return &pb.DataUser{
		IdUser:        User.IDUser.String(),
		Fullname: User.Fullname,
		Email: User.Email,
	}
}
