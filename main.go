package main

import (
	//"log"
	"fmt"
	"strings"
	"time"

	"gobodyworkplans/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var opts = []string{
	"New Blood",
	"Good Behaviour",
	"Veterano",
	"Solitary Confinement",
	"Supermax",
}

var weekdayStrs = []string{
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
	"Sunday",
}

var exerciseStrs = []string{
	"Pushups",
	"Pullups",
	"Leg Raises",
	"Squats",
	"Bridges",
	"Handstand Pushups",
}

var userexe string

func DrawDays(days *widget.Accordion, opt string, t time.Time, toggleFn func()) {
	weekday := t.Weekday().String()

	plan, found := data.Plans[opt]
	if !found {
		fyne.LogError(fmt.Sprintf("failed to find plan: %q", opt), nil)
		opt = opts[0]
		plan = data.Plans[opt]
	}

	days.Items = days.Items[:0]

	for n, day := range weekdayStrs {
		title := day
		detail := container.NewVBox()
		if len(plan.Days[day]) > 0 {
			for _, workout := range plan.Days[day] {
				lt := fmt.Sprintf("   •    %s", workout.Name)
				if len(workout.Sets) == 1 {
					s := ""
					if workout.Sets[0] > 1 {
						s = "s"
					}
					lt += fmt.Sprintf(" (%d set%s)", workout.Sets[0], s)
				} else if len(workout.Sets) == 2 {
					lt += fmt.Sprintf(" (%d-%d sets)", workout.Sets[0], workout.Sets[1])
				}

				// exercise view on tap
				name := workout.Name
				wod := NewTapLabel(lt, func(_ *TapLabel) {
					// make sure exercise exists
					_, exists := data.Exercises[name]
					if !exists {
						return
					}
					// keep change
					userexe = name
					// and switch to exercise view
					toggleFn()
				})

				detail.Add(wod)
			}
		} else {
			title += ": Off"

			wod := widget.NewLabel("   •    Rest")

			detail.Add(wod)
		}

		if day == weekday {
			title += " (TODAY)"
		}

		daycon := widget.NewAccordionItem(title, detail)
		days.Append(daycon)

		if day == weekday {
			days.Open(n)
		}
	}
}

func prettyDuration(d time.Duration) string {
	return strings.ReplaceAll(d.String(), "m0s", "m")
}

func main() {
	optsidx := make(map[string]int)
	for n, opt := range opts {
		optsidx[opt] = n
	}

	var toggleFn func()

	a := app.NewWithID("io.github.sdassow.bodyworkplans")
	a.SetIcon(resourceFaviconSvg)
	w := a.NewWindow("BodyWorkPlans")

	// selected workout plan
	useropt := a.Preferences().StringWithFallback("WorkoutPlan", opts[0])

	// selected exercise
	userexe = a.Preferences().StringWithFallback("WorkoutExercise", exerciseStrs[0])

	// each day of the week to click on
	days := widget.NewAccordion()

	// dropdown to choose from plans
	choices := widget.NewSelect(opts, func(opt string) {
		useropt = opt

		DrawDays(days, opt, time.Now(), toggleFn)

		a.Preferences().SetString("WorkoutPlan", opt)
	})

	// selected step for an exercise
	userstep := -1

	var exerciseList *widget.List

	// dropdown to choose from exercises
	exerciseDropdown := widget.NewSelect(exerciseStrs, func(opt string) {
		userexe = opt

		// load step for the given exercise
		userstep = a.Preferences().IntWithFallback(fmt.Sprintf("WorkoutStep:%s", userexe), -1)

		// clear any previously selected list item
		exerciseList.UnselectAll()
		// select/scroll to correct list item
		if userstep > -1 {
			exerciseList.Select(userstep)
		} else {
			exerciseList.ScrollToTop()
		}

		a.Preferences().SetString("WorkoutExercise", opt)
	})

	// list of exercises (step 1-10)
	exerciseList = widget.NewList(
		func() int {
			return 10
		},
		func() fyne.CanvasObject {
			number := widget.NewLabel("#")

			easy := widget.NewLabel("easy")
			normal := widget.NewLabelWithStyle("normal", fyne.TextAlignCenter, fyne.TextStyle{})
			hard := widget.NewLabelWithStyle("hard", fyne.TextAlignTrailing, fyne.TextStyle{})
			levels := container.New(
				layout.NewBorderLayout(nil, nil, easy, hard),
				easy,
				normal,
				hard,
			)

			body := container.New(
				NewVBoxLayout(),
				widget.NewLabel("name"),
				levels,
			)

			return container.New(
				layout.NewBorderLayout(nil, nil, number, nil),
				number,
				body,
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			con := obj.(*fyne.Container)

			number := con.Objects[0].(*widget.Label)
			body := con.Objects[1].(*fyne.Container)

			name := body.Objects[0].(*widget.Label)
			levels := body.Objects[1].(*fyne.Container)

			exe, found := data.Exercises[userexe]
			if !found {
				return
			}

			steps := exe.Steps

			number.SetText(fmt.Sprintf("%d", id+1))
			name.SetText(steps[id].Name)

			// easy, normal, hard
			for m, goal := range steps[id].Goals {
				level := levels.Objects[m].(*widget.Label)

				if goal.Repetitions > 0 {
					level.SetText(fmt.Sprintf("%d×%d", goal.Amount, goal.Repetitions))
				} else {
					level.SetText(fmt.Sprintf("%d×%s", goal.Amount, prettyDuration(goal.Duration)))
				}
			}
		},
	)
	exerciseList.OnSelected = func(id widget.ListItemID) {
		a.Preferences().SetInt(fmt.Sprintf("WorkoutStep:%s", userexe), id)
	}

	exerciseView := container.New(
		layout.NewBorderLayout(nil, nil, nil, nil),
		exerciseList,
	)

	var toggle *widget.Button
	var menu *fyne.Container
	var content *fyne.Container

	mode := "plan"
	toggleFn = func() {
		if mode == "plan" {
			mode = "info"
			toggle.Icon = theme.HomeIcon()

			menu.Objects[1] = exerciseDropdown
			content.Objects[1] = exerciseView
			exerciseDropdown.SetSelected(userexe)
			exerciseDropdown.Refresh()
		} else {
			mode = "plan"
			toggle.Icon = theme.InfoIcon()

			menu.Objects[1] = choices
			content.Objects[1] = days
		}

		content.Refresh()
	}

	toggle = widget.NewButtonWithIcon("", theme.InfoIcon(), toggleFn)
	menu = container.New(
		layout.NewBorderLayout(nil, nil, toggle, nil),
		toggle,
		choices,
	)

	content = container.New(
		layout.NewBorderLayout(menu, nil, nil, nil),
		menu,
		days,
	)

	w.SetContent(content)

	if !a.Driver().Device().IsMobile() {
		w.Resize(fyne.Size{200, 450})
	}

	// keep last selected plan
	choices.SetSelected(useropt)

	DrawDays(days, useropt, time.Now(), toggleFn)

	w.ShowAndRun()
}
