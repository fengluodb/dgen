package parser

type enumStat struct {
	name    string
	members []string
}

type messageStat struct {
	name    string
	members []messageMember
}

type messageMember struct {
	seq      int
	optional bool
	typ      interface{}
	name     string
}

type serviceStat struct {
	name    string
	members []serviceMember
}

type serviceMember struct {
	name string
	req  string
	resp string
}

type mapType struct {
	key string
	val interface{}
}

type listType struct {
	ele interface{}
}
