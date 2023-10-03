package pod

type Message interface {
	isPodMessage()
}

type AddedPodMsg struct {
	Pod Item
}

type RemovedPodMsg struct {
	Pod Item
}

type ModifiedPodMsg struct {
	Pod Item
}

type ListMsg struct {
	Pods            []Item
	ResourceVersion string
}

type NextEventMsg struct{}

type WatchingPodsMsg struct{}

type ErrorMsg struct {
	Err error
}

func (m AddedPodMsg) isPodMessage()     {}
func (m RemovedPodMsg) isPodMessage()   {}
func (m WatchingPodsMsg) isPodMessage() {}
func (m NextEventMsg) isPodMessage()    {}
func (m ErrorMsg) isPodMessage()        {}
