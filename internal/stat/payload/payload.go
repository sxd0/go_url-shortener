package payload

type GetStatResponse struct {
	Period string `json:"period"`
	Sum    int    `json:"sum"`
}
