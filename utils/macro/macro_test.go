package macro

import "testing"

var ma Macro

func TestMacro_ParserDefineMacro(t *testing.T) {
	ma.SetMacroSplit(",")
	var str = []string{"MY_DEFINE_CONSTANT=1", "version=${file,`../../version`}", "name=BuilderX-Go", "${define,`defineName`,`defineValue`}", "v2=${using,`version`}"}
	ma.SetDefineContext(str)
	t.Logf(ma.GetDefine("MY_DEFINE_CONSTANT"))
	t.Logf(ma.GetDefine("version"))
	t.Logf(ma.GetDefine("defineName"))
	t.Logf(ma.GetDefine("v2"))
	t.Logf("\n")
	t.Logf(ma.ParserMacro("${using,`MY_DEFINE_CONSTANT`}"))
	t.Logf(ma.ParserMacro("${using,`version`}"))
	t.Logf(ma.ParserMacro("${using,`defineName`}"))
	t.Logf(ma.ParserMacro("${using,`v2`}"))
}

func TestMacro_IsDefineMacro(t *testing.T) {
	ma.SetMacroSplit(",")
	str := "${define,`defineName`,`defineValue`}"
	str2 := "v2=${using,`version`}"
	t.Log(ma.IsDefineMacro(str))
	t.Log(ma.IsDefineMacro(str2))
}
