package args

import "reflect"

var factories map[reflect.Kind]optionFactory

func init() {
	factories = map[reflect.Kind]optionFactory{
		reflect.Bool: factoryFn(booleanOptionFactory),
		reflect.Int: singleValuedOptionFactoryGenerator(
			reflect.TypeOf(IntOption(0)),
			func(opt *reflect.Value, val interface{}) {
				opt.SetInt(int64(val.(int)))
			}),
		reflect.String: singleValuedOptionFactoryGenerator(
			reflect.TypeOf(StringOption("")),
			func(opt *reflect.Value, val interface{}) {
				opt.SetString(val.(string))
			}),
		reflect.Slice: factoryFn(multiOptionsFactory),
	}
}

type optionFactory interface {
	create(value interface{}) (interface{}, error)
}

type factoryFn func(interface{}) (interface{}, error)

func (cf factoryFn) create(value interface{}) (interface{}, error) {
	return cf(value)
}

func singleValuedOptionFactoryGenerator(typ reflect.Type, setter func(opt *reflect.Value, val interface{})) factoryFn {
	optVal := reflect.New(typ).Elem()
	f := func(v interface{}) (interface{}, error) {
		setter(&optVal, v)
		return optVal.Interface(), nil
	}
	return factoryFn(f)
}

func booleanOptionFactory(value interface{}) (interface{}, error) {
	return BooleanOption(value.(bool)), nil
}

func multiOptionsFactory(value interface{}) (interface{}, error) {
	supportedFlagsOptions := map[reflect.Kind]interface{}{
		reflect.Bool:   BooleanOption(false),
		reflect.Int:    IntOption(0),
		reflect.String: StringOption(""),
	}

	vals := value.([]interface{})
	for _, val := range vals {
		kind := reflect.ValueOf(val).Kind()
		var err error
		supportedFlagsOptions[kind], err = factories[kind].create(val)
		if err != nil {
			return nil, err
		}
	}

	return MultiOptions{
		BooleanOption: supportedFlagsOptions[reflect.Bool].(BooleanOption),
		IntOption:     supportedFlagsOptions[reflect.Int].(IntOption),
		StringOption:  supportedFlagsOptions[reflect.String].(StringOption),
	}, nil
}
