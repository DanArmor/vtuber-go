package selected

import (
	"encoding/json"
	"errors"
)

type Selected string

const (
	Yes Selected = "yes"
	No  Selected = "no"
	All Selected = "all"
)

func (Selected) Values() (kinds []string) {
	for _, s := range []Selected{Yes, No, All} {
		kinds = append(kinds, string(s))
	}
	return kinds
}

func (r Selected) String() string {
	switch r {
	case Yes:
		return "yes"
	case No:
		return "no"
	case All:
		return "all"
	default:
		return ""
	}
}

func ParseDifficulty(s string) (Selected, error) {
	switch s {
	case "yes":
		return Yes, nil
	case "no":
		return No, nil
	case "all":
		return All, nil
	default:
		return Yes, errors.New("cant parse difficulty from string")
	}
}

func (r Selected) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *Selected) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if *r, err = ParseDifficulty(s); err != nil {
		return err
	}
	return nil
}
