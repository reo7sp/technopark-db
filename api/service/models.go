package service

type StatusModel struct {
	UsersCount   int64 `json:"user"`
	ForumsCount  int64 `json:"forum"`
	ThreadsCount int64 `json:"thread"`
	PostsCount   int64 `json:"post"`
}
