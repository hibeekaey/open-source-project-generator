package generators

import (
	"fmt"
	"path/filepath"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// DocumentationGenerator handles documentation file generation
type DocumentationGenerator struct {
	fsOps FileSystemOperationsInterface
}

// NewDocumentationGenerator creates a new documentation generator
func NewDocumentationGenerator(fsOps FileSystemOperationsInterface) *DocumentationGenerator {
	return &DocumentationGenerator{
		fsOps: fsOps,
	}
}

// GenerateDocumentationFiles creates documentation files
func (dg *DocumentationGenerator) GenerateDocumentationFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate CONTRIBUTING.md
	contributingContent := fmt.Sprintf(`# Contributing to %s

Thank you for your interest in contributing to %s! This document provides guidelines and information for contributors.

## Development Setup

1. Clone the repository:
   `+"```bash"+`
   git clone %s
   cd %s
   `+"```"+`

2. Set up the development environment:
   `+"```bash"+`
   make setup
   `+"```"+`

3. Start development servers:
   `+"```bash"+`
   make dev
   `+"```"+`

## Code Style

- Follow the existing code style in each component
- Run linting before submitting: `+"```bash"+`make lint`+"```"+`
- Ensure all tests pass: `+"```bash"+`make test`+"```"+`

## Pull Request Process

1. Fork the repository
2. Create a feature branch: `+"```bash"+`git checkout -b feature/your-feature`+"```"+`
3. Make your changes and add tests
4. Ensure all tests pass and code is properly formatted
5. Submit a pull request with a clear description

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/).

## License

By contributing to %s, you agree that your contributions will be licensed under the %s License.
`, config.Name, config.Name, "", config.Name, config.Name, config.License)

	contributingPath := filepath.Join(projectPath, "CONTRIBUTING.md")
	if err := dg.fsOps.WriteFile(contributingPath, []byte(contributingContent), 0644); err != nil {
		return fmt.Errorf("failed to create CONTRIBUTING.md: %w", err)
	}

	// Generate SECURITY.md
	securityContent := fmt.Sprintf(`# Security Policy

## Supported Versions

We provide security updates for the following versions of %s:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability in %s, please report it responsibly:

1. **Do not** create a public GitHub issue for security vulnerabilities
2. Email us at: security@%s.com (replace with your actual security contact)
3. Include detailed information about the vulnerability
4. Allow us time to address the issue before public disclosure

## Security Best Practices

When contributing to %s:

- Keep dependencies up to date
- Follow secure coding practices
- Use environment variables for sensitive configuration
- Validate all user inputs
- Use HTTPS for all external communications

## Response Timeline

- We will acknowledge receipt of vulnerability reports within 48 hours
- We aim to provide an initial assessment within 7 days
- We will work to resolve critical vulnerabilities within 30 days

Thank you for helping keep %s secure!
`, config.Name, config.Name, config.Organization, config.Name, config.Name)

	securityPath := filepath.Join(projectPath, "SECURITY.md")
	if err := dg.fsOps.WriteFile(securityPath, []byte(securityContent), 0644); err != nil {
		return fmt.Errorf("failed to create SECURITY.md: %w", err)
	}

	// Generate LICENSE file
	licenseContent := dg.generateLicenseContent(config)
	licensePath := filepath.Join(projectPath, "LICENSE")
	if err := dg.fsOps.WriteFile(licensePath, []byte(licenseContent), 0644); err != nil {
		return fmt.Errorf("failed to create LICENSE: %w", err)
	}

	return nil
}

// generateLicenseContent generates license content based on the selected license
func (dg *DocumentationGenerator) generateLicenseContent(config *models.ProjectConfig) string {
	year := "2024"
	author := config.Author
	if author == "" {
		author = config.Organization
	}

	switch config.License {
	case "MIT":
		return fmt.Sprintf(`MIT License

Copyright (c) %s %s

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.`, year, author)

	case "Apache-2.0":
		return fmt.Sprintf(`Apache License
Version 2.0, January 2004
http://www.apache.org/licenses/

Copyright %s %s

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.`, year, author)

	case "GPL-3.0":
		return fmt.Sprintf(`GNU GENERAL PUBLIC LICENSE
Version 3, 29 June 2007

Copyright (C) %s %s

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.`, year, author)

	case "BSD-3-Clause":
		return fmt.Sprintf(`BSD 3-Clause License

Copyright (c) %s, %s
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
   contributors may be used to endorse or promote products derived from
   this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.`, year, author)

	default:
		return fmt.Sprintf(`Copyright (c) %s %s

All rights reserved.`, year, author)
	}
}
