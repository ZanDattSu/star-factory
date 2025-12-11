package model

type Session struct {
	UUID        string `redis:"uuid"`
	UserUUID    string `redis:"user_uuid"`
	Login       string `redis:"login"`
	Email       string `redis:"email"`
	CreatedAtNs int64  `redis:"created_at_ns"`
	UpdatedAtNs int64  `redis:"updated_at_ns"`
	ExpiresAtNs int64  `redis:"expires_at_ns"`
}
