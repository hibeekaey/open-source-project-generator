# TODO Resolution Report

**Total TODOs Found:** 2954

## Summary

- **Resolved:** 3
- **Documented for Future Work:** 945
- **Removed (Obsolete):** 0
- **False Positives:** 2006

## Resolved TODOs

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:216
**Original:** .Message), "send email") {
**Action:** Replaced with documentation
**Description:** TODO replaced with proper implementation guidance

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:232
**Original:** Send email with reset token",
**Action:** Replaced with documentation
**Description:** TODO replaced with proper implementation guidance

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:216
**Original:** Send email with reset token
**Action:** Replaced with documentation
**Description:** TODO replaced with proper implementation guidance

## Documented for Future Work

### /Users/inertia/cuesoft/working-area/cmd/todo-scanner/main.go:84
**Message:** s: %d\n",
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/cmd/todo-scanner/main.go:85
**Message:** s, report.Summary.PerformanceTODOs,
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/docker-compose.yml:32
**Message:** 
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/app/app.go:201
**Message:** 
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/app.go:216
**Message:** , info, warn, error)")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/app/errors.go:128
**Message:** {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/errors.go:136
**Message:** {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/errors.go:137
**Message:** ("Stack trace: %s", err.Stack)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/errors.go:124
**Message:** ("Caused by: %v", err.Cause)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/errors.go:129
**Message:** ("Stack trace: %s", err.Stack)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/errors.go:123
**Message:** && err.Cause != nil {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/logger.go:16
**Message:** LogLevel = iota
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/logger.go:83
**Message:** (msg string, args ...interface{}) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/logger.go:85
**Message:** Logger.Printf(msg, args...)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/logger.go:114
**Message:** 
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/logger.go:84
**Message:** {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/logger.go:28
**Message:** Logger *log.Logger
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/logger.go:113
**Message:** ":
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/logger.go:69
**Message:** Logger = log.New(multiWriter, "DEBUG ", log.Ldate|log.Ltime|log.Lshortfile)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/app/logger.go:82
**Message:** logs a debug message
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer.go:302
**Message:** 
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer.go:301
**Message:** ") || strings.Contains(message, "fix") {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer.go:47
**Message:** 
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer.go:277
**Message:** Type == "fixme" || todoType == "bug" || strings.Contains(message, "security") {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:214
**Message:** ", CategoryBug},
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:185
**Message:** ", "security issue", PriorityHigh},
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:25
**Message:** This is a security vulnerability
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:187
**Message:** ", "security vulnerability", PriorityHigh},
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:186
**Message:** ", "critical bug", PriorityHigh},
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/integration_test.go:322
**Message:** %s", todo.Message)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/integration_test.go:327
**Message:** ")
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/integration_test.go:318
**Message:** .Category == CategorySecurity {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/integration_test.go:109
**Message:** This is a security issue
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:276
**Message:** .Category == CategorySecurity {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:444
**Message:** 
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:445
**Message:** "
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:32
**Message:** s processes all remaining TODOs after security implementation
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:51
**Message:** ",
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:271
**Message:** ",
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:269
**Message:** ") ||
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:505
**Message:** "
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:504
**Message:** 
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:476
**Message:** TODOs,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:473
**Message:** s,
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:460
**Message:** s": %d
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:360
**Message:** s:** %d\n", report.Summary.BugTODOs))
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:357
**Message:** s))
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:305
**Message:** TODOs++
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:304
**Message:** 
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:299
**Message:** s++
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:274
**Message:** 
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:268
**Message:** category
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:46
**Message:** TODOs         int
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:43
**Message:** s    int
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:145
**Message:** ", "memory leak", "", PriorityCritical},
**Reason:** Performance optimization - requires benchmarking and analysis
**Category:** Performance

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:223
**Message:** s: 1,
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:102
**Message:** = false
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:105
**Message:** .Category == CategorySecurity {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:106
**Message:** = true
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:113
**Message:** {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:114
**Message:** ")
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:123
**Message:** .Priority == PriorityCritical {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:143
**Message:** ", "security vulnerability", "", PriorityCritical},
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:144
**Message:** ", "security issue here", "", PriorityCritical},
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:49
**Message:** Update API documentation`,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:42
**Message:** Memory leak in this function
**Reason:** Performance optimization - requires benchmarking and analysis
**Category:** Performance

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:172
**Message:** Fix security hole", "", CategorySecurity},
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:178
**Message:** ", "", "", CategoryBug},
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:28
**Message:** This is a security vulnerability
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:214
**Message:** ",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:216
**Message:** Security vulnerability",
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:339
**Message:** , got %d", summary.CriticalTODOs)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:254
**Message:** ",
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:327
**Message:** , Priority: PriorityHigh},
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:335
**Message:** s != 2 {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:79
**Message:** s in main.go, security.go, and README.md
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:336
**Message:** s, got %d", summary.SecurityTODOs)
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:333
**Message:** s immediately", securityCount))
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:231
**Message:** s := cu.filterTodosByCategory(analysis.TODOs, CategorySecurity)
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:232
**Message:** s) > 0 {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:235
**Message:** s and vulnerabilities",
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:237
**Message:** s)) * 3 * time.Minute,
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:238
**Message:** s(securityTodos),
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:265
**Message:** s := cu.filterTodosExcludeCategory(analysis.TODOs, CategorySecurity)
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:300
**Message:** 
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:301
**Message:** "
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:330
**Message:** sByCategory(analysis.TODOs, CategorySecurity))
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:450
**Message:** , Message: "Bug fix"},
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:462
**Message:** s, got %d", len(nonSecurityTodos))
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:454
**Message:** s := utils.filterTodosByCategory(todos, CategorySecurity)
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:455
**Message:** s) != 2 {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:410
**Message:** , "Bug"},
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:367
**Message:** s, duplicates)
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:323
**Message:** ",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:456
**Message:** s, got %d", len(securityTodos))
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:170
**Message:** message")
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:83
**Message:** ",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:460
**Message:** s := utils.filterTodosExcludeCategory(todos, CategorySecurity)
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:461
**Message:** s) != 2 {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:466
**Message:** s)+len(nonSecurityTodos) != len(todos) {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/cli/integration_test.go:132
**Message:** info (use env var for dev)
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/integration/version_consistency_test.go:1003
**Message:** ging
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/models/security.go:79
**Message:** ging while maintaining security.
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/models/security_test.go:521
**Message:** mpty(t, err.Error(), "Security error should have message")
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/models/security_test.go:522
**Message:** mpty(t, err.Component, "Security error should have component")
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/models/security_test.go:523
**Message:** mpty(t, err.Operation, "Security error should have operation")
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/models/security_test.go:524
**Message:** mpty(t, err.Remediation, "Security error should have remediation")
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/authentication_security_validation_test.go:282
**Message:** xpected, result)
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/authentication_security_validation_test.go:78
**Message:** xpectedContent: []string{"SECURITY FIX"},
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/cors_security_validation_test.go:62
**Message:** xpectedContent: []string{"SECURITY FIX"},
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/cors_security_validation_test.go:71
**Message:** xpectedContent: []string{"SECURITY FIX"},
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/fixes.go:105
**Message:** information exposure",
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/fixes.go:329
**Message:** Exposure(line string) string {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/fixes.go:343
**Message:** info (use env var for dev)"
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/fixes.go:102
**Message:** Information Exposure",
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/fixes.go:106
**Message:** information in production",
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/fixes.go:331
**Message:** ") ||
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/fixes.go:330
**Message:** information in production
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/fixes.go:107
**Message:** Exposure,
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/fixes.go:103
**Message:** |trace|stack).*(true|enabled|on)(.*)`),
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/lint_rules.go:201
**Message:** Information Exposure",
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/lint_rules.go:207
**Message:** information should not be enabled in production environments",
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/lint_rules.go:208
**Message:** information in production environments",
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/lint_rules.go:203
**Message:** |trace|stack).*(?:true|enabled|on)`,
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/lint_rules.go:202
**Message:** |trace|stack).*(?:true|enabled|on)`),
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/lint_rules.go:210
**Message:** ", "information-leakage", "production", "security"},
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/lint_rules.go:206
**Message:** information may be exposed in production",
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/patterns.go:118
**Message:** |trace|stack).*(?:true|enabled|on)`),
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/patterns.go:122
**Message:** information in production environments",
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/patterns.go:117
**Message:** Information Exposure",
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/patterns.go:121
**Message:** information may be exposed in production",
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/patterns_test.go:91
**Message:** Enabled",
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/patterns_test.go:92
**Message:** true`,
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/security_regression_prevention_test.go:217
**Message:** os.Getenv("DEBUG") == "true"`,
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/security_regression_prevention_test.go:316
**Message:** true
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/security_regression_prevention_test.go:590
**Message:** ",
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/security_regression_test.go:390
**Message:** information exposure",
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/security_regression_test.go:55
**Message:** os.Getenv("DEBUG") == "true"`,
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/security_regression_test.go:56
**Message:** should not be flagged",
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/security_regression_test.go:158
**Message:** enabled in production",
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/security_regression_test.go:54
**Message:** config",
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/security_regression_test.go:162
**Message:** enabled should always be detected",
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/security_regression_test.go:159
**Message:** true`,
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/security_regression_test.go:534
**Message:** true
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/security_validation_test.go:281
**Message:** config",
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/security_validation_test.go:282
**Message:** os.Getenv("DEBUG") == "true"`,
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/template/docker_build_test.go:222
**Message:** .log*
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/pkg/template/docker_build_test.go:221
**Message:** .log*
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/pkg/template/processor_test.go:946
**Message:** info (use env var for dev)
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/template/processor_test.go:956
**Message:** info (use env var for dev)
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/template/processor_test.go:960
**Message:** info (use env var for dev)
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/template/processor_test.go:979
**Message:** info (use env var for dev)
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/template/processor_test.go:911
**Message:** info (use env var for dev)
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/validation/setup.go:379
**Message:** ", "build"); err != nil {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/setup.go:360
**Message:** ",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/setup.go:355
**Message:** "); err != nil {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/vercel_validator.go:463
**Message:** ",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/vercel_validator.go:679
**Message:** ",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/vercel_validator.go:435
**Message:** logging
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/vercel_validator.go:469
**Message:** ",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/vercel_validator.go:479
**Message:** ",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/vercel_validator.go:487
**Message:** ",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/vercel_validator.go:686
**Message:** ",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/vercel_validator.go:455
**Message:** ",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/vercel_validator.go:668
**Message:** ",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/vercel_validator.go:446
**Message:** ",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/vercel_validator.go:437
**Message:** ",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/vercel_validator.go:493
**Message:** ",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/security-audit-config.yml:149
**Message:** information exposure"
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer.go:280
**Message:** Type == "hack" || strings.Contains(message, "performance") {
**Reason:** Performance optimization - requires benchmarking and analysis
**Category:** Performance

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:189
**Message:** ", "performance optimization", PriorityMedium},
**Reason:** Performance optimization - requires benchmarking and analysis
**Category:** Performance

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:28
**Message:** Temporary workaround
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:188
**Message:** ", "performance issue", PriorityMedium},
**Reason:** Performance optimization - requires benchmarking and analysis
**Category:** Performance

### /Users/inertia/cuesoft/working-area/internal/cleanup/integration_test.go:112
**Message:** Temporary workaround
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:280
**Message:** .Category == CategoryPerformance {
**Reason:** Performance optimization - requires benchmarking and analysis
**Category:** Performance

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:108
**Message:** .Category == CategoryPerformance {
**Reason:** Performance optimization - requires benchmarking and analysis
**Category:** Performance

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:146
**Message:** ", "performance problem", "", PriorityHigh},
**Reason:** Performance optimization - requires benchmarking and analysis
**Category:** Performance

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:30
**Message:** Quick fix for performance issue
**Reason:** Performance optimization - requires benchmarking and analysis
**Category:** Performance

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:174
**Message:** performance", "", "", CategoryPerformance},
**Reason:** Performance optimization - requires benchmarking and analysis
**Category:** Performance

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:337
**Message:** sByCategory(analysis.TODOs, CategoryPerformance))
**Reason:** Performance optimization - requires benchmarking and analysis
**Category:** Performance

### /Users/inertia/cuesoft/working-area/cmd/todo-scanner/main.go:36
**Message:** ", "FIXME", "HACK", "XXX", "BUG", "NOTE", "OPTIMIZE"},
**Reason:** Performance optimization - requires benchmarking and analysis
**Category:** Performance

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer.go:292
**Message:** ") {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:32
**Message:** This function needs refactoring
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Refactor

### /Users/inertia/cuesoft/working-area/internal/cleanup/backup_test.go:138
**Message:** The cleanup logic in the current implementation is simplified
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Refactor

### /Users/inertia/cuesoft/working-area/internal/cleanup/integration_test.go:116
**Message:** This function needs refactoring
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Refactor

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:223
**Message:** was addressed during cleanup",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Refactor

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:66
**Message:** ", "FIXME", "HACK", "XXX", "BUG", "NOTE", "OPTIMIZE"},
**Reason:** Performance optimization - requires benchmarking and analysis
**Category:** Performance

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:222
**Message:** ") ||
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:251
**Message:** ") ||
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:224
**Message:** Type == "optimize" {
**Reason:** Performance optimization - requires benchmarking and analysis
**Category:** Performance

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:147
**Message:** ", "optimize this function", "", PriorityMedium},
**Reason:** Performance optimization - requires benchmarking and analysis
**Category:** Performance

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:44
**Message:** This could be faster
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:393
**Message:** s(todos []TODOItem) []CleanupTask {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Refactor

### /Users/inertia/cuesoft/working-area/pkg/models/security_test.go:407
**Message:** mpty(t, result.Warnings, "Should have warnings for disabled secure cleanup")
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/.github/workflows/release.yml:273
**Message:** s
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/.github/workflows/release.yml:336
**Message:** s.md
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/.github/workflows/release.yml:276
**Message:** s.md << EOF
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/.github/workflows/release.yml:274
**Message:** s
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/cmd/cleanup/main.go:96
**Message:** .File, todo.Line, todo.Type, todo.Message)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/cmd/cleanup/main.go:90
**Message:** /FIXME Comments:")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/cmd/cleanup/main.go:93
**Message:** s)-i)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/cmd/cleanup/main.go:89
**Message:** s) > 0 {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/cmd/cleanup/main.go:91
**Message:** = range analysis.TODOs {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/cmd/import-analyzer/main.go:23
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/cmd/security-linter/main.go:112
**Message:** xist(err) {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/cmd/todo-resolver/main.go:47
**Message:** s: %v", err)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/cmd/todo-resolver/main.go:16
**Message:** -resolution-report.md", "Output file for the resolution report")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/cmd/todo-resolver/main.go:45
**Message:** s()
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/cmd/todo-resolver/main.go:29
**Message:** s in: %s\n", absRoot)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/cmd/todo-resolver/main.go:40
**Message:** resolver
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/cmd/todo-resolver/main.go:44
**Message:** s...")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/cmd/todo-resolver/main.go:43
**Message:** s
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/cmd/todo-resolver/main.go:68
**Message:** This was a dry run. No files were modified.")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/cmd/todo-resolver/main.go:59
**Message:** Resolution Complete!\n")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/cmd/todo-resolver/main.go:35
**Message:** ResolverConfig{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/cmd/todo-resolver/main.go:41
**Message:** Resolver(absRoot, config)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/cmd/todo-resolver/main.go:60
**Message:** s found: %d\n", report.TotalFound)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/cmd/todo-scanner/main.go:55
**Message:** s, report.FilesScanned)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/cmd/todo-scanner/main.go:86
**Message:** s, report.Summary.BugTODOs)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/cmd/todo-scanner/main.go:79
**Message:** Analysis Summary:\n")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/cmd/todo-scanner/main.go:80
**Message:** s: %d\n", report.Summary.TotalTODOs)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/cmd/todo-scanner/main.go:82
**Message:** s, report.Summary.HighTODOs,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/cmd/todo-scanner/main.go:33
**Message:** ScanConfig{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/cmd/todo-scanner/main.go:41
**Message:** Scanner(config)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/cmd/todo-scanner/main.go:45
**Message:** /FIXME comments...")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/cmd/todo-scanner/main.go:54
**Message:** /FIXME comments in %d files\n",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/cmd/todo-scanner/main.go:83
**Message:** s, report.Summary.LowTODOs)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/app.go:1551
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/app.go:1561
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/app.go:1891
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/app.go:1618
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/app.go:1609
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/app.go:1456
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/app.go:1442
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/app.go:1332
**Message:** Full version store integration will be completed in task 6.2")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/app.go:599
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/app.go:762
**Message:** Individual configuration setting will be implemented in a future version")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/app.go:773
**Message:** Configuration reset will be implemented in a future version")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/app.go:823
**Message:** Full template update integration will be completed in task 6.2")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/app.go:992
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/app.go:828
**Message:** Template update requires version storage integration")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/app/app_test.go:270
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer.go:109
**Message:** Regex.FindStringSubmatch(line); matches != nil {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer.go:20
**Message:** Item struct {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer.go:91
**Message:** Regex := regexp.MustCompile(`(?i)(TODO|FIXME|HACK|XXX|BUG|NOTE)[\s:]*(.*)`)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer.go:88
**Message:** Comments(rootDir string) ([]TODOItem, error) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer.go:119
**Message:** s = append(todos, todo)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer.go:275
**Message:** Type = strings.ToLower(todoType)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer.go:273
**Message:** Type, message string) Priority {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer.go:110
**Message:** = TODOItem{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer.go:126
**Message:** s, err
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer.go:89
**Message:** s []TODOItem
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:65
**Message:** .Line)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:52
**Message:** s, got %d", expectedTodos, len(todos))
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:64
**Message:** .Line <= 0 {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:60
**Message:** .File != testFile {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:58
**Message:** Types[todo.Type] = true
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:9
**Message:** Comments(t *testing.T) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:17
**Message:** comments
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:181
**Message:** Type string
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:198
**Message:** Type, test.message, test.expected, priority)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:57
**Message:** = range todos {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:34
**Message:** This is just a note
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:71
**Message:** Types[expectedType] {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:69
**Message:** ", "FIXME", "HACK", "XXX", "NOTE"}
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:42
**Message:** comments
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:56
**Message:** Types := make(map[string]bool)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:46
**Message:** comments: %v", err)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:195
**Message:** Type, test.message)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:191
**Message:** ", "documentation", PriorityLow},
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:55
**Message:** s
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:61
**Message:** .File)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:190
**Message:** ", "add feature", PriorityLow},
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:51
**Message:** s) != expectedTodos {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:23
**Message:** Implement proper error handling
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:50
**Message:** s := 5 // TODO, FIXME, HACK, XXX, NOTE
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/analyzer_test.go:44
**Message:** s, err := analyzer.AnalyzeTODOComments(tempDir)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/backup.go:165
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/backup_test.go:52
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/integration_test.go:308
**Message:** s: %v", err)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/integration_test.go:306
**Message:** s, err := analyzer.AnalyzeTODOComments(tempDir)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/integration_test.go:118
**Message:** This is just a note
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/integration_test.go:305
**Message:** analysis
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/integration_test.go:154
**Message:** Add utility functions
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/integration_test.go:107
**Message:** Add proper error handling
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/integration_test.go:141
**Message:** Add application logic
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/integration_test.go:312
**Message:** comments")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/integration_test.go:321
**Message:** .Priority == PriorityHigh && todo.Type == "FIXME" {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/integration_test.go:311
**Message:** s) == 0 {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/integration_test.go:129
**Message:** Implement proper testing
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/integration_test.go:50
**Message:** s) == 0 {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/integration_test.go:51
**Message:** comments in test project")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/integration_test.go:317
**Message:** = range todos {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/integration_test.go:315
**Message:** categories and priorities
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/manager.go:201
**Message:** /FIXME comments", len(todos))
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/manager.go:194
**Message:** comments
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/manager.go:195
**Message:** /FIXME comments...")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/manager.go:196
**Message:** s, err := m.analyzer.AnalyzeTODOComments(m.projectRoot)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/manager.go:198
**Message:** comments: %w", err)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/manager.go:123
**Message:** s:   []string{},
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/manager.go:76
**Message:** sResolved     int
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/manager.go:200
**Message:** s = todos
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/manager.go:270
**Message:** s        []TODOItem
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/manager.go:27
**Message:** s   []string
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/manager.go:75
**Message:** sAnalyzed     int
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/manager.go:280
**Message:** /FIXME comments: %d\n"+
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/manager.go:284
**Message:** s),
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:127
**Message:** .File, ".md") ||
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:306
**Message:** part from the line
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:370
**Message:** struct {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:371
**Message:** Item
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:375
**Message:** resolution
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:376
**Message:** ResolutionReport) GenerateReport() string {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:379
**Message:** Resolution Report\n\n")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:380
**Message:** s Found:** %d\n\n", report.TotalFound))
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:389
**Message:** s
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:391
**Message:** s\n\n")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:400
**Message:** s
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:365
**Message:** Item
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:364
**Message:** struct {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:363
**Message:** represents a TODO that was removed as obsolete
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:357
**Message:** Item
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:356
**Message:** struct {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:355
**Message:** represents a TODO that was documented for future work
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:350
**Message:** Item
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:349
**Message:** struct {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:348
**Message:** represents a TODO that was resolved
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:345
**Message:** 
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:344
**Message:** 
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:343
**Message:** 
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:411
**Message:** s
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:342
**Message:** 
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:340
**Message:** ResolutionReport struct {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:339
**Message:** ResolutionReport contains the results of TODO resolution
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:228
**Message:** Resolver) resolveEmailTODO(todo TODOItem) (*ResolvedTODO, error) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:227
**Message:** resolves email-related TODOs by adding proper documentation
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:236
**Message:** .Line-1 >= len(lines) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:221
**Message:** ,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:237
**Message:** .Line)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:220
**Message:** {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:217
**Message:** (todo)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:240
**Message:** with proper documentation
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:215
**Message:** s
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:214
**Message:** Resolver) resolveTODO(todo TODOItem) (*ResolvedTODO, error) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:213
**Message:** resolves a TODO by implementing or fixing it
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:241
**Message:** Email sending should be implemented based on your email service provider"
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:205
**Message:** .Message), pattern) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:198
**Message:** s should be documented for future work
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:334
**Message:** ,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:333
**Message:** {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:196
**Message:** Resolver) canResolveFeatureTODO(todo TODOItem) bool {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:195
**Message:** checks if a feature TODO can be resolved now
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:242
**Message:** .Line-1] = newComment
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:9
**Message:** Resolver handles the resolution and documentation of remaining TODOs
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:10
**Message:** Resolver struct {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:12
**Message:** ResolverConfig
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:15
**Message:** ResolverConfig holds configuration for TODO resolution
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:187
**Message:** .Message), pattern) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:16
**Message:** ResolverConfig struct {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:329
**Message:** ) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:327
**Message:** ) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:21
**Message:** Resolver creates a new TODO resolver
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:325
**Message:** comment"
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:22
**Message:** Resolver(projectRoot string, config *TODOResolverConfig) *TODOResolver {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:178
**Message:** Resolver) isObsoleteTODO(todo TODOItem) bool {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:369
**Message:** represents a false positive TODO
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:24
**Message:** ResolverConfig{}
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:324
**Message:** Resolver) markAsFalsePositive(todo TODOItem) *FalsePositiveTODO {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:323
**Message:** as a false positive
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:230
**Message:** .File)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:319
**Message:** ",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:177
**Message:** checks if the TODO is obsolete or already completed
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:174
**Message:** ActionDocument
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:318
**Message:** ,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:172
**Message:** s are intentional placeholders for generated projects
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:317
**Message:** {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:26
**Message:** Resolver{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:312
**Message:** .File, []byte(newContent), 0644)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:171
**Message:** Resolver) handleTemplateTODO(todo TODOItem) TODOAction {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:307
**Message:** .Line-1] = strings.Replace(originalLine, todo.Context, "", 1)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:33
**Message:** Resolver) ResolveRemainingTODOs() (*TODOResolutionReport, error) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:170
**Message:** determines action for template TODOs
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:34
**Message:** Scanner(DefaultTODOScanConfig())
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:88
**Message:** ActionRemove
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:166
**Message:** Resolver) isTemplateFile(filePath string) bool {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:157
**Message:** .Context, ref) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:36
**Message:** s
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:39
**Message:** s: %w", err)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:42
**Message:** ResolutionReport{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:43
**Message:** s),
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:246
**Message:** .File, []byte(newContent), 0644)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:44
**Message:** , 0),
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:45
**Message:** , 0),
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:304
**Message:** .Line-1], lines[todo.Line:]...)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:46
**Message:** , 0),
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:302
**Message:** .Context) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:301
**Message:** .Line-1]
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:300
**Message:** line if it's a standalone comment
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:47
**Message:** , 0),
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:297
**Message:** .Line)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:50
**Message:** 
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:51
**Message:** = range report.TODOs {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:251
**Message:** {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:144
**Message:** Resolver) isLegitimateCodeReference(todo TODOItem) bool {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:52
**Message:** Action(todo)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:55
**Message:** ActionResolve:
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:296
**Message:** .Line-1 >= len(lines) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:290
**Message:** .File)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:56
**Message:** (todo)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:58
**Message:** at %s:%d: %w", todo.File, todo.Line, err)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:62
**Message:** ActionDocument:
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:288
**Message:** Resolver) removeTODO(todo TODOItem) (*RemovedTODO, error) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:114
**Message:** that can be resolved
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:287
**Message:** removes an obsolete TODO
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:115
**Message:** (todo) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:63
**Message:** (todo)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:66
**Message:** ActionRemove:
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:67
**Message:** (todo)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:69
**Message:** at %s:%d: %w", todo.File, todo.Line, err)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:73
**Message:** ActionIgnore:
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:74
**Message:** )
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:82
**Message:** Action represents the action to take for a TODO
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:116
**Message:** ActionResolve
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:83
**Message:** Action int
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:252
**Message:** ,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:273
**Message:** for generated projects"
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:272
**Message:** .File) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:133
**Message:** s in comments about TODOs
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:254
**Message:** replaced with proper implementation guidance",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:86
**Message:** ActionResolve TODOAction = iota
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:258
**Message:** documents a TODO for future development
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:413
**Message:** s\n\n")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:87
**Message:** ActionDocument
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:271
**Message:** Resolver) getDocumentationReason(todo TODOItem) string {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:270
**Message:** is being documented
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:263
**Message:** ,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:111
**Message:** ActionRemove
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:110
**Message:** (todo) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:106
**Message:** (todo)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:105
**Message:** .File) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:101
**Message:** ActionIgnore
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:100
**Message:** ) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:262
**Message:** {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:260
**Message:** )
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:259
**Message:** Resolver) documentTODO(todo TODOItem) *DocumentedTODO {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:120
**Message:** ActionDocument
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:123
**Message:** is actually a false positive
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:124
**Message:** Resolver) isFalsePositive(todo TODOItem) bool {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:126
**Message:** .File, "/docs/") ||
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:96
**Message:** ActionIgnore
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:95
**Message:** ) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:93
**Message:** Resolver) determineTODOAction(todo TODOItem) TODOAction {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:92
**Message:** Action decides what action to take for a given TODO
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:128
**Message:** .File, "/.kiro/specs/") ||
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:129
**Message:** .File, "/scripts/") {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:89
**Message:** ActionIgnore
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:98
**Message:** TODOItem{File: "/test/docs/README.md"},
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:365
**Message:** s",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:363
**Message:** s",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:357
**Message:** s Found:** 10",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:356
**Message:** Resolution Report",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:342
**Message:** Item{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:340
**Message:** {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:332
**Message:** Item{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:330
**Message:** {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:319
**Message:** Item{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:317
**Message:** {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:308
**Message:** Item{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:306
**Message:** {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:304
**Message:** ResolutionReport{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:303
**Message:** ResolutionReport_GenerateReport(t *testing.T) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:295
**Message:** )
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:286
**Message:** TODOItem{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:285
**Message:** ",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:279
**Message:** TODOItem{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:278
**Message:** ",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:272
**Message:** TODOItem{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:268
**Message:** for generated projects",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:265
**Message:** TODOItem{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:260
**Message:** TODOItem
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:256
**Message:** Resolver("/test", &TODOResolverConfig{})
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:255
**Message:** Resolver_getDocumentationReason(t *testing.T) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:250
**Message:** Email sending should be implemented") {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:237
**Message:** () failed: %v", err)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:235
**Message:** (todo)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:228
**Message:** = TODOItem{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:226
**Message:** Resolver(tmpDir, &TODOResolverConfig{})
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:208
**Message:** Resolver_resolveEmailTODO(t *testing.T) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:174
**Message:** Resolver("/test", &TODOResolverConfig{})
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:173
**Message:** Resolver_isTemplateFile(t *testing.T) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:165
**Message:** )
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:158
**Message:** TODOItem{Context: "// TODO: implement this feature"},
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:157
**Message:** comment",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:152
**Message:** check",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:143
**Message:** TODOItem
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:139
**Message:** Resolver("/test", &TODOResolverConfig{})
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:138
**Message:** Resolver_isLegitimateCodeReference(t *testing.T) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:130
**Message:** )
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:123
**Message:** TODOItem{File: "/test/pkg/template/engine.go"},
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:117
**Message:** comment about TODOs",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:113
**Message:** TODOItem{File: "/test/scripts/audit.sh"},
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:108
**Message:** TODOItem{File: "/test/.kiro/specs/feature/design.md"},
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:103
**Message:** TODOItem{File: "/test/CONTRIBUTING.md"},
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:93
**Message:** TODOItem
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:89
**Message:** Resolver("/test", &TODOResolverConfig{})
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:88
**Message:** Resolver_isFalsePositive(t *testing.T) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:82
**Message:** Action() = %v, want %v", action, tt.expected)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:80
**Message:** Action(tt.todo)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:74
**Message:** ActionDocument,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:69
**Message:** TODOItem{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:68
**Message:** for documentation",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:65
**Message:** ActionResolve,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:61
**Message:** TODOItem{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:60
**Message:** ",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:57
**Message:** ActionRemove,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:52
**Message:** TODOItem{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:48
**Message:** ActionDocument,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:44
**Message:** TODOItem{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:43
**Message:** ",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:40
**Message:** ActionIgnore,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:36
**Message:** TODOItem{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:32
**Message:** ActionIgnore,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:30
**Message:** comments",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:28
**Message:** TODOItem{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:24
**Message:** ActionIgnore,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:20
**Message:** TODOItem{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:16
**Message:** Action
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:15
**Message:** TODOItem
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:11
**Message:** Resolver("/test", &TODOResolverConfig{})
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:10
**Message:** Resolver_determineTODOAction(t *testing.T) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:452
**Message:** s": %d,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:34
**Message:** Summary
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:350
**Message:** s))
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:349
**Message:** s))
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:348
**Message:** s:** %d\n", report.Summary.TotalTODOs))
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:341
**Message:** /FIXME Analysis Report\n\n")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:337
**Message:** Scanner) generateMarkdownReport(report *TODOReport) (string, error) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:352
**Message:** s))
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:325
**Message:** Scanner) GenerateReport(report *TODOReport) (string, error) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:317
**Message:** s++
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:358
**Message:** s))
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:315
**Message:** s++
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:359
**Message:** s))
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:313
**Message:** s++
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:373
**Message:** s := report.PriorityBreakdown[priority]
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:311
**Message:** s++
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:374
**Message:** s) == 0 {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:378
**Message:** s (%d)\n\n", priorityNames[priority], len(todos)))
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:309
**Message:** .Priority {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:380
**Message:** = range todos {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:381
**Message:** .File, todo.Line, todo.Type))
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:382
**Message:** .Message))
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:383
**Message:** .Context))
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:384
**Message:** .Category)))
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:303
**Message:** s++
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:389
**Message:** s
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:301
**Message:** s++
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:390
**Message:** s\n\n")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:397
**Message:** s := range report.FileBreakdown {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:398
**Message:** s)})
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:409
**Message:** s\n", fc.file, fc.count))
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:297
**Message:** .Category {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:416
**Message:** Scanner) generateTextReport(report *TODOReport) (string, error) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:294
**Message:** s++
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:419
**Message:** /FIXME Analysis Report\n")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:293
**Message:** = range todos {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:291
**Message:** Summary{}
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:290
**Message:** Scanner) generateSummary(todos []TODOItem) *TODOSummary {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:427
**Message:** s: %d\n", report.Summary.TotalTODOs))
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:429
**Message:** s, report.Summary.HighTODOs,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:430
**Message:** s, report.Summary.LowTODOs))
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:432
**Message:** = range report.TODOs {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:434
**Message:** .Priority), todo.File, todo.Line, todo.Type))
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:435
**Message:** .Message))
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:436
**Message:** .Context))
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:443
**Message:** Scanner) generateJSONReport(report *TODOReport) (string, error) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:462
**Message:** _count": %d
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:468
**Message:** s,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:233
**Message:** Scanner) determineCategory(message, context, filePath string) Category {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:469
**Message:** s,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:232
**Message:** item
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:470
**Message:** s,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:471
**Message:** s,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:472
**Message:** s,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:474
**Message:** s,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:475
**Message:** s,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:210
**Message:** Type == "hack" ||
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:477
**Message:** s)), nil
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:200
**Message:** Type == "fixme" || todoType == "bug" ||
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:481
**Message:** Scanner) priorityToString(p Priority) string {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:197
**Message:** Type = strings.ToLower(todoType)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:496
**Message:** Scanner) categoryToString(c Category) string {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:194
**Message:** Scanner) determinePriority(todoType, message, context string) Priority {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:193
**Message:** item
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:15
**Message:** Scanner struct {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:184
**Message:** Scanner) isTextFile(path string) bool {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:16
**Message:** ScanConfig
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:174
**Message:** Scanner) shouldSkipFile(path string) bool {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:19
**Message:** ScanConfig holds configuration for TODO scanning
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:170
**Message:** s, scanner.Err()
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:20
**Message:** ScanConfig struct {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:166
**Message:** s = append(todos, todo)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:27
**Message:** Report represents a comprehensive report of all TODO items
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:157
**Message:** = TODOItem{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:28
**Message:** Report struct {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:156
**Message:** Regex.FindStringSubmatch(line); matches != nil {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:33
**Message:** s             []TODOItem
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:148
**Message:** s []TODOItem
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:141
**Message:** Scanner) scanFile(filePath string, todoRegex *regexp.Regexp) ([]TODOItem, error) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:140
**Message:** comments
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:351
**Message:** s))
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:134
**Message:** s[i].File < report.TODOs[j].File
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:35
**Message:** Item
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:132
**Message:** s[i].Priority > report.TODOs[j].Priority // Higher priority first
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:36
**Message:** Item
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:131
**Message:** s[i].Priority != report.TODOs[j].Priority {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:37
**Message:** Item
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:130
**Message:** s, func(i, j int) bool {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:40
**Message:** Summary provides summary statistics
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:129
**Message:** s by priority and then by file
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:41
**Message:** Summary struct {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:127
**Message:** s)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:42
**Message:** s       int
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:116
**Message:** .File] = append(report.FileBreakdown[todo.File], todo)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:115
**Message:** .Priority] = append(report.PriorityBreakdown[todo.Priority], todo)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:44
**Message:** s int
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:45
**Message:** s     int
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:114
**Message:** .Category] = append(report.CategoryBreakdown[todo.Category], todo)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:47
**Message:** s    int
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:113
**Message:** = range todos {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:110
**Message:** s = append(report.TODOs, todos...)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:48
**Message:** s        int
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:105
**Message:** s, err := ts.scanFile(path, todoRegex)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:49
**Message:** s      int
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:84
**Message:** Regex := regexp.MustCompile(fmt.Sprintf(`(?i)(%s)[\s:]*(.*)`, keywords))
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:50
**Message:** s         int
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:82
**Message:** detection
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:53
**Message:** Scanner creates a new TODO scanner
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:79
**Message:** Item),
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:78
**Message:** Item),
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:54
**Message:** Scanner(config *TODOScanConfig) *TODOScanner {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:77
**Message:** Item),
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:76
**Message:** s:             []TODOItem{},
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:73
**Message:** Report{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:72
**Message:** Scanner) ScanProject(rootDir string) (*TODOReport, error) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:71
**Message:** comments
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:56
**Message:** ScanConfig()
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:58
**Message:** Scanner{config: config}
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:63
**Message:** ScanConfig{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:62
**Message:** ScanConfig() *TODOScanConfig {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:61
**Message:** ScanConfig returns default configuration
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:129
**Message:** {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:124
**Message:** = true
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:122
**Message:** = range report.TODOs {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:130
**Message:** ")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:121
**Message:** = false
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:117
**Message:** ")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:116
**Message:** {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:134
**Message:** Scanner_DeterminePriority(t *testing.T) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:109
**Message:** = true
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:104
**Message:** = range report.TODOs {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:103
**Message:** = false
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:135
**Message:** Scanner(DefaultTODOScanConfig())
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:97
**Message:** s in %s, but didn't", expectedFile)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:91
**Message:** in vendor directory, should be skipped: %s", todo.File)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:138
**Message:** Type string
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:90
**Message:** .File, "vendor/") {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:86
**Message:** .File)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:148
**Message:** ", "add feature", "", PriorityLow},
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:84
**Message:** = range report.TODOs {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:80
**Message:** s in vendor/ directory
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:149
**Message:** ", "remember to update", "", PriorityLow},
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:76
**Message:** comments, but found none")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:153
**Message:** Type, test.message, test.context)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:75
**Message:** s == 0 {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:66
**Message:** Scanner(DefaultTODOScanConfig())
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:155
**Message:** type '%s' with message '%s', expected priority %v, got %v",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:50
**Message:** This should be ignored`,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:48
**Message:** Add installation instructions
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:156
**Message:** Type, test.message, test.expected, priority)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:39
**Message:** This needs immediate attention
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:161
**Message:** Scanner_DetermineCategory(t *testing.T) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:37
**Message:** Implement proper authentication
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:25
**Message:** Add proper error handling
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:162
**Message:** Scanner(DefaultTODOScanConfig())
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:19
**Message:** patterns
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:13
**Message:** _scanner_test")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:11
**Message:** Scanner_ScanProject(t *testing.T) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:192
**Message:** Scanner_GenerateMarkdownReport(t *testing.T) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:193
**Message:** Scanner(DefaultTODOScanConfig())
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:196
**Message:** Report{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:201
**Message:** s: []TODOItem{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:205
**Message:** ",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:207
**Message:** Add error handling",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:221
**Message:** Summary{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:222
**Message:** s:    2,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:224
**Message:** s:  1,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:225
**Message:** s: 1,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:226
**Message:** s:   1,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:228
**Message:** Item),
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:229
**Message:** Item),
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:230
**Message:** Item),
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:234
**Message:** = range report.TODOs {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:235
**Message:** .Category] = append(report.CategoryBreakdown[todo.Category], todo)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:236
**Message:** .Priority] = append(report.PriorityBreakdown[todo.Priority], todo)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:237
**Message:** .File] = append(report.FileBreakdown[todo.File], todo)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:248
**Message:** /FIXME Analysis Report",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:250
**Message:** s:** 2",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:253
**Message:** s",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:265
**Message:** Scanner_ShouldSkipFile(t *testing.T) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:266
**Message:** ScanConfig{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:269
**Message:** Scanner(config)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:291
**Message:** Scanner_IsTextFile(t *testing.T) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:292
**Message:** ScanConfig{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:295
**Message:** Scanner(config)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:319
**Message:** Scanner_GenerateSummary(t *testing.T) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:320
**Message:** Scanner(DefaultTODOScanConfig())
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:322
**Message:** s := []TODOItem{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:330
**Message:** s)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:332
**Message:** s != 5 {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:333
**Message:** s, got %d", summary.TotalTODOs)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:338
**Message:** s != 1 {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:341
**Message:** s != 2 {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner_test.go:342
**Message:** s, got %d", summary.HighTODOs)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:199
**Message:** Time += 2 * time.Minute
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:216
**Message:** Time + duplicateTime + unusedTime + importTime
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:72
**Message:** /FIXME Comments\n\n")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:75
**Message:** Item)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:76
**Message:** = range analysis.TODOs {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:77
**Message:** .Category] = append(categories[todo.Category], todo)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:80
**Message:** s := range categories {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:82
**Message:** s)))
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:84
**Message:** = range todos {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:85
**Message:** .Priority)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:87
**Message:** .Type, priority, todo.File, todo.Line, todo.Message))
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:192
**Message:** (varies by priority)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:193
**Message:** Time := time.Duration(0)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:194
**Message:** = range analysis.TODOs {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:195
**Message:** .Priority {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:197
**Message:** Time += 5 * time.Minute
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:70
**Message:** /FIXME Details
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:201
**Message:** Time += 1 * time.Minute
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:264
**Message:** s
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:266
**Message:** s) > 0 {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:268
**Message:** s",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:269
**Message:** /FIXME comments",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:271
**Message:** s)) * 1 * time.Minute,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:272
**Message:** s(remainingTodos),
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:340
**Message:** s for optimization opportunities", perfCount))
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:360
**Message:** s) + len(analysis.UnusedCode) + len(analysis.Duplicates) + len(analysis.ImportIssues)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:373
**Message:** sByCategory(todos []TODOItem, category Category) []TODOItem {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:374
**Message:** Item
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:375
**Message:** = range todos {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:376
**Message:** .Category == category {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:377
**Message:** )
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:383
**Message:** sExcludeCategory(todos []TODOItem, category Category) []TODOItem {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:384
**Message:** Item
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:385
**Message:** = range todos {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:386
**Message:** .Category != category {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:387
**Message:** )
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:395
**Message:** = range todos {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:397
**Message:** ",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:398
**Message:** .Type, todo.Message),
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:399
**Message:** .File,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:400
**Message:** .Line,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:401
**Message:** .Priority,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:71
**Message:** s) > 0 {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils.go:203
**Message:** Time += 30 * time.Second
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:321
**Message:** s: []TODOItem{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:382
**Message:** s + 1 duplicate + 1 unused + 1 import
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:197
**Message:** ", Description: "Fixed TODO"},
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:329
**Message:** ",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:236
**Message:** ", Description: "Fixed TODO"},
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:273
**Message:** s:        []TODOItem{},
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:282
**Message:** s",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:284
**Message:** s: []TODOItem{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:149
**Message:** /FIXME Comments",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Documentation

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:71
**Message:** s: []TODOItem{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:208
**Message:** ", Description: "Fixed TODO"},
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:443
**Message:** s(t *testing.T) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:446
**Message:** s := []TODOItem{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:467
**Message:** s should account for all original TODOs")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:122
**Message:** ",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/utils_test.go:75
**Message:** ",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/validator.go:333
**Message:** xist(err) || !info.IsDir() {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/validator.go:319
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/cleanup/validator_test.go:168
**Message:** The current parsing implementation is simplified and counts
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/config/manager_test.go:720
**Message:** Current validation may not catch path traversal - this is a test for future enhancement
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/internal/config/manager_test.go:422
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/cli/template_analysis.go:22
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/filesystem/integration_test.go:313
**Message:** xist {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/filesystem/integration_test.go:297
**Message:** xist := []string{
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/integration/comprehensive_integration_test.go:130
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/integration/security_integration_test.go:269
**Message:** xist(err) {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/integration/version_consistency_test.go:110
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/integration/version_consistency_test.go:569
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/integration/version_consistency_test.go:935
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/models/errors_test.go:122
**Message:** mpty(t, err.Component, "Component should be specified")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/models/errors_test.go:124
**Message:** mpty(t, err.Remediation, "Remediation should be provided")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/models/errors_test.go:123
**Message:** mpty(t, err.Operation, "Operation should be specified")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/models/errors_test.go:121
**Message:** mpty(t, err.Error(), "Error message should not be empty")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Bug

### /Users/inertia/cuesoft/working-area/pkg/models/security_error_examples_test.go:63
**Message:** qual(t, randomBytes, randomBytes2, "Generated bytes should be different")
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/models/security_error_examples_test.go:178
**Message:** mpty(t, err.Remediation, "Error %s should have remediation guidance", err.Error())
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/models/security_error_examples_test.go:259
**Message:** mpty(t, id1)
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/models/security_error_examples_test.go:260
**Message:** mpty(t, id2)
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/models/security_error_examples_test.go:261
**Message:** qual(t, id1, id2, "Generated IDs should be unique")
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/models/security_test.go:448
**Message:** mpty(t, result.Warnings, "Should have warnings for disabled entropy check")
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/models/security_test.go:490
**Message:** This is testing the logic, not the validator framework integration
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/models/security_test.go:352
**Message:** mpty(t, result.Errors, "Should have errors for dangerous directory")
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/models/security_test.go:370
**Message:** mpty(t, result.Warnings, "Should have warnings for weak random length")
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/models/security_test.go:423
**Message:** mpty(t, result.Warnings, "Should have warnings for low entropy")
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/models/security_test.go:435
**Message:** mpty(t, result.Warnings, "Should have warnings for alphanumeric format")
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/models/security_test.go:270
**Message:** mpty(t, result.Warnings, "Should have compatibility warnings")
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/models/security_test.go:322
**Message:** mpty(t, result.Warnings, "Should have warnings for insecure permissions")
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/models/security_test.go:218
**Message:** mpty(t, result.Errors, "Should have validation errors")
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/models/security_test.go:104
**Message:** mpty(t, result.Errors, "Should have validation errors")
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/reporting/audit.go:315
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/reporting/audit.go:213
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/reporting/audit_test.go:37
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/reporting/audit_test.go:320
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/reporting/generator.go:336
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/security/authentication_security_validation_test.go:234
**Message:** xpected []string
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/authentication_security_validation_test.go:18
**Message:** xpectedContent []string
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/authentication_security_validation_test.go:38
**Message:** xpectedContent: []string{},
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/authentication_security_validation_test.go:48
**Message:** xpectedContent: []string{},
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/authentication_security_validation_test.go:58
**Message:** xpectedContent: []string{},
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/authentication_security_validation_test.go:68
**Message:** xpectedContent: []string{},
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/authentication_security_validation_test.go:88
**Message:** xpectedContent: []string{},
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/authentication_security_validation_test.go:116
**Message:** xpected := range tt.notExpectedContent {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/authentication_security_validation_test.go:117
**Message:** xpected) {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/authentication_security_validation_test.go:118
**Message:** xpected, result)
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/authentication_security_validation_test.go:28
**Message:** xpectedContent: []string{`"none"`},
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/authentication_security_validation_test.go:241
**Message:** xpected: []string{`"none"`},
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/authentication_security_validation_test.go:248
**Message:** xpected: []string{},
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/authentication_security_validation_test.go:255
**Message:** xpected: []string{},
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/authentication_security_validation_test.go:280
**Message:** xpected := range tt.notExpected {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/authentication_security_validation_test.go:281
**Message:** xpected) {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/cors_security_validation_test.go:17
**Message:** xpectedContent []string
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/cors_security_validation_test.go:35
**Message:** xpectedContent: []string{`'null'`},
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/cors_security_validation_test.go:44
**Message:** xpectedContent: []string{},
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/cors_security_validation_test.go:53
**Message:** xpectedContent: []string{},
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/cors_security_validation_test.go:80
**Message:** xpectedContent: []string{},
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/cors_security_validation_test.go:117
**Message:** xpected := range tt.notExpectedContent {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/cors_security_validation_test.go:118
**Message:** xpected) {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/cors_security_validation_test.go:119
**Message:** xpected, result)
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/cors_security_validation_test.go:26
**Message:** xpectedContent: []string{`"null"`},
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/fileops.go:330
**Message:** xist(err) {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/fileops_test.go:294
**Message:** xist(err) {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/fileops_test.go:283
**Message:** xist(err) {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/fileops_test.go:328
**Message:** xist(err) {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/fixes.go:386
**Message:** crypto/rand has different API - use rand.Read([]byte) instead of rand.Int()"
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/integration_test.go:53
**Message:** xist(err) {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/integration_test.go:217
**Message:** xist(err) {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/integration_test.go:228
**Message:** xist(err) {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/integration_test.go:324
**Message:** xist(err) {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/linter.go:366
**Message:** "
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/linter_test.go:249
**Message:** xist(err) {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/linter_test.go:238
**Message:** xist(err) {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/linter_test.go:260
**Message:** xist(err) {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/security_regression_test.go:342
**Message:** Some of these may not be detected by current patterns
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/security/security_regression_test.go:500
**Message:** Some fixes may not apply to all formats (e.g., YAML/JSON)
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/template/engine_test.go:147
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/template/import_detector.go:111
**Message:** xist": "os",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/template/processor.go:236
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/template/processor.go:315
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/template/processor_test.go:166
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/template/processor_test.go:189
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/template/processor_test.go:258
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/template/processor_test.go:264
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/template/processor_test.go:270
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/template/processor_test.go:360
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/template/processor_test.go:466
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/template/processor_test.go:991
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/template/processor_test.go:1013
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/template/processor_test.go:1026
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/template/scanner.go:95
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/template/scanner_test.go:14
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/template/scanner_test.go:86
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/template/scanner_test.go:93
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/template/template_compilation_integration_test.go:41
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/template/template_compilation_integration_test.go:60
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/template/template_compilation_integration_test.go:76
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/template/template_compilation_verification_test.go:38
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/template/template_compilation_verification_test.go:150
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/template/template_compilation_verification_test.go:321
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/template/template_compilation_verification_test.go:341
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/template/template_edge_cases_test.go:511
**Message:** Add time.Sleep() here later
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/template/template_fixes_test_suite.go:90
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/engine.go:46
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/engine.go:264
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/performance_test.go:80
**Message:** mpty(t, projectName)
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/project_types.go:43
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/project_types.go:64
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/project_types.go:114
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/project_types.go:138
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/project_types.go:267
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/project_types.go:387
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/project_types.go:435
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/project_types.go:489
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/setup.go:121
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/setup.go:159
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/setup.go:294
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/setup.go:320
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/template_validator.go:36
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/template_validator.go:226
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/template_validator.go:280
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/template_validator.go:323
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/template_validator.go:332
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/vercel_validator.go:90
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/vercel_validator.go:116
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/vercel_validator.go:221
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/vercel_validator.go:308
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/vercel_validator.go:355
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/validation/vercel_validator.go:364
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/version/cache.go:165
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/version/compatibility.go:13
**Message:** s        string            `json:"notes,omitempty"`
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/version/compatibility.go:122
**Message:** s as warnings
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/version/compatibility.go:123
**Message:** s != "" {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/version/compatibility.go:127
**Message:** s,
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/version/compatibility.go:193
**Message:** s: "Next.js 15 requires React 18 or later",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/version/compatibility.go:203
**Message:** s: "Next.js 14 requires React 18 or later",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/version/compatibility.go:233
**Message:** s: "Gin v1.9+ requires Go 1.19 or later",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/version/compatibility.go:252
**Message:** s: "Kotlin 2.0 requires Android Gradle Plugin 8.0+",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/version/compatibility_test.go:20
**Message:** s: "React 18 requires React DOM 18+",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/version/github_client_test.go:31
**Message:** s...",
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/version/npm_registry_test.go:252
**Message:** Current implementation only returns latest version
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/version/npm_registry_test.go:266
**Message:** LatestVersion will be whatever the mock server returns
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/version/security_integration_test.go:85
**Message:** xist(err) {
**Reason:** Security enhancement - requires careful implementation and testing
**Category:** Security

### /Users/inertia/cuesoft/working-area/pkg/version/storage.go:49
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/version/storage.go:276
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/version/storage.go:298
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/version/storage_test.go:67
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/version/storage_test.go:576
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/version/storage_test.go:744
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/version/template_updater.go:30
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/version/template_updater.go:66
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/version/template_updater.go:82
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/version/template_updater.go:100
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/version/template_updater.go:123
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/version/template_updater.go:151
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/pkg/version/template_updater.go:179
**Message:** Template restoration will be fully implemented in a future version")
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/test/integration/template_generation_test.go:166
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/test/integration/template_generation_test.go:203
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/test/integration/template_generation_test.go:208
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

### /Users/inertia/cuesoft/working-area/test/integration/template_generation_test.go:1059
**Message:** xist(err) {
**Reason:** Feature enhancement - requires design and implementation planning
**Category:** Feature

## False Positives

These were identified as false positives and ignored:

- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/design.md:326** - s were identified: (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/design.md:232** - ging (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/design.md:333** - s for secure random generation and file operations (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/design.md:324** - s (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/requirements.md:89** - ging issues, I want consistent error handling and logging throughout the application, so that I can effectively troubleshoot problems. (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/requirements.md:20** - s are found THEN the system SHALL either implement the security features or document why they are deferred (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/tasks.md:120** - logging and temporary print statements (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/tasks.md:16** - s (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/CONTRIBUTING.md:492** - reports and feature requests (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/CONTRIBUTING.md:457** - s, please include: (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/CONTRIBUTING.md:3** - report, and feature suggestion. (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/CONTRIBUTING.md:31** - fixes or new features (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/CONTRIBUTING.md:118** - Fixes (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/CONTRIBUTING.md:190** - fix (non-breaking change which fixes an issue) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/CONTRIBUTING.md:342** - fix (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/CONTRIBUTING.md:455** - Reports (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/CONTRIBUTING.md:27** - Reports**: Help us identify and fix bugs (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/DISTRIBUTION.md:234** - build with symbols (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/cmd/import-detector/README.md:36** - ging (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/cmd/security-fixer/README.md:11** - information exposure (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/cmd/security-scanner/README.md:11** - information and detailed error messages that could leak sensitive data (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/cmd/security-scanner/README.md:88** - information enabled in production (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/AUDIT.md:279** - ging (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/CLI_USAGE.md:224** - , info, warn, error) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/CLI_USAGE.md:412** - ging purposes. Log levels can be controlled with the `--log-level` flag. (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/EXAMPLES.md:332** - generate (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/EXAMPLES.md:331** - logging (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/EXAMPLES.md:325** - issues with verbose logging: (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/INSTALLATION.md:254** - , info, warn, error) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/SECURITY_CODING_GUIDE.md:525** - bool   `json:"debug"` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/SECURITY_PATTERNS_EXAMPLES.md:493** - Config(config *DatabaseConfig) { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/SECURITY_PATTERNS_EXAMPLES.md:111** - ging (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/SECURITY_PATTERNS_EXAMPLES.md:581** - ConfigSecure(config *DatabaseConfig) { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/SECURITY_PATTERNS_EXAMPLES.md:938** - ging and forensics (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/SECURITY_RATIONALE.md:486** - .Stack(), // Shows code structure (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/SECURITY_RATIONALE.md:501** - .Stack()), (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/SECURITY_UTILITIES_USAGE.md:671** - ging (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/SECURITY_UTILITIES_USAGE.md:682** - SecurityOperation(operation string, details interface{}) { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/SECURITY_UTILITIES_USAGE.md:683** - _SECURITY") == "true" { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/SECURITY_UTILITIES_USAGE.md:684** - %s - %+v", operation, details) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/TEMPLATE_SECURITY_CHECKLIST.md:130** - ging (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/TEMPLATE_VALIDATION_CHECKLIST.md:115** - ging (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/TROUBLESHOOTING.md:166** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/TROUBLESHOOTING.md:721** - logging for detailed troubleshooting: (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/TROUBLESHOOTING.md:725** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/TROUBLESHOOTING.md:728** - --verbose (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/TROUBLESHOOTING.md:742** - s and feature requests](https://github.com/open-source-template-generator/generator/issues) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/internal/cleanup/README.md:135** - s:   []string{"IMPORTANT:", "CRITICAL:"}, (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/internal/cleanup/README.md:276** - s or unexpected behavior (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/internal/cleanup/README.md:321** - Mode (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/internal/cleanup/README.md:195** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/internal/cleanup/README.md:323** - logging for detailed troubleshooting: (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/internal/cleanup/README.md:187** - , BUG, security mentions (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/scripts/validate-templates/README.md:139** - template generation issues: (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/scripts/validate-templates/README.md:137** - Mode (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/scripts/validate-templates/README.md:81** - ging) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1662** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1356** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7298** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7140** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6528** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:20** - s:** 36 (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:24** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:27** - ging` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:36** - s` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6519** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6510** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6501** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6492** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6483** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6474** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6465** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6456** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6447** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6438** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6429** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6420** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6411** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6402** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6393** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6384** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6375** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6366** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6357** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6348** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6339** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6330** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6321** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6312** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6303** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6294** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6285** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6276** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6267** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6258** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6249** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6240** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6231** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6222** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6213** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6204** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6195** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6186** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6177** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6168** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6159** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6150** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6096** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6087** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6078** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6069** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6060** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6051** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6042** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6033** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6024** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6015** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6006** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5997** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5988** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5979** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5970** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5939** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5898** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5732** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5723** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5714** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5701** - Comments", (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5525** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5485** - comments", (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5296** - Comments\n\n") (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5282** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5278** - Details (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5120** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5116** - Analysis Report", (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5003** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4652** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4508** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4436** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4432** - Analysis Report\n") (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4427** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4238** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4234** - Analysis Report\n\n") (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4229** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4220** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4076** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4072** - " || todoType == "bug" || (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3977** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3788** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3505** - comments: %d\n"+ (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3487** - comments", len(todos)) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3451** - comments...") (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3401** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3293** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3194** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3190** - ", "HACK", "XXX", "NOTE"} (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3095** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3050** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2969** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2902** - , HACK, XXX, BUG, NOTE (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2866** - /HACK comment detection with priority and category classification (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2848** - comments (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2685** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2676** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2659** - comments regularly (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2654** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2650** - TODOs) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2596** - comments in %d files\n", (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2587** - comments...") (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2559** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2515** - Comments:") (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2452** - /HACK comments (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2443** - comments (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2434** - comments and resolve or document them appropriately (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2425** - comment resolution (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2407** - , HACK, XXX (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2389** - comments (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2332** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2270** - ", "HACK", "XXX", "BUG", "NOTE", "OPTIMIZE"}, (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2225** - ", "HACK", "XXX", "BUG", "NOTE", "OPTIMIZE"}, (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1428** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1494** - Exposure(line string) string {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2007** - information exposure"` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2004** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1998** - template generation issues:` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1995** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1989** - Mode` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1986** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1980** - ging)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1977** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1971** - Implement actual security checking` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1962** - Implement actual security checking using govulncheck or similar` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1953** - ",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1950** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1944** - ",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1941** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1935** - ",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1932** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1926** - ",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1923** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1917** - ",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1914** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1908** - ",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1905** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1899** - ",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1896** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1890** - ",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1887** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1881** - ",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1878** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1872** - ",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1869** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1863** - ",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1860** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1854** - logging` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1851** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1845** - ", "build"); err != nil {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1842** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1836** - ",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1833** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1827** - "); err != nil {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1824** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1818** - info (use env var for dev)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1815** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1809** - info (use env var for dev)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1806** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1800** - info (use env var for dev)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1797** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1791** - info (use env var for dev)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1788** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1782** - info (use env var for dev)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1779** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1773** - .log*` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1770** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1764** - .log*` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1761** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1755** - os.Getenv("DEBUG") == "true"`,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1753** - ") == "true"`, (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1752** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1746** - config",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1743** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1737** - true` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1734** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1728** - information exposure",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1725** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1719** - enabled should always be detected",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1716** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1710** - true`,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1707** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1701** - enabled in production",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1698** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1692** - should not be flagged",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1689** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1683** - os.Getenv("DEBUG") == "true"`,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1681** - ") == "true"`, (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1680** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1674** - config",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1671** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1665** - ",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1656** - true` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1653** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1647** - os.Getenv("DEBUG") == "true"`,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1645** - ") == "true"`, (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1644** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1638** - true`,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1635** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1629** - Enabled",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1626** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1620** - information in production environments",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1617** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1611** - information may be exposed in production",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1608** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1602** - |trace|stack).*(?:true|enabled|on)`),` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1599** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1593** - Information Exposure",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1590** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1584** - ", "information-leakage", "production", "security"},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1581** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1575** - information in production environments",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1572** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1566** - information should not be enabled in production environments",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1563** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1557** - information may be exposed in production",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1554** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1548** - |trace|stack).*(?:true|enabled|on)`,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1545** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1539** - |trace|stack).*(?:true|enabled|on)`),` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1536** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1530** - Information Exposure",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1527** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1521** - info (use env var for dev)"` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1518** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1512** - ") ||` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1509** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1503** - information in production` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1500** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1491** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1485** - Exposure,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1482** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1476** - information in production",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1473** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1467** - information exposure",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1464** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1458** - |trace|stack).*(true|enabled|on)(.*)`),` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1455** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1449** - Information Exposure",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1446** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1440** - xpectedContent: []string{"SECURITY FIX"},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1437** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1431** - xpectedContent: []string{"SECURITY FIX"},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1422** - xpected, result)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1419** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1413** - xpectedContent: []string{"SECURITY FIX"},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1410** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1404** - mpty(t, err.Remediation, "Security error should have remediation")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1401** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1395** - mpty(t, err.Operation, "Security error should have operation")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1392** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1386** - mpty(t, err.Component, "Security error should have component")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1383** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1377** - mpty(t, err.Error(), "Security error should have message")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1374** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1368** - ging while maintaining security.` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1365** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1359** - ging` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7357** - s (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1350** - info (use env var for dev)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1347** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1341** - s)+len(nonSecurityTodos) != len(todos) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1339** - s) != len(todos) { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1332** - s, got %d", len(nonSecurityTodos))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1330** - s)) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1323** - s) != 2 {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1314** - s := utils.filterTodosExcludeCategory(todos, CategorySecurity)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1312** - sExcludeCategory(todos, CategorySecurity) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1305** - s, got %d", len(securityTodos))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1303** - s)) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1296** - s) != 2 {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1287** - s := utils.filterTodosByCategory(todos, CategorySecurity)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1285** - sByCategory(todos, CategorySecurity) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1280** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1278** - , Message: "Bug fix"},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1276** - fix"}, (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1275** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1271** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1269** - , "Bug"},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1267** - "}, (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1266** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1260** - s, duplicates)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1251** - ",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1248** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1242** - message")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1233** - ",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1230** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1224** - s immediately", securityCount))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1215** - sByCategory(analysis.TODOs, CategorySecurity))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1213** - s, CategorySecurity)) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1206** - "` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1203** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1197** - ` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1194** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1188** - s := cu.filterTodosExcludeCategory(analysis.TODOs, CategorySecurity)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1186** - sExcludeCategory(analysis.TODOs, CategorySecurity) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1179** - s(securityTodos),` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1177** - s), (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1170** - s)) * 3 * time.Minute,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1161** - s and vulnerabilities",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1152** - s) > 0 {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1143** - s := cu.filterTodosByCategory(analysis.TODOs, CategorySecurity)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1141** - sByCategory(analysis.TODOs, CategorySecurity) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1134** - , got %d", summary.CriticalTODOs)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1125** - s, got %d", summary.SecurityTODOs)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1123** - s) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1116** - s != 2 {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1107** - , Priority: PriorityHigh},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1098** - ",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1089** - s: 1,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1080** - Security vulnerability",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1071** - ",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1064** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1062** - ", "", "", CategoryBug},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1060** - }, (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1053** - Fix security hole", "", CategorySecurity},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1044** - ", "memory leak", "", PriorityCritical},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1035** - ", "security issue here", "", PriorityCritical},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1026** - ", "security vulnerability", "", PriorityCritical},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1017** - .Priority == PriorityCritical {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1008** - ")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:999** - {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:990** - = true` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:981** - .Category == CategorySecurity {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:972** - = false` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:963** - s in main.go, security.go, and README.md` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:954** - Update API documentation`,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:945** - Memory leak in this function` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:936** - This is a security vulnerability` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:927** - "` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:918** - ` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:909** - TODOs,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:900** - s,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:891** - s": %d` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:884** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:882** - s:** %d\n", report.Summary.BugTODOs))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:880** - TODOs)) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:873** - s))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:864** - TODOs++` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:855** - ` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:846** - s++` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:837** - ` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:828** - ") ||` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:819** - category` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:810** - TODOs         int` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:801** - s    int` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:792** - ")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:783** - %s", todo.Message)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:780** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:774** - .Category == CategorySecurity {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:765** - This is a security issue` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:762** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:758** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:756** - ", CategoryBug},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:754** - }, (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:753** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:747** - ", "security vulnerability", PriorityHigh},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:740** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:738** - ", "critical bug", PriorityHigh},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:736** - ", PriorityHigh}, (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:735** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:729** - ", "security issue", PriorityHigh},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:726** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:720** - This is a security vulnerability` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:717** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:711** - ` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:708** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:704** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:702** - ") || strings.Contains(message, "fix") {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:699** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:693** - Type == "fixme" || todoType == "bug" || strings.Contains(message, "security") {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:691** - " || todoType == "bug" || strings.Contains(message, "security") { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:684** - ` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:681** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:675** - logging for detailed troubleshooting:` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:672** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:666** - Mode` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:663** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:657** - s or unexpected behavior` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:654** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:648** - ` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:645** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:639** - , BUG, security mentions` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:637** - , security mentions (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:636** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:630** - s:   []string{"IMPORTANT:", "CRITICAL:"},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:621** - ` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:618** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:612** - ":` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:609** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:603** - Logger.Printf(msg, args...)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:600** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:594** - {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:591** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:585** - (msg string, args ...interface{}) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:582** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:578** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:576** - logs a debug message` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:574** - message (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:573** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:569** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:567** - Logger = log.New(multiWriter, "DEBUG ", log.Ldate|log.Ltime|log.Lshortfile)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:565** - ", log.Ldate|log.Ltime|log.Lshortfile) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:564** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:558** - Logger *log.Logger` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:555** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:549** - LogLevel = iota` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:546** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:540** - ("Stack trace: %s", err.Stack)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:537** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:531** - {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:528** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:522** - ("Stack trace: %s", err.Stack)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:519** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:513** - {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:510** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:504** - ("Caused by: %v", err.Cause)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:501** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:495** - && err.Cause != nil {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:492** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:488** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:486** - , info, warn, error)")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:483** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:477** - ` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:474** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:468** - s and feature requests](https://github.com/open-source-template-generator/generator/issues)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:465** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:459** - --verbose` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:456** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:450** - ` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:447** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:441** - logging for detailed troubleshooting:` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:438** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:432** - ` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:429** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:423** - ging` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:420** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:414** - ging` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:411** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:405** - %s - %+v", operation, details)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:402** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:396** - _SECURITY") == "true" {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:393** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:387** - SecurityOperation(operation string, details interface{}) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:384** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:378** - ging` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:375** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:369** - .Stack()),` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:366** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:360** - .Stack(), // Shows code structure` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:357** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:351** - ging and forensics` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:348** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:342** - ConfigSecure(config *DatabaseConfig) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:339** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:333** - Config(config *DatabaseConfig) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:330** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:324** - ging` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:321** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:315** - bool   `json:"debug"`` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:313** - "` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:312** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:306** - , info, warn, error)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:303** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:297** - generate` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:294** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:288** - logging` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:285** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:279** - issues with verbose logging:` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:276** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:270** - ging purposes. Log levels can be controlled with the `--log-level` flag.` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:267** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:261** - , info, warn, error)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:258** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:252** - ging` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:249** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:243** - ` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:240** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:234** - s, report.Summary.PerformanceTODOs,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:225** - s: %d\n",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:216** - information enabled in production` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:213** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:207** - information and detailed error messages that could leak sensitive data` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:204** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:198** - information exposure` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:195** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:189** - ging` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:186** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:180** - build with symbols` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:177** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:171** - reports and feature requests` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:168** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:162** - s, please include:` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:159** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:153** - Reports` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:150** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:144** - fix` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:141** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:135** - fix (non-breaking change which fixes an issue)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:132** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:126** - Fixes` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:123** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:117** - fixes or new features` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:114** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:108** - Reports**: Help us identify and fix bugs` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:106** - s (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:105** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:99** - report, and feature suggestion.` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:96** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:90** - logging and temporary print statements` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:87** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:81** - s` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:72** - ging issues, I want consistent error handling and logging throughout the application, so that I can effectively troubleshoot problems.` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:69** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:63** - s are found THEN the system SHALL either implement the security features or document why they are deferred` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:54** - s for secure random generation and file operations` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:45** - s were identified:` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:29** - Implement actual security checking using govulncheck or similar" (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:30** - s for secure implementations: (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:37** - , Documentation, Refactor (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:69** - /FIXME comments and created a comprehensive categorized report. The next subtask (2.2) can now proceed to implement the security-related TODOs that have been identified and catalogued. (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:28** - Implement actual security checking" (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:27** - s Identified (as per requirements) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:25** - s:** 36 TODOs (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:22** - s (including the key ones identified in requirements) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/tasks.md:129** - performance and fix inefficiencies (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/tasks.md:136** - template processing performance (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/CHANGELOG.md:34** - d template processing performance (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/CLI_USAGE.md:298** - d for performance (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/CLI_USAGE.md:344** - d for performance and SEO with modern design patterns (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/internal/cleanup/README.md:188** - , performance mentions (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2153** - sByCategory(analysis.TODOs, CategoryPerformance))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2108** - Temporary workaround` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2078** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2087** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2090** - ", "performance issue", PriorityMedium},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2151** - s, CategoryPerformance)) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2027** - template processing performance` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2072** - Type == "hack" || strings.Contains(message, "performance") {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2099** - ", "performance optimization", PriorityMedium},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2045** - d for performance` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2105** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2144** - performance", "", "", CategoryPerformance},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2117** - Quick fix for performance issue` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2081** - Temporary workaround` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4081** - " || (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2070** - " || strings.Contains(message, "performance") { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2018** - performance and fix inefficiencies` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2036** - d template processing performance` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2126** - .Category == CategoryPerformance {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2054** - d for performance and SEO with modern design patterns` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2135** - ", "performance problem", "", PriorityHigh},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2063** - , performance mentions` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2060** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/design.md:366** - string operations (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/design.md:361** - d file I/O patterns (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/tasks.md:144** - memory allocations in hot paths (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/tasks.md:138** - template parsing and rendering (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/tasks.md:132** - string operations and memory allocations (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/CHANGELOG.md:75** - d cross-platform build process (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/CHANGELOG.md:72** - d memory usage in version caching (10,000 ops in ~3ms) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2272** - ", "FIXME", "HACK", "XXX", "BUG", "NOTE", "OPTIMIZE"},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2173** - string operations` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2033** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2042** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2051** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2024** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2335** - mpty(t, result.Warnings, "Should have warnings for disabled secure cleanup")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2308** - This could be faster` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2305** - _scanner_test.go:44 - OPTIMIZE (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2299** - ") ||` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2296** - _scanner.go:251 - OPTIMIZE (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2290** - Type == "optimize" {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2141** - _scanner_test.go:174 - OPTIMIZE (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2288** - " { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2326** - s(todos []TODOItem) []CleanupTask {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2281** - ") ||` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2278** - _scanner.go:222 - OPTIMIZE (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2161** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2164** - d file I/O patterns` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2170** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2263** - This function needs refactoring` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2179** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2324** - s []TODOItem) []CleanupTask { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2254** - The cleanup logic in the current implementation is simplified` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2182** - string operations and memory allocations` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2245** - This function needs refactoring` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2188** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2315** - this function", "", PriorityMedium}, (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2236** - ") {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2233** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2227** - ", "FIXME", "HACK", "XXX", "BUG", "NOTE", "OPTIMIZE"},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2317** - ", "optimize this function", "", PriorityMedium},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2015** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2218** - d cross-platform build process` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2215** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2209** - d memory usage in version caching (10,000 ops in ~3ms)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2206** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2200** - memory allocations in hot paths` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2191** - template parsing and rendering` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2197** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:36** - , FIXME, HACK, XXX, BUG, NOTE, OPTIMIZE keywords (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/design.md:289** - /FIXME comment resolution (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/design.md:186** - s    []string (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/design.md:147** - , FIXME, HACK, XXX (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/design.md:144** - Item struct { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/design.md:63** - /FIXME comments (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/design.md:52** - Comments() ([]TODOItem, error) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/requirements.md:15** - /FIXME comments and resolve or document them appropriately (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/tasks.md:10** - /FIXME/HACK comments (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/tasks.md:23** - s where appropriate (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/tasks.md:22** - s (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/tasks.md:24** - s that should remain for future development (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/tasks.md:25** - s (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/tasks.md:11** - -style comments across the codebase (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/.kiro/specs/project-cleanup/tasks.md:9** - /FIXME comments (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/AUDIT.md:238** - /FIXME comments regularly (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/INSTALLATION.md:358** - s (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/SECURITY_PATTERNS_EXAMPLES.md:733** - xist(err) { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/TEMPLATE_SECURITY_CHECKLIST.md:387** - **: This checklist should be customized based on specific template types and organizational requirements. Not all items may be applicable to every template, but all applicable items should be completed before deployment. (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/TEMPLATE_VALIDATION_CHECKLIST.md:145** - comments include context and assignee (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/docs/TEMPLATE_VALIDATION_CHECKLIST.md:151** - d (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/internal/cleanup/README.md:181** - comments by: (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/internal/cleanup/README.md:189** - s and notes (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/internal/cleanup/README.md:147** - s**: TODO patterns that should not be automatically resolved (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/internal/cleanup/README.md:179** - Analysis (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/internal/cleanup/README.md:67** - /FIXME/HACK comment detection with priority and category classification (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/internal/cleanup/README.md:54** - s, err := analyzer.AnalyzeTODOComments(rootDir) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/internal/cleanup/README.md:183** - , FIXME, HACK, XXX, BUG, NOTE (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/internal/cleanup/README.md:53** - /FIXME comments (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/internal/cleanup/analyzer.go:19** - Item represents a TODO/FIXME comment found in code (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/internal/cleanup/analyzer.go:87** - Comments scans for TODO/FIXME comments in Go files (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:330** - )" (Legitimate code reference (e.g., context.TODO))
- **/Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:179** - s that reference already implemented features (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:153** - comments without issues", // PR template (Legitimate code reference (e.g., context.TODO))
- **/Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:152** - ", (Legitimate code reference (e.g., context.TODO))
- **/Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:146** - .Context, "context.TODO") { (Legitimate code reference (e.g., context.TODO))
- **/Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:145** - is a legitimate Go standard library function (Legitimate code reference (e.g., context.TODO))
- **/Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:136** - .Context, "Check for TODO") { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:135** - .Message, "TODO/FIXME") || (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:134** - .Message, "TODO comments") || (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver.go:99** - ) (Legitimate code reference (e.g., context.TODO))
- **/Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:153** - TODOItem{Context: "- [ ] No TODO comments without issues"}, (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:148** - TODOItem{Context: `"context.TODO": "context",`}, (Legitimate code reference (e.g., context.TODO))
- **/Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:147** - reference", (Legitimate code reference (e.g., context.TODO))
- **/Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:118** - TODOItem{Message: "Check for TODO comments"}, (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/internal/cleanup/todo_resolver_test.go:38** - ": "context",`, (Legitimate code reference (e.g., context.TODO))
- **/Users/inertia/cuesoft/working-area/internal/cleanup/todo_scanner.go:14** - Scanner provides enhanced TODO/FIXME comment analysis (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/pkg/template/import_detector.go:130** - ":         "context", (Legitimate code reference (e.g., context.TODO))
- **/Users/inertia/cuesoft/working-area/scripts/check-imports/main.go:225** - that this check should be done (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/scripts/template-import-scanner/main.go:86** - xist": "os", (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/scripts/template-import-scanner/main.go:152** - ":         "context", (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/scripts/validate-templates/parser.go:128** - We're not returning the error here because it would stop the inspection (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/scripts/validate-templates/parser.go:169** - ", "WithCancel", "WithTimeout", "WithValue"}, (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3147** - Types[todo.Type] = true` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4863** - _scanner_test.go:134 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2460** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2463** - -style comments across the codebase` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2469** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2472** - s` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2478** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2481** - s where appropriate` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2487** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2490** - s that should remain for future development` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2496** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2499** - s` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2505** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2508** - s) > 0 {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2514** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2451** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2517** - /FIXME Comments:")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2523** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2524** - s { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2526** - = range analysis.TODOs {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2532** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2535** - s)-i)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2541** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2542** - .Line, todo.Type, todo.Message) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2544** - .File, todo.Line, todo.Type, todo.Message)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2550** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2553** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2445** - /FIXME comments` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2562** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2568** - -scanner/main.go:33 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2571** - ScanConfig{` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2577** - -scanner/main.go:41 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2580** - Scanner(config)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2586** - -scanner/main.go:45 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2442** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2589** - /FIXME comments...")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2595** - -scanner/main.go:54 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2436** - /FIXME comments and resolve or document them appropriately` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2598** - /FIXME comments in %d files\n",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2604** - -scanner/main.go:55 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2607** - s, report.FilesScanned)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2613** - -scanner/main.go:79 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2616** - Analysis Summary:\n")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2622** - -scanner/main.go:80 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2623** - s) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2625** - s: %d\n", report.Summary.TotalTODOs)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2631** - -scanner/main.go:82 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2632** - s, (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2634** - s, report.Summary.HighTODOs,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2640** - -scanner/main.go:83 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2641** - s) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2643** - s, report.Summary.LowTODOs)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2649** - -scanner/main.go:86 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2433** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2652** - s, report.Summary.BugTODOs)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2427** - /FIXME comment resolution` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2658** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2424** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2661** - /FIXME comments regularly` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2667** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2670** - s` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2418** - s    []string` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2679** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2415** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2688** - **: This checklist should be customized based on specific template types and organizational requirements. Not all items may be applicable to every template, but all applicable items should be completed before deployment.` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2694** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2697** - comments include context and assignee` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2703** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2706** - d` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2712** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2715** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2721** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2724** - Individual configuration setting will be implemented in a future version")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2730** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2733** - Configuration reset will be implemented in a future version")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2739** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2742** - Full template update integration will be completed in task 6.2")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2748** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2751** - Template update requires version storage integration")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2757** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2760** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2766** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2769** - Full version store integration will be completed in task 6.2")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2775** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2778** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2784** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2787** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2793** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2796** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2802** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2805** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2811** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2814** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2820** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2823** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2829** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2832** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2838** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2841** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2847** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2409** - , FIXME, HACK, XXX` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2850** - /FIXME comments` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2856** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2857** - Comments(rootDir) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2859** - s, err := analyzer.AnalyzeTODOComments(rootDir)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2865** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2406** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2868** - /FIXME/HACK comment detection with priority and category classification` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2874** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2875** - patterns that should not be automatically resolved (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2877** - s**: TODO patterns that should not be automatically resolved` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2883** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2886** - Analysis` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2892** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2895** - comments by:` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2901** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2400** - Item struct {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2904** - , FIXME, HACK, XXX, BUG, NOTE` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2910** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2911** - s (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2913** - s and notes` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2919** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2920** - /FIXME comment found in code (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2922** - Item represents a TODO/FIXME comment found in code` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2928** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2931** - Item struct {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2937** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2938** - /FIXME comments in Go files (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2940** - Comments scans for TODO/FIXME comments in Go files` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2946** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2947** - Item, error) { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2949** - Comments(rootDir string) ([]TODOItem, error) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2955** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2956** - Item (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2958** - s []TODOItem` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2964** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2965** - |FIXME|HACK|XXX|BUG|NOTE)[\s:]*(.*)`) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2967** - Regex := regexp.MustCompile(`(?i)(TODO|FIXME|HACK|XXX|BUG|NOTE)[\s:]*(.*)`)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2397** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2973** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2976** - Regex.FindStringSubmatch(line); matches != nil {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2982** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2983** - Item{ (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2985** - = TODOItem{` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2991** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2992** - s, todo) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2994** - s = append(todos, todo)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3000** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3003** - s, err` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3009** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3012** - Type, message string) Priority {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3018** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3019** - Type) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3021** - Type = strings.ToLower(todoType)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3027** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3030** - Comments(t *testing.T) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3036** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3039** - comments` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3045** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3048** - Implement proper error handling` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2391** - /FIXME comments` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3054** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3055** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3057** - This is just a note` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3063** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3066** - comments` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3072** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3073** - Comments(tempDir) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3075** - s, err := analyzer.AnalyzeTODOComments(tempDir)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3081** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3084** - comments: %v", err)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3090** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3091** - , FIXME, HACK, XXX, NOTE (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3093** - s := 5 // TODO, FIXME, HACK, XXX, NOTE` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2388** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3099** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3100** - s { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3102** - s) != expectedTodos {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3108** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3109** - s, len(todos)) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3111** - s, got %d", expectedTodos, len(todos))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3117** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3120** - s` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3126** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3129** - Types := make(map[string]bool)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3135** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3136** - s { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3138** - = range todos {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3144** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3145** - .Type] = true (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2382** - Comments() ([]TODOItem, error)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3153** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3156** - .File != testFile {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3162** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3165** - .File)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3171** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3174** - .Line <= 0 {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3180** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3183** - .Line)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3189** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2380** - Item, error) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3192** - ", "FIXME", "HACK", "XXX", "NOTE"}` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2379** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3198** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3201** - Types[expectedType] {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3207** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3210** - Type string` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3216** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3219** - ", "add feature", PriorityLow},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3225** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3228** - ", "documentation", PriorityLow},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3234** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3237** - Type, test.message)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3243** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3246** - Type, test.message, test.expected, priority)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3252** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3255** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3261** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3264** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3270** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3273** - s) == 0 {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3279** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3282** - comments in test project")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3288** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3291** - Add proper error handling` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2373** - s.md` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3297** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3298** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3300** - This is just a note` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3306** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3309** - Implement proper testing` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3315** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3318** - Add application logic` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3324** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3327** - Add utility functions` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3333** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3336** - analysis` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3342** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3343** - Comments(tempDir) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3345** - s, err := analyzer.AnalyzeTODOComments(tempDir)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3351** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3354** - s: %v", err)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3360** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3363** - s) == 0 {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3369** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3372** - comments")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3378** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3381** - categories and priorities` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3387** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3388** - s { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3390** - = range todos {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3396** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3397** - .Type == "FIXME" { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3399** - .Priority == PriorityHigh && todo.Type == "FIXME" {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2370** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3405** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3408** - s   []string` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3414** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3417** - sAnalyzed     int` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3423** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3426** - sResolved     int` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3432** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3435** - s:   []string{},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3441** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3444** - comments` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3450** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2364** - s.md << EOF` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3453** - /FIXME comments...")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3459** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3460** - Comments(m.projectRoot) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3462** - s, err := m.analyzer.AnalyzeTODOComments(m.projectRoot)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3468** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3471** - comments: %w", err)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3477** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3478** - s (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3480** - s = todos` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3486** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2361** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3489** - /FIXME comments", len(todos))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3495** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3496** - Item (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3498** - s        []TODOItem` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3504** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2355** - s` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3507** - /FIXME comments: %d\n"+` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3513** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3516** - s),` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3522** - _scanner.go:14 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3523** - /FIXME comment analysis (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3525** - Scanner provides enhanced TODO/FIXME comment analysis` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3531** - _scanner.go:15 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3534** - Scanner struct {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3540** - _scanner.go:16 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3543** - ScanConfig` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3549** - _scanner.go:19 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3550** - scanning (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3552** - ScanConfig holds configuration for TODO scanning` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3558** - _scanner.go:20 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3561** - ScanConfig struct {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3567** - _scanner.go:27 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3568** - items (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3570** - Report represents a comprehensive report of all TODO items` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3576** - _scanner.go:28 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3579** - Report struct {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3585** - _scanner.go:33 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3586** - Item (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3588** - s             []TODOItem` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3594** - _scanner.go:34 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3597** - Summary` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3603** - _scanner.go:35 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3606** - Item` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3612** - _scanner.go:36 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3615** - Item` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3621** - _scanner.go:37 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3624** - Item` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3630** - _scanner.go:40 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3633** - Summary provides summary statistics` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3639** - _scanner.go:41 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3642** - Summary struct {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3648** - _scanner.go:42 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3651** - s       int` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3657** - _scanner.go:44 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3660** - s int` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3666** - _scanner.go:45 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3669** - s     int` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3675** - _scanner.go:47 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3678** - s    int` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3684** - _scanner.go:48 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3687** - s        int` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3693** - _scanner.go:49 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3696** - s      int` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3702** - _scanner.go:50 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3705** - s         int` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3711** - _scanner.go:53 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3712** - scanner (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3714** - Scanner creates a new TODO scanner` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3720** - _scanner.go:54 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3721** - ScanConfig) *TODOScanner { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3723** - Scanner(config *TODOScanConfig) *TODOScanner {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3729** - _scanner.go:56 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3732** - ScanConfig()` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3738** - _scanner.go:58 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3741** - Scanner{config: config}` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3747** - _scanner.go:61 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3750** - ScanConfig returns default configuration` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3756** - _scanner.go:62 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3757** - ScanConfig { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3759** - ScanConfig() *TODOScanConfig {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3765** - _scanner.go:63 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3768** - ScanConfig{` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3774** - _scanner.go:71 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3777** - comments` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3783** - _scanner.go:72 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3784** - Report, error) { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3786** - Scanner) ScanProject(rootDir string) (*TODOReport, error) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2352** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3792** - _scanner.go:73 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3795** - Report{` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3801** - _scanner.go:76 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3802** - Item{}, (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3804** - s:             []TODOItem{},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3810** - _scanner.go:77 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3813** - Item),` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3819** - _scanner.go:78 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3822** - Item),` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3828** - _scanner.go:79 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3831** - Item),` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3837** - _scanner.go:82 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3840** - detection` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3846** - _scanner.go:84 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3849** - Regex := regexp.MustCompile(fmt.Sprintf(`(?i)(%s)[\s:]*(.*)`, keywords))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3855** - _scanner.go:105 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3856** - Regex) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3858** - s, err := ts.scanFile(path, todoRegex)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3864** - _scanner.go:110 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3865** - s, todos...) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3867** - s = append(report.TODOs, todos...)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3873** - _scanner.go:113 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3874** - s { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3876** - = range todos {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3882** - _scanner.go:114 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3883** - .Category], todo) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3885** - .Category] = append(report.CategoryBreakdown[todo.Category], todo)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3891** - _scanner.go:115 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3892** - .Priority], todo) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3894** - .Priority] = append(report.PriorityBreakdown[todo.Priority], todo)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3900** - _scanner.go:116 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3901** - .File], todo) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3903** - .File] = append(report.FileBreakdown[todo.File], todo)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3909** - _scanner.go:127 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3912** - s)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3918** - _scanner.go:129 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3921** - s by priority and then by file` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3927** - _scanner.go:130 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3930** - s, func(i, j int) bool {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3936** - _scanner.go:131 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3937** - s[j].Priority { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3939** - s[i].Priority != report.TODOs[j].Priority {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3945** - _scanner.go:132 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3946** - s[j].Priority // Higher priority first (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3948** - s[i].Priority > report.TODOs[j].Priority // Higher priority first` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3954** - _scanner.go:134 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3955** - s[j].File (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3957** - s[i].File < report.TODOs[j].File` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3963** - _scanner.go:140 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3966** - comments` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3972** - _scanner.go:141 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3973** - Regex *regexp.Regexp) ([]TODOItem, error) { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3975** - Scanner) scanFile(filePath string, todoRegex *regexp.Regexp) ([]TODOItem, error) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2346** - s` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3981** - _scanner.go:148 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3982** - Item (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3984** - s []TODOItem` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3990** - _scanner.go:156 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3993** - Regex.FindStringSubmatch(line); matches != nil {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:3999** - _scanner.go:157 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4000** - Item{ (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4002** - = TODOItem{` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4008** - _scanner.go:166 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4009** - s, todo) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4011** - s = append(todos, todo)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4017** - _scanner.go:170 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4020** - s, scanner.Err()` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4026** - _scanner.go:174 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4029** - Scanner) shouldSkipFile(path string) bool {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4035** - _scanner.go:184 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4038** - Scanner) isTextFile(path string) bool {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4044** - _scanner.go:193 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4047** - item` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4053** - _scanner.go:194 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4054** - Type, message, context string) Priority { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4056** - Scanner) determinePriority(todoType, message, context string) Priority {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4062** - _scanner.go:197 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4063** - Type) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4065** - Type = strings.ToLower(todoType)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4071** - _scanner.go:200 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2343** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4074** - Type == "fixme" || todoType == "bug" ||` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2341** - s (556) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4080** - _scanner.go:210 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2323** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4083** - Type == "hack" ||` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4089** - _scanner.go:232 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4092** - item` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4098** - _scanner.go:233 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4101** - Scanner) determineCategory(message, context, filePath string) Category {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4107** - _scanner.go:290 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4108** - s []TODOItem) *TODOSummary { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4110** - Scanner) generateSummary(todos []TODOItem) *TODOSummary {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4116** - _scanner.go:291 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4119** - Summary{}` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4125** - _scanner.go:293 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4126** - s { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4128** - = range todos {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4134** - _scanner.go:294 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4137** - s++` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4143** - _scanner.go:297 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4146** - .Category {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4152** - _scanner.go:301 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4155** - s++` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4161** - _scanner.go:303 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4164** - s++` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4170** - _scanner.go:309 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4173** - .Priority {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4179** - _scanner.go:311 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4182** - s++` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4188** - _scanner.go:313 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4191** - s++` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4197** - _scanner.go:315 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4200** - s++` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4206** - _scanner.go:317 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4209** - s++` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4215** - _scanner.go:325 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4216** - Report) (string, error) { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4218** - Scanner) GenerateReport(report *TODOReport) (string, error) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2314** - _scanner_test.go:147 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4224** - _scanner.go:337 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4225** - Report) (string, error) { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4227** - Scanner) generateMarkdownReport(report *TODOReport) (string, error) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2287** - _scanner.go:224 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4233** - _scanner.go:341 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2269** - _scanner.go:66 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4236** - /FIXME Analysis Report\n\n")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2260** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4242** - _scanner.go:348 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4243** - s)) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4245** - s:** %d\n", report.Summary.TotalTODOs))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4251** - _scanner.go:349 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4254** - s))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4260** - _scanner.go:350 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4263** - s))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4269** - _scanner.go:351 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4272** - s))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4278** - _scanner.go:352 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4281** - s))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4287** - _scanner.go:358 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4290** - s))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4296** - _scanner.go:359 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4299** - s))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4305** - _scanner.go:373 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4308** - s := report.PriorityBreakdown[priority]` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4314** - _scanner.go:374 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4317** - s) == 0 {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4323** - _scanner.go:378 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4324** - s))) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4326** - s (%d)\n\n", priorityNames[priority], len(todos)))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4332** - _scanner.go:380 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4333** - s { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4335** - = range todos {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4341** - _scanner.go:381 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4342** - .Line, todo.Type)) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4344** - .File, todo.Line, todo.Type))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4350** - _scanner.go:382 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4353** - .Message))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4359** - _scanner.go:383 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4362** - .Context))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4368** - _scanner.go:384 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4371** - .Category)))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4377** - _scanner.go:389 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4380** - s` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4386** - _scanner.go:390 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4389** - s\n\n")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4395** - _scanner.go:397 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4398** - s := range report.FileBreakdown {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4404** - _scanner.go:398 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4407** - s)})` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4413** - _scanner.go:409 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4416** - s\n", fc.file, fc.count))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4422** - _scanner.go:416 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4423** - Report) (string, error) { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4425** - Scanner) generateTextReport(report *TODOReport) (string, error) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2251** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4431** - _scanner.go:419 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2242** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4434** - /FIXME Analysis Report\n")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2224** - -scanner/main.go:36 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4440** - _scanner.go:427 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4441** - s)) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4443** - s: %d\n", report.Summary.TotalTODOs))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4449** - _scanner.go:429 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4450** - s, (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4452** - s, report.Summary.HighTODOs,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4458** - _scanner.go:430 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4459** - s)) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4461** - s, report.Summary.LowTODOs))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4467** - _scanner.go:432 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4468** - s { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4470** - = range report.TODOs {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4476** - _scanner.go:434 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4477** - .File, todo.Line, todo.Type)) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4479** - .Priority), todo.File, todo.Line, todo.Type))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4485** - _scanner.go:435 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4488** - .Message))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4494** - _scanner.go:436 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4497** - .Context))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4503** - _scanner.go:443 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4504** - Report) (string, error) { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4506** - Scanner) generateJSONReport(report *TODOReport) (string, error) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2150** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4512** - _scanner.go:452 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4515** - s": %d,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4521** - _scanner.go:462 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4524** - _count": %d` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4530** - _scanner.go:468 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4533** - s,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4539** - _scanner.go:469 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4542** - s,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4548** - _scanner.go:470 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4551** - s,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4557** - _scanner.go:471 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4560** - s,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4566** - _scanner.go:472 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4569** - s,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4575** - _scanner.go:474 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4578** - s,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4584** - _scanner.go:475 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4587** - s,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4593** - _scanner.go:477 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4596** - s)), nil` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4602** - _scanner.go:481 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4605** - Scanner) priorityToString(p Priority) string {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4611** - _scanner.go:496 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4614** - Scanner) categoryToString(c Category) string {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4620** - _scanner_test.go:11 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4623** - Scanner_ScanProject(t *testing.T) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4629** - _scanner_test.go:13 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4632** - _scanner_test")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4638** - _scanner_test.go:19 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4641** - patterns` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4647** - _scanner_test.go:25 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4650** - Add proper error handling` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2132** - _scanner_test.go:146 - HACK (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4656** - _scanner_test.go:37 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4659** - Implement proper authentication` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4665** - _scanner_test.go:39 - XXX (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4668** - This needs immediate attention` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4674** - _scanner_test.go:48 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4677** - Add installation instructions` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4683** - _scanner_test.go:50 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4686** - This should be ignored`,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4692** - _scanner_test.go:66 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4693** - ScanConfig()) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4695** - Scanner(DefaultTODOScanConfig())` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4701** - _scanner_test.go:75 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4704** - s == 0 {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4710** - _scanner_test.go:76 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4713** - comments, but found none")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4719** - _scanner_test.go:80 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4722** - s in vendor/ directory` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4728** - _scanner_test.go:84 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4729** - s { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4731** - = range report.TODOs {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4737** - _scanner_test.go:86 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4740** - .File)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4746** - _scanner_test.go:90 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4749** - .File, "vendor/") {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4755** - _scanner_test.go:91 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4756** - .File) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4758** - in vendor directory, should be skipped: %s", todo.File)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4764** - _scanner_test.go:97 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4767** - s in %s, but didn't", expectedFile)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4773** - _scanner_test.go:103 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4776** - = false` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4782** - _scanner_test.go:104 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4783** - s { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4785** - = range report.TODOs {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4791** - _scanner_test.go:109 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4794** - = true` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4800** - _scanner_test.go:116 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4803** - {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4809** - _scanner_test.go:117 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4812** - ")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4818** - _scanner_test.go:121 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4821** - = false` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4827** - _scanner_test.go:122 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4828** - s { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4830** - = range report.TODOs {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4836** - _scanner_test.go:124 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4839** - = true` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4845** - _scanner_test.go:129 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4848** - {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4854** - _scanner_test.go:130 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4857** - ")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2454** - /FIXME/HACK comments` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4866** - Scanner_DeterminePriority(t *testing.T) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4872** - _scanner_test.go:135 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4873** - ScanConfig()) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4875** - Scanner(DefaultTODOScanConfig())` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4881** - _scanner_test.go:138 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4884** - Type string` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4890** - _scanner_test.go:148 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4893** - ", "add feature", "", PriorityLow},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4899** - _scanner_test.go:149 - NOTE (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4902** - ", "remember to update", "", PriorityLow},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4908** - _scanner_test.go:153 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4911** - Type, test.message, test.context)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4917** - _scanner_test.go:155 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4920** - type '%s' with message '%s', expected priority %v, got %v",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4926** - _scanner_test.go:156 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4929** - Type, test.message, test.expected, priority)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4935** - _scanner_test.go:161 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4938** - Scanner_DetermineCategory(t *testing.T) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4944** - _scanner_test.go:162 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4945** - ScanConfig()) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4947** - Scanner(DefaultTODOScanConfig())` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4953** - _scanner_test.go:192 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4956** - Scanner_GenerateMarkdownReport(t *testing.T) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4962** - _scanner_test.go:193 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4963** - ScanConfig()) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4965** - Scanner(DefaultTODOScanConfig())` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4971** - _scanner_test.go:196 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4974** - Report{` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4980** - _scanner_test.go:201 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4981** - Item{ (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4983** - s: []TODOItem{` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4989** - _scanner_test.go:205 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4992** - ",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:4998** - _scanner_test.go:207 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5001** - Add error handling",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2123** - _scanner_test.go:108 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5007** - _scanner_test.go:221 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5010** - Summary{` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5016** - _scanner_test.go:222 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5019** - s:    2,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5025** - _scanner_test.go:224 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5028** - s:  1,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5034** - _scanner_test.go:225 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5037** - s: 1,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5043** - _scanner_test.go:226 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5046** - s:   1,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5052** - _scanner_test.go:228 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5055** - Item),` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5061** - _scanner_test.go:229 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5064** - Item),` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5070** - _scanner_test.go:230 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5073** - Item),` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5079** - _scanner_test.go:234 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5080** - s { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5082** - = range report.TODOs {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5088** - _scanner_test.go:235 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5089** - .Category], todo) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5091** - .Category] = append(report.CategoryBreakdown[todo.Category], todo)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5097** - _scanner_test.go:236 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5098** - .Priority], todo) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5100** - .Priority] = append(report.PriorityBreakdown[todo.Priority], todo)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5106** - _scanner_test.go:237 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5107** - .File], todo) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5109** - .File] = append(report.FileBreakdown[todo.File], todo)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5115** - _scanner_test.go:248 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2114** - _scanner_test.go:30 - HACK (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5118** - /FIXME Analysis Report",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2096** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5124** - _scanner_test.go:250 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5127** - s:** 2",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5133** - _scanner_test.go:253 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5136** - s",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5142** - _scanner_test.go:265 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5145** - Scanner_ShouldSkipFile(t *testing.T) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5151** - _scanner_test.go:266 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5154** - ScanConfig{` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5160** - _scanner_test.go:269 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5163** - Scanner(config)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5169** - _scanner_test.go:291 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5172** - Scanner_IsTextFile(t *testing.T) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5178** - _scanner_test.go:292 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5181** - ScanConfig{` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5187** - _scanner_test.go:295 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5190** - Scanner(config)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5196** - _scanner_test.go:319 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5199** - Scanner_GenerateSummary(t *testing.T) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5205** - _scanner_test.go:320 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5206** - ScanConfig()) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5208** - Scanner(DefaultTODOScanConfig())` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5214** - _scanner_test.go:322 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5215** - Item{ (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5217** - s := []TODOItem{` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5223** - _scanner_test.go:330 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5226** - s)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5232** - _scanner_test.go:332 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5235** - s != 5 {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5241** - _scanner_test.go:333 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5242** - s) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5244** - s, got %d", summary.TotalTODOs)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5250** - _scanner_test.go:338 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5253** - s != 1 {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5259** - _scanner_test.go:341 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5262** - s != 2 {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5268** - _scanner_test.go:342 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5269** - s) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5271** - s, got %d", summary.HighTODOs)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5277** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2069** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5280** - /FIXME Details` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2013** - s (16) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5286** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5289** - s) > 0 {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5295** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1968** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5298** - /FIXME Comments\n\n")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5304** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5307** - Item)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5313** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5314** - s { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5316** - = range analysis.TODOs {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5322** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5323** - .Category], todo) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5325** - .Category] = append(categories[todo.Category], todo)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5331** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5334** - s := range categories {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5340** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5343** - s)))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5349** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5350** - s { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5352** - = range todos {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5358** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5361** - .Priority)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5367** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5368** - .File, todo.Line, todo.Message)) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5370** - .Type, priority, todo.File, todo.Line, todo.Message))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5376** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5379** - (varies by priority)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5385** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5388** - Time := time.Duration(0)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5394** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5395** - s { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5397** - = range analysis.TODOs {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5403** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5406** - .Priority {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5412** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5415** - Time += 5 * time.Minute` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5421** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5424** - Time += 2 * time.Minute` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5430** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5433** - Time += 1 * time.Minute` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5439** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5442** - Time += 30 * time.Second` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5448** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5451** - Time + duplicateTime + unusedTime + importTime` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5457** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5460** - s` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5466** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5469** - s) > 0 {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5475** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5478** - s",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5484** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1959** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5487** - /FIXME comments",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5493** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5496** - s)) * 1 * time.Minute,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5502** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5503** - s), (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5505** - s(remainingTodos),` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5511** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5514** - s for optimization opportunities", perfCount))` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5520** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5523** - s) + len(analysis.UnusedCode) + len(analysis.Duplicates) + len(analysis.ImportIssues)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:2159** - s (20) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5529** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5530** - s []TODOItem, category Category) []TODOItem { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5532** - sByCategory(todos []TODOItem, category Category) []TODOItem {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5538** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5541** - Item` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5547** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5548** - s { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5550** - = range todos {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5556** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5559** - .Category == category {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5565** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5568** - )` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5574** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5575** - s []TODOItem, category Category) []TODOItem { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5577** - sExcludeCategory(todos []TODOItem, category Category) []TODOItem {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5583** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5586** - Item` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5592** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5593** - s { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5595** - = range todos {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5601** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5604** - .Category != category {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5610** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5613** - )` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5619** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5620** - s { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5622** - = range todos {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5628** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5631** - ",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5637** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5638** - .Message), (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5640** - .Type, todo.Message),` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5646** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5649** - .File,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5655** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5658** - .Line,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5664** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5667** - .Priority,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5673** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5674** - Item{ (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5676** - s: []TODOItem{` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5682** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5685** - ",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5691** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5694** - ",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5700** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1338** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5703** - /FIXME Comments",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5709** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5710** - "}, (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5712** - ", Description: "Fixed TODO"},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1329** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5718** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5719** - "}, (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5721** - ", Description: "Fixed TODO"},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1320** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5727** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5728** - "}, (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5730** - ", Description: "Fixed TODO"},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1311** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5736** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5737** - Item{}, (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5739** - s:        []TODOItem{},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5745** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5748** - s",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5754** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5755** - Item{ (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5757** - s: []TODOItem{` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5763** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5764** - Item{ (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5766** - s: []TODOItem{` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5772** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5775** - ",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5781** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5784** - s + 1 duplicate + 1 unused + 1 import` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5790** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5793** - s(t *testing.T) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5799** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5800** - Item{ (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5802** - s := []TODOItem{` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5808** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5809** - s") (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5811** - s should account for all original TODOs")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5817** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5820** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5826** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5829** - xist(err) || !info.IsDir() {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5835** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5838** - The current parsing implementation is simplified and counts` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5844** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5847** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5853** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5856** - Current validation may not catch path traversal - this is a test for future enhancement` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5862** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5865** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5871** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5874** - xist := []string{` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5880** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5883** - xist {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5889** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5892** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1302** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5901** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5907** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5910** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5916** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5919** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5925** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5928** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5934** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5937** - mpty(t, err.Error(), "Error message should not be empty")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1293** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5943** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5946** - mpty(t, err.Component, "Component should be specified")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5952** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5955** - mpty(t, err.Operation, "Operation should be specified")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5961** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5964** - mpty(t, err.Remediation, "Remediation should be provided")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1284** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5973** - qual(t, randomBytes, randomBytes2, "Generated bytes should be different")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1257** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5982** - mpty(t, err.Remediation, "Error %s should have remediation guidance", err.Error())` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1239** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:5991** - mpty(t, id1)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1221** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6000** - mpty(t, id2)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1212** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6009** - qual(t, id1, id2, "Generated IDs should be unique")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1185** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6018** - mpty(t, result.Errors, "Should have validation errors")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1176** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6027** - mpty(t, result.Errors, "Should have validation errors")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1167** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6036** - mpty(t, result.Warnings, "Should have compatibility warnings")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1158** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6045** - mpty(t, result.Warnings, "Should have warnings for insecure permissions")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1149** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6054** - mpty(t, result.Errors, "Should have errors for dangerous directory")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1140** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6063** - mpty(t, result.Warnings, "Should have warnings for weak random length")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1132** - s) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6072** - mpty(t, result.Warnings, "Should have warnings for low entropy")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1131** - _scanner_test.go:339 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6081** - mpty(t, result.Warnings, "Should have warnings for alphanumeric format")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1122** - _scanner_test.go:336 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6090** - mpty(t, result.Warnings, "Should have warnings for disabled entropy check")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1113** - _scanner_test.go:335 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6099** - This is testing the logic, not the validator framework integration` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6105** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6108** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6114** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6117** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6123** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6126** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6132** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6135** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6141** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6144** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1104** - _scanner_test.go:327 - BUG (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6153** - xpectedContent []string` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1095** - _scanner_test.go:254 - FIXME (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6162** - xpectedContent: []string{`"none"`},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1086** - _scanner_test.go:223 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6171** - xpectedContent: []string{},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1077** - _scanner_test.go:216 - FIXME (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6180** - xpectedContent: []string{},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1068** - _scanner_test.go:214 - FIXME (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6189** - xpectedContent: []string{},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1059** - _scanner_test.go:178 - BUG (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6198** - xpectedContent: []string{},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1050** - _scanner_test.go:172 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6207** - xpectedContent: []string{},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1041** - _scanner_test.go:145 - BUG (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6214** - xpectedContent { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6216** - xpected := range tt.notExpectedContent {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1032** - _scanner_test.go:144 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6225** - xpected) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1023** - _scanner_test.go:143 - FIXME (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6234** - xpected, result)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1014** - _scanner_test.go:123 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6243** - xpected []string` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1005** - _scanner_test.go:114 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6252** - xpected: []string{`"none"`},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:996** - _scanner_test.go:113 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6261** - xpected: []string{},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:987** - _scanner_test.go:106 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6270** - xpected: []string{},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:978** - _scanner_test.go:105 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6277** - xpected { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6279** - xpected := range tt.notExpected {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:969** - _scanner_test.go:102 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6288** - xpected) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:960** - _scanner_test.go:79 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6297** - xpectedContent []string` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:951** - _scanner_test.go:49 - FIXME (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6306** - xpectedContent: []string{`"null"`},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:942** - _scanner_test.go:42 - BUG (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6315** - xpectedContent: []string{`'null'`},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:933** - _scanner_test.go:28 - FIXME (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6324** - xpectedContent: []string{},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:924** - _scanner.go:505 - BUG (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6333** - xpectedContent: []string{},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:915** - _scanner.go:504 - BUG (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6342** - xpectedContent: []string{},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:907** - s, (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6349** - xpectedContent { (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6351** - xpected := range tt.notExpectedContent {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:906** - _scanner.go:476 - BUG (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6360** - xpected) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:897** - _scanner.go:473 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6369** - xpected, result)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:888** - _scanner.go:460 - BUG (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6378** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:879** - _scanner.go:360 - BUG (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6387** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:870** - _scanner.go:357 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6396** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:862** - s++ (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6405** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:861** - _scanner.go:305 - BUG (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6414** - crypto/rand has different API - use rand.Read([]byte) instead of rand.Int()"` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:852** - _scanner.go:304 - BUG (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6423** - Replace with secure random generation"` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:843** - _scanner.go:299 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6432** - Replace with secure ID generation"` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:834** - _scanner.go:274 - BUG (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6441** - Replace with secure temp file creation"` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:825** - _scanner.go:269 - BUG (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6450** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:816** - _scanner.go:268 - BUG (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6459** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:808** - s         int (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6468** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:807** - _scanner.go:46 - BUG (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6477** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:798** - _scanner.go:43 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6486** - "` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:789** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6495** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:781** - .Message) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6504** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:771** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6513** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:744** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6522** - Some of these may not be detected by current patterns` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:690** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6531** - Some fixes may not apply to all formats (e.g., YAML/JSON)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6537** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6540** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6546** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6549** - xist": "os",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6555** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6558** - ":         "context",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6564** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6567** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6573** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6576** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6582** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6585** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6591** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6594** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6600** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6603** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6609** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6612** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6618** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6621** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6627** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6630** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6636** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6639** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6645** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6648** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6654** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6657** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6663** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6666** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6672** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6675** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6681** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6684** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6690** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6693** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6699** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6702** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6708** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6711** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6717** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6720** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6726** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6729** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6735** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6738** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6744** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6747** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6753** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6756** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6762** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6765** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6771** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6774** - Add time.Sleep() here later` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6780** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6783** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6789** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6792** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6798** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6801** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6807** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6810** - mpty(t, projectName)` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6816** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6819** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6825** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6828** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6834** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6837** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6843** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6846** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6852** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6855** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6861** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6864** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6870** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6873** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6879** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6882** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6888** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6891** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6897** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6900** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6906** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6909** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6915** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6918** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6924** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6927** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6933** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6936** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6942** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6945** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6951** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6954** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6960** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6963** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6969** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6972** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6978** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6981** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6987** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6990** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6996** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:6999** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7005** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7008** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7014** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7017** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7023** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7026** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7032** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7033** - s,omitempty"` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7035** - s        string            `json:"notes,omitempty"`` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7041** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7044** - s as warnings` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7050** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7053** - s != "" {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7059** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7062** - s,` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7068** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7071** - s: "Next.js 15 requires React 18 or later",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7077** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7080** - s: "Next.js 14 requires React 18 or later",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7086** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7089** - s: "Gin v1.9+ requires Go 1.19 or later",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7095** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7098** - s: "Kotlin 2.0 requires Android Gradle Plugin 8.0+",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7104** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7107** - s: "React 18 requires React DOM 18+",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7113** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7116** - s...",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7122** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7125** - Current implementation only returns latest version` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7131** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7134** - LatestVersion will be whatever the mock server returns` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:627** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7143** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7149** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7152** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7158** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7161** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7167** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7170** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7176** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7179** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7185** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7188** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7194** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7197** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7203** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7206** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7212** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7215** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7221** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7224** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7230** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7233** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7239** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7242** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7248** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7251** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7257** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7260** - Template restoration will be fully implemented in a future version")` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7266** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7269** - that this check should be done` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7275** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7278** - xist": "os",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7284** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7287** - ":         "context",` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7293** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7296** - We're not returning the error here because it would stop the inspection` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:232** - s, (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7302** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7305** - ", "WithCancel", "WithTimeout", "WithValue"},` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7311** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7314** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7320** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7323** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7329** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7332** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7338** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7341** - xist(err) {` (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7347** - s (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7349** - _scanner.go:** 141 TODOs (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7350** - _scanner_test.go:** 102 TODOs (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7351** - s (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7352** - s (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7353** - s (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7354** - s (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7355** - s (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7356** - s (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:231** - -scanner/main.go:85 - TODO (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:7358** - s (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:1** - /FIXME Analysis Report (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:9** - s:** 813 (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:22** - s (221) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:33** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:42** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:51** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:222** - -scanner/main.go:84 - BUG (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:60** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-analysis-report.md:78** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:7** - _scanner.go` - Enhanced TODO scanner with categorization (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:41** -  (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:14** - /FIXME Comments Found:** 813 (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:4** - /FIXME comment scanner for the project codebase. (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:8** - -scanner/main.go` - Command-line tool for running scans (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:31** - Replace with secure random generation" (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:32** - Replace with secure ID generation" (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:33** - Replace with secure temp file creation" (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:23** - s (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:24** - s (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:9** - _scanner_test.go` - Comprehensive test suite (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:44** - -analysis-report.md` (comprehensive markdown report) (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:50** - -scanner (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:53** - -scanner -output report.md (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:56** - -scanner -verbose (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:59** - -scanner -format json (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:60** - -scanner -format text (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:64** - -style comments across the codebase (Documentation or specification file)
- **/Users/inertia/cuesoft/working-area/todo-scan-summary.md:1** - /FIXME Scan Summary - Task 2.1 Completion (Documentation or specification file)

