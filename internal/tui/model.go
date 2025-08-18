package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"github.com/docker/cagent/pkg/chat"
	"github.com/docker/cagent/pkg/history"
	"github.com/docker/cagent/pkg/runtime"
	"github.com/docker/cagent/pkg/session"
)

// model represents the application state
type model struct {
	// TUI components
	chatViewport viewport.Model
	toolViewport viewport.Model
	textInput    textinput.Model
	spinner      spinner.Model
	renderer     *glamour.TermRenderer

	// Content state
	chatContent    string // rendered chat content
	rawChatContent string // raw markdown chat content
	toolContent    string // rendered tool content
	err            error

	// App state
	ready        bool
	showInput    bool
	width        int
	height       int
	chatHeight   int
	toolHeight   int
	userScrolled bool // Track if user has manually scrolled
	isWorking    bool // Track if LLM is actively working

	// Tool call tracking
	activeToolCalls    map[string]ToolCall
	completedToolCalls []ToolCall

	// Business logic
	rt         *runtime.Runtime
	sess       *session.Session
	responseCh chan string
	toolCh     chan any // Channel for tool events
	history    *history.History
}

// NewModel creates and initializes a new model
func NewModel(rt *runtime.Runtime, sess *session.Session) (*model, error) {
	// Initialize text input
	ti := textinput.New()
	ti.Placeholder = "Enter your message..."
	ti.Focus()
	ti.CharLimit = 0
	ti.Prompt = inputPromptStyle.Render("> ")

	hist, err := history.New()
	if err != nil {
		return nil, err
	}

	// Create viewports with smooth scrolling
	chatVp := viewport.New(0, 0)
	chatVp.MouseWheelEnabled = true
	chatVp.MouseWheelDelta = 1 // Reduced from 3 to 1 for smoother scrolling

	toolVp := viewport.New(0, 0)
	toolVp.MouseWheelEnabled = true
	toolVp.MouseWheelDelta = 1

	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500"))

	return &model{
		chatViewport:       chatVp,
		toolViewport:       toolVp,
		textInput:          ti,
		spinner:            s,
		rt:                 rt,
		sess:               sess,
		responseCh:         make(chan string, 100),
		toolCh:             make(chan any, 100),
		history:            hist,
		activeToolCalls:    make(map[string]ToolCall),
		completedToolCalls: make([]ToolCall, 0),
	}, nil
}

func (m *model) updateDimensions(width, height int) {
	m.width = width
	m.height = height

	// Calculate heights
	headerHeight := 1
	footerHeight := 3
	spacingHeight := 1 // Space between the two viewports
	availableHeight := height - headerHeight - footerHeight - spacingHeight

	// Allocate space: 70% for chat, 30% for tools (minimum 5 lines for tools including borders)
	m.chatHeight = int(float64(availableHeight) * 0.7)
	m.toolHeight = availableHeight - m.chatHeight
	if m.toolHeight < 5 {
		m.toolHeight = 5
		m.chatHeight = availableHeight - 5
	}

	// Update chat viewport
	m.chatViewport.Width = width - 2         // Account for borders
	m.chatViewport.Height = m.chatHeight - 2 // Account for borders
	m.chatViewport.Style = chatViewportStyle.
		Width(width).
		Height(m.chatHeight)

	// Update tool viewport
	m.toolViewport.Width = width - 2
	m.toolViewport.Height = m.toolHeight - 2
	m.toolViewport.Style = toolViewportStyle.
		Width(width).
		Height(m.toolHeight).
		PaddingLeft(2).
		MarginLeft(0)

	// Update text input width
	m.textInput.Width = width - 2

	// Update renderer width
	var err error
	m.renderer, err = glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width-4), // Account for borders and padding
	)
	if err != nil {
		m.err = err
	}
}

// renderChatContent renders the raw markdown chat content
func (m *model) renderChatContent() error {
	rendered, err := m.renderer.Render(m.rawChatContent)
	if err != nil {
		return err
	}
	m.chatContent = rendered
	m.chatViewport.SetContent(m.chatContent)
	return nil
}

// renderToolContent renders the tool calls content
func (m *model) renderToolContent() {
	var content strings.Builder

	// Show active tool calls
	if len(m.activeToolCalls) > 0 {
		content.WriteString(toolCallStyle.Render("🔧 Active Tool Calls:") + "\n")
		for _, toolCall := range m.activeToolCalls {
			elapsed := time.Since(toolCall.StartTime).Truncate(time.Second)
			content.WriteString(fmt.Sprintf("%s %s(%s) - %v\n",
				m.spinner.View(),
				toolCall.Name,
				truncateWithEllipsis(toolCall.Arguments, 40),
				elapsed))
		}
		content.WriteString("\n")
	}

	// Show recently completed tool calls (last 3)
	if len(m.completedToolCalls) > 0 {
		content.WriteString(toolCompletedStyle.Render("✅ Recent Completions:") + "\n")
		start := max(len(m.completedToolCalls)-3, 0)
		for i := start; i < len(m.completedToolCalls); i++ {
			toolCall := m.completedToolCalls[i]
			content.WriteString(fmt.Sprintf("✓ %s - %s\n",
				toolCall.Name,
				truncateWithEllipsis(toolCall.Response, 50)))
		}
	}

	if content.Len() == 0 {
		content.WriteString("No active tool calls")
	}

	m.toolContent = content.String()
	m.toolViewport.SetContent(m.toolContent)
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.spinner.Tick,
	)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case showInputMsg:
		m.showInput = true
		return m, nil

	case tea.WindowSizeMsg:
		if !m.ready {
			m.updateDimensions(msg.Width, msg.Height)
			m.ready = true

			return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
				return showInputMsg{}
			})
		}

		m.updateDimensions(msg.Width, msg.Height)
		if err := m.renderChatContent(); err != nil {
			m.err = err
		}
		m.renderToolContent()

	case tea.KeyMsg:
		if !m.showInput {
			return m, nil
		}

		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyUp:
			if msg.Alt {
				// Alt+Up for slow scrolling up
				m.chatViewport.ScrollUp(1)
				m.userScrolled = true
				return m, nil
			}
			m.textInput.SetValue(m.history.Previous())
			return m, nil
		case tea.KeyDown:
			if msg.Alt {
				// Alt+Down for slow scrolling down
				m.chatViewport.ScrollDown(1)
				// Check if we're at the bottom
				if m.chatViewport.AtBottom() {
					m.userScrolled = false
				}
				return m, nil
			}
			m.textInput.SetValue(m.history.Next())
			return m, nil
		case tea.KeyPgUp:
			// Page up for faster scrolling
			m.chatViewport.HalfPageUp()
			m.userScrolled = true
			return m, nil
		case tea.KeyPgDown:
			// Page down for faster scrolling
			m.chatViewport.HalfPageDown()
			// Check if we're at the bottom
			if m.chatViewport.AtBottom() {
				m.userScrolled = false
			}
			return m, nil
		case tea.KeyEnter:
			if strings.TrimSpace(m.textInput.Value()) == "" {
				return m, nil
			}
			cmd := m.handleUserInput()
			return m, cmd
		}

	case responseMsg:
		m.rawChatContent += msg.content
		if err := m.renderChatContent(); err != nil {
			m.err = err
		}
		// Only auto-scroll to bottom if user hasn't manually scrolled up
		if !m.userScrolled {
			m.chatViewport.GotoBottom()
		}
		return m, tea.Tick(time.Millisecond*10, func(t time.Time) tea.Msg {
			return readResponseMsg{}
		})

	case readResponseMsg:
		return m, readResponse(m.responseCh)

	case errorMsg:
		m.err = error(msg)
		return m, nil

	case toolCallMsg:
		m.activeToolCalls[msg.toolCall.ID] = msg.toolCall
		m.renderToolContent()
		return m, readToolEvents(m.toolCh)

	case toolCompleteMsg:
		if toolCall, exists := m.activeToolCalls[msg.id]; exists {
			// Move to completed
			toolCall.IsActive = false
			toolCall.IsCompleted = true
			toolCall.Response = msg.response
			m.completedToolCalls = append(m.completedToolCalls, toolCall)

			// Remove from active
			delete(m.activeToolCalls, msg.id)
			m.renderToolContent()
		}
		return m, readToolEvents(m.toolCh)

	case workStartMsg:
		m.isWorking = true
		return m, readToolEvents(m.toolCh)

	case workEndMsg:
		m.isWorking = false
		return m, readToolEvents(m.toolCh)

	case spinner.TickMsg:
		var spinnerCmd tea.Cmd
		m.spinner, spinnerCmd = m.spinner.Update(msg)
		if len(m.activeToolCalls) > 0 {
			m.renderToolContent()
		}
		return m, spinnerCmd
	}

	// Handle viewport updates and track user scrolling
	var chatVpCmd, toolVpCmd tea.Cmd

	// Store previous position to detect user scrolling
	prevChatY := m.chatViewport.YOffset

	m.chatViewport, chatVpCmd = m.chatViewport.Update(msg)
	m.toolViewport, toolVpCmd = m.toolViewport.Update(msg)

	// Detect if user manually scrolled (position changed via user input, not programmatically)
	if m.chatViewport.YOffset != prevChatY {
		// Check if user scrolled up from bottom
		maxScroll := len(strings.Split(m.chatContent, "\n")) - m.chatViewport.Height
		if maxScroll < 0 {
			maxScroll = 0
		}
		if m.chatViewport.YOffset < maxScroll {
			m.userScrolled = true
		} else {
			m.userScrolled = false
		}
	}

	if chatVpCmd != nil {
		cmds = append(cmds, chatVpCmd)
	}
	if toolVpCmd != nil {
		cmds = append(cmds, toolVpCmd)
	}

	// Handle textinput updates if input is shown
	if m.showInput {
		var tiCmd tea.Cmd
		m.textInput, tiCmd = m.textInput.Update(msg)
		if tiCmd != nil {
			cmds = append(cmds, tiCmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// handleUserInput processes user input and returns appropriate commands
func (m *model) handleUserInput() tea.Cmd {
	input := m.textInput.Value()
	m.textInput.Reset()

	if err := m.history.Add(input); err != nil {
		m.err = err
	}

	userMsg := fmt.Sprintf("\n\n**You**: %s\n", input)
	m.rawChatContent += userMsg
	if err := m.renderChatContent(); err != nil {
		m.err = err
	}
	// Reset scroll state and go to bottom for new user input
	m.userScrolled = false
	m.chatViewport.GotoBottom()

	m.sess.Messages = append(m.sess.Messages, session.NewAgentMessage(m.rt.CurrentAgent(), &chat.Message{
		Role:    chat.MessageRoleUser,
		Content: input,
	}))

	m.responseCh = make(chan string, 100)
	m.toolCh = make(chan any, 100)

	return tea.Batch(
		processStream(m.rt, m.sess, m.responseCh, m.toolCh),
		readResponse(m.responseCh),
		readToolEvents(m.toolCh),
	)
}

func (m *model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// Build header
	headerText := "🤖 AI Chat"
	if m.isWorking || len(m.activeToolCalls) > 0 {
		headerText += " " + m.spinner.View() + " Working..."
	}
	header := headerStyle.Render(headerText)

	// Build chat viewport
	chatView := m.chatViewport.View()

	// Build tool viewport
	toolView := m.toolViewport.View()

	// Build footer with status and input
	var footer string
	if m.err != nil {
		footer = errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	} else {
		status := statusStyle.Render("🤖 Ready\n")
		input := ""
		if m.showInput {
			input = "\n" + m.textInput.View() + "\n"
		}
		footer = footerStyle.Render(status + input)
	}

	// Combine all sections
	return appStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			chatView,
			"", // Empty line for spacing
			toolView,
			footer,
		),
	)
}
