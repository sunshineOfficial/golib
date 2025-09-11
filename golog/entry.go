package golog

type LogEntry interface {
	WithTags(tags ...Tag) LogEntry
	Write()
}

type Entry struct {
	writeFunc func(Entry)
	message   string
	tags      []Tag
	level     MessageLevel
}

func (l Entry) WithTags(tags ...Tag) LogEntry {
	l.tags = append(l.tags, tags...)
	return l
}

func (l Entry) Write() {
	l.writeFunc(l)
}
