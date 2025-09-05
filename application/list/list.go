package list

const LikedQuestionsListName = "Liked questions"

// A list of questions
type List struct {
	// The unique identifier for the list.
	ID string `json:"id" firestore:"id" binding:"required"`
	// The human readable name for the list
	Name string `json:"name" firestore:"name" binding:"required"`
}

func (l *List) String() string {
	return "List(" +
		"Id=" + l.ID +
		", Name=" + l.Name +
		")"
}
