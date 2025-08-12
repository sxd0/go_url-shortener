package kafka

import "time"

type Kind string

const (
	LinkVisitedKind Kind = "link.visited"
	LinkCreatedKind Kind = "link.created"
)

type Event struct {
	Version int32     `json:"version"`
	Kind    Kind      `json:"kind"`
	LinkID  uint      `json:"link_id,omitempty"`
	UserID  uint      `json:"user_id,omitempty"`
	Ts      time.Time `json:"ts"`
}
