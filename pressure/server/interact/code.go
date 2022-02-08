package interact

var Success = &code{0, "操作成功"}
var Failure = &code{1, "操作失败"}
var Loginless = &code{2, "没有登陆"}

type code struct {
	value    int
	describe string
}

func (c *code) Value() int {
	return c.value
}

func (c *code) Describe() string {
	return c.describe
}

type DevNumber struct {
	DevType    string
	SubDevType string
	Number     string
}
