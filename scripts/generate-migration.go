package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Template struct {
	Name         string
	Path         string
	TemplateBody string
	Config       map[string]interface{}
}

func main() {
	// Flags
	changedTemplates := flag.String("changed", "", "JSON array of changed template paths")
	deletedTemplates := flag.String("deleted", "", "JSON array of deleted template names")
	allTemplates := flag.Bool("all", false, "Process all templates")
	validateOnly := flag.Bool("validate-only", false, "Only validate, don't generate")
	output := flag.String("output", "", "Output SQL file path")
	version := flag.String("version", time.Now().Format("20060102150405"), "Migration version")
	commitSHA := flag.String("commit-sha", "", "Source commit SHA")

	flag.Parse()

	// Parse deleted template names
	var deletedNames []string
	if *deletedTemplates != "" {
		if err := json.Unmarshal([]byte(*deletedTemplates), &deletedNames); err != nil {
			fatal("Failed to parse deleted templates JSON: %v", err)
		}
	}

	// Load templates
	var templates []Template
	var err error

	if *allTemplates {
		templates, err = scanAllTemplates("../templates")
	} else if *changedTemplates != "" {
		var changedPaths []string
		if err := json.Unmarshal([]byte(*changedTemplates), &changedPaths); err != nil {
			fatal("Failed to parse changed templates JSON: %v", err)
		}
		templates, err = loadChangedTemplates(changedPaths)
	} else if len(deletedNames) > 0 {
		// Only deletions, no changed templates
		templates = []Template{}
	} else {
		fatal("Error: --changed, --deleted, or --all flag required")
	}

	if err != nil {
		fatal("Failed to load templates: %v", err)
	}

	if len(templates) == 0 && len(deletedNames) == 0 {
		fmt.Println("⚠️  No templates to process")

		if *output != "" {
			emptySQL := generateEmptyMigration(*version, *commitSHA)
			writeFile(*output, emptySQL)
			fmt.Printf("✅ Empty migration file created: %s\n", *output)
		}
		return
	}

	// Validate templates
	if err := validateTemplates(templates); err != nil {
		fatal("Validation failed: %v", err)
	}

	if len(templates) > 0 {
		fmt.Printf("✅ Validated %d templates\n", len(templates))
	}
	if len(deletedNames) > 0 {
		fmt.Printf("⚠️  Detected %d deleted templates: %v\n", len(deletedNames), deletedNames)
	}

	if *validateOnly {
		fmt.Println("✅ Validation complete (--validate-only mode)")
		return
	}

	// Generate migration SQL
	sql := generateMigrationSQL(templates, deletedNames, *version, *commitSHA)

	// Write output
	if *output != "" {
		writeFile(*output, sql)
		fmt.Printf("✅ Migration file generated: %s\n", *output)
		fmt.Printf("   Templates updated: %d\n", len(templates))
		fmt.Printf("   Templates deleted: %d\n", len(deletedNames))
	} else {
		fmt.Println(sql)
	}
}

func scanAllTemplates(baseDir string) ([]Template, error) {
	var templates []Template

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read templates directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		templateDir := filepath.Join(baseDir, entry.Name())
		tmpl, err := loadTemplate(templateDir)
		if err != nil {
			return nil, fmt.Errorf("failed to load template at %s: %w", templateDir, err)
		}
		templates = append(templates, tmpl)
	}

	return templates, nil
}

func loadChangedTemplates(paths []string) ([]Template, error) {
	var templates []Template

	for _, path := range paths {
		// path: "templates/vue" or "templates/react"
		tmpl, err := loadTemplate(path)
		if err != nil {
			return nil, fmt.Errorf("failed to load template at %s: %w", path, err)
		}
		templates = append(templates, tmpl)
	}

	return templates, nil
}

func loadTemplate(path string) (Template, error) {
	// Read template.tmpl
	tmplPath := filepath.Join(path, "template.tmpl")
	templateBody, err := os.ReadFile(tmplPath)
	if err != nil {
		return Template{}, fmt.Errorf("template.tmpl not found: %w", err)
	}

	// Read config.json
	configPath := filepath.Join(path, "config.json")
	var config map[string]interface{}
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return Template{}, fmt.Errorf("config.json not found: %w", err)
	}
	if err := json.Unmarshal(configData, &config); err != nil {
		return Template{}, fmt.Errorf("invalid config.json: %w", err)
	}

	// Read metadata.json
	metadataPath := filepath.Join(path, "metadata.json")
	var metadata map[string]interface{}
	metadataData, err := os.ReadFile(metadataPath)
	if err != nil {
		return Template{}, fmt.Errorf("metadata.json not found: %w", err)
	}
	if err := json.Unmarshal(metadataData, &metadata); err != nil {
		return Template{}, fmt.Errorf("invalid metadata.json: %w", err)
	}

	// Merge config and metadata
	finalConfig := make(map[string]interface{})
	for k, v := range metadata {
		finalConfig[k] = v
	}
	for k, v := range config {
		finalConfig[k] = v
	}

	return Template{
		Name:         metadata["name"].(string),
		Path:         path,
		TemplateBody: string(templateBody),
		Config:       finalConfig,
	}, nil
}

func validateTemplates(templates []Template) error {
	names := make(map[string]bool)

	for _, tmpl := range templates {
		// Check for duplicate names
		if names[tmpl.Name] {
			return fmt.Errorf("duplicate template name: %s", tmpl.Name)
		}
		names[tmpl.Name] = true

		// Validate required fields
		if tmpl.Name == "" {
			return fmt.Errorf("template at %s has empty name", tmpl.Path)
		}
		if tmpl.TemplateBody == "" {
			return fmt.Errorf("template %s has empty body", tmpl.Name)
		}

		// Validate config structure
		requiredFields := []string{"description", "version", "categories"}
		for _, field := range requiredFields {
			if _, ok := tmpl.Config[field]; !ok {
				return fmt.Errorf("template %s missing '%s' in metadata", tmpl.Name, field)
			}
		}
	}

	return nil
}

func generateMigrationSQL(templates []Template, deletedNames []string, version, commitSHA string) string {
	var b strings.Builder

	// Header
	b.WriteString("-- Generated by container-go-template\n")
	b.WriteString(fmt.Sprintf("-- Version: %s\n", version))
	b.WriteString(fmt.Sprintf("-- Generated at: %s\n", time.Now().Format(time.RFC3339)))
	if commitSHA != "" {
		b.WriteString(fmt.Sprintf("-- Source commit: %s\n", commitSHA))
	}
	b.WriteString(fmt.Sprintf("-- Changed templates: %d\n", len(templates)))
	b.WriteString(fmt.Sprintf("-- Deleted templates: %d\n\n", len(deletedNames)))

	// List changed template names
	if len(templates) > 0 {
		b.WriteString("-- Updated templates:\n")
		for _, tmpl := range templates {
			categories := tmpl.Config["categories"].([]interface{})
			categoryStr := ""
			if len(categories) > 0 {
				categoryStr = categories[0].(string)
			}
			b.WriteString(fmt.Sprintf("--   - %s (%s)\n", tmpl.Name, categoryStr))
		}
		b.WriteString("\n")
	}

	// List deleted template names
	if len(deletedNames) > 0 {
		b.WriteString("-- Deleted templates:\n")
		for _, name := range deletedNames {
			b.WriteString(fmt.Sprintf("--   - %s\n", name))
		}
		b.WriteString("\n")
	}

	// Phase 0: Mark deleted templates as inactive
	if len(deletedNames) > 0 {
		b.WriteString("-- Phase 0: Mark deleted templates as inactive\n")
		b.WriteString("-- These templates were removed from the repository\n")
		b.WriteString("-- Existing containers using these templates are not affected\n")
		b.WriteString("UPDATE TEMPLATES\n")
		b.WriteString("SET status = 'inactive'\n")
		b.WriteString("WHERE name IN (")

		escapedNames := make([]string, len(deletedNames))
		for i, name := range deletedNames {
			escapedNames[i] = fmt.Sprintf("'%s'", escapeSQLString(name))
		}
		b.WriteString(strings.Join(escapedNames, ", "))
		b.WriteString(")\n")
		b.WriteString("AND status = 'active';\n\n")
	}

	// Phase 1: Deprecate existing templates
	if len(templates) > 0 {
		b.WriteString("-- Phase 1: Deprecate existing versions of changed templates\n")
		b.WriteString("UPDATE TEMPLATES\n")
		b.WriteString("SET status = 'deprecated'\n")
		b.WriteString("WHERE name IN (")

		names := make([]string, len(templates))
		for i, tmpl := range templates {
			names[i] = fmt.Sprintf("'%s'", escapeSQLString(tmpl.Name))
		}
		b.WriteString(strings.Join(names, ", "))
		b.WriteString(");\n\n")
	}

	// Phase 2: Insert new template versions
	if len(templates) > 0 {
		b.WriteString("-- Phase 2: Insert new versions of changed templates\n")
		b.WriteString("INSERT INTO TEMPLATES (name, template_body, template_config, status) VALUES\n\n")

		for i, tmpl := range templates {
			b.WriteString(fmt.Sprintf("-- %s\n", tmpl.Name))
			b.WriteString(fmt.Sprintf("('%s',\n", escapeSQLString(tmpl.Name)))
			b.WriteString(fmt.Sprintf("'%s',\n", escapeSQLString(tmpl.TemplateBody)))

			configJSON, _ := json.Marshal(tmpl.Config)
			b.WriteString(fmt.Sprintf("'%s',\n", escapeSQLString(string(configJSON))))
			b.WriteString("'active')")

			if i < len(templates)-1 {
				b.WriteString(",\n\n")
			} else {
				b.WriteString(";\n")
			}
		}
	}

	return b.String()
}

func generateEmptyMigration(version, commitSHA string) string {
	var b strings.Builder

	b.WriteString("-- Generated by container-go-template\n")
	b.WriteString(fmt.Sprintf("-- Version: %s\n", version))
	b.WriteString(fmt.Sprintf("-- Generated at: %s\n", time.Now().Format(time.RFC3339)))
	if commitSHA != "" {
		b.WriteString(fmt.Sprintf("-- Source commit: %s\n", commitSHA))
	}
	b.WriteString("\n")
	b.WriteString("-- No template changes detected\n")
	b.WriteString("-- This migration file is intentionally empty\n")

	return b.String()
}

func escapeSQLString(s string) string {
	// Escape backslashes and single quotes for SQL
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "'", "''")
	return s
}

func writeFile(path, content string) {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		fatal("Failed to write file %s: %v", path, err)
	}
}

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "❌ "+format+"\n", args...)
	os.Exit(1)
}
