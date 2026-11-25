package health

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	bookingv1 "github.com/AntowQa/ai-for-developers-project-386/api/gen/booking/v1"
)

type Server struct {
	bookingv1.UnimplementedHealthServiceServer
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Server {
	return &Server{pool: pool}
}

func (s *Server) Check(ctx context.Context, _ *bookingv1.CheckRequest) (*bookingv1.CheckResponse, error) {
	if err := s.pool.Ping(ctx); err != nil {
		return nil, status.Error(codes.Unavailable, "database unavailable")
	}
	return &bookingv1.CheckResponse{Status: "ok"}, nil
}
