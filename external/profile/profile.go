package profile

const (
	address = "localhost:9001"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

type ProfileId struct {
	Value string
}
