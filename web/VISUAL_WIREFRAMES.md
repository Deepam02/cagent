# Message Action Buttons - Visual Wireframes

## Desktop Layout Detailed Wireframes

### Standard Message Layout (Desktop)
```
┌────────────────────────────────────────────────────────────────────────────┐
│                                   Desktop View (≥1024px)                    │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                              AGENT MESSAGE                          │   │
│  │ ┌─────────────────────────────────────────────────────────────────┐ │   │
│  │ │ 🤖 Assistant    [assistant]              [Copy] [Replay]        │ │   │
│  │ │ ↑                ↑                        ↑       ↑             │ │   │
│  │ │ Icon           Badge                   32x32px  32x32px         │ │   │
│  │ │                                      6px gap                    │ │   │
│  │ │                                   Hidden by default            │ │   │
│  │ │                                   Visible on hover             │ │   │
│  │ └─────────────────────────────────────────────────────────────────┘ │   │
│  │                                                                     │   │
│  │ Here's a comprehensive response with **markdown** formatting        │   │
│  │ that includes code blocks and other rich content:                   │   │
│  │                                                                     │   │
│  │ ```javascript                                                       │   │
│  │ const example = "This is a code block example";                    │   │
│  │ console.log(example);                                               │   │
│  │ ```                                                                 │   │
│  │                                                                     │   │
│  │ - Bullet point 1                                                    │   │
│  │ - Bullet point 2                                                    │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                               USER MESSAGE                          │   │
│  │ ┌─────────────────────────────────────────────────────────────────┐ │   │
│  │ │ 👤 User    [user]                         [Copy]               │ │   │
│  │ │ ↑           ↑                             ↑                     │ │   │
│  │ │ Icon      Badge                        32x32px                  │ │   │
│  │ │                                     Only Copy button            │ │   │
│  │ │                                   (No Replay for user)          │ │   │
│  │ └─────────────────────────────────────────────────────────────────┘ │   │
│  │                                                                     │   │
│  │ What is the current weather in San Francisco?                       │   │
│  │                                                                     │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└────────────────────────────────────────────────────────────────────────────┘
```

### Hover State Animation (Desktop)
```
Before Hover:
┌─────────────────────────────────────────────────────────────┐
│ 🤖 Assistant    [assistant]                                 │ ← Actions hidden
│                                                             │
│ Message content here...                                     │
└─────────────────────────────────────────────────────────────┘

During Hover (0.2s transition):
┌─────────────────────────────────────────────────────────────┐
│ 🤖 Assistant    [assistant]           [Copy] [Replay]       │ ← Actions fade in
│                                       ^^^^^^ ^^^^^^^        │
│ Message content here...              Opacity: 0 → 1        │
└─────────────────────────────────────────────────────────────┘

Button Hover Effects:
[Copy]     →    [Copy]      (translateY: 0 → -1px, shadow appears)
 ↑               ↑
Normal        Hovered
```

## Mobile Layout Detailed Wireframes

### Mobile Message Layout (≤768px)
```
┌─────────────────────────────────────────┐
│          Mobile View (≤768px)           │
│                                         │
│ ┌─────────────────────────────────────┐ │
│ │          AGENT MESSAGE              │ │
│ │ ┌─────────────────────────────────┐ │ │
│ │ │ 🤖 Assistant [assistant] [⋯][↻] │ │ │
│ │ │ ↑             ↑         ↑  ↑   │ │ │
│ │ │ Icon        Badge    44px 44px  │ │ │
│ │ │                     8px gap     │ │ │
│ │ │                  Always visible │ │ │
│ │ └─────────────────────────────────┘ │ │
│ │                                   │ │
│ │ Here's a mobile-optimized         │ │
│ │ response with **markdown**         │ │
│ │ formatting that wraps nicely      │ │
│ │ on smaller screens.               │ │
│ │                                   │ │
│ │ ```js                             │ │
│ │ const mobile = "example";         │ │
│ │ ```                               │ │
│ │                                   │ │
│ └─────────────────────────────────────┘ │
│                                         │
│ ┌─────────────────────────────────────┐ │
│ │           USER MESSAGE              │ │
│ │ ┌─────────────────────────────────┐ │ │
│ │ │ 👤 User [user]            [⋯]  │ │ │
│ │ │ ↑        ↑                ↑    │ │ │
│ │ │ Icon   Badge           44px     │ │ │
│ │ │                    Only Copy    │ │ │
│ │ └─────────────────────────────────┘ │ │
│ │                                   │ │
│ │ What's the weather?               │ │
│ │                                   │ │
│ └─────────────────────────────────────┘ │
│                                         │
└─────────────────────────────────────────┘
```

### Touch Target Specifications (Mobile)
```
Button Touch Target Analysis:

┌─────────────────┐ ┌─────────────────┐
│       [⋯]       │ │       [↻]       │
│    44px x 44px  │ │    44px x 44px  │
│                 │ │                 │
│ ┌─────────────┐ │ │ ┌─────────────┐ │
│ │    Icon     │ │ │ │    Icon     │ │
│ │  20px x 20px│ │ │ │  20px x 20px│ │
│ │   Centered  │ │ │ │   Centered  │ │
│ └─────────────┘ │ │ └─────────────┘ │
│                 │ │                 │
│   8px margin    │ │   8px margin    │
│   all around    │ │   all around    │
│                 │ │                 │
└─────────────────┘ └─────────────────┘
    Copy Button        Replay Button
  (All messages)     (User messages only)
```

## State Visual Specifications

### Copy Button States
```
Idle State:
┌──────────┐
│    📋    │  ← Copy icon (outline)
│          │    Color: text-gray-500 (light) / text-[#d8dee9] (dark)
│          │    Background: rgba(229,231,235,0.1) / rgba(76,86,106,0.1)
└──────────┘

Hover State:
┌──────────┐
│    📋    │  ← Same icon, enhanced colors
│          │    Color: text-gray-700 (light) / text-[#eceff4] (dark)  
│          │    Background: rgba(229,231,235,0.8) / rgba(76,86,106,0.3)
│          │    Transform: translateY(-1px)
│          │    Shadow: 0 2px 4px rgba(0,0,0,0.1) / 0 2px 8px rgba(0,0,0,0.3)
└──────────┘

Success State (2 seconds):
┌──────────┐
│    ✓     │  ← Check icon (solid)
│          │    Color: text-green-500 (light) / text-[#a3be8c] (dark)
│          │    Background: rgba(34,197,94,0.1) / rgba(163,190,140,0.2)
│          │    Border: rgba(34,197,94,0.2) / rgba(163,190,140,0.3)
│          │    Animation: success-pulse (scale 1 → 1.1 → 1)
└──────────┘

Error State (2 seconds):
┌──────────┐
│    ⚠️     │  ← AlertCircle icon
│          │    Color: text-red-500 (light) / text-[#bf616a] (dark)
│          │    Background: rgba(239,68,68,0.1) / rgba(191,97,106,0.2)
│          │    Border: rgba(239,68,68,0.2) / rgba(191,97,106,0.3)
└──────────┘
```

### Replay Button States
```
Idle State (User messages only):
┌──────────┐
│    ↻     │  ← RotateCcw icon (outline)
│          │    Color: text-gray-500 (light) / text-[#d8dee9] (dark)
│          │    Background: rgba(229,231,235,0.1) / rgba(76,86,106,0.1)
└──────────┘

Hover State:
┌──────────┐
│    ↻     │  ← Same icon, enhanced colors
│          │    Color: text-gray-700 (light) / text-[#eceff4] (dark)
│          │    Background: rgba(229,231,235,0.8) / rgba(76,86,106,0.3)
│          │    Transform: translateY(-1px)
│          │    Shadow: 0 2px 4px rgba(0,0,0,0.1) / 0 2px 8px rgba(0,0,0,0.3)
└──────────┘

Loading State:
┌──────────┐
│    ⟳     │  ← Loader2 icon with spin animation
│          │    Color: text-gray-500 (light) / text-[#d8dee9] (dark)
│          │    Opacity: 0.6
│          │    Cursor: not-allowed
│          │    Animation: spin 1s linear infinite
└──────────┘

Hidden State (Agent messages):
┌──────────┐
│          │  ← Not rendered (display: none)
│   NONE   │    Only Copy button visible for agent messages
│          │
└──────────┘
```

## Responsive Breakpoint Wireframes

### Tablet View (768px - 1023px)
```
┌─────────────────────────────────────────────────────────────┐
│                    Tablet View (768px - 1023px)             │
│                                                             │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ 🤖 Assistant    [assistant]           [Copy] [Replay]   │ │
│ │ ↑                ↑                     ↑       ↑       │ │
│ │ Icon           Badge                36x36px  36x36px   │ │
│ │                                    6px gap             │ │
│ │                                 Always visible         │ │
│ └─────────────────────────────────────────────────────────┘ │
│                                                             │
│ Tablet-optimized content that provides a balance between   │
│ mobile touch-friendliness and desktop information density. │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Compact Desktop (1024px - 1279px)
```
┌───────────────────────────────────────────────────────────────────┐
│                Compact Desktop (1024px - 1279px)                  │
│                                                                   │
│ ┌───────────────────────────────────────────────────────────────┐ │
│ │ 🤖 Assistant    [assistant]              [Copy] [Replay]      │ │
│ │ ↑                ↑                        ↑       ↑           │ │
│ │ Icon           Badge                   32x32px  32x32px       │ │
│ │                                       4px gap                │ │
│ │                                   Hover to reveal            │ │
│ └───────────────────────────────────────────────────────────────┘ │
│                                                                   │
│ Standard desktop experience with hover-based action discovery.    │
│ Optimized for precision pointer input (mouse/trackpad).          │
│                                                                   │
└───────────────────────────────────────────────────────────────────┘
```

### Large Desktop (≥1280px)
```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Large Desktop (≥1280px)                              │
│                                                                             │
│ ┌─────────────────────────────────────────────────────────────────────────┐ │
│ │ 🤖 Assistant    [assistant]                    [Copy] [Replay]          │ │
│ │ ↑                ↑                              ↑       ↑               │ │
│ │ Icon           Badge                         32x32px  32x32px           │ │
│ │                                             4px gap                     │ │
│ │                                         Hover to reveal                 │ │
│ └─────────────────────────────────────────────────────────────────────────┘ │
│                                                                             │
│ Maximum information density with generous spacing for comfortable reading   │
│ and precise interaction. Actions remain subtle until needed.                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Animation Sequence Diagrams

### Copy Button Success Flow
```
User Clicks Copy Button:

Step 1: Initial Click
┌──────────┐  Click Event   ┌──────────┐
│    📋    │ ────────────→  │    📋    │
│  Idle    │                │ Active   │
└──────────┘                └──────────┘
                               ↓
Step 2: Copy to Clipboard (50ms delay)
                         ┌──────────┐
                         │    ✓     │
                         │ Success  │  ← success-pulse animation
                         └──────────┘
                               ↓
Step 3: Hold Success State (2000ms)
                         ┌──────────┐
                         │    ✓     │  ← Visible feedback
                         │ Success  │    Green colors
                         └──────────┘
                               ↓
Step 4: Return to Idle (200ms transition)
                         ┌──────────┐
                         │    📋    │
                         │  Idle    │  ← Back to original state
                         └──────────┘
```

### Replay Button Loading Flow
```
User Clicks Replay Button:

Step 1: Initial Click
┌──────────┐  Click Event   ┌──────────┐
│    ↻     │ ────────────→  │    ⟳     │
│  Idle    │                │ Loading  │  ← spin animation starts
└──────────┘                └──────────┘
                               ↓
Step 2: API Call in Progress (variable duration)
                         ┌──────────┐
                         │    ⟳     │  ← Continuous spin
                         │ Loading  │    Disabled state
                         └──────────┘    opacity: 0.6
                               ↓
Step 3: API Complete (success/error)
                         ┌──────────┐
                         │    ↻     │
                         │  Idle    │  ← Back to interactive state
                         └──────────┘
```

### Hover Discovery Animation (Desktop Only)
```
Message Container Hover Sequence:

Initial State (no hover):
┌─────────────────────────────────────────────────────────┐
│ 🤖 Assistant [assistant]                                │ ← opacity: 0 (actions hidden)
│                                                         │
│ Message content...                                      │
└─────────────────────────────────────────────────────────┘

Hover Start (mouse enters):
┌─────────────────────────────────────────────────────────┐
│ 🤖 Assistant [assistant]           [Copy] [Replay]      │ ← opacity: 0 → 1
│                                    ^^^^^^ ^^^^^^^       │   (200ms transition)
│ Message content...                                      │
└─────────────────────────────────────────────────────────┘

Full Hover State:
┌─────────────────────────────────────────────────────────┐
│ 🤖 Assistant [assistant]           [Copy] [Replay]      │ ← opacity: 1 (fully visible)
│                                                         │   Ready for interaction
│ Message content...                                      │
└─────────────────────────────────────────────────────────┘

Hover End (mouse leaves):
┌─────────────────────────────────────────────────────────┐
│ 🤖 Assistant [assistant]                                │ ← opacity: 1 → 0
│                                                         │   (200ms transition)
│ Message content...                                      │
└─────────────────────────────────────────────────────────┘
```

## Accessibility Visual Indicators

### Focus States
```
Keyboard Focus Indicators:

Copy Button Focus:
┌────────────────┐
│  ┌──────────┐  │  ← 2px blue outline
│  │    📋    │  │    offset: 2px
│  │          │  │    Color: #3b82f6 (light) / #81a1c1 (dark)
│  └──────────┘  │
└────────────────┘

Replay Button Focus:
┌────────────────┐
│  ┌──────────┐  │  ← Same outline treatment
│  │    ↻     │  │    Consistent focus styling
│  │          │  │
│  └──────────┘  │
└────────────────┘
```

### Screen Reader Announcements
```
Visual Feedback for Screen Reader States:

Copy Success Announcement:
┌─────────────────────────────────────────┐
│ "Message copied to clipboard"           │ ← aria-live="polite"
│                                         │   Announced to screen reader
│ [Not visually displayed to sighted users] │   but invisible on screen
│                                         │
└─────────────────────────────────────────┘

Copy Error Announcement:
┌─────────────────────────────────────────┐
│ "Failed to copy message"                │ ← aria-live="polite"
│                                         │   Error feedback
│ [Not visually displayed to sighted users] │   for screen reader users
│                                         │
└─────────────────────────────────────────┘

Button Label Announcements:
┌─────────────────────────────────────────┐
│ Copy Button: "Copy message content      │
│              to clipboard"              │
│                                         │
│ Replay Button: "Resend this message"    │
│                                         │
│ Loading State: "Replaying message..."   │
└─────────────────────────────────────────┘
```

This comprehensive visual wireframe specification provides exact pixel measurements, color codes, animation timings, and accessibility considerations for implementing the message action buttons in the React chat interface.