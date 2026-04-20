package db

// TODO: migrate to komodo-forge-sdk-go managed client once the SDK ships
// a SQLite or RDS package. Until then, wire the driver directly here.
//
// Suggested driver: modernc.org/sqlite (pure Go, no CGo).
// Add to go.mod:  require modernc.org/sqlite v1.34.5
// Bootstrap:      db.Open("file:stats.db?cache=shared&mode=rwc")

// Client wraps the SQLite connection pool.
type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Open() error {
	return nil
}

func (c *Client) Query() error {
	// TODO: implement
	return nil
}

func (c *Client) Transaction() error {
	// TODO: implement
	return nil
}

func (c *Client) Close() error {
	return nil
}
