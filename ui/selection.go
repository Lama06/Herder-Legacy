package ui

import (
	"fmt"

	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/hajimehoshi/ebiten/v2"
)

type SelectionConfig[T any] struct {
	Position Position
	Text     string
	Value    T
	Values   []T
	Callback func(T)
}

type Selection[T any] struct {
	config SelectionConfig[T]
	button *Button
	value  T
}

func NewSelection[T any](herderLegacy herderlegacy.HerderLegacy, config SelectionConfig[T]) *Selection[T] {
	selection := Selection[T]{
		config: config,
		button: NewButton(ButtonConfig{
			Position: config.Position,
		}),
	}

	selection.button.SetCallback(func() {
		previousScreen := herderLegacy.CurrentScreen()

		widgets := make([]ListScreenWidget, len(config.Values))
		for i, möglichkeit := range config.Values {
			möglichkeit := möglichkeit
			widgets[i] = ListScreenButtonWidget{
				Text: fmt.Sprint(möglichkeit),
				Callback: func() {
					selection.SetValue(möglichkeit)
					herderLegacy.OpenScreen(previousScreen)
				},
			}
		}

		herderLegacy.OpenScreen(NewListScreen(herderLegacy, ListScreenConfig{
			Title: config.Text,
			CancelAction: func() herderlegacy.Screen {
				return previousScreen
			},

			Widgets: widgets,
		}))
	})

	selection.SetValue(config.Value)

	return &selection
}

func (s *Selection[T]) Value() T {
	return s.value
}

func (s *Selection[T]) SetValue(value T) {
	s.value = value
	s.button.SetText(fmt.Sprintf("%v: %v", s.config.Text, value))
	if s.config.Callback != nil {
		s.config.Callback(value)
	}
}

func (s *Selection[T]) Draw(screen *ebiten.Image) {
	s.button.Draw(screen)
}

func (s *Selection[T]) Update() {
	s.button.Update()
}
