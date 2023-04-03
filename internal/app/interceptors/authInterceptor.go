package interceptors

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"golang.org/x/exp/slices"

	"protected-storage-server/internal/security"
)

type AuthInterceptor struct {
	jwtManager *security.JWTManager
}

func NewAuthInterceptor(jwtManager *security.JWTManager) *AuthInterceptor {
	return &AuthInterceptor{jwtManager}
}

func accessibleMethods() []string {
	const servicePath = "/server.GrpcService/"

	return []string{
		servicePath + "CreateUser",
		servicePath + "LoginUser",
	}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Info().Msgf("AuthInterceptor intercept method %s", info.FullMethod)

		if slices.Contains(accessibleMethods(), info.FullMethod) {
			return handler(ctx, req)
		}

		err := interceptor.authorize(ctx)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	_, err := interceptor.jwtManager.VerifyTokenAndExtractUserID(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	return nil
}
