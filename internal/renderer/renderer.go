package renderer

import (
    "fmt"
    "strings"
)

func h1(text string) string { return "# " + text + "\n" }
func h2(text string) string { return "## " + text + "\n" }

func bullet(items []string) string {
    if len(items) == 0 {
        return ""
    }
    var b strings.Builder
    for _, it := range items {
        b.WriteString("- " + it + "\n")
    }
    return b.String()
}

func getSlice(m map[string]interface{}, key string) []interface{} {
    v, ok := m[key]
    if !ok {
        return nil
    }
    if s, ok := v.([]interface{}); ok {
        return s
    }
    return nil
}

func toStrings(v interface{}) []string {
    arr, ok := v.([]interface{})
    if !ok {
        return nil
    }
    out := make([]string, 0, len(arr))
    for _, x := range arr {
        if s, ok := x.(string); ok && strings.TrimSpace(s) != "" {
            out = append(out, s)
        }
    }
    return out
}

func sectionHeadline(cfg map[string]interface{}) string {
    if v, ok := cfg["headline"].(string); ok && strings.TrimSpace(v) != "" {
        return h1(v)
    }
    return ""
}

func sectionBadges(cfg map[string]interface{}) string {
    outCfg, _ := cfg["output"].(map[string]interface{})
    if outCfg == nil || outCfg["include_badges"] == nil || outCfg["include_badges"] == false {
        return ""
    }
    badges := getSlice(cfg, "badges")
    if len(badges) == 0 {
        return ""
    }
    lines := make([]string, 0, len(badges))
    for _, b := range badges {
        row, _ := b.(map[string]interface{})
        label, _ := row["label"].(string)
        value, _ := row["value"].(string)
        label = strings.TrimSpace(label)
        value = strings.TrimSpace(value)
        if label != "" && value != "" {
            lines = append(lines, fmt.Sprintf("**%s:** %s", label, value))
        }
    }
    if len(lines) == 0 {
        return ""
    }
    return strings.Join(lines, "\n") + "\n\n"
}

func sectionSkills(cfg map[string]interface{}) string {
    skills := getSlice(cfg, "skills")
    if len(skills) == 0 {
        return ""
    }
    var out []string
    out = append(out, h2("Skills"))
    for _, g := range skills {
        row, _ := g.(map[string]interface{})
        name, _ := row["group"].(string)
        items := toStrings(row["items"])
        if name == "" {
            name = "Skills"
        }
        out = append(out, fmt.Sprintf("**%s:**", name))
        out = append(out, bullet(items))
    }
    return strings.Join(out, "\n")
}

func sectionProjects(cfg map[string]interface{}) string {
    projects := getSlice(cfg, "projects")
    if len(projects) == 0 {
        return ""
    }
    var out []string
    out = append(out, h2("Projects"))
    for _, p := range projects {
        row, _ := p.(map[string]interface{})
        name, _ := row["name"].(string)
        desc, _ := row["description"].(string)
        tags := toStrings(row["tags"])
        highlights := toStrings(row["highlights"])
        repo, _ := row["repo"].(string)

        if name == "" {
            name = "Project"
        }
        out = append(out, fmt.Sprintf("**Title:** %s", name))
        if strings.TrimSpace(desc) != "" {
            out = append(out, fmt.Sprintf("**Summary:** %s", desc))
        }
        if len(tags) > 0 {
            out = append(out, "**Tags:** "+strings.Join(tags, ", "))
        }
        if len(highlights) > 0 {
            out = append(out, "**Highlights:**")
            out = append(out, bullet(highlights))
        }
        if strings.TrimSpace(repo) != "" {
            out = append(out, fmt.Sprintf("**Repository:** %s", repo))
        }
        out = append(out, "")
    }
    return strings.Join(out, "\n")
}

func sectionHighlights(cfg map[string]interface{}) string {
    highlights := toStrings(cfg["highlights"])
    if len(highlights) == 0 {
        return ""
    }
    var out []string
    out = append(out, h2("Highlights"))
    out = append(out, bullet(highlights))
    return strings.Join(out, "\n")
}

func RenderMarkdown(cfg map[string]interface{}) string {
    title := "Profile overview"
    if outCfg, ok := cfg["output"].(map[string]interface{}); ok {
        if t, ok := outCfg["title"].(string); ok && strings.TrimSpace(t) != "" {
            title = t
        }
    }
    parts := []string{
        sectionHeadline(cfg),
        sectionBadges(cfg),
        h2(title),
        sectionSkills(cfg),
        sectionProjects(cfg),
        sectionHighlights(cfg),
    }
    // join non-empty parts with single newlines
    var out []string
    for _, p := range parts {
        if strings.TrimSpace(p) != "" {
            out = append(out, strings.TrimSpace(p))
        }
    }
    return strings.Join(out, "\n\n") + "\n"
}
