# Janus - Cross-Domain Project Scaffolder

A CLI tool that works like `create-next-app` but supports multiple domains (web, backend, Python, etc.) and works completely offline using local templates.

## 🚀 Quick Start

```bash
# Build the CLI
go build -o janus main.go

# Make it executable (Linux/Mac)
chmod +x janus

# Move to PATH (optional)
mv janus /usr/local/bin/janus
```

## 📖 Usage

```bash
# Scaffold a new project
janus <domain> [subdomain]

# List available offline templates
janus -o templates

# Show help
janus -h

# Show version
janus -v
```

## 🎯 Examples

```bash
# Create a Next.js project
janus nextjs
# Creates ./nextjs-app/

# Create a Python Flask project
janus python flask
# Creates ./python-flask-app/

# Create a Node.js Express project
janus node express
# Creates ./node-express-app/

# List all available templates
janus -o templates
```

## 📁 Template Structure

Templates are stored at: `~/.janus/templates/`

Example structure:
```
~/.janus/templates/
├── nextjs/
│   ├── package.json
│   ├── pages/
│   └── components/
├── python-flask/
│   ├── requirements.txt
│   ├── app.py
│   └── templates/
└── node-express/
    ├── package.json
    ├── server.js
    └── routes/
```

## ✨ Features

- **Offline First**: Works completely offline with local templates
- **Cross-Platform**: Single binary works on Linux, Mac, Windows
- **Smart Dependencies**: Auto-detects and installs dependencies
  - `package.json` → runs `npm install`
  - `requirements.txt` → runs `pip install`
  - `go.mod` → runs `go mod tidy`
  - `Cargo.toml` → runs `cargo fetch`
- **Safe Operations**: Checks for existing directories before scaffolding
- **Auto Git Init**: Initializes git repository if git is available
- **No External Dependencies**: Uses only Go standard library

## 🔧 How It Works

1. **Template Storage**: Templates are stored locally at `~/.janus/templates/`
2. **Scaffolding**: Copies template folder to `./[domain]-[subdomain]-app`
3. **Dependencies**: Automatically installs dependencies if package manager files exist
4. **Git Init**: Initializes git repository for version control

## 📦 Installation

### From Source
```bash
git clone <repository-url>
cd janus
go build -o janus main.go
```

### Create Templates Directory
```bash
mkdir -p ~/.janus/templates
```

### Add Your Templates
Copy your project templates to `~/.janus/templates/[template-name]/`

## 🎨 Creating Templates

1. Create a directory under `~/.janus/templates/`
2. Add your project boilerplate files
3. Include dependency files (`package.json`, `requirements.txt`, etc.)
4. Use descriptive names like `nextjs`, `python-flask`, `node-express`

Example template creation:
```bash
mkdir -p ~/.janus/templates/nextjs
cd ~/.janus/templates/nextjs
# Add your Next.js boilerplate files here
```

## 🔍 Supported Project Types

Janus can scaffold any type of project. Common examples:

- **Web**: `nextjs`, `react`, `vue`, `angular`
- **Backend**: `node-express`, `python-flask`, `python-django`, `go-gin`
- **Mobile**: `react-native`, `flutter`
- **Desktop**: `electron`, `tauri`
- **Data**: `python-jupyter`, `r-project`

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📝 License

MIT License - see LICENSE file for details

## 🐛 Troubleshooting

### Templates not found
- Check if `~/.janus/templates/` exists
- Verify template names match exactly
- Run `janus -o templates` to list available templates

### Dependency installation fails
- Ensure package managers are installed (`npm`, `pip`, `go`, `cargo`)
- Check internet connection for dependency downloads
- Verify dependency files are correct (`package.json`, `requirements.txt`)

### Permission errors
- Ensure write permissions in current directory
- On Linux/Mac, make binary executable with `chmod +x janus`

## 📊 Roadmap

- [ ] Interactive template selection
- [ ] Online template download fallback
- [ ] Template versioning
- [ ] Custom template variables/placeholders
- [ ] Template validation
- [ ] Configuration file support
