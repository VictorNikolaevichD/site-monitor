package messaging

import "context"

type EventPublisher interface {
	Publish(ctx context.Context, evennt SiteCheckEvent) error
	Close() error
}
