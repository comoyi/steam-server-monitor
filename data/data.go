package data

type Data struct {
	Counter   int
	ChCounter chan struct{}
}

func New() *Data {
	data := &Data{}
	return data
}

func (d *Data) Init() error {
	d.ChCounter = make(chan struct{})
	return nil
}
