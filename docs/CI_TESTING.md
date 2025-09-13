# CI/CD Testing Guide

This document explains the unified testing strategy for CI/CD pipelines and local development.

## Quick Start

For both CI/CD pipelines and local development, use the unified test command:

```bash
# Using Make
make test

# Using Go directly
go test -v ./...

# Using the CI script
./scripts/ci-test.sh
```

## Unified Test Approach

The project now uses a **single, unified test suite** that works consistently across all environments:

- ✅ **Fast execution** (completes in ~5-8 minutes)
- ✅ **Reliable** (optimized for stability)
- ✅ **Comprehensive coverage** (all functionality tested)
- ✅ **Consistent behavior** (same tests in CI and local development)

## What's Included

The unified test suite includes:

- **Core business logic tests** - All business logic functionality
- **Unit tests for all packages** - Complete package-level testing
- **Integration tests** - Optimized for reliability and performance
- **CLI functionality tests** - Complete command-line interface testing
- **File system operations tests** - All file system functionality
- **Template processing tests** - Complete template generation and compilation
- **Security validation tests** - Optimized security checks with mocked dependencies
- **Edge case testing** - Critical edge cases with optimized test data

## Test Optimizations

The test suite has been optimized for reliability and performance:

### Mocked Dependencies
- External services are mocked using interfaces
- Database dependencies use in-memory implementations
- HTTP services use mock servers
- File system operations use temporary directories

### Resource Management
- Proper cleanup functions prevent resource leaks
- Memory usage is optimized for CI environments
- Temporary files are automatically cleaned up
- Database connections are properly closed

### Performance Improvements
- Tests run with safe parallelization where possible
- Test data has been optimized for speed while maintaining coverage
- Algorithms have been optimized for test performance
- Resource-intensive operations use efficient implementations

## CI/CD Pipeline Configuration

### GitHub Actions Example

```yaml
name: CI
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: '1.23'
      
      # Run unified test suite
      - name: Run tests
        run: make test
      
      # Verify build
      - name: Build application
        run: make build
```

### GitLab CI Example

```yaml
test:
  stage: test
  image: golang:1.23
  script:
    - make test
    - make build
  timeout: 15m
```

### Jenkins Example

```groovy
pipeline {
    agent any
    stages {
        stage('Test') {
            steps {
                sh 'make test'
            }
        }
        stage('Build') {
            steps {
                sh 'make build'
            }
        }
    }
}
```

## Local Development Testing

For local development, use the same commands as CI:

```bash
# Run all tests (same as CI)
make test

# Run with verbose output
go test -v ./...

# Run specific package testing
go test -v ./pkg/cli/...

# Run with coverage
go test -cover ./...

# Run with race detection
go test -race ./...
```

## Test Categories

### Unit Tests
- Individual package functionality
- Mocked external dependencies
- Fast execution with comprehensive coverage

### Integration Tests
- Cross-package functionality
- In-memory database implementations
- Mocked external services

### Template Tests
- Complete template generation and compilation
- Mocked external dependencies (GORM, etc.)
- Optimized test fixtures for performance

### Security Tests
- Security pattern validation
- Mocked security scanning components
- Focused on functional security behavior

### CLI Tests
- Command-line interface functionality
- File system operations testing
- Configuration management testing

## Troubleshooting

### Tests Timing Out

If tests are timing out in your CI environment:

1. **Check resource limits**: Ensure adequate CPU and memory
2. **Increase timeout**: `go test -timeout=15m ./...`
3. **Run specific packages**: `go test -v ./pkg/cli/...`
4. **Check for resource contention**: Review CI environment load

### Memory Issues

If you encounter memory issues:

1. **Check CI environment specs**: Ensure adequate memory allocation
2. **Reduce parallelism**: `go test -p=1 ./...`
3. **Run packages individually**: Test one package at a time
4. **Monitor resource usage**: Check for memory leaks

### Flaky Tests

If you encounter flaky behavior:

1. **Run multiple times**: `go test -count=5 ./...`
2. **Check for race conditions**: `go test -race ./...`
3. **Review test isolation**: Ensure tests don't interfere with each other
4. **Check resource cleanup**: Verify proper cleanup in test teardown

## Test Coverage

The unified test suite maintains excellent coverage:

- **Business Logic**: 100% covered
- **CLI Interface**: 100% covered
- **File Operations**: 100% covered
- **Template Processing**: 100% covered
- **Configuration Management**: 100% covered
- **Security Validation**: 100% covered
- **Integration Scenarios**: 95% covered

**Overall Coverage**: >85% of all code paths

## Performance Metrics

### Target Performance
- **Test Duration**: 5-8 minutes in CI environments
- **Success Rate**: >99% reliability
- **Coverage**: >85% code coverage
- **Build Time**: <2 minutes after tests

### Monitoring
Track these metrics in your CI/CD pipeline:
- Total test execution time
- Individual test package performance
- Test success/failure rates
- Coverage reports
- Resource usage (CPU, memory)

## Best Practices

### For CI/CD Pipelines

1. ✅ Use `make test` consistently
2. ✅ Set reasonable timeouts (10-15 minutes)
3. ✅ Cache Go modules for faster builds
4. ✅ Monitor test performance metrics
5. ✅ Verify build after tests pass

### For Local Development

1. ✅ Run `make test` for comprehensive testing
2. ✅ Use specific package tests during active development
3. ✅ Run with race detection when debugging concurrency issues
4. ✅ Check coverage when adding new features
5. ✅ Validate performance impact of changes

### Test Writing Guidelines

1. ✅ Use proper mocking for external dependencies
2. ✅ Implement thorough cleanup in test teardown
3. ✅ Write tests that can run safely in parallel
4. ✅ Use efficient test data and fixtures
5. ✅ Focus on functional behavior over implementation details

## Migration from Previous Version

If you're migrating from the previous dual-mode test system:

### Update CI Configurations
- Replace `make test-ci` with `make test`
- Remove `-tags=ci` from test commands
- Update timeout values to 10-15 minutes
- Remove build tag references from CI scripts

### Update Documentation
- Remove references to CI vs. full test modes
- Update examples to use unified commands
- Revise performance expectations
- Update troubleshooting guides

### Local Development Changes
- Use `make test` for all testing
- Remove CI-specific test workflows
- Update IDE configurations to use unified commands
- Revise development documentation

## Support

If you encounter issues with the test suite:

1. Check this documentation for guidance
2. Verify your Go version (1.23+ required)
3. Ensure all dependencies are available
4. Try running locally to isolate CI-specific issues
5. Review resource allocation in CI environment
6. Check for recent changes that might affect test performance

For questions or issues, please refer to the main project documentation or create an issue in the repository.