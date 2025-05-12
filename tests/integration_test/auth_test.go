package server_test

import (
	"MamangRust/paymentgatewaygrpc/internal/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerTestSuite) TestAuthRegister() {
	s.Run("Success Register", func() {

		registerRequest := &pb.RegisterRequest{
			Firstname:       "John",
			Lastname:        "Doe",
			Email:           "test@example.com",
			Password:        "password123",
			ConfirmPassword: "password123",
		}

		registerResp, err := s.authClient.RegisterUser(s.ctx, registerRequest)

		myexpected := &pb.ApiResponseRegister{
			Status:  "success",
			Message: "Registration successful",
			Data: &pb.UserResponse{
				Id:        1,
				Firstname: "John",
				Lastname:  "Doe",
				Email:     "test@example.com",
			},
		}

		s.NoError(err)
		s.NotNil(registerResp)
		s.Equal("success", registerResp.Status)
		s.Equal("Registration successful", registerResp.Message)

		s.Equal(myexpected.Status, registerResp.Status)
		s.Equal(myexpected.Message, registerResp.Message)
		s.Equal(myexpected.Data.Id, registerResp.Data.Id)
		s.Equal(myexpected.Data.Firstname, registerResp.Data.Firstname)
		s.Equal(myexpected.Data.Lastname, registerResp.Data.Lastname)
		s.Equal(myexpected.Data.Email, registerResp.Data.Email)
	})

	s.Run("Failed Register", func() {
		registerRequest := &pb.RegisterRequest{
			Firstname:       "John",
			Lastname:        "Doe",
			Email:           "test@example.com",
			Password:        "password123",
			ConfirmPassword: "password123",
		}

		res, err := s.authClient.RegisterUser(s.ctx, registerRequest)

		statuf := status.Errorf(
			codes.InvalidArgument,
			"%v",
			&pb.ErrorResponse{
				Status:  "error",
				Message: "Registration failed: ",
			},
		)

		s.Nil(res)
		s.NotNil(err)

		s.Equal(statuf, err)
	})

}

func (s *ServerTestSuite) TestAuthLogin() {
	s.Run("Success Login", func() {

		loginRequest := &pb.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		loginResp, err := s.authClient.LoginUser(s.ctx, loginRequest)

		s.NoError(err)
		s.NotNil(loginResp)
		s.Equal("success", loginResp.Status)
		s.Equal("Login successful", loginResp.Message)

	})

	s.Run("Failed Login", func() {
		loginRequest := &pb.LoginRequest{
			Email:    "test@example.com",
			Password: "wrong-password",
		}

		loginResp, err := s.authClient.LoginUser(s.ctx, loginRequest)

		s.Nil(loginResp)
		s.NotNil(err)

		statuf := status.Errorf(
			codes.Unauthenticated,
			"%v",
			&pb.ErrorResponse{
				Status:  "error",
				Message: "Login failed: ",
			},
		)

		s.Equal(statuf.Error(), err.Error())
	})
}
