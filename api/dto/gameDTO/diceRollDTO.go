package gameDTO

import "fmt"

func ErrParamIsRequired(name, typ string) error {
	return fmt.Errorf("param %s (type: %s) is required", name, typ)
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
type RollResultRequest struct {
	Total    int   `json:"total"`
	Sides    int   `json:"sides"`
	NumDices int   `json:"numDices"`
	Bonuses  []int `json:"bonuses"`
}

func (r *RollResultRequest) Validate() error {
	if r.Sides == 0 {
		return ErrParamIsRequired("sides", "int")
	}

	if r.NumDices == 0 {
		return ErrParamIsRequired("numDices", "int")
	}

	return nil
}

type RollResultResponse struct {
	Rolls      []int  `json:"rolls"`
	Bonuses    []int  `json:"bonuses"`
	SumOfBonus int    `json:"sum_of_bonuses"`
	SumOfRolls int    `json:"sum_of_rolls"`
	Total      int    `json:"total"`
	UserName   string `json:"userName"`
}
