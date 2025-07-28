package payload

type GetStatResponse struct {
	Date   string `json:"date"`
	LinkId uint   `json:"link_id"`
	Clicks int    `json:"clicks"`
}
