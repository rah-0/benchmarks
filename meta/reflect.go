package meta

type Person struct {
	Name string
}

func (p *Person) SetName(name string) {
	p.Name = name
}

func (p *Person) GetName() string {
	return p.Name
}
