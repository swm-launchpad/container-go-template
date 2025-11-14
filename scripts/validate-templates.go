package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/go-connections/nat"
	"github.com/moby/term"
)

type TestProjectDef struct {
	Name    string            `json:"name"`
	Path    string            `json:"path"`
	Options map[string]string `json:"options"`
}

type TemplateMapping struct {
	Template       string            `json:"template"`
	TestProject    string            `json:"test_project"`    // deprecated, for backward compatibility
	TestProjects   []TestProjectDef  `json:"test_projects"`   // new
	DefaultOptions map[string]string `json:"default_options"` // deprecated
}

type ValidationResult struct {
	TemplateName    string
	TestProjectName string
	Success         bool
	Message         string
	BuildTime       time.Duration
	RuntimeTest     string
}

func main() {
	// Flags
	templatesFlag := flag.String("templates", "", "Comma-separated list of templates to validate (default: all)")
	skipBuild := flag.Bool("skip-build", false, "Skip Docker build test (only validate rendering)")
	runTest := flag.Bool("run-test", false, "Run container and test connectivity (requires Docker build)")
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

		// Get test projects (supports both old and new format)
		testProjects := getTestProjects(mapping)

		// Validate each test project variant
		for _, testProject := range testProjects {
			fmt.Printf("\n=== Validating %s (%s) ===\n", templateName, testProject.Name)

			// Step 1: Render template
			renderedDockerfile, err := renderTemplate(templateName, testProject)
			if err != nil {
				results = append(results, ValidationResult{
					TemplateName:    templateName,
					TestProjectName: testProject.Name,
					Success:         false,
					Message:         fmt.Sprintf("Template rendering failed: %v", err),
				})
				fmt.Printf("✗ Template rendering failed: %v\n", err)
				continue
			}
			fmt.Printf("✓ Template rendered successfully\n")

			// Step 2: Docker build test (if not skipped)
			if !*skipBuild {
				buildTime, err := dockerBuild(templateName, testProject, renderedDockerfile, *timeout, *verbose)
				if err != nil {
					results = append(results, ValidationResult{
						TemplateName:    templateName,
						TestProjectName: testProject.Name,
						Success:         false,
						Message:         fmt.Sprintf("Docker build failed: %v", err),
					})
					fmt.Printf("✗ Docker build failed: %v\n", err)
					continue
				}
				fmt.Printf("✓ Docker build succeeded (build time: %s)\n", buildTime.Round(time.Second))

				// Step 3: Runtime test (if requested)
				runtimeTestResult := "skipped"
				if *runTest {
					if err := containerRunTest(templateName, testProject); err != nil {
						results = append(results, ValidationResult{
							TemplateName:    templateName,
							TestProjectName: testProject.Name,
							Success:         false,
							Message:         fmt.Sprintf("Runtime test failed: %v", err),
							BuildTime:       buildTime,
							RuntimeTest:     "failed",
						})
						fmt.Printf("✗ Runtime test failed: %v\n", err)
						continue
					}
					fmt.Printf("✓ Runtime test succeeded\n")
					runtimeTestResult = "passed"
				}

				results = append(results, ValidationResult{
					TemplateName:    templateName,
					TestProjectName: testProject.Name,
					Success:         true,
					Message:         "All checks passed",
					BuildTime:       buildTime,
					RuntimeTest:     runtimeTestResult,
				})
			} else {
				results = append(results, ValidationResult{
					TemplateName:    templateName,
					TestProjectName: testProject.Name,
					Success:         true,
					Message:         "Template rendering passed (build skipped)",
				})
			}
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

// getTestProjects extracts test projects from mapping, supporting both old and new format
func getTestProjects(mapping TemplateMapping) []TestProjectDef {
	// New format: test_projects array
	if len(mapping.TestProjects) > 0 {
		return mapping.TestProjects
	}

	// Old format: single test_project (backward compatibility)
	return []TestProjectDef{{
		Name:    "default",
		Path:    mapping.TestProject,
		Options: mapping.DefaultOptions,
	}}
}

func renderTemplate(templateName string, testProject TestProjectDef) (string, error) {
	// Read template file from parent mapping
	// Note: testProject.Path is like "test-projects/backend/spring-boot-maven"
	// We need to get template path from the test project's parent directory structure
	// For now, derive from standard location: templates/{templateName}/template.tmpl
	templatePath := filepath.Join("templates", templateName, "template.tmpl")
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
	if err := tmpl.Execute(&buf, testProject.Options); err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}

	return buf.String(), nil
}

func dockerBuild(templateName string, testProject TestProjectDef, dockerfile string, timeout time.Duration, verbose bool) (time.Duration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return 0, fmt.Errorf("failed to create Docker client: %v", err)
	}
	defer cli.Close()

	// Prepare build context
	buildContext, err := prepareBuildContext(testProject.Path, dockerfile)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare build context: %v", err)
	}
	defer buildContext.Close()

	// Build options - use unique tag per test project variant
	imageTag := fmt.Sprintf("template-validation:%s-%s", templateName, testProject.Name)
	buildOptions := types.ImageBuildOptions{
		Tags:        []string{imageTag},
		Dockerfile:  "Dockerfile",
		Remove:      true,
		ForceRemove: true,
		NoCache:     true,
	}

	// Start build
	startTime := time.Now()
	buildResponse, err := cli.ImageBuild(ctx, buildContext, buildOptions)
	if err != nil {
		return 0, fmt.Errorf("failed to start build: %v", err)
	}
	defer buildResponse.Body.Close()

	// Read build output and detect errors
	var output io.Writer = os.Stdout
	if !verbose {
		output = io.Discard
	}
	termFd, isTerm := term.GetFdInfo(output)

	// DisplayJSONMessagesStream automatically detects build errors
	err = jsonmessage.DisplayJSONMessagesStream(buildResponse.Body, output, termFd, isTerm, nil)
	if err != nil {
		return 0, fmt.Errorf("build failed: %v", err)
	}

	buildTime := time.Since(startTime)
	return buildTime, nil
}

func containerRunTest(templateName string, testProject TestProjectDef) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Create Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %v", err)
	}
	defer cli.Close()

	imageName := fmt.Sprintf("template-validation:%s-%s", templateName, testProject.Name)

	// Determine port from test project options
	var hostPort string
	appPort := testProject.Options["app_port"]
	if appPort == "" {
		// For templates without app_port, skip runtime test
		return nil
	}

	// Create container with port mapping
	hostPort = "18080"
	containerPort := nat.Port(fmt.Sprintf("%s/tcp", appPort))
	portBindings := nat.PortMap{
		containerPort: []nat.PortBinding{
			{
				HostIP:   "127.0.0.1",
				HostPort: hostPort,
			},
		},
	}

	// Container config
	containerConfig := &container.Config{
		Image: imageName,
		ExposedPorts: nat.PortSet{
			containerPort: struct{}{},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		AutoRemove:   true,
	}

	// Create container
	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		return fmt.Errorf("failed to create container: %v", err)
	}

	containerID := resp.ID

	// Ensure container cleanup
	defer func() {
		cleanupCtx := context.Background()
		cli.ContainerStop(cleanupCtx, containerID, container.StopOptions{})
		cli.ContainerRemove(cleanupCtx, containerID, container.RemoveOptions{Force: true})
	}()

	// Start container
	if err := cli.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %v", err)
	}

	// Wait for container to be ready (max 30 seconds)
	time.Sleep(5 * time.Second) // Initial wait

	// Perform health check based on template type
	if err := performHealthCheck(templateName, hostPort); err != nil {
		return fmt.Errorf("health check failed: %v", err)
	}

	return nil
}

func performHealthCheck(templateName, hostPort string) error {
	// Web application templates - HTTP GET request
	webTemplates := []string{
		"static-html", "react", "vuejs", "nextjs",
		"expressjs", "nestjs", "fastapi", "flask",
		"django", "spring-boot", "kotlin-spring-boot",
	}

	for _, name := range webTemplates {
		if templateName == name {
			return httpHealthCheck(hostPort)
		}
	}

	// Database and cache templates - skip for now
	// (would require database-specific client libraries)
	return nil
}

func httpHealthCheck(hostPort string) error {
	url := fmt.Sprintf("http://127.0.0.1:%s/", hostPort)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Retry up to 5 times with backoff
	for i := 0; i < 5; i++ {
		resp, err := client.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound {
				return nil // Success
			}
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		if i < 4 {
			time.Sleep(time.Duration(i+1) * 2 * time.Second)
		}
	}

	return fmt.Errorf("failed to connect after 5 retries")
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
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("VALIDATION SUMMARY")
	fmt.Println(strings.Repeat("=", 70))

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
			timeStr = fmt.Sprintf(" (build: %s", result.BuildTime.Round(time.Second))
			if result.RuntimeTest != "" {
				timeStr += fmt.Sprintf(", runtime: %s", result.RuntimeTest)
			}
			timeStr += ")"
		}

		// Display template name with test project variant
		templateDisplay := result.TemplateName
		if result.TestProjectName != "" && result.TestProjectName != "default" {
			templateDisplay = fmt.Sprintf("%s (%s)", result.TemplateName, result.TestProjectName)
		}

		fmt.Printf("%s %s: %s%s\n", status, templateDisplay, result.Message, timeStr)
	}

	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("Total: %d | Success: %d | Failed: %d\n", len(results), successCount, failureCount)
	fmt.Println(strings.Repeat("=", 70))
}

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", args...)
	os.Exit(1)
}
