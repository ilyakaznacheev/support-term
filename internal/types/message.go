package types

// Question consists of question by some user
type Question struct {
	ID       int64
	UserName string
	Text     string
}

// Answer consist or answer to some question
type Answer struct {
	ID      int64
	SupName string
	Text    string
}
