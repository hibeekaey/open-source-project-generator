# Documentation

Welcome to the Open Source Project Generator documentation. This tool uses a **tool-orchestration architecture** that delegates project creation to industry-standard CLI tools.

## 📚 Documentation

### For Users

- **[Getting Started](GETTING_STARTED.md)** - Installation, quick start, and basic usage
- **[CLI Commands](CLI_COMMANDS.md)** - Complete command reference
- **[Configuration Guide](CONFIGURATION.md)** - Configuration file format and options
- **[Examples](EXAMPLES.md)** - Real-world project configurations
- **[Troubleshooting](TROUBLESHOOTING.md)** - Common issues and solutions

### For Developers

- **[Architecture](ARCHITECTURE.md)** - System design and architecture overview
- **[Adding Bootstrap Tools](ADDING_TOOLS.md)** - Guide for adding new tool support
- **[API Reference](API_REFERENCE.md)** - Developer API documentation
- **[Contributing](../CONTRIBUTING.md)** - How to contribute to the project

## 🚀 Quick Links

| I want to... | Go to... |
|--------------|----------|
| Install and use the generator | [Getting Started](GETTING_STARTED.md) |
| See all available commands | [CLI Commands](CLI_COMMANDS.md) |
| Create a configuration file | [Configuration Guide](CONFIGURATION.md) |
| See example projects | [Examples](EXAMPLES.md) |
| Fix an issue | [Troubleshooting](TROUBLESHOOTING.md) |
| Understand the architecture | [Architecture](ARCHITECTURE.md) |
| Add support for a new tool | [Adding Bootstrap Tools](ADDING_TOOLS.md) |
| Contribute code | [Contributing](../CONTRIBUTING.md) |

## 🎯 What is Tool-Orchestration?

Instead of maintaining templates manually, this generator:

1. **Discovers** available bootstrap tools on your system (like `npx`, `go`, `gradle`)
2. **Executes** these tools to generate projects using their official CLIs
3. **Maps** the generated output to a standardized directory structure
4. **Integrates** components together with Docker Compose, scripts, etc.

**Benefits:**

- ✅ Always up-to-date dependencies (no manual template maintenance)
- ✅ Industry-standard project structures
- ✅ Leverages community expertise
- ✅ Graceful fallback when tools unavailable
- ✅ Offline support with caching

## 📖 Documentation Structure

```text
docs/
├── README.md                   # This file - Documentation index
├── GETTING_STARTED.md          # Installation, quick start, tool requirements
├── CLI_COMMANDS.md             # Complete CLI command reference
├── CONFIGURATION.md            # Configuration file format and options
├── EXAMPLES.md                 # Real-world configuration examples
├── TROUBLESHOOTING.md          # Common issues and solutions
├── ARCHITECTURE.md             # System architecture and design
├── ADDING_TOOLS.md             # Guide for adding new bootstrap tools
├── API_REFERENCE.md            # Developer API documentation
└── DOCUMENTATION_SUMMARY.md    # Documentation overview and statistics
```

### Documentation Coverage

- ✅ **9 documentation files** covering all aspects
- ✅ **150+ pages** of comprehensive documentation
- ✅ **10+ examples** for common project types
- ✅ **All 5 CLI commands** fully documented
- ✅ **All 4 component types** explained
- ✅ **10 core components** architecture documented

## 🆘 Getting Help

- **Documentation**: Start with [Getting Started](GETTING_STARTED.md)
- **Issues**: [GitHub Issues](https://github.com/cuesoftinc/open-source-project-generator/issues)
- **Discussions**: [GitHub Discussions](https://github.com/cuesoftinc/open-source-project-generator/discussions)
- **Email**: <support@cuesoft.io>

---

**Ready to start?** Begin with the [Getting Started Guide](GETTING_STARTED.md)!
