# continuous Delivery & Deployment Manager (D2M)

D2M is a robust continuous delivery and deployment management tool that provides a seamless interface for managing deployments across your projects. It features both a CLI interface and a web panel for deployment management and monitoring.

## Features

- Daemon-based deployment management
- Web-based deployment monitoring panel
- Real-time deployment logs
- Branch and commit tracking
- Deployment status monitoring
- TCP-based inter-process communication
- SQLite database for persistent storage

## Components

### CLI Application

The CLI interface provides commands for:
- Initialization
- Deployment updates
- Log management
- Daemon control

### Web Panel

Built with:
- React
- TypeScript
- Vite
- Modern UI components

Features:
- Deployment monitoring dashboard
- Real-time status updates
- Detailed deployment logs
- Branch and commit information
- Time-based tracking

## Installation

1. Clone the repository:
```bash
git clone https://github.com/Techzy-Programmer/d2m.git
```

2. Build the CLI and Web Panel:
```bash
cd web/panel
npm install
npm run build

go build ./app/cli
```

## Usage
```bash
./d2m init # Initialize the deployment manager
./d2m --help # Display help information
```

## Web Panel
Web panel is automatically started when the daemon is initialized.

Access the panel at `http://localhost:[configured-port]`.

or if running on a VPS or remote server, access the panel at `http://[server-ip]:[configured-port]`.
