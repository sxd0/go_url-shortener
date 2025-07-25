package event

const EventLinkVisited = "link.visited"

type LinkVisitedEvent struct {
	LinkID uint
	UserID uint
}
