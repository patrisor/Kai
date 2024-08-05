package gui_elements

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// Custom button that can handle mouse events.
type HoldableImageButton struct {
    widget.BaseWidget
    image     *canvas.Image
    onPress   func()
    onRelease func()
    size      fyne.Size
}

// Method creates a new HoldableImageButton.
func NewHoldableImageButton(
    image fyne.Resource, 
    size fyne.Size, 
    onPress, 
    onRelease func(),
) *HoldableImageButton {
    button := &HoldableImageButton{
        image:     canvas.NewImageFromResource(image),
        onPress:   onPress,
        onRelease: onRelease,
        size:      size,
    }
    button.image.SetMinSize(size)
    button.image.FillMode = canvas.ImageFillContain
    button.ExtendBaseWidget(button)
    return button
}

// Method is called when the mouse button is pressed.
func (h *HoldableImageButton) MouseDown(*desktop.MouseEvent) {
    if h.onPress != nil {
        h.onPress()
    }
}

// Method is called when the mouse button is released.
func (h *HoldableImageButton) MouseUp(*desktop.MouseEvent) {
    if h.onRelease != nil {
        h.onRelease()
    }
}

// Method returns the renderer for this widget.
func (h *HoldableImageButton) CreateRenderer() fyne.WidgetRenderer {
    return widget.NewSimpleRenderer(h.image)
}

// Method returns the minimum size of the widget.
func (h *HoldableImageButton) MinSize() fyne.Size {
    return h.size
}

// Method resizes the widget.
func (h *HoldableImageButton) Resize(size fyne.Size) {
    h.size = size
    h.image.Resize(size)
    h.BaseWidget.Resize(size)
}