package optiongen

//go:generate optiongen --option_with_struct_name=false --new_func=NewTestConfig --xconf=true --empty_composite_nil=true --usage_tag_name=usage
func ConfigOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"OptionWithStructName": false, // @MethodComment(should the option func with struct name?)
		"NewFunc":              "",    // @MethodComment(new function name)
		"XConf@xconf#xconf":    false, // @MethodComment(should gen xconf tag?)
		"UsageTagName":         "",    // @MethodComment(usage tag name)
		"EmptyCompositeNil":    false, // @MethodComment(should empty slice or map to be nil default?)
		"Debug":                false, // @MethodComment(debug will print more detail info)
	}
}