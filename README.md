# optiongen
[![GoDoc](https://godoc.org/github.com/timestee/optiongen?status.svg)](https://godoc.org/github.com/timestee/optiongen)
[![Go Report Card](https://goreportcard.com/badge/github.com/timestee/optiongen)](https://goreportcard.com/report/github.com/timestee/optiongen)[![Sourcegraph](https://sourcegraph.com/github.com/timestee/optiongen/-/badge.svg)](https://sourcegraph.com/github.com/timestee/optiongen?badge)



optiongen is a fork of [XSAM/optionGen](https://github.com/XSAM/optionGen), a tool to generate go Struct option for test, mock or more flexible. The purpose of this fork is to provide more powerful and flexible option generation. 

## Functional Options
Functional options are an idiomatic way of creating APIs with options on types. The initial idea for this design pattern can be found in an article published by Rob Pike called [Self-referential functions and the design of options](https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html).

## Install
Install using go install, and this will build the optionGen binary in $GOPATH/bin.
```bash
go install github.com/timestee/optiongen/cmd/optiongen@master
```

optionGen require [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports) to format code which is generated. So you may confirm that `goimports` has been installed

```bash
go install golang.org/x/tools/cmd/goimports@latest
```


## 快速开始
```golang
type WatchError = func(loaderName string, confPath string, watchErr error)

//go:generate optiongen --option_with_struct_name=true --xconf=true --usage_tag_name=usage --xconf=true
func RedisOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"Endpoints":      []string{"192.168.0.1", "192.168.0.2"},
		"Cluster":        true,
		"TimeoutsStruct": (Timeouts)(Timeouts{}),
	}
}
//go:generate optiongen
func XXXXXXOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"Endpoints":        []string{"10.0.0.1", "10.0.0.2"},
		"ReadTimeout":      time.Duration(time.Second),
		"TypeMapIntString": map[int]string{1: "a", 2: "b"},
		"TypeSliceInt64":   []int64{1, 2, 3, 4},
		"TypeBool":         false,
		"MapRedis":         (map[string]*Redis)(map[string]*Redis{"test": NewRedis()}),
		// annotation@Redis(getter="RedisVisitor")
		"Redis":              (*Redis)(NewRedis()), // 辅助指定类型为*Redis
		"OnWatchError":       WatchError(nil),      // 辅助指定类型为WatchError
		"OnWatchErrorNotNil": func(loaderName string, confPath string, watchErr error) {},
		"TypeSliceDuratuon":  []time.Duration([]time.Duration{time.Second, time.Minute, time.Hour}), // 辅助指定类型为WatchError
	}
}
```
- 命令提示：`//go:generate optiongen `
- 声明XXX`OptionDeclareWithDefault`, optiongen会识别后缀为`OptionDeclareWithDefault`的func声明做进一步处理,XXX`OptionDeclareWithDefault`的函数约定如上代码段所示。
- 在XXX`OptionDeclareWithDefault`函数返回的` map[string]interface{}`中声明字段名称,默认值
  - 基础类型如上例中的`Endpoints`,`TypeBool`等基础类型或者基础类型的slice，map等可以直接以字面值给出默认值
  - 函数类型如果非空可以直接以字面值给出默认值，如`OnWatchErrorNotNil`,但`OnWatchError`的默认值为`nil`，则需要辅助性的给出类型定义`WatchError`
  - `time.Duration`,`[]time.Duration`,`map[string]*Redis`此类的非基础类型的slice或者map都需要辅助指明类型
  > 在使用过程中可以尝试运行，如`optiongen`遇到无法处理的类型会给出错误提示，如将`TypeSliceDuratuon`默认值写为`[]time.Duration{time.Second, time.Minute, time.Hour}`会得到如下错误提示：
	`panic: optionGen "TypeSliceDuratuon" got type []time.Duration support basic types only`
- 运行`go generate`,会有如下输出
  ```shell
  🚀  optiongen running => /xxxxxx/github/optiongen/example/config.go:159 [XXXXXXOptionDeclareWithDefault] ...
  ```

## 使用帮助
### optiongen支持的参数
可以在`//go:generate optiongen`中根据具体需求加入参数调整代码生成行为，支持的参数可以运行:`optiongen --help`查看。以上文提到的`XXXXXXOptionDeclareWithDefault`为例。
- `--debug`, bool类型，默认false，是否打开调试模式，调试模式下会输出详尽的运行日志，主要用于开发调试
- `--new_func`,生成的Struct的New方法名称,如不指定则默认为默认为`NewXXXXXX`,可以指定为如`NewConf`
- `--new_func_return`,生成的Struct的New方法的返回类型，默认 `pointer`
	- `pointer`, 返回类型指针: `*XXXXXX`，比如定义类型为`confOptionDeclareWithDefault`首字符小写，生成的配置为`conf`不会被导出，此时将返回类型设定为`interface`或`visitor`更为恰当。
	- `interface`,返回类型接口: `XXXXXXInterface`
	- `visitor`,返回类型访问接口: `XXXXXXVisitor`
- `--option_prefix`,生成的Option方法前缀
    - 默认设置下，上例中的`Endpoints`字段生成的Option方法签名为`WithEndpoints`
	- 为了避免方法签名冲突，也更为明确方法的含义(一个package中定义了多个Option结构)，可以指定Option方法前缀如:`WithServer`,则生成的Option方法签名为`WithServerEndpoints`
- `--option_with_struct_name`,默认false，功能与`--option_prefix`类似，Option名称是否携带Struct名称
    - 如设定为true,在不设定`--option_prefix`的情况下,上例中的`Endpoints`生成的Option签名为：`WithXXXXXXEndpoints`
	- 设定`--option_prefix`时，该参数无效
- `--option_return_previous`， bool类型，默认true，生成的Option方法是否返回原始值
	返回原始值的Option签名为:
	```golang
	// WithReadTimeout option func for ReadTimeout
	func WithReadTimeout(v time.Duration) XXXXXXOption {
		return func(cc *XXXXXX) XXXXXXOption {
			previous := cc.ReadTimeout
			cc.ReadTimeout = v
			return WithReadTimeout(previous)
		}
	}
	// 可以在一些测试场景下将参数设定为需要的值，在退出当前场景后将配置恢复
	func TestApplyOption(xx XXXXXXInterface) {
		old := xx.ApplyOption(WithReadTimeout(time.Second))
		defer xx.ApplyOption(old...)
		// ...
	}
	```

	不返回原始值的Option签名:
	```golang
	// WithReadTimeout option func for ReadTimeout
	func WithReadTimeout(v time.Duration) XXXXXXOption {
		return func(cc *XXXXXX)  {
			cc.ReadTimeout = v
		}
	}
	```
- `--xconf`,bool类型，默认false，是否生成[XConf](https://github.com/sandwich-go/xconf)支持
	- 如设定为true，会生成`XConf`所需的TAG,更新逻辑等,此时可以将optiongen作为配置存在
	- 如有配置热更需求如对接ETCD,Apollo,文件系统等，`XConf`解析时传入`AtomicETCD`,`XConf`会自动解析配置，自动在配置更新后将最新的配置更新到`AtomicETCD`返回的指针中。	
	```golang
	//go:generate optiongen --option_prefix=WithETCD --xconf=true  --usage_tag_name=usage --option_return_previous=false
	func ETCDOptionDeclareWithDefault() interface{} {
		return map[string]interface{}{
			// annotation@Endpoints(comment="etcd地址")
			"Endpoints": []string{"10.0.0.1", "10.0.0.2"},
			// annotation@TimeoutsPointer(comment="timeout设置")
			"TimeoutsPointer": (*Timeouts)(&Timeouts{}),
			// annotation@writeTimeout(private="true",arg=1)
			"writeTimeout": time.Duration(time.Second),
			// annotation@Redis(getter="RedisVisitor")
			"Redis": (*Redis)(NewRedis()),
	}
	// ETCD should use NewETCD to initialize it
	type ETCD struct {
		// annotation@Endpoints(comment="etcd地址")
		Endpoints []string `xconf:"endpoints" usage:"etcd地址"`
		// annotation@TimeoutsPointer(comment="timeout设置")
		TimeoutsPointer *Timeouts `xconf:"timeouts_pointer" usage:"timeout设置"`
	}
	// AtomicSetFunc used for XConf
	func (cc *ETCD) AtomicSetFunc() func(interface{}) { return AtomicETCDSet }
	// InstallCallbackOnAtomicETCDSet install callback
	func InstallCallbackOnAtomicETCDSet(callback func(cc ETCDInterface) bool) { ... }
	// AtomicETCDSet atomic setter for *ETCD
	func AtomicETCDSet(update interface{}) {...}
	// AtomicETCD return atomic *ETCDVisitor
	func AtomicETCD() ETCDVisitor {...}
	```
- `--usage_tag_name`，字符串类型，生成的usage标签名称，默认空，如指定:`usage`会生生成usage信息，`XConf`会将这部分信息展现在`xonf.Usage`以及`FlagSet.Usage`中。

### optiongen支持的标注
`optiongen`支持通过标注的方式对字段级的代码生成进行更为灵活的控制，目前支持通过在注释中定位标注格式如下：
```golang
// annotation@Redis(getter="RedisVisitor") 
// annotation@ReadTimeout(private="true", xconf="read_timeout_user_define_name")
// annotation@TypeMapStringIntNotLeaf(xconf="type_map_string_int_not_leaf,notleaf")
// annotation@ReadTimeout(arg=1,tag_json=",omitempty")
```
- `private`,指定字段为私有字段，不生成Option，不会影响字段本身的访问属性，字段本身的访问属性设定通过首字符大小写决定，如上例`ETCDOptionDeclareWithDefault`的`writeTimeout`字段。
- `arg`，指定arg参数的字段不会生成Option方法，并会作为New方法的参数存在
  - 如上例指定`writeTimeout`的arg，则生成的New方法为:`NewETCD(writeTimeout time.Duration, opts ...ETCDOption) `
  - 允许设定多个arg，指定参数的index即可，index不可重复 
- `xconf`，自定义xconf标签
- `inline`,将字段inline
- `getter`,生成的Get接口返回值类型，默认为定义时指定的类型，可通过该方式指定返回类型对应的接口，如上例中`Redis`的访问接口返回为:`RedisVisitor`
- `option`, 指定该字段生成的option方法名称，覆盖`--option_prefix`和`--option_with_struct_name`规则，
- `deprecated`,字符串，指定字段为deprecated，在Option以及Get方法上都会生成`//Deprecated`注释，如果启用了xconf支持，会一并在xconf标签中生成deprecated支持。
- `tag_{name}`,其中`{name}`为tag名称，如json，例如`tag_json=",omitempty"`
