package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

type TemplateMapping struct {
	Template       string            `json:"template"`
	TestProject    string            `json:"test_project"`
	DefaultOptions map[string]string `json:"default_options"`
}

type ValidationResult struct {
	TemplateName string
	Success      bool
	Message      string
	BuildTime    time.Duration
}

func main() {
	// Flags
	templatesFlag := flag.String("templates", "", "Comma-separated list of templates to validate (default: all)")
	skipBuild := flag.Bool("skip-build", false, "Skip Docker build test (only validate rendering)")
	timeout := flag.Duration("timeout", 10*time.Minute, "Build timeout per template")
	verbose := flag.Bool("verbose", false, "Show detailed build logs")

	flag.Parse()

	// Change to repository root
	if err := os.Chdir(".."); err != nil {
		fatal("Failed to change directory: %v", err)
	}

	// Load mappings
	mappings, err := loadMappings("test-projects/mappings.json")
	if err != nil {
		fatal("Failed to load mappings: %v", err)
	}

	// Filter templates if specified
	var templatesToValidate []string
	if *templatesFlag != "" {
		templatesToValidate = strings.Split(*templatesFlag, ",")
	} else {
		for name := range mappings {
			templatesToValidate = append(templatesToValidate, name)
		}
	}

	// Validate each template
	var results []ValidationResult
	for _, templateName := range templatesToValidate {
		mapping, ok := mappings[templateName]
		if !ok {
			results = append(results, ValidationResult{
				TemplateName: templateName,
				Success:      false,
				Message:      "Template not found in mappings",
			})
			continue
		}

		fmt.Printf("\n=== Validating %s ===\n", templateName)

		// Step 1: Render template
		renderedDockerfile, err := renderTemplate(templateName, mapping)
		if err != nil {
			results = append(results, ValidationResult{
				TemplateName: templateName,
				Success:      false,
				Message:      fmt.Sprintf("Template rendering failed: %v", err),
			})
			fmt.Printf("✗ Template rendering failed: %v\n", err)
			continue
		}
		fmt.Printf("✓ Template rendered successfully\n")

		// Step 2: Docker build test (if not skipped)
		if !*skipBuild {
			buildTime, err := dockerBuild(templateName, mapping, renderedDockerfile, *timeout, *verbose)
			if err != nil {
				results = append(results, ValidationResult{
					TemplateName: templateName,
					Success:      false,
					Message:      fmt.Sprintf("Docker build failed: %v", err),
				})
				fmt.Printf("✗ Docker build failed: %v\n", err)
				continue
			}
			results = append(results, ValidationResult{
				TemplateName: templateName,
				Success:      true,
				Message:      "All checks passed",
				BuildTime:    buildTime,
			})
			fmt.Printf("✓ Docker build succeeded (build time: %s)\n", buildTime.Round(time.Second))
		} else {
			results = append(results, ValidationResult{
				TemplateName: templateName,
				Success:      true,
				Message:      "Template rendering passed (build skipped)",
			})
		}
	}

	// Print summary
	printSummary(results)

	// Exit with error if any validation failed
	for _, result := range results {
		if !result.Success {
			os.Exit(1)
		}
	}
}

func loadMappings(path string) (map[string]TemplateMapping, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var mappings map[string]TemplateMapping
	if err := json.Unmarshal(data, &mappings); err != nil {
		return nil, err
	}

	return mappings, nil
}

func renderTemplate(templateName string, mapping TemplateMapping) (string, error) {
	// Read template file
	templatePath := filepath.Join(mapping.Template, "template.tmpl")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %v", err)
	}

	// Parse and execute template
	tmpl, err := template.New(templateName).Parse(string(templateContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, mapping.DefaultOptions); err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}

	return buf.String(), nil
}

func dockerBuild(templateName string, mapping TemplateMapping, dockerfile string, timeout time.Duration, verbose bool) (time.Duration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return 0, fmt.Errorf("failed to create Docker client: %v", err)
	}
	defer cli.Close()

	// Prepare build context
	buildContext, err := prepareBuildContext(mapping.TestProject, dockerfile)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare build context: %v", err)
	}
	defer buildContext.Close()

	// Build options
	buildOptions := types.ImageBuildOptions{
		Tags:       []string{fmt.Sprintf("template-validation:%s", templateName)},
		Dockerfile: "Dockerfile",
		Remove:     true,
		ForceRemove: true,
		NoCache:    true,
	}

	// Start build
	startTime := time.Now()
	buildResponse, err := cli.ImageBuild(ctx, buildContext, buildOptions)
	if err != nil {
		return 0, fmt.Errorf("failed to start build: %v", err)
	}
	defer buildResponse.Body.Close()

	// Read build output
	if verbose {
		_, err = io.Copy(os.Stdout, buildResponse.Body)
	} else {
		_, err = io.Copy(io.Discard, buildResponse.Body)
	}
	if err != nil {
		return 0, fmt.Errorf("build failed: %v", err)
	}

	buildTime := time.Since(startTime)
	return buildTime, nil
}

func prepareBuildContext(testProjectPath, dockerfile string) (io.ReadCloser, error) {
	// Create temporary directory for build context
	tempDir, err := os.MkdirTemp("", "docker-build-*")
	if err != nil {
		return nil, err
	}
	// Note: tempDir will be cleaned up by the OS
	// We cannot defer RemoveAll here because TarWithOptions reads files asynchronously

	// Copy test project to temp dir
	if err := copyDir(testProjectPath, tempDir); err != nil {
		os.RemoveAll(tempDir)
		return nil, err
	}

	// Write Dockerfile to temp dir
	dockerfilePath := filepath.Join(tempDir, "Dockerfile")
	if err := os.WriteFile(dockerfilePath, []byte(dockerfile), 0644); err != nil {
		os.RemoveAll(tempDir)
		return nil, err
	}

	// Create tar archive
	tarReader, err := archive.TarWithOptions(tempDir, &archive.TarOptions{})
	if err != nil {
		os.RemoveAll(tempDir)
		return nil, err
	}

	// Wrap the reader to clean up tempDir when closed
	return &contextCloser{
		ReadCloser: tarReader,
		cleanup: func() error {
			return os.RemoveAll(tempDir)
		},
	}, nil
}

// contextCloser wraps an io.ReadCloser and calls cleanup when closed
type contextCloser struct {
	io.ReadCloser
	cleanup func() error
}

func (c *contextCloser) Close() error {
	err1 := c.ReadCloser.Close()
	err2 := c.cleanup()
	if err1 != nil {
		return err1
	}
	return err2
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		// Copy file
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(targetPath, data, info.Mode())
	})
}

func printSummary(results []ValidationResult) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("VALIDATION SUMMARY")
	fmt.Println(strings.Repeat("=", 60))

	successCount := 0
	failureCount := 0

	for _, result := range results {
		status := "✓"
		if !result.Success {
			status = "✗"
			failureCount++
		} else {
			successCount++
		}

		timeStr := ""
		if result.BuildTime > 0 {
			timeStr = fmt.Sprintf(" (build time: %s)", result.BuildTime.Round(time.Second))
		}

		fmt.Printf("%s %s: %s%s\n", status, result.TemplateName, result.Message, timeStr)
	}

	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Total: %d | Success: %d | Failed: %d\n", len(results), successCount, failureCount)
	fmt.Println(strings.Repeat("=", 60))
}

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", args...)
	os.Exit(1)
}
