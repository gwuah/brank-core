package queue

type JobIdentifier string

func (i JobIdentifier) String() string {
	return string(i)
}
