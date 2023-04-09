package ads

type Ad struct {
	ID        int64
	Title     string `validate:"len:1-100"`
	Text      string `validate:"len:1-500"`
	AuthorID  int64
	Published bool
}
