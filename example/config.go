package example

import (
	"log"
	"time"
)

type SubTest struct {
	HTTPAddress string         `xconf:"http_address"`
	MapNotLeaf  map[string]int `xconf:"map_not_leaf,notleaf"`
	Map2        map[string]int `xconf:"map2"`
	Map3        map[string]int `xconf:"map3"`
	Slice2      []int64        `xconf:"slice2"`
}

// Google Public DNS provides two distinct DoH APIs at these endpoints
// Using the GET method can reduce latency, as it is cached more effectively.
// RFC 8484 GET requests must have a ?dns= query parameter with a Base64Url encoded DNS message. The GET method is the only method supported for the JSON API.

//go:generate optiongen --option_with_struct_name=false --new_func=NewFuncNameSpecified --xconf=true --usage_tag_name=usage --new_func_return=interface
func ConfigOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		// test comment 1
		// annotation@TestNil(option="WithTTTTTTTT")
		"TestNil":           nil, // test comment 3
		"TestInt":           32,  // @MethodComment(这里是函数注释1,"test") @MethodComment(这里是函数注释2)
		"TestInt64":         int64(32),
		"TestSliceInt":      []int{1, 2, 3},
		"TestSliceInt64":    []int64{1, 2, 3},
		"TestSliceString":   []string{"test1", "test2"},
		"TestSliceBool":     []bool{false, true},
		"TestSliceIntNil":   []int(nil),
		"TestSliceByte":     []byte(nil),
		"TestSliceIntEmpty": []int{},
		"TestHTTPPort":      "",

		"TestEmptyMap":        map[int]int{},
		"TestMapIntInt":       map[int]int{1: 1, 2: 2, 3: 3},
		"TestMapIntString":    map[int]string{1: "test"},
		"TestMapStringInt":    map[string]int{"test": 1},
		"TestMapStringString": map[string]string{"test": "test"},

		"TestString": "Meow",
		"Food":       (*string)(nil),
		"Walk": func() {
			log.Println("Walking")
		},
		"TestNilFunc": (func())(nil), // 中文1
		// annotation@TestParamterBool(arg=1)
		"TestParamterBool": false, // reserved parameter 1
		// annotation@TestParamterStr(arg=22)
		"TestParamterStr": "", // reserved parameter 2
		// annotation@TestProtected(private="true")
		"TestProtected": []byte(nil),
		// annotation@fOO(inline="true")
		"fOO": (*FOO)(nil),
		// annotation@Paths(inline="true")
		"Paths":   Paths(Paths{}),
		"SubTest": (*SubTest)(&SubTest{}),
		"SpecSub": (*spec)(NewSpec()), // annotation@SpecSub(getter="SpecVisitor")
	}
}

type FOO struct {
	Usernames_Passwords map[string]string
}
type Paths struct {
	Path1 string
	Path2 string
}

// HTTP parsing and communication with DNS resolver was successful, and the response body content is a DNS response in either binary or JSON encoding,
// depending on the query endpoint, Accept header and GET parameters.

//go:generate optiongen --option_prefix=WithServer --option_return_previous=false --xconf=true
func specOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		// test comment 5
		// test comment 6
		// annotation@TestNil1(comment="method commnet", private="true", xconf="test_nil1",tag_json=",omitempty")
		"TestNil1":  nil,   // test comment 1
		"TestBool1": false, // test comment 2
		"TestInt1":  32,    // @MethodComment(这里是函数注释3) @MethodComment(这里是函数注释4)

		"TestNilFunc1":   (func())(nil), // 中文2
		"TestReserved2_": []byte(nil),   // sql.DB对外暴露出了其运行时的状态db.DBStats，sql.DB在关闭，创建，释放连接时候，会维护更新这个状态。
		// 我们可以通过prometheus来收集连接池状态，然后在grafana面板上配置指标，使指标可以动态的展示。
		"TestReserved2Inner1": 1,
	}
}

//go:generate optiongen --option_return_previous=true
func OnlyAppendOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"Address": []string{"10.0.0.1:6379", "10.0.0.2:6379"},
	}
}

type Timeouts struct {
	ReadTimeout  time.Duration `xconf:"read_timeout" default:"5s"`
	WriteTimeout time.Duration `xconf:"write_timeout" default:"10s"`
	ConnTimeout  time.Duration `xconf:"conn_timeout" default:"20s"`
}

//go:generate optiongen --option_with_struct_name=true  --usage_tag_name=usage --option_return_previous=false
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
}

//go:generate optiongen --option_with_struct_name=false  --xconf=true --empty_composite_nil=true --usage_tag_name=usage --xconf=true
func AllConfigOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"TypeBool":     false,
		"TypeString":   "a",
		"TypeDuration": time.Duration(time.Second),

		"TypeFloat32": float32(32.32),
		"TypeFloat64": float32(64.64),

		"TypeInt":    32,
		"TypeUint":   32,
		"TypeInt8":   int8(8),
		"TypeUint8":  uint8(8),
		"TypeInt16":  int16(16),
		"TypeUint16": uint16(16),
		"TypeInt32":  int32(32),
		"TypeUint32": uint32(32),
		"TypeInt64":  int64(64),
		"TypeUint64": uint64(64),

		"TypeSliceInt":      []int{1, 2, 3, 4},
		"TypeSliceUint":     []uint{1, 2, 3, 4},
		"TypeSliceInt8":     []int8{1, 2, 3, 4},
		"TypeSliceUint8":    []uint8{1, 2, 3, 4},
		"TypeSliceInt16":    []int16{1, 2, 3, 4},
		"TypeSliceUin16":    []uint16{1, 2, 3, 4},
		"TypeSliceInt32":    []int32{1, 2, 3, 4},
		"TypeSliceUint32":   []uint32{1, 2, 3, 4},
		"TypeSliceInt64":    []int64{1, 2, 3, 4},
		"TypeSliceUint64":   []uint64{1, 2, 3, 4},
		"TypeSliceString":   []string{"a", "b", "c"},
		"TypeSliceFloat32":  []float32{1.32, 2.32, 3.32, 4.32},
		"TypeSliceFloat64":  []float64{1.64, 2.64, 3.64, 4.64},
		"TypeSliceDuratuon": []time.Duration([]time.Duration{time.Second, time.Minute, time.Hour}),
		// annotation@TypeMapStringIntNotLeaf(xconf="type_map_string_int_not_leaf,notleaf")
		"TypeMapStringIntNotLeaf": map[string]int{"a": 1, "b": 2},
		"TypeMapStringInt":        map[string]int{"a": 1, "b": 2},
		"TypeMapIntString":        map[int]string{1: "a", 2: "b"},
		"TypeMapStringString":     map[string]string{"a": "a", "b": "b"},
		"TypeMapIntInt":           map[int]int{1: 1, 2: 2},
		"TypeMapStringDuration":   map[string]time.Duration(map[string]time.Duration{"read": time.Second, "write": time.Second * 5}),
		// annotation@Redis(getter="RedisVisitor")
		"Redis":         (*Redis)(NewRedis()),
		"ETCD":          (*ETCD)(NewETCD(time.Second)),
		"TestInterface": (interface{})(nil),
	}
}

type WatchError = func(loaderName string, confPath string, watchErr error)

//go:generate optiongen --option_with_struct_name=true --xconf=true --usage_tag_name=usage --xconf=true
func RedisOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"Endpoints": []string{"192.168.0.1", "192.168.0.2"},
		// annotation@Address(slice_only_append="true")
		"Address":        []string{"10.0.0.1:6379", "10.0.0.2:6379"},
		"Cluster":        true,
		"TimeoutsStruct": (Timeouts)(Timeouts{}),
	}
}

var optionUsage = `
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
}
`

//go:generate optiongen --option_with_struct_name=true --debug=false --xconf=true --usage_tag_name=usage
func XXXXXXOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"OptionUsage":      string(optionUsage),
		"Endpoints":        []string{"10.0.0.1", "10.0.0.2"},
		"ReadTimeout":      time.Duration(time.Second),
		"TypeMapIntString": map[int]string{1: "a", 2: "b"},
		"TypeSliceInt64":   []int64{1, 2, 3, 4},
		"TypeBool":         false,
		"MapRedis":         (map[string]*Redis)(map[string]*Redis{"test": NewRedis()}),
		// annotation@Redis(getter="RedisVisitor",deprecated="use MapRedis intead")
		"Redis":              (*Redis)(NewRedis()), // 辅助指定类型为*Redis
		"OnWatchError":       WatchError(nil),      // 辅助指定类型为WatchError
		"OnWatchErrorNotNil": func(loaderName string, confPath string, watchErr error) {},
		"TypeSliceDuratuon":  []time.Duration([]time.Duration{time.Second, time.Minute, time.Hour}), // 辅助指定类型为WatchError
	}
}
