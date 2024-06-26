package optiongen

const templateTextWithPreviousSupport = `

{{- range $_, $comment := $.ClassComments }}
{{ $comment }}
{{- end }}

// {{ $.ClassName }} should use {{ $.ClassNewFuncName }} to initialize it
type {{ $.ClassName }} struct {
	{{- range $index, $option := $.ClassOptionInfo }}
		{{- range $_, $comment := $option.LastRowComments }}
			{{ unescaped $comment }}
 		{{- end }}
		{{- if $option.Inline }}
		{{ $option.Type }} {{unescaped $option.TagString}} {{ unescaped $option.SameRowComment }} 
		{{- else }}
		{{ $option.Name }} {{ $option.Type }} {{unescaped $option.TagString}} {{ unescaped $option.SameRowComment }} 
		{{- end }}
	{{- end }}
}

// {{ $.ClassNewFuncName }} new {{ $.ClassName }}
{{ $.ClassNewFuncSignature }} {
	cc := newDefault{{ $.ClassNameTitle }}()
	{{- range $index, $option := $.ClassOptionInfo }}
	{{- if eq $option.ArgIndex 0}}
	{{- else}}
	cc.{{$option.Name}} = {{$option.NameAsParameter}}
	{{- end}}
	{{- end }}
	for _, opt := range opts  {
		opt.Apply(cc)
	}
	if watchDog{{$.ClassNameTitle}} != nil {
		watchDog{{$.ClassNameTitle}}(cc)
	}
	return cc
}

{{- if $.OptionReturnPrevious }}
// ApplyOption apply multiple new option and return the old ones
// sample: 
// old := cc.ApplyOption(WithTimeout(time.Second))
// defer cc.ApplyOption(old...)
func (cc *{{ $.ClassName }}) ApplyOption(opts... {{$.ClassOptionTypeName }}) []{{$.ClassOptionTypeName }}{
	var previous []{{$.ClassOptionTypeName }}
	for _, opt := range opts  {
		previous = append(previous,opt.Apply(cc))
	}
	return previous
}
{{- else}}
// ApplyOption apply multiple new option
func (cc *{{ $.ClassName }}) ApplyOption(opts... {{$.ClassOptionTypeName }}){
	for _, opt := range opts  {
		opt.Apply(cc)
	}
}
{{- end }}

// {{ $.ClassOptionFunctionTypeName }} option func
{{- if $.OptionReturnPrevious }}
type {{ $.ClassOptionTypeName }} interface {
	Apply(cc *{{$.ClassName}}) {{ $.ClassOptionTypeName }}
}
var _ {{ $.ClassOptionTypeName }} = {{ $.ClassOptionFunctionTypeName }}(nil)
type {{ $.ClassOptionFunctionTypeName }} func(cc *{{$.ClassName}}) {{ $.ClassOptionFunctionTypeName }}
func (f {{ $.ClassOptionFunctionTypeName }}) Apply(cc *{{$.ClassName}}) {{ $.ClassOptionTypeName }} {
	return f(cc)
}
{{- else}}
type {{ $.ClassOptionTypeName }} interface {
	Apply(cc *{{$.ClassName}})
}
var _ {{ $.ClassOptionTypeName }} = {{ $.ClassOptionFunctionTypeName }}(nil)
type {{ $.ClassOptionFunctionTypeName }} func(cc *{{$.ClassName}})
func (f {{ $.ClassOptionFunctionTypeName }}) Apply(cc *{{$.ClassName}}) {
	f(cc)
}



{{- end }}

{{ range $index, $option := $.ClassOptionInfo }}
{{- if eq $option.GenOptionFunc true }}
	{{- if eq $option.Slice true }}
		{{- if not $option.OnlyAppend }}
			{{- if $.OptionReturnPrevious }}
{{ unescaped $option.OptionComment }}
func {{$option.OptionFuncName}}(v ...{{$option.SliceElemType}}) {{ $.ClassOptionFunctionTypeName }}   { return func(cc *{{$.ClassName}}) {{ $.ClassOptionFunctionTypeName }} {
	previous := cc.{{$option.Name}}
	cc.{{$option.Name}} = v
return {{$option.OptionFuncName}}(previous...)
} }
			{{- else }}
{{ unescaped $option.OptionComment }}
func {{$option.OptionFuncName}}(v ...{{$option.SliceElemType}}) {{ $.ClassOptionFunctionTypeName }}   { return func(cc *{{$.ClassName}})  {
	cc.{{$option.Name}} = v
} }
			{{- end }}
		{{- end }}
	{{ unescaped $option.AppendComment }}
		{{- if $.OptionReturnPrevious }}
func {{$option.AppendFuncName}}(v ...{{$option.SliceElemType}}) {{ $.ClassOptionFunctionTypeName }}   { return func(cc *{{$.ClassName}}) {{ $.ClassOptionFunctionTypeName }} {
	previous := cc.{{$option.Name}}
	cc.{{$option.Name}} = append(cc.{{$option.Name}}, v...)
			{{- $OptionFuncName := $option.OptionFuncName}}
			{{- if $option.OnlyAppend }}
				{{- $OptionFuncName = $option.AppendFuncName}}
			{{- end }}
	return {{$OptionFuncName}}(previous...)
} }
		{{- else }}
func {{$option.AppendFuncName}}(v ...{{$option.SliceElemType}}) {{ $.ClassOptionFunctionTypeName }}   { return func(cc *{{$.ClassName}})  {
	cc.{{$option.Name}} = append(cc.{{$option.Name}}, v...)
} }
		{{- end }}
	{{- else }}
	{{ unescaped $option.OptionComment }}
		{{- if $.OptionReturnPrevious }}
func {{$option.OptionFuncName}}(v {{$option.Type}}) {{ $.ClassOptionFunctionTypeName }}   { return func(cc *{{$.ClassName}}) {{ $.ClassOptionFunctionTypeName }} {	
	previous := cc.{{$option.Name}}
	cc.{{$option.Name}} = v
return {{$option.OptionFuncName}}(previous)
} }
		{{- else }}
func {{$option.OptionFuncName}}(v {{$option.Type}}) {{ $.ClassOptionFunctionTypeName }}   { return func(cc *{{$.ClassName}})  {
	cc.{{$option.Name}} = v
} }
		{{- end }}
	{{- end }}
{{- end }}

{{ end }}

// Install{{$.ClassNameTitle}}WatchDog the installed func will called when {{ $.ClassNewFuncName }}  called
func Install{{$.ClassNameTitle}}WatchDog(dog func(cc *{{$.ClassName}})) { watchDog{{$.ClassNameTitle}} = dog }

// watchDog{{$.ClassNameTitle}} global watch dog
var watchDog{{$.ClassNameTitle}} func(cc *{{$.ClassName}})

// set{{ $.ClassNameTitle }}DefaultValue default {{ $.ClassName }} value
func set{{ $.ClassNameTitle }}DefaultValue (cc *{{ $.ClassName }}) {
{{- range $index, $option := $.ClassOptionInfo }}
	{{- if eq $option.GenOptionFunc false }}
		{{- if eq $option.FieldType 0 }}
	cc.{{$option.Name}} = {{ $option.Type }} {{ $option.Body}}
		{{- else }}
	cc.{{$option.Name}} = {{ $option.Body}}
		{{- end }}
	{{- end }}
{{- end }}	
	for _, opt := range [...]{{ $.ClassOptionFunctionTypeName }} {
{{- range $index, $option := $.ClassOptionInfo }}
	{{- if eq $option.GenOptionFunc true }}
		{{- if eq $option.Slice true }}
			{{- $OptionFuncName := $option.OptionFuncName}}
			{{- if $option.OnlyAppend }}
				{{- $OptionFuncName = $option.AppendFuncName}}
			{{- end }}
			{{- if eq $option.FieldType 0 }}
				{{$OptionFuncName}}({{ $option.Type }} {{ $option.Body}}...),
			{{- else }}
				{{$OptionFuncName}}({{ $option.Body}}...),
			{{- end }}
		{{- else }}
			{{- if eq $option.FieldType 0 }}
				{{$option.OptionFuncName}}({{ $option.Type }} {{ $option.Body}}),
			{{- else }}
				{{$option.OptionFuncName}}({{ $option.Body}}),
			{{- end }}
		{{- end }}
	{{- end }}
{{- end }}
	}  {
		opt(cc)
	}
}

// newDefault{{ $.ClassNameTitle }} new default {{ $.ClassName }} 
func newDefault{{ $.ClassNameTitle }} () *{{ $.ClassName }} {
	cc := &{{ $.ClassName }}{}
	set{{ $.ClassNameTitle }}DefaultValue(cc)
	return cc
}

{{- if $.XConf }}

// AtomicSetFunc used for XConf
func (cc *{{ $.ClassName }}) AtomicSetFunc() func(interface{}) { return Atomic{{ $.ClassNameTitle }}Set }

// atomic{{ $.ClassName }} global *{{ $.ClassName }} holder
var atomic{{ $.ClassNameTitle }} unsafe.Pointer

// onAtomic{{ $.ClassNameTitle }}Set global call back when  Atomic{{ $.ClassNameTitle }}Set called by XConf.
// use {{ $.InterfaceName }}.ApplyOption to modify the updated cc
// if passed in cc not valid, then return false, cc will not set to atomic{{ $.ClassNameTitle }}
var onAtomic{{ $.ClassNameTitle }}Set func(cc {{ $.InterfaceName }}) bool

// InstallCallbackOnAtomic{{ $.ClassNameTitle }}Set install callback
func InstallCallbackOnAtomic{{ $.ClassNameTitle }}Set(callback func(cc {{ $.InterfaceName }}) bool) { onAtomic{{ $.ClassNameTitle }}Set = callback}

// Atomic{{ $.ClassNameTitle }}Set atomic setter for *{{ $.ClassName }}
func Atomic{{ $.ClassNameTitle }}Set(update interface{}) {
	cc := update.(*{{ $.ClassName }})
	if onAtomic{{ $.ClassNameTitle }}Set != nil && !onAtomic{{ $.ClassNameTitle }}Set(cc) {
		return
	}
	atomic.StorePointer(&atomic{{ $.ClassNameTitle }}, (unsafe.Pointer)(cc))
}

// Atomic{{ $.ClassNameTitle }} return atomic *{{ $.VisitorName }}
func Atomic{{ $.ClassNameTitle }}() {{ $.VisitorName }} {
	current := (*{{ $.ClassName }})(atomic.LoadPointer(&atomic{{ $.ClassNameTitle }}))
	if current == nil {
		defaultOne := newDefault{{ $.ClassNameTitle }}()
		if watchDog{{$.ClassNameTitle}} != nil {
			watchDog{{$.ClassNameTitle}}(defaultOne)
		}
		atomic.CompareAndSwapPointer(&atomic{{ $.ClassNameTitle }}, nil, (unsafe.Pointer)(defaultOne))
		return (*{{ $.ClassName }})(atomic.LoadPointer(&atomic{{ $.ClassNameTitle }}))
	}
	return current
}
{{- end}}


// all getter func
{{- range $index, $option := $.ClassOptionInfo }}{{- if eq $option.GenVisitFunc true }}{{ unescaped $option.VisitFuncComment }}
func (cc *{{ $.ClassName }}) {{$option.VisitFuncName}}() {{ $option.VisitFuncReturnType }} { return cc.{{$option.Name}} }
{{- end }}
{{- end }}

// {{ $.VisitorName }} visitor interface for {{ $.ClassName }}
type {{ $.VisitorName }} interface {
	{{- range $index, $option := $.ClassOptionInfo }}{{- if eq $option.GenVisitFunc true }}{{ unescaped $option.VisitFuncComment }}
	{{$option.VisitFuncName}}() {{ $option.VisitFuncReturnType }} 
	{{- end }}
	{{- end }}
}

// {{ $.InterfaceName }} visitor + ApplyOption interface for {{ $.ClassName }}
type {{ $.InterfaceName }} interface {
	{{ $.VisitorName }}
	{{- if $.OptionReturnPrevious }}
	ApplyOption(... {{$.ClassOptionTypeName }}) []{{$.ClassOptionTypeName }} 
	{{- else }}
	ApplyOption(... {{$.ClassOptionTypeName }})
	{{- end }}
}
`
