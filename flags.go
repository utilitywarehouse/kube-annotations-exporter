package main

import "strings"

type StringSliceFlag struct {
	Values []string

	changed bool
}

func (s *StringSliceFlag) String() string {
	return strings.Join(s.Values, ",")
}

func (s *StringSliceFlag) Set(value string) error {
	if !s.changed {
		s.Values = []string{}
		s.changed = true
	}
	s.Values = append(s.Values, strings.Split(value, ",")...)

	return nil
}

func (s *StringSliceFlag) StringSlice() []string {
	return s.Values
}
