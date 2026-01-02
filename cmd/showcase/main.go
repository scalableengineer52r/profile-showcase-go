package main

import (
    "flag"
    "fmt"
    "os"

    "gopkg.in/yaml.v3"

    "profile-showcase-go/internal/renderer"
)

func loadConfig(path string) (map[string]interface{}, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("read config: %w", err)
    }
    var cfg map[string]interface{}
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("parse yaml: %w", err)
    }
    if cfg == nil {
        cfg = map[string]interface{}{}
    }
    return cfg, nil
}

func main() {
    configPath := flag.String("config", "", "Path to config.yml")
    outputPath := flag.String("output", "OUTPUT.md", "Path to write rendered Markdown")
    flag.Parse()

    if *configPath == "" {
        fmt.Fprintln(os.Stderr, "missing --config path")
        os.Exit(1)
    }

    cfg, err := loadConfig(*configPath)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    md := renderer.RenderMarkdown(cfg)
    if err := os.WriteFile(*outputPath, []byte(md), 0o644); err != nil {
        fmt.Fprintln(os.Stderr, "write output:", err)
        os.Exit(1)
    }
    fmt.Println("Written:", *outputPath)
}
