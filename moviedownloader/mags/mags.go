package mags

type TargetItem struct {
	Name       string
	Session    int
	Episode    int
	Size       int
	FileName   string
	Link       string
	Resolution string
}

type MagFinder interface {
	ListAvailableItems() ([]*TargetItem, error)
}
