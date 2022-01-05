// Code generated by optiongen. DO NOT EDIT.
// optiongen: github.com/timestee/optiongen

package optiongen

import (
	"sync/atomic"
	"unsafe"
)

type Config struct {
	OptionWithStructName bool   `xconf:"option_with_struct_name" usage:"should the option func with struct name?"`
	NewFunc              string `xconf:"new_func" usage:"new function name"`
	XConf                bool   `xconf:"xconf" usage:"should gen xconf tag?"`
	UsageTagName         string `xconf:"usage_tag_name" usage:"usage tag name"`
	EmptyCompositeNil    bool   `xconf:"empty_composite_nil" usage:"should empty slice or map to be nil default?"`
	Debug                bool   `xconf:"debug" usage:"debug will print more detail info"`
}

func (cc *Config) SetOption(opt ConfigOption) {
	_ = opt(cc)
}

func (cc *Config) ApplyOption(opts ...ConfigOption) {
	for _, opt := range opts {
		_ = opt(cc)
	}
}

func (cc *Config) GetSetOption(opt ConfigOption) ConfigOption {
	return opt(cc)
}

type ConfigOption func(cc *Config) ConfigOption

// should the option func with struct name?
func WithOptionWithStructName(v bool) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.OptionWithStructName
		cc.OptionWithStructName = v
		return WithOptionWithStructName(previous)
	}
}

// new function name
func WithNewFunc(v string) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.NewFunc
		cc.NewFunc = v
		return WithNewFunc(previous)
	}
}

// should gen xconf tag?
func WithXConf(v bool) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.XConf
		cc.XConf = v
		return WithXConf(previous)
	}
}

// usage tag name
func WithUsageTagName(v string) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.UsageTagName
		cc.UsageTagName = v
		return WithUsageTagName(previous)
	}
}

// should empty slice or map to be nil default?
func WithEmptyCompositeNil(v bool) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.EmptyCompositeNil
		cc.EmptyCompositeNil = v
		return WithEmptyCompositeNil(previous)
	}
}

// debug will print more detail info
func WithDebug(v bool) ConfigOption {
	return func(cc *Config) ConfigOption {
		previous := cc.Debug
		cc.Debug = v
		return WithDebug(previous)
	}
}

func NewTestConfig(opts ...ConfigOption) *Config {
	cc := newDefaultConfig()

	for _, opt := range opts {
		_ = opt(cc)
	}
	if watchDogConfig != nil {
		watchDogConfig(cc)
	}
	return cc
}

func InstallConfigWatchDog(dog func(cc *Config)) {
	watchDogConfig = dog
}

var watchDogConfig func(cc *Config)

func newDefaultConfig() *Config {

	cc := &Config{}

	for _, opt := range [...]ConfigOption{
		WithOptionWithStructName(false),
		WithNewFunc(""),
		WithXConf(false),
		WithUsageTagName(""),
		WithEmptyCompositeNil(false),
		WithDebug(false),
	} {
		_ = opt(cc)
	}

	return cc
}

func (cc *Config) AtomicSetFunc() func(interface{}) { return AtomicConfigSet }

var atomicConfig unsafe.Pointer

func AtomicConfigSet(update interface{}) {
	atomic.StorePointer(&atomicConfig, (unsafe.Pointer)(update.(*Config)))
}

func AtomicConfig() ConfigInterface {
	current := (*Config)(atomic.LoadPointer(&atomicConfig))
	if current == nil {
		atomic.CompareAndSwapPointer(&atomicConfig, nil, (unsafe.Pointer)(newDefaultConfig()))
		return (*Config)(atomic.LoadPointer(&atomicConfig))
	}
	return current
}

// all getter func
func (cc *Config) GetOptionWithStructName() bool { return cc.OptionWithStructName }
func (cc *Config) GetNewFunc() string            { return cc.NewFunc }
func (cc *Config) GetXConf() bool                { return cc.XConf }
func (cc *Config) GetUsageTagName() string       { return cc.UsageTagName }
func (cc *Config) GetEmptyCompositeNil() bool    { return cc.EmptyCompositeNil }
func (cc *Config) GetDebug() bool                { return cc.Debug }

// interface for Config
type ConfigInterface interface {
	GetOptionWithStructName() bool
	GetNewFunc() string
	GetXConf() bool
	GetUsageTagName() string
	GetEmptyCompositeNil() bool
	GetDebug() bool
}