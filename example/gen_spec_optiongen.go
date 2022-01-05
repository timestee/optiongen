// Code generated by optiongen. DO NOT EDIT.
// optiongen: github.com/timestee/optiongen

package example

// HTTP parsing and communication with DNS resolver was successful, and the response body content is a DNS response in either binary or JSON encoding,
// depending on the query endpoint, Accept header and GET parameters.

type Spec struct {
	// test comment 5
	// test comment 6
	TestNil1       interface{} // test comment 1
	TestBool1      bool        // test comment 2
	TestInt1       int
	TestNilFunc1   func() // 中文2
	TestReserved2_ []byte // sql.DB对外暴露出了其运行时的状态db.DBStats，sql.DB在关闭，创建，释放连接时候，会维护更新这个状态。
	// 我们可以通过prometheus来收集连接池状态，然后在grafana面板上配置指标，使指标可以动态的展示。
	TestReserved2Inner1 int
}

func (cc *Spec) SetOption(opt SpecOption) {
	_ = opt(cc)
}

func (cc *Spec) ApplyOption(opts ...SpecOption) {
	for _, opt := range opts {
		_ = opt(cc)
	}
}

func (cc *Spec) GetSetOption(opt SpecOption) SpecOption {
	return opt(cc)
}

type SpecOption func(cc *Spec) SpecOption

func WithTestNil1(v interface{}) SpecOption {
	return func(cc *Spec) SpecOption {
		previous := cc.TestNil1
		cc.TestNil1 = v
		return WithTestNil1(previous)
	}
}

func WithTestBool1(v bool) SpecOption {
	return func(cc *Spec) SpecOption {
		previous := cc.TestBool1
		cc.TestBool1 = v
		return WithTestBool1(previous)
	}
}

// 这里是函数注释3
// 这里是函数注释4
func WithTestInt1(v int) SpecOption {
	return func(cc *Spec) SpecOption {
		previous := cc.TestInt1
		cc.TestInt1 = v
		return WithTestInt1(previous)
	}
}

func WithTestNilFunc1(v func()) SpecOption {
	return func(cc *Spec) SpecOption {
		previous := cc.TestNilFunc1
		cc.TestNilFunc1 = v
		return WithTestNilFunc1(previous)
	}
}

func WithTestReserved2Inner1(v int) SpecOption {
	return func(cc *Spec) SpecOption {
		previous := cc.TestReserved2Inner1
		cc.TestReserved2Inner1 = v
		return WithTestReserved2Inner1(previous)
	}
}

func NewSpec(opts ...SpecOption) *Spec {
	cc := newDefaultSpec()

	for _, opt := range opts {
		_ = opt(cc)
	}
	if watchDogSpec != nil {
		watchDogSpec(cc)
	}
	return cc
}

func InstallSpecWatchDog(dog func(cc *Spec)) {
	watchDogSpec = dog
}

var watchDogSpec func(cc *Spec)

func newDefaultSpec() *Spec {

	cc := &Spec{
		TestReserved2_: nil,
	}

	for _, opt := range [...]SpecOption{
		WithTestNil1(nil),
		WithTestBool1(false),
		WithTestInt1(32),
		WithTestNilFunc1(nil),
		WithTestReserved2Inner1(1),
	} {
		_ = opt(cc)
	}

	return cc
}

// all getter func
func (cc *Spec) GetTestNil1() interface{}    { return cc.TestNil1 }
func (cc *Spec) GetTestBool1() bool          { return cc.TestBool1 }
func (cc *Spec) GetTestInt1() int            { return cc.TestInt1 }
func (cc *Spec) GetTestNilFunc1() func()     { return cc.TestNilFunc1 }
func (cc *Spec) GetTestReserved2_() []byte   { return cc.TestReserved2_ }
func (cc *Spec) GetTestReserved2Inner1() int { return cc.TestReserved2Inner1 }

// interface for Spec
type SpecVisitor interface {
	GetTestNil1() interface{}
	GetTestBool1() bool
	GetTestInt1() int
	GetTestNilFunc1() func()
	GetTestReserved2_() []byte
	GetTestReserved2Inner1() int
}
