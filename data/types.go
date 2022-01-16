package data

//go:generate go run gen.go
//go:generate go run genplans.go

import (
	"encoding/json"
	"fmt"
	"time"
)

type Goal struct {
	Amount      int
	Repetitions int
	Duration    time.Duration
}

func (g *Goal) UnmarshalJSON(b []byte) error {
	parts := []json.RawMessage{}

	//fmt.Printf("wag\n")
	if err := json.Unmarshal(b, &parts); err != nil {
		return err
	}

	if len(parts) != 2 { //|| len(parts) > 3 {
		return fmt.Errorf("need 2 parts: %v", parts)
	}

	//fmt.Printf("hmpf\n")
	if err := json.Unmarshal(parts[0], &g.Amount); err != nil {
		return err
	}

	l := len(parts[1])
	if parts[1][0] == '"' && parts[1][l-1] == '"' {
		// try duration first
		if dur, err := time.ParseDuration(string(parts[1][1 : l-1])); err != nil {
			return fmt.Errorf("gna: %v (%v)", err, string(parts[1]))
		} else {
			g.Duration = dur
		}
		return nil
	}

	if err := json.Unmarshal(parts[1], &g.Repetitions); err != nil {
		return fmt.Errorf("wtf: %v (%v)", err, string(parts[1]))
	}

	return nil
}

type Step struct {
	Name  string
	Goals []Goal
}

type Exercise struct {
	Name  string
	Steps []Step
}

/*

 */
type Workout struct {
	Name string `json:"name"`
	Sets []int  `json:"sets"`
}

/*
type Week struct {
	Monday []Workout
	Tuesday []Workout
	Wednesday []Workout
	Thursday []Workout
	Friday []Workout
	Saturday []Workout
	Sunday []Workout
}
*/

type Days map[string][]Workout

type Plan struct {
	Name string
	//Days map[string][]Workout `json:"days"`
	Days Days `json:"days"`
}
