# Template TODO Documentation

## Overview

This document explains the intentional TODO comments found in template files. These are **NOT** issues that need to be resolved, but rather placeholders and guidance for developers using the templates.

## Template TODOs by Category

### 1. Android Mobile Template

#### Data Extraction Rules Configuration

**File**: `templates/mobile/android-kotlin/app/src/main/res/xml/data_extraction_rules.xml.tmpl`
**Line**: 4
**TODO**: `<!-- TODO: Use <include> and <exclude> to control what is backed up.`

**Purpose**: Guides developers to configure Android's data extraction rules for cloud backup and device transfer.

**Implementation Guidance**:

- Developers should uncomment and configure the `<include>` and `<exclude>` tags
- Specify which app data should be backed up to the cloud
- Define what data should be transferred between devices
- See [Android Data Extraction Rules](https://developer.android.com/guide/topics/data/autobackup#IncludingFiles) for details

### 2. Go Gin Backend Template

#### Email Service Integration

**File**: `templates/backend/go-gin/internal/services/auth_service.go.tmpl`
**Line**: 188
**TODO**: `// TODO: Send email with reset token`

**Purpose**: Placeholder for email service integration in password reset functionality.

**Implementation Guidance**:

- Integrate with email service provider (SendGrid, AWS SES, etc.)
- Implement email template for password reset
- Add email configuration to application config
- Handle email sending errors appropriately
- Consider rate limiting for password reset emails

**Example Implementation**:

```go
// Send password reset email
emailService := s.emailService // Inject email service
err = emailService.SendPasswordResetEmail(user.Email, resetToken)
if err != nil {
    // Log error but don't reveal to user for security
    s.logger.Error("Failed to send password reset email", "error", err)
}
```

### 3. Pull Request Template

#### Code Quality Checklist

**File**: `templates/base/.github/PULL_REQUEST_TEMPLATE.md.tmpl`
**Line**: 244
**TODO**: `- [ ] No TODO comments without issues`

**Purpose**: Checklist item to ensure developers don't leave untracked TODO comments in their code.

**Implementation Guidance**:

- All TODO comments should be linked to GitHub issues
- Use format: `// TODO(#123): Description` where 123 is the issue number
- Remove TODO comments when implementing the feature
- Document any intentional TODOs (like these template ones)

## Template TODO Best Practices

### For Template Maintainers

1. **Clear Documentation**: Always document why a TODO exists in templates
2. **Implementation Guidance**: Provide clear instructions on how to resolve the TODO
3. **Examples**: Include code examples where helpful
4. **Links**: Reference relevant documentation or guides

### For Template Users

1. **Review All TODOs**: When using a template, review all TODO comments
2. **Implement Required Features**: Address TODOs that are required for your use case
3. **Remove Completed TODOs**: Remove TODO comments after implementation
4. **Track Remaining TODOs**: Create issues for TODOs you plan to address later

## Template TODO Format

For consistency, template TODOs should follow this format:

```
// TEMPLATE_TODO: Brief description
// Implementation: Detailed guidance on how to implement
// Documentation: Link to relevant docs
// Example: Code example if helpful
```

## Maintenance

This documentation should be updated when:

- New template TODOs are added
- Existing template TODOs are modified
- Templates are removed or deprecated
- Implementation guidance changes

## Related Documentation

- [Template Development Guide](docs/TEMPLATE_MAINTENANCE.md)
- [Template Security Guide](docs/TEMPLATE_SECURITY_GUIDE.md)
- [Template Validation Checklist](docs/TEMPLATE_VALIDATION_CHECKLIST.md)
