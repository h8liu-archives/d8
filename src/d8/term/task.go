package term

type Task interface {
	Run(c Cursor)
}
