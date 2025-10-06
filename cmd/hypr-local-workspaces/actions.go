package main

func NewAction(hyprctl hyprctl, dispatcher dispatcher) *Action {
	return &Action{
		hyprctl:    hyprctl,
		dispatcher: dispatcher,
	}
}

func (a *Action) GoToWorkspace(targetIndex int) error {
	return nil
}

func (a *Action) MoveToWorkspace(targetIndex int, all bool) error {
	return nil
}

func (a *Action) CycleWorkspace(dir string) error {
	return nil
}

func (a *Action) InitWorkspaces() error {
	return nil
}
