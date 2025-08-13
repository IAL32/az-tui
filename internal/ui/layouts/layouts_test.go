package layouts

import (
	"strings"
	"testing"
)

func TestLayoutSystem_Creation(t *testing.T) {
	// Test basic layout system creation
	ls := NewLayoutSystem(80, 24)
	if ls == nil {
		t.Fatal("NewLayoutSystem returned nil")
	}

	// Test dimensions
	contentW, contentH := ls.GetContentDimensions(LayoutOptions{})
	if contentW <= 0 || contentH <= 0 {
		t.Errorf("Invalid content dimensions: %dx%d", contentW, contentH)
	}
}

func TestLayoutSystem_TableLayout(t *testing.T) {
	ls := NewLayoutSystem(80, 24)

	statusContext := StatusContext{
		Mode:     "apps",
		Loading:  false,
		Counters: map[string]int{"App": 5},
	}

	helpContext := HelpContext{
		Mode:    "apps",
		ShowAll: false,
	}

	tableContent := "Sample Table Content"
	result := ls.CreateTableLayout(tableContent, statusContext, helpContext)

	if result == "" {
		t.Error("CreateTableLayout returned empty string")
	}

	if !strings.Contains(result, tableContent) {
		t.Error("Result does not contain table content")
	}
}

func TestLayoutSystem_LoadingLayout(t *testing.T) {
	ls := NewLayoutSystem(80, 24)

	statusContext := StatusContext{
		Mode:    "apps",
		Loading: true,
	}

	helpContext := HelpContext{
		Mode: "apps",
	}

	message := "Loading apps..."
	result := ls.CreateLoadingLayout(message, statusContext, helpContext)

	if result == "" {
		t.Error("CreateLoadingLayout returned empty string")
	}

	if !strings.Contains(result, message) {
		t.Error("Result does not contain loading message")
	}
}

func TestLayoutSystem_ErrorLayout(t *testing.T) {
	ls := NewLayoutSystem(80, 24)

	statusContext := StatusContext{
		Mode:  "apps",
		Error: nil,
	}

	helpContext := HelpContext{
		Mode: "apps",
	}

	errorMsg := "Connection failed"
	helpMsg := "Press r to retry"
	result := ls.CreateErrorLayout(errorMsg, helpMsg, statusContext, helpContext)

	if result == "" {
		t.Error("CreateErrorLayout returned empty string")
	}

	if !strings.Contains(result, errorMsg) {
		t.Error("Result does not contain error message")
	}
}

func TestLayoutSystem_ModalLayout(t *testing.T) {
	ls := NewLayoutSystem(80, 24)

	statusContext := StatusContext{
		Mode: "apps",
	}

	helpContext := HelpContext{
		Mode: "apps",
	}

	content := "Are you sure?"
	result := ls.CreateModalLayout(content, statusContext, helpContext)

	if result == "" {
		t.Error("CreateModalLayout returned empty string")
	}

	if !strings.Contains(result, content) {
		t.Error("Result does not contain modal content")
	}
}

func TestLayoutSystem_DimensionUpdates(t *testing.T) {
	ls := NewLayoutSystem(80, 24)

	// Test dimension updates
	ls.SetDimensions(100, 30)

	contentW, contentH := ls.GetContentDimensions(LayoutOptions{})
	if contentW <= 80 || contentH <= 24 {
		t.Error("Dimensions were not updated correctly")
	}
}

func TestThemeManager_DefaultTheme(t *testing.T) {
	theme := DefaultTheme()

	if theme.Colors.Primary == "" {
		t.Error("Default theme has empty primary color")
	}

	if theme.Colors.Success == "" {
		t.Error("Default theme has empty success color")
	}

	if theme.Colors.Error == "" {
		t.Error("Default theme has empty error color")
	}
}

func TestThemeManager_StyleAccess(t *testing.T) {
	ls := NewLayoutSystem(80, 24)
	theme := ls.GetTheme()

	// Test that we can get styles
	errorStyle := theme.GetStyle("error")
	// Just test that the style exists and can be used
	rendered := errorStyle.Render("test")
	if rendered == "" {
		t.Error("Error style failed to render content")
	}

	statusStyle := theme.GetStyle("statusBar")
	// statusStyle should exist (even if empty)
	_ = statusStyle
}

func TestTemplateManager_DefaultTemplates(t *testing.T) {
	ls := NewLayoutSystem(80, 24)

	templates := ls.ListTemplates()
	expectedTemplates := []string{"table", "loading", "error", "modal", "context"}

	for _, expected := range expectedTemplates {
		found := false
		for _, template := range templates {
			if template == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected template '%s' not found", expected)
		}
	}
}

func TestLayoutSystem_CustomTemplate(t *testing.T) {
	ls := NewLayoutSystem(80, 24)

	// Register a custom template
	customTemplate := LayoutTemplate{
		Name:        "custom",
		Description: "Custom test template",
		Options: LayoutOptions{
			CenterContent: true,
		},
		ContentFunc: func(content string, state LayoutState) string {
			return "CUSTOM: " + content
		},
	}

	ls.RegisterTemplate(customTemplate)

	// Test that we can retrieve it
	template, exists := ls.GetTemplate("custom")
	if !exists {
		t.Error("Custom template was not registered")
	}

	if template.Name != "custom" {
		t.Error("Retrieved template has wrong name")
	}

	// Test rendering with custom template
	result := ls.RenderTemplate("custom", "test content", LayoutOptions{})
	if !strings.Contains(result, "CUSTOM: test content") {
		t.Error("Custom template was not applied correctly")
	}
}
