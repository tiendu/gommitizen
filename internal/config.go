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

// Option represents a selectable option for a form item.
type Option struct {
    Name string `json:"name"`
    Desc string `json:"desc"`
}

// Item represents a field in the commit form.
type Item struct {
    Name       string   `json:"name"`
    Desc       string   `json:"desc"`
    Form       string   `json:"form"`
    Options    []Option `json:"options,omitempty"`
    Required   bool     `json:"required,omitempty"`
    Default    string   `json:"default,omitempty"`
    Hint       string   `json:"hint,omitempty"`
    Validation string   `json:"validation,omitempty"`
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

// LoadConfig loads the configuration from the file "configs/default.json".
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

// getHighlightColors returns the foreground and background colors for a change type option.
func getHighlightColors(option string) (fg, bg string) {
    switch strings.ToLower(option) {
    case "feat":
        return "white", "green"
    case "fix":
        return "white", "red"
    case "docs":
        return "black", "cyan"
    case "style":
        return "black", "yellow"
    case "refactor":
        return "white", "blue"
    case "perf":
        return "black", "magenta"
    case "test":
        return "white", "black"
    case "chore":
        return "white", "cyan"
    case "revert":
        return "black", "green"
    case "wip":
        return "black", "white"
    default:
        return "white", "black"
    }
}

// CollectUserInput prompts the user for input based on the configuration.
func CollectUserInput(cfg Config) map[string]string {
    reader := bufio.NewReader(os.Stdin)
    userInput := make(map[string]string)

    for _, item := range cfg.Message.Items {
        var input string
        for {
            // Print the field description using bold and color.
            // (Using our beautify helpers: Bold and Color.)
            fmt.Println(Bold(Color(item.Desc, "cyan")))
            if item.Hint != "" {
                fmt.Println("Hint:", item.Hint)
            }

            if item.Form == "select" {
                // List the options.
                for idx, option := range item.Options {
                    // For the "type" field, use our Highlight function.
                    if strings.ToLower(item.Name) == "type" {
                        fg, bg := getHighlightColors(option.Name)
                        fmt.Printf("%s) %s: %s\n",
                            Color(strconv.Itoa(idx+1), "yellow"),
                            Highlight(option.Name, fg, bg),
                            option.Desc)
                    } else {
                        fmt.Printf("%s) %s: %s\n",
                            Color(strconv.Itoa(idx+1), "yellow"),
                            option.Name,
                            option.Desc)
                    }
                }
                fmt.Print(Color("Enter option number: ", "green"))
            } else {
                // For input or multiline fields.
                prompt := "Enter value"
                if item.Required {
                    prompt += " (required)"
                }
                if item.Default != "" {
                    prompt += fmt.Sprintf(" (default: %s)", item.Default)
                }
                prompt += ": "
                fmt.Print(Color(prompt, "green"))
            }

            rawInput, err := reader.ReadString('\n')
            if err != nil {
                fmt.Printf("Error reading input: %v\n", err)
                continue
            }
            input = strings.TrimSpace(rawInput)
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

