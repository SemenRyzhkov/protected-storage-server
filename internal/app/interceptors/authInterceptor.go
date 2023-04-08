package interceptors

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"golang.org/x/exp/slices"

	"protected-storage-server/internal/security"
)

// AuthInterceptor type
type AuthInterceptor struct {
	jwtManager *security.JWTManager
}

// NewAuthInterceptor constructor
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

// Unary is a gRPC auth interceptor
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
	accessToken, err := interceptor.jwtManager.ExtractJWTFromContext(ctx)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, err.Error())
	}

	log.Printf("Access token is %s", accessToken)

	token, err := interceptor.jwtManager.ParseToken(accessToken)
	if err != nil || !token.Valid {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	return nil
}
