# Simple Git

A simplified Git implementation built from scratch in Go for educational purposes. This project demonstrates the core concepts and internal workings of Git's version control system.

## Why Build This?

> *"What I cannot create, I do not understand."* - Richard Feynman

Building Git from scratch helps us understand:

- **How version control actually works** under the hood
- **Content-addressable storage** and why it's revolutionary
- **The elegance of Git's object model** (blobs, trees, commits)
- **How branches are just lightweight pointers** to commits
- **Why Git is so fast and efficient**

This isn't about reinventing the wheel but about **deeply understanding** how the wheel works!

## Architecture Overview

Simple Git implements the core Git concepts with a clean, readable codebase:

### Core Components
1. **Object Database** (`objects.go`)
   - Content-addressable storage using SHA-1 hashing
   - Three object types: blobs (files), trees (directories), commits (snapshots)
   - Compressed storage using zlib

2. **Index/Staging Area** (`index.go`)
   - Tracks which files will be included in the next commit
   - Binary format similar to Git's actual index
   - Efficient change detection

3. **References** (`refs.go`)
   - Branches as lightweight pointers to commits
   - HEAD pointer for current branch/commit tracking
   - Simple branch creation and switching

4. **Commands** (`commands.go`)
   - User-facing Git commands (init, add, commit, etc.)
   - Orchestrates the interaction between components

### How It All Works Together

```
Working Directory → Index (staging) → Object Database
       ↑                                      ↓
       └─────── References (branches) ────────┘
```

## User Journey Example

Let's trace through what happens when you use Simple Git:

### 1. Initialize Repository
```bash
simple-git init
```
- Creates `.git/` directory structure
- Sets up object database
- Creates HEAD pointing to master branch

### 2. Stage Changes
```bash
echo "Hello World" > hello.txt
simple-git add hello.txt
```
- Reads file content
- Calculates SHA-1 hash of content
- Stores content as blob object in `.git/objects/`
- Updates index to track this file

### 3. Create Commit
```bash
simple-git commit -m "Initial commit"
```
- Reads all staged files from index
- Creates tree objects representing directory structure
- Creates commit object pointing to root tree
- Updates branch reference to point to new commit

### 4. View History
```bash
simple-git log
```
- Follows commit chain from current branch
- Displays commit information and messages

## Installation & Setup

### Prerequisites
- Go 1.17 or higher
- Make (optional, for convenience)

### Build from Source

1. **Clone and build:**
   ```bash
   git clone https://github.com/mutasim77/lets-build-x/simple-git
   cd simple-git
   go build -o simple-git .
   ```

2. **Or use Make:**
   ```bash
   make build
   ```

### Quick Start

1. **Initialize a repository:**
   ```bash
   ./simple-git init
   ```

2. **Create and add files:**
   ```bash
   echo "Hello Simple Git" > file.txt
   ./simple-git add file.txt
   ```

3. **Create your first commit:**
   ```bash
   ./simple-git commit -m "feat: initial commit"
   ```

4. **Check status and history:**
   ```bash
   ./simple-git status
   ./simple-git log
   ```

## Available Commands

| Command             | Description                   | Example                              |
| ------------------- | ----------------------------- | ------------------------------------ |
| `init`              | Initialize new repository     | `simple-git init`                    |
| `add <files>`       | Stage files for commit        | `simple-git add file.txt src/`       |
| `commit -m <msg>`   | Create new commit             | `simple-git commit -m "Add feature"` |
| `status`            | Show working directory status | `simple-git status`                  |
| `log`               | Show commit history           | `simple-git log`                     |
| `branch [name]`     | Create or list branches       | `simple-git branch feature`          |
| `checkout <branch>` | Switch branches               | `simple-git checkout feature`        |
| `help`              | Show help information         | `simple-git help`                    |

## Development Commands (Makefile)

```bash
# Build the project
make build

# Run with arguments
make run ARGS="init"
make run ARGS="add file.txt"
make run ARGS="commit -m 'message'"

# Quick demo
make demo

# Clean build artifacts
make clean
```

## Learning Path
### Recommended Reading Order

1. **`main.go`** - Start here to understand the overall interface
2. **`objects.go`** - Core concept: how Git stores everything
3. **`index.go`** - How staging works
4. **`refs.go`** - How branches and HEAD work
5. **`commands.go`** - How it all comes together
6. **`utils.go`** - Supporting utilities

### Key Concepts to Understand

1. **Content-Addressable Storage**
   - Everything identified by SHA-1 hash of content
   - Automatic deduplication
   - Immutable objects

2. **The Three Object Types**
   - **Blob**: File contents
   - **Tree**: Directory listings (points to blobs and other trees)
   - **Commit**: Snapshot metadata (points to a tree + parent commits)

3. **References Are Just Files**
   - Branches are files containing commit hashes
   - HEAD is a file pointing to current branch
   - Lightweight and fast

4. **The Index/Staging Area**
   - Intermediate state between working directory and repository
   - Allows crafting commits carefully
   - Tracks file metadata and hashes

## Deep Dive: What Makes Git Special?

### Content-Addressable Storage
```
Content: "Hello World"
↓
SHA-1: 5e1c309dae7f45e0f39b1bf3ac3cd9db12e7d689
↓
Storage: .git/objects/5e/1c309dae7f45e0f39b1bf3ac3cd9db12e7d689
```

### Object Relationships
```
Commit Object
├── tree abc123...
├── parent def456...
├── author Alice <alice@example.com>
└── message "Add feature"

Tree Object (abc123...)
├── blob 111222... "README.md"
├── blob 333444... "main.go"
└── tree 555666... "src/"
```

### Branch Pointer Magic
```
master → commit abc123
feature → commit def456
HEAD → refs/heads/master
```

## 🧪 Testing Your Understanding

Try these experiments:

1. **Explore the object database:**
   ```bash
   find .git/objects -type f
   # See the actual stored objects!
   ```

2. **Look at the index file:**
   ```bash
   xxd .git/index | head
   # Binary format showing staged files
   ```

3. **Examine references:**
   ```bash
   cat .git/HEAD
   cat .git/refs/heads/master
   # See how branches are just files with commit hashes
   ```

## What's Not Implemented (Yet!)

This is a learning implementation focusing on core concepts. Missing features include:

- **Merging and conflict resolution**
- **Remote repositories and networking**
- **Rebasing and cherry-picking**
- **Submodules and worktrees**
- **Pack files and delta compression**
- **Hooks and advanced configuration**

Each omission is intentional to keep the focus on understanding Git's fundamental architecture.

## Contributing

Found a bug or want to add a feature? Contributions welcome!

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Further Reading

- [Git Internals - Git Objects](https://git-scm.com/book/en/v2/Git-Internals-Git-Objects)
- [Building Git - Online Book](https://shop.jcoglan.com/building-git/)
- [Git from the Bottom Up](https://jwiegley.github.io/git-from-the-bottom-up/)

## License

MIT License - see LICENSE file for details.

---
**Remember**: The goal isn't to replace Git, but to understand it deeply. Happy learning! ❤️