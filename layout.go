package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
)

// Declare conformity with Layout interface
var _ fyne.Layout = (*myLayout)(nil)

const padding float32 = -4

type myLayout struct{}

// newMyLayout returns a vertical box layout for stacking a number of child
// canvas objects or widgets top to bottom. The objects are always displayed
// at their vertical MinSize. Use a different layout if the objects are intended
// to be larger then their vertical MinSize.
func newMyLayout() fyne.Layout {
	return &myLayout{}
}

func (g *myLayout) isSpacer(obj fyne.CanvasObject) bool {
	if !obj.Visible() {
		return false // invisible spacers don't impact layout
	}

	spacer, ok := obj.(layout.SpacerObject)
	if !ok {
		return false
	}

	return spacer.ExpandVertical()
}

// Layout is called to pack all child objects into a specified size.
// For a MyLayout this will pack objects into a single column where each item
// is full width but the height is the minimum required.
// Any spacers added will pad the view, sharing the space if there are two or more.
func (g *myLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	spacers := 0
	visibleObjects := 0
	// Size taken up by visible objects
	total := float32(0)

	for _, child := range objects {
		if !child.Visible() {
			continue
		}
		if g.isSpacer(child) {
			spacers++
			continue
		}

		visibleObjects++
		total += child.MinSize().Height
	}

	// Amount of space not taken up by visible objects and inter-object padding
	extra := size.Height - total - (padding * float32(visibleObjects-1))

	// Spacers split extra space equally
	spacerSize := float32(0)
	if spacers > 0 {
		spacerSize = extra / float32(spacers)
	}

	x, y := float32(0), float32(0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if g.isSpacer(child) {
			y += spacerSize
			continue
		}
		child.Move(fyne.NewPos(x, y))

		height := child.MinSize().Height
		y += padding + height
		child.Resize(fyne.NewSize(size.Width, height))
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For a BoxLayout this is the width of the widest item and the height is
// the sum of of all children combined with padding between each.
func (g *myLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	addPadding := false
	for _, child := range objects {
		if !child.Visible() || g.isSpacer(child) {
			continue
		}

		childMin := child.MinSize()
		minSize.Width = fyne.Max(childMin.Width, minSize.Width)
		minSize.Height += childMin.Height
		if addPadding {
			minSize.Height += padding
		}
		addPadding = true
	}
	return minSize
}
