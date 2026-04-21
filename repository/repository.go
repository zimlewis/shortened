package repository

import "context"

type Repository interface {
	AddShortenedLink(ctx context.Context, shortened string, full string) error
	GetShortenedResult(ctx context.Context, shortened string) (string, error)
	IncreaseLinkClick(ctx context.Context, shortened string) (int, error)
	GetClickedCount(ctx context.Context, shortened string) (uint32, error)
}
