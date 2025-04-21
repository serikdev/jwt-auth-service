package entities

import "time"

type Session struct {
	ID          string
	UserGUID    string
	RefreshHash string
	JTI         string
	IPAddress   string
	ExpiresAt   time.Time
	Used        bool
	CreatedAt   time.Time
}
