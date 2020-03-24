package utils

import "github.com/kardianos/service"

func WinMain(name, desc string, run func()) {
	info := NewConfig(name, desc)
	prg := &Program{Main: run}
	s, _ := service.New(prg, info)
	s.Run()
}

func NewConfig(name, desc string) *service.Config {
	return &service.Config{
		Name:        name,
		DisplayName: name,
		Description: desc,
	}
}

type Program struct {
	Main func()
}

func (p *Program) Start(s service.Service) error {
	go p.Main()
	return nil
}

func (Program) Stop(s service.Service) error {
	return nil
}
