//go:build !ci

package security

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestSecurityHeaderIntegration tests comprehensive security header implementation across different frameworks
func TestSecurityHeaderIntegration(t *testing.T) {
	testCases := []struct {
		name            string
		framework       string
		template        string
		expectedHeaders []string
	}{
		{
			name:      "Go Gin API Template",
			framework: "gin",
			template: `package main

import "github.com/gin-gonic/gin"

func setupRoutes() *gin.Engine {
	r := gin.Default()
	
	r.GET("/api/users", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.JSON(200, gin.H{"users": []string{}})
	})
	
	return r
}`,
			expectedHeaders: []string{
				"X-Content-Type-Options",
				"X-Frame-Options",
				"X-XSS-Protection",
				"nosniff",
				"DENY",
			},
		},
		{
			name:      "Node.js Express API Template",
			framework: "express",
			template: `const express = require('express');
const app = express();

app.get('/api/users', (req, res) => {
  res.setHeader('Content-Type', 'application/json');
  res.json({ users: [] });
});

module.exports = app;`,
			expectedHeaders: []string{
				"X-Content-Type-Options",
				"X-Frame-Options",
				"X-XSS-Protection",
				"nosniff",
				"DENY",
			},
		},
		{
			name:      "Go HTTP Server Template",
			framework: "http",
			template: `package main

import (
	"net/http"
	"encoding/json"
)

func userHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"users": []string{}})
}`,
			expectedHeaders: []string{
				"X-Content-Type-Options",
				"nosniff",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "template."+tc.framework+".tmpl")

			err := os.WriteFile(testFile, []byte(tc.template), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Apply security fixes
			fixer := NewFixer()
			options := FixerOptions{
				DryRun:       false,
				Verbose:      false,
				FixType:      "headers",
				CreateBackup: false,
			}

			result, err := fixer.FixFile(testFile, options)
			if err != nil {
				t.Fatalf("Security fix failed: %v", err)
			}

			// Read the fixed content
			fixedContent, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Failed to read fixed file: %v", err)
			}

			fixedStr := string(fixedContent)

			// Verify security headers were added
			for _, header := range tc.expectedHeaders {
				if !strings.Contains(fixedStr, header) {
					t.Errorf("Expected security header %q to be added for %s framework", header, tc.framework)
				}
			}

			// Verify fixes were reported
			headerFixFound := false
			for _, fix := range result.FixedIssues {
				if fix.IssueType == MissingSecurityHeader {
					headerFixFound = true
					break
				}
			}

			if !headerFixFound {
				t.Error("Expected security header fix to be reported")
			}
		})
	}
}

// TestCORSIntegrationAcrossFrameworks tests CORS fixes across different web frameworks
func TestCORSIntegrationAcrossFrameworks(t *testing.T) {
	testCases := []struct {
		name          string
		framework     string
		template      string
		vulnerability string
	}{
		{
			name:      "Gin CORS Null Origin",
			framework: "gin",
			template: `func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "null")
		c.Next()
	}
}`,
			vulnerability: "null",
		},
		{
			name:      "Express CORS Wildcard",
			framework: "express",
			template: `app.use((req, res, next) => {
  res.setHeader('Access-Control-Allow-Origin', '*');
  next();
});`,
			vulnerability: "wildcard",
		},
		{
			name:      "HTTP Server CORS Null",
			framework: "http",
			template: `func corsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "null")
}`,
			vulnerability: "null",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "cors_template."+tc.framework+".tmpl")

			err := os.WriteFile(testFile, []byte(tc.template), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// First scan for CORS vulnerabilities
			scanner := NewScanner()
			issues, err := scanner.ScanFile(testFile)
			if err != nil {
				t.Fatalf("Security scan failed: %v", err)
			}

			// Verify CORS vulnerability was detected
			corsIssueFound := false
			for _, issue := range issues {
				if issue.IssueType == CORSVulnerability {
					corsIssueFound = true
					break
				}
			}

			if !corsIssueFound {
				t.Error("Expected CORS vulnerability to be detected")
			}

			// Apply CORS fixes
			fixer := NewFixer()
			options := FixerOptions{
				DryRun:       false,
				Verbose:      false,
				FixType:      "cors",
				CreateBackup: false,
			}

			result, err := fixer.FixFile(testFile, options)
			if err != nil {
				t.Fatalf("CORS fix failed: %v", err)
			}

			// Verify CORS fix was applied
			corsFixFound := false
			for _, fix := range result.FixedIssues {
				if fix.IssueType == CORSVulnerability {
					corsFixFound = true
					break
				}
			}

			if !corsFixFound {
				t.Error("Expected CORS fix to be applied")
			}

			// Read the fixed content
			fixedContent, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Failed to read fixed file: %v", err)
			}

			fixedStr := string(fixedContent)

			// Verify vulnerability was removed/fixed
			switch tc.vulnerability {
			case "null":
				if strings.Contains(fixedStr, `"null"`) {
					t.Error("CORS null origin vulnerability should have been fixed")
				}
				if !strings.Contains(fixedStr, "SECURITY FIX") {
					t.Error("Expected security fix comment to be added")
				}
			case "wildcard":
				if strings.Contains(fixedStr, `'*'`) && !strings.Contains(fixedStr, "isAllowedOrigin") {
					t.Error("CORS wildcard should have been replaced with origin validation")
				}
			}
		})
	}
}

// TestAuthenticationIntegration tests authentication security improvements across different auth patterns
func TestAuthenticationIntegration(t *testing.T) {
	testCases := []struct {
		name     string
		template string
		authType string
	}{
		{
			name: "JWT None Algorithm Vulnerability",
			template: `package auth

import "github.com/golang-jwt/jwt/v4"

func createToken(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	return token.SignedString(jwt.UnsafeAllowNoneSignatureType)
}`,
			authType: "jwt_none",
		},
		{
			name: "JWT Missing Expiration",
			template: `const jwt = require('jsonwebtoken');

function createToken(payload) {
	return jwt.sign(payload, process.env.JWT_SECRET);
}`,
			authType: "jwt_no_exp",
		},
		{
			name: "Insecure Cookie Configuration",
			template: `package session

import "net/http"

func setSessionCookie(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:  "session",
		Value: sessionID,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
}`,
			authType: "insecure_cookie",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "auth_template.tmpl")

			err := os.WriteFile(testFile, []byte(tc.template), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Scan for authentication vulnerabilities
			scanner := NewScanner()
			issues, err := scanner.ScanFile(testFile)
			if err != nil {
				t.Fatalf("Security scan failed: %v", err)
			}

			// Verify authentication vulnerability was detected
			authIssueFound := false
			for _, issue := range issues {
				if issue.IssueType == WeakAuthentication {
					authIssueFound = true
					break
				}
			}

			if !authIssueFound {
				t.Error("Expected authentication vulnerability to be detected")
			}

			// Apply authentication fixes
			fixer := NewFixer()
			options := FixerOptions{
				DryRun:       false,
				Verbose:      false,
				FixType:      "auth",
				CreateBackup: false,
			}

			result, err := fixer.FixFile(testFile, options)
			if err != nil {
				t.Fatalf("Authentication fix failed: %v", err)
			}

			// Read the fixed content
			fixedContent, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Failed to read fixed file: %v", err)
			}

			fixedStr := string(fixedContent)

			// Verify specific authentication fixes
			switch tc.authType {
			case "jwt_none":
				if strings.Contains(fixedStr, "SigningMethodNone") {
					t.Error("JWT none algorithm should have been fixed")
				}
				if !strings.Contains(fixedStr, "HS256") {
					t.Error("Expected secure JWT algorithm to be used")
				}
			case "jwt_no_exp":
				if !strings.Contains(fixedStr, "expir") {
					t.Error("Expected JWT expiration guidance to be added")
				}
			case "insecure_cookie":
				if !strings.Contains(fixedStr, "HttpOnly") || !strings.Contains(fixedStr, "Secure") {
					t.Error("Expected secure cookie flags guidance to be added")
				}
			}

			// Verify authentication fix was reported
			authFixFound := false
			for _, fix := range result.FixedIssues {
				if fix.IssueType == WeakAuthentication {
					authFixFound = true
					break
				}
			}

			if !authFixFound {
				t.Error("Expected authentication fix to be reported")
			}
		})
	}
}

// TestSQLInjectionIntegration tests SQL injection prevention across different database patterns
func TestSQLInjectionIntegration(t *testing.T) {
	testCases := []struct {
		name     string
		template string
		dbType   string
	}{
		{
			name: "Go SQL String Concatenation",
			template: `package db

import "database/sql"

func getUserByID(db *sql.DB, userID string) error {
	query := "SELECT * FROM users WHERE id = " + userID
	_, err := db.Query(query)
	return err
}`,
			dbType: "go_sql",
		},
		{
			name: "Node.js SQL Template Literal",
			template: `const mysql = require('mysql2');

function getUserByID(connection, userID) {
	const query = ` + "`SELECT * FROM users WHERE id = ${userID}`" + `;
	return connection.query(query);
}`,
			dbType: "node_sql",
		},
		{
			name: "Go SQL Format String",
			template: `package db

import (
	"database/sql"
	"fmt"
)

func updateUser(db *sql.DB, userID, name string) error {
	query := fmt.Sprintf("UPDATE users SET name = '%s' WHERE id = %s", name, userID)
	_, err := db.Exec(query)
	return err
}`,
			dbType: "go_format",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "db_template.tmpl")

			err := os.WriteFile(testFile, []byte(tc.template), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Scan for SQL injection vulnerabilities
			scanner := NewScanner()
			issues, err := scanner.ScanFile(testFile)
			if err != nil {
				t.Fatalf("Security scan failed: %v", err)
			}

			// Verify SQL injection vulnerability was detected
			sqlIssueFound := false
			for _, issue := range issues {
				if issue.IssueType == SQLInjectionRisk {
					sqlIssueFound = true
					break
				}
			}

			if !sqlIssueFound {
				t.Error("Expected SQL injection vulnerability to be detected")
			}

			// Apply SQL injection fixes
			fixer := NewFixer()
			options := FixerOptions{
				DryRun:       false,
				Verbose:      false,
				FixType:      "sql",
				CreateBackup: false,
			}

			result, err := fixer.FixFile(testFile, options)
			if err != nil {
				t.Fatalf("SQL injection fix failed: %v", err)
			}

			// Verify SQL injection fix was applied
			sqlFixFound := false
			for _, fix := range result.FixedIssues {
				if fix.IssueType == SQLInjectionRisk {
					sqlFixFound = true
					break
				}
			}

			if !sqlFixFound {
				t.Error("Expected SQL injection fix to be applied")
			}

			// Read the fixed content
			fixedContent, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Failed to read fixed file: %v", err)
			}

			fixedStr := string(fixedContent)

			// Verify security guidance was added
			if !strings.Contains(fixedStr, "SECURITY FIX") {
				t.Error("Expected security fix comment to be added")
			}

			if !strings.Contains(fixedStr, "parameterized") {
				t.Error("Expected parameterized query guidance to be added")
			}
		})
	}
}

// TestEndToEndSecurityWorkflow tests the complete security validation workflow
func TestEndToEndSecurityWorkflow(t *testing.T) {
	// Create a realistic template directory structure
	tmpDir := t.TempDir()

	templates := map[string]string{
		"backend/middleware/cors.go.tmpl": `package middleware

import "github.com/gin-gonic/gin"

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "null")
		c.Header("Content-Type", "application/json")
		c.Next()
	}
}`,
		"backend/auth/jwt.go.tmpl": `package auth

import "github.com/golang-jwt/jwt/v4"

func CreateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"user_id": userID,
	})
	return token.SignedString(jwt.UnsafeAllowNoneSignatureType)
}`,
		"backend/db/users.go.tmpl": `package db

import "database/sql"

func GetUser(db *sql.DB, userID string) (*User, error) {
	query := "SELECT * FROM users WHERE id = " + userID
	row := db.QueryRow(query)
	// ... rest of implementation
	return nil, nil
}`,
		"frontend/api/client.js.tmpl": `const API_BASE = 'http://api.example.com';

function makeRequest(endpoint, data) {
	return fetch(API_BASE + endpoint, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify(data)
	});
}`,
	}

	// Create template files
	for filePath, content := range templates {
		fullPath := filepath.Join(tmpDir, filePath)
		dir := filepath.Dir(fullPath)

		err := os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}

		err = os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create template file %s: %v", filePath, err)
		}
	}

	// Step 1: Scan entire directory for security issues
	scanner := NewScanner()
	report, err := scanner.ScanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("Directory security scan failed: %v", err)
	}

	// Verify comprehensive security issues were detected
	if report.ScannedFiles != 4 {
		t.Errorf("Expected 4 files to be scanned, got %d", report.ScannedFiles)
	}

	if len(report.Issues) == 0 {
		t.Error("Expected security issues to be detected across templates")
	}

	// Verify we have different types of security issues
	issueTypes := make(map[SecurityIssueType]bool)
	for _, issue := range report.Issues {
		issueTypes[issue.IssueType] = true
	}

	expectedTypes := []SecurityIssueType{
		CORSVulnerability,
		WeakAuthentication,
		SQLInjectionRisk,
	}

	for _, expectedType := range expectedTypes {
		if !issueTypes[expectedType] {
			t.Errorf("Expected to find %s issues in template directory", expectedType)
		}
	}

	// Step 2: Apply security fixes to entire directory
	fixer := NewFixer()
	options := FixerOptions{
		DryRun:       false,
		Verbose:      true,
		FixType:      "all",
		CreateBackup: true,
	}

	fixResult, err := fixer.FixDirectory(tmpDir, options)
	if err != nil {
		t.Fatalf("Directory security fix failed: %v", err)
	}

	// Verify fixes were applied across multiple files
	if len(fixResult.FixedIssues) == 0 {
		t.Error("Expected security fixes to be applied across template directory")
	}

	// Verify backups were created
	if fixResult.BackupsCreated == 0 {
		t.Error("Expected backup files to be created")
	}

	// Step 3: Re-scan to verify issues were resolved
	newReport, err := scanner.ScanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("Post-fix directory scan failed: %v", err)
	}

	// Verify critical issues were resolved
	criticalIssues := newReport.CountBySeverity(SeverityCritical)
	if criticalIssues > 0 {
		t.Errorf("Expected critical security issues to be resolved, but %d remain", criticalIssues)
	}

	// Step 4: Verify specific fixes in each file
	fixedFiles := map[string][]string{
		"backend/middleware/cors.go.tmpl": {"SECURITY FIX", "X-Content-Type-Options"},
		"backend/auth/jwt.go.tmpl":        {"HS256", "SECURITY FIX"},
		"backend/db/users.go.tmpl":        {"parameterized", "SECURITY FIX"},
	}

	for filePath, expectedContent := range fixedFiles {
		fullPath := filepath.Join(tmpDir, filePath)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			t.Fatalf("Failed to read fixed file %s: %v", filePath, err)
		}

		contentStr := string(content)
		for _, expected := range expectedContent {
			if !strings.Contains(contentStr, expected) {
				t.Errorf("File %s should contain %q after security fixes", filePath, expected)
			}
		}
	}

	t.Logf("End-to-end security workflow completed successfully:")
	t.Logf("- Scanned %d files", report.ScannedFiles)
	t.Logf("- Found %d security issues", len(report.Issues))
	t.Logf("- Applied %d security fixes", len(fixResult.FixedIssues))
	t.Logf("- Created %d backup files", fixResult.BackupsCreated)
	t.Logf("- Remaining critical issues: %d", criticalIssues)
}
