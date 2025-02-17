package internal

import (
    "bufio"
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "text/template"
)

// ANSI style constants
const (
    DefaultColorReset  = "\033[0m"
    DefaultColorGreen  = "\033[32m"
    DefaultColorYellow = "\033[33m"
    DefaultColorCyan   = "\033[36m"
    DefaultBold        = "\033[1m"
)

// Global style variables (adjustable in one place)
var (
    PromptColor      = DefaultColorGreen  // Color for input prompts.
    DescriptionColor = DefaultColorCyan   // Color for field descriptions.
    OptionColor      = DefaultColorYellow // Color for option numbers.
    BoldStyle        = DefaultBold        // Bold style.
    ColorReset       = DefaultColorReset  // Reset style.
)

// Option represents a selectable option for a form item.
type Option struct {
    Name string `json:"name"`
    Desc string `json:"desc"`
}

// Item represents a field in the commit form.
type Item struct {
    Name     string   `json:"name"`
    Desc     string   `json:"desc"`
    Form     string   `json:"form"`
    Options  []Option `json:"options,omitempty"`
    Required bool     `json:"required,omitempty"`
    // Optional default value.
    Default string `json:"default,omitempty"`
    // Optional hint to guide the user.
    Hint string `json:"hint,omitempty"`
    // Optional validation rule (for documentation purposes).
    Validation string `json:"validation,omitempty"`
}

// MessageConfig holds the form definition and commit message template.
type MessageConfig struct {
    Items    []Item `json:"items"`
    Template string `json:"template"`
}

// Config is the root configuration structure.
type Config struct {
    Message MessageConfig `json:"message"`
}

// LoadConfig loads the configuration from the file "configs/default.json"
// relative to the current working directory.
func LoadConfig() (Config, error) {
    var cfg Config
    configPath := filepath.Join("configs", "default.json")

    data, err := ioutil.ReadFile(configPath)
    if err != nil {
        return cfg, fmt.Errorf("failed to read config file %s: %v", configPath, err)
    }

    if err := json.Unmarshal(data, &cfg); err != nil {
        return cfg, fmt.Errorf("failed to parse config file %s: %v", configPath, err)
    }

    log.Printf("Loaded config from %s\n", configPath)
    return cfg, nil
}

// LoadDefaultConfig returns a built-in default configuration as a fallback.
func LoadDefaultConfig() Config {
    // A minimal built-in config (fallback) without repetition.
    defaultJSON := `{
        "message": {
            "items": [
                {
                    "name": "type",
                    "desc": "Select the type of change (required):",
                    "form": "select",
                    "options": [
                        { "name": "feat", "desc": "A new feature" },
                        { "name": "fix", "desc": "A bug fix" },
                        { "name": "docs", "desc": "Documentation only changes" },
                        { "name": "style", "desc": "Changes that do not affect the meaning of the code" },
                        { "name": "refactor", "desc": "A code change that neither fixes a bug nor adds a feature" },
                        { "name": "perf", "desc": "A code change that improves performance" },
                        { "name": "test", "desc": "Adding missing tests" },
                        { "name": "chore", "desc": "Changes to build process or auxiliary tools" },
                        { "name": "revert", "desc": "Revert to a commit" },
                        { "name": "WIP", "desc": "Work in progress" }
                    ],
                    "required": true,
                    "hint": "Choose one of the available change types."
                },
                {
                    "name": "scope",
                    "desc": "Scope (optional): Specify the area affected (e.g., users, db, poll)",
                    "form": "input",
                    "default": ""
                },
                {
                    "name": "subject",
                    "desc": "Subject (required): Concise description in imperative, lower case, no final dot",
                    "form": "input",
                    "required": true,
                    "validation": "max:100"
                },
                {
                    "name": "body",
                    "desc": "Body (optional): Detailed motivation for the change",
                    "form": "multiline",
                    "default": ""
                },
                {
                    "name": "footer",
                    "desc": "Footer (optional): Information about breaking changes or related issues",
                    "form": "multiline",
                    "default": ""
                }
            ],
            "template": "{{.type}}{{if .scope}} ({{.scope}}){{end}}: {{.subject}}{{if .body}}\n\n{{.body}}{{end}}{{if .footer}}\n\n{{.footer}}{{end}}"
        }
    }`
    var cfg Config
    if err := json.Unmarshal([]byte(defaultJSON), &cfg); err != nil {
        log.Fatalf("Error parsing built-in default config: %v", err)
    }
    return cfg
}

// CollectUserInput prompts the user for input based on the configuration.
func CollectUserInput(cfg Config) map[string]string {
    reader := bufio.NewReader(os.Stdin)
    userInput := make(map[string]string)

    for _, item := range cfg.Message.Items {
        var input string

        for {
            // Print the field description in bold and in the designated color.
            fmt.Println(BoldStyle + DescriptionColor + item.Desc + ColorReset)
            // Optionally, display a hint if available.
            if item.Hint != "" {
                fmt.Println("Hint:", item.Hint)
            }

            // For select fields, list the options.
            if item.Form == "select" {
                for idx, option := range item.Options {
                    // For the "type" field, highlight option names in bold.
                    if strings.ToLower(item.Name) == "type" {
                        fmt.Printf("%s%d%s) %s%s%s: %s\n", OptionColor, idx+1, ColorReset, BoldStyle, option.Name, ColorReset, option.Desc)
                    } else {
                        fmt.Printf("%s%d%s) %s: %s\n", OptionColor, idx+1, ColorReset, option.Name, option.Desc)
                    }
                }
                // Prompt for the option number.
                fmt.Print(PromptColor + "Enter option number: " + ColorReset)
            } else {
                // For input and multiline fields.
                prompt := "Enter value: "
                if item.Required {
                    prompt = "Enter value (required): "
                }
                // If a default value exists, show it.
                if item.Default != "" {
                    prompt = fmt.Sprintf("%s (default: %s): ", prompt, item.Default)
                }
                fmt.Print(PromptColor + prompt + ColorReset)
            }

            rawInput, err := reader.ReadString('\n')
            if err != nil {
                fmt.Printf("Error reading input: %v\n", err)
                continue
            }
            input = strings.TrimSpace(rawInput)

            // Use default if input is empty and default exists.
            if input == "" && item.Default != "" {
                input = item.Default
            }

            if input == "" && item.Required {
                fmt.Println("This field is required.")
                continue
            }

            break
        }

        // For select fields, convert numeric input to the corresponding option name.
        if item.Form == "select" && input != "" {
            choice, err := strconv.Atoi(input)
            if err != nil || choice < 1 || choice > len(item.Options) {
                fmt.Println("Please enter a valid option number.")
            } else {
                input = item.Options[choice-1].Name
            }
        }

        userInput[item.Name] = input
    }

    return userInput
}

// RenderTemplate renders the commit message template using the provided data.
func RenderTemplate(cfg Config, data map[string]string) (string, error) {
    tmpl, err := template.New("commitMessage").Parse(cfg.Message.Template)
    if err != nil {
        return "", err
    }

    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", err
    }

    return buf.String(), nil
}

