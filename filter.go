package hapi

type Filters []Filter

type Filter struct {
	Field    string
	Operator FilterOperator
	Values   Values
}

func (f Filters) GetFromField(field string) (list []Filter) {
	for _, filter := range f {
		if filter.Field == field {
			list = append(list, filter)
		}
	}

	return
}

func (f Filters) GetFromFields(fields []string) (list []Filter) {
	for _, field := range fields {
		list = append(list, f.GetFromField(field)...)
	}

	return
}

func (f Filters) GetFirstFromField(field string) Filter {
	for _, filter := range f {
		if filter.Field == field {
			return filter
		}
	}

	return Filter{}
}
