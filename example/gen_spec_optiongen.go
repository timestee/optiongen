// Code generated by optiongen. DO NOT EDIT.
// optiongen: github.com/timestee/optiongen

package example

// HTTP parsing and communication with DNS resolver was successful, and the response body content is a DNS response in either binary or JSON encoding,
// depending on the query endpoint, Accept header and GET parameters.

// Spec struct
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

// SetOption apply single option
func (cc *Spec) SetOption(opt SpecOption) {
	_ = opt(cc)
}

// ApplyOption apply mutiple options
func (cc *Spec) ApplyOption(opts ...SpecOption) {
	for _, opt := range opts {
		_ = opt(cc)
	}
}

// GetSetOption apply new option and return the old optuon
// sample:
// old := cc.GetSetOption(WithTimeout(time.Second))
// defer cc.SetOption(old)
func (cc *Spec) GetSetOption(opt SpecOption) SpecOption {
	return opt(cc)
}

// SpecOption option func
type SpecOption func(cc *Spec) SpecOption

// WithTestNil1 option func for TestNil1
func WithTestNil1(v interface{}) SpecOption {
	return func(cc *Spec) SpecOption {
		previous := cc.TestNil1
		cc.TestNil1 = v
		return WithTestNil1(previous)
	}
}

// WithTestBool1 option func for TestBool1
func WithTestBool1(v bool) SpecOption {
	return func(cc *Spec) SpecOption {
		previous := cc.TestBool1
		cc.TestBool1 = v
		return WithTestBool1(previous)
	}
}

// 这里是函数注释3
// 这里是函数注释4
// WithTestInt1 option func for TestInt1
func WithTestInt1(v int) SpecOption {
	return func(cc *Spec) SpecOption {
		previous := cc.TestInt1
		cc.TestInt1 = v
		return WithTestInt1(previous)
	}
}

// WithTestNilFunc1 option func for TestNilFunc1
func WithTestNilFunc1(v func()) SpecOption {
	return func(cc *Spec) SpecOption {
		previous := cc.TestNilFunc1
		cc.TestNilFunc1 = v
		return WithTestNilFunc1(previous)
	}
}

// WithTestReserved2Inner1 option func for TestReserved2Inner1
func WithTestReserved2Inner1(v int) SpecOption {
	return func(cc *Spec) SpecOption {
		previous := cc.TestReserved2Inner1
		cc.TestReserved2Inner1 = v
		return WithTestReserved2Inner1(previous)
	}
}

// NewSpec(opts... SpecOption) new Spec
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

// InstallSpecWatchDog the installed func will called when NewSpec(opts... SpecOption)  called
func InstallSpecWatchDog(dog func(cc *Spec)) {
	watchDogSpec = dog
}

// watchDogSpec global watch dog
var watchDogSpec func(cc *Spec)

// newDefaultSpec new default Spec
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
// GetTestNil1 return TestNil1
func (cc *Spec) GetTestNil1() interface{} { return cc.TestNil1 }

// GetTestBool1 return TestBool1
func (cc *Spec) GetTestBool1() bool { return cc.TestBool1 }

// GetTestInt1 return TestInt1
func (cc *Spec) GetTestInt1() int { return cc.TestInt1 }

// GetTestNilFunc1 return TestNilFunc1
func (cc *Spec) GetTestNilFunc1() func() { return cc.TestNilFunc1 }

// GetTestReserved2_ return TestReserved2_
func (cc *Spec) GetTestReserved2_() []byte { return cc.TestReserved2_ }

// GetTestReserved2Inner1 return TestReserved2Inner1
func (cc *Spec) GetTestReserved2Inner1() int { return cc.TestReserved2Inner1 }

// SpecVisitor visitor interface for Spec
type SpecVisitor interface {
	GetTestNil1() interface{}
	GetTestBool1() bool
	GetTestInt1() int
	GetTestNilFunc1() func()
	GetTestReserved2_() []byte
	GetTestReserved2Inner1() int
}
