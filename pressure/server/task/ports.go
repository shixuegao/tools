package task

type Ports struct {
	signatures map[string]*info
	ports      map[int]int
}

type info struct {
	startPort int
	portCount int
}

func newPorts() *Ports {
	return &Ports{
		signatures: make(map[string]*info),
		ports:      make(map[int]int),
	}
}

func (p *Ports) addPorts(signature string, startPort, portCount int) bool {
	for i := 0; i < portCount; i++ {
		port := startPort + i
		if _, ok := p.ports[port]; ok {
			return false
		}
	}
	p.signatures[signature] = &info{startPort, portCount}
	for i := 0; i < portCount; i++ {
		port := startPort + i
		p.ports[port] = 1
	}
	return true
}

func (p *Ports) removePorts(signatrue string) {
	s := p.signatures[signatrue]
	if nil == s {
		return
	}
	for i := 0; i < s.portCount; i++ {
		port := s.startPort + i
		delete(p.ports, port)
	}
}
