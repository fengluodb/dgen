package parser

type EnumStat struct {
	Name    string
	Members []string
}

type MessageStat struct {
	Name    string
	Members []*MessageMember
}

type MessageMember struct {
	Seq      int
	Optional bool
	Type     interface{}
	Name     string
}

type ServiceStat struct {
	Name    string
	Members []ServiceMember
}

type ServiceMember struct {
	Name string
	Req  string
	Resp string
}

type MapType struct {
	Key string
	Val interface{}
}

type ListType struct {
	Ele interface{}
}
