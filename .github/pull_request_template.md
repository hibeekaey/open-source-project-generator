# Pull Request

## Description

<!-- Provide a brief description of the changes in this PR -->

## Type of Change

<!-- Mark the relevant option with an "x" -->

- [ ] 🐛 Bug fix (non-breaking change which fixes an issue)
- [ ] ✨ New feature (non-breaking change which adds functionality)
- [ ] 💥 Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] 📚 Documentation update
- [ ] 🔧 Configuration change
- [ ] 🧹 Code cleanup/refactoring
- [ ] ⚡ Performance improvement
- [ ] 🔒 Security improvement
- [ ] 🧪 Test improvement

## Related Issues

<!-- Link to related issues using "Fixes #123" or "Closes #123" or "Related to #123" -->
<!-- Remove unused lines below -->

- Fixes #
- Closes #
- Related to #

## Changes Made

<!-- Provide a detailed list of changes made -->

### Added

-

### Changed

-

### Removed

-

### Fixed

-

## Testing

<!-- Describe the tests you ran to verify your changes -->

- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed
- [ ] All existing tests pass

### Test Coverage
<!-- If applicable, mention test coverage changes -->

- [ ] Test coverage maintained or improved
- [ ] New tests added for new functionality
- [ ] Existing tests updated if needed

### Build Verification

- [ ] `make build` succeeds
- [ ] `make test` passes
- [ ] `make lint` passes
- [ ] `make security-scan` passes
- [ ] `make pre-commit` passes (optional but recommended)
- [ ] `make test-coverage` shows adequate coverage

## Screenshots/Demo

<!-- If applicable, add screenshots or demo links -->

## Checklist

<!-- Mark completed items with an "x" -->

### Code Quality

- [ ] Code follows the project's style guidelines
- [ ] Self-review of code completed
- [ ] Code is properly commented
- [ ] No hardcoded values or magic numbers
- [ ] Error handling is appropriate
- [ ] Logging is appropriate

### Architecture

- [ ] Dependencies flow inward (presentation → business → infrastructure)
- [ ] Interfaces defined in `pkg/interfaces/` (if new component)
- [ ] Error handling uses `pkg/errors/` categorization (SecurityError, ValidationError, etc.)
- [ ] No global state or singletons introduced
- [ ] Docker containers use UID 1001 (if Docker changes)

### Documentation

- [ ] README updated (if needed)
- [ ] API documentation updated (if needed)
- [ ] Code comments added/updated
- [ ] Changelog updated (if needed)

### Security

- [ ] No sensitive information exposed
- [ ] Input validation implemented
- [ ] User paths sanitized via `pkg/security/SanitizePath()`
- [ ] Errors use categorized types from `pkg/errors/`
- [ ] Security best practices followed (see SECURITY.md)
- [ ] Dependencies are secure and up-to-date
- [ ] `make security-scan` passes

### Performance

- [ ] No performance regressions introduced
- [ ] Memory usage is appropriate
- [ ] Database queries are optimized (if applicable)
- [ ] Caching is implemented where appropriate

### Compatibility

- [ ] Backward compatibility maintained
- [ ] Cross-platform compatibility verified
- [ ] Browser compatibility verified (if applicable)

### Template Changes (if applicable)

- [ ] Template metadata (`template.yaml`) updated
- [ ] Template tested with generation
- [ ] Template security validation passes
- [ ] Template documentation updated

## Breaking Changes

<!-- If this PR contains breaking changes, describe them here -->

## Migration Guide

<!-- If applicable, provide migration instructions -->

## Additional Notes

<!-- Any additional information that reviewers should know -->

## Reviewers

<!-- Tag specific reviewers if needed: @username -->

## Deployment Notes

<!-- Any special deployment considerations -->

## Rollback Plan

<!-- How to rollback if this change causes issues -->

---

**By submitting this pull request, I confirm that:**

- [ ] I have read and agree to the [Code of Conduct](CODE_OF_CONDUCT.md)
- [ ] I have performed a self-review of my own code
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings or errors
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] Any dependent changes have been merged and published
