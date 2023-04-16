package ads

type Ad struct {
	ID        int64
	Title     string `validate:"len:100"`
	Text      string `validate:"len:500"`
	AuthorID  int64
	Published bool
}
