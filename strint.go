package main

import "strconv"

type strint int64

func (s *strint) UnmarshalJSON(in []byte) error {
	data := in
	if in[0] == '"' {
		data = in[1 : len(in)-1]
	}

	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*s = strint(i)

	return nil
}

func (s *strint) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(*s), 64)), nil
}

func (s strint) String() string {
	return strconv.FormatInt(int64(s), 10)
}
