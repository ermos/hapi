package hapi

import "strconv"

type Values []Value

func (v Values) First() Value {
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

func (v Values) ToString() []string {
	var res []string
	for _, value := range v {
		res = append(res, value.String())
	}
	return res
}

type Value string

func (v Value) String() string {
	return string(v)
}

func (v Value) Int() int {
	res, err := strconv.Atoi(v.String())
	if err != nil {
		return 0
	}
	return res
}

func (v Value) Int64() int64 {
	res, err := strconv.ParseInt(v.String(), 10, 64)
	if err != nil {
		return 0
	}
	return res
}

func (v Value) Float64() float64 {
	res, err := strconv.ParseFloat(v.String(), 64)
	if err != nil {
		return 0
	}
	return res
}

func (v Value) Bool() bool {
	res, err := strconv.ParseBool(v.String())
	if err != nil {
		return false
	}
	return res
}
