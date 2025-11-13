# GoLand Docker Setup Guide

This guide explains how to set up and use Docker with GoLand IDE for the Nevarol project.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Initial Setup](#initial-setup)
- [Docker Integration](#docker-integration)
- [Running Docker Compose from GoLand](#running-docker-compose-from-goland)
- [Debugging with Docker](#debugging-with-docker)
- [Database Tools Integration](#database-tools-integration)
- [Common Workflows](#common-workflows)
- [Troubleshooting](#troubleshooting)

## Prerequisites

1. **GoLand IDE**: Download and install from [JetBrains](https://www.jetbrains.com/go/)
2. **Docker Desktop**: Install from [Docker's website](https://www.docker.com/products/docker-desktop)
   - For Windows: Docker Desktop for Windows
   - For macOS: Docker Desktop for Mac
   - For Linux: Docker Engine and Docker Compose
3. **Go SDK**: GoLand will help you set this up, but you can also install Go 1.24.5+ manually

## Initial Setup

### 1. Open the Project in GoLand

1. Launch GoLand
2. Click **File → Open**
3. Navigate to your `nevarol` project directory
4. Click **OK**

### 2. Configure Go SDK

1. Go to **File → Settings** (Windows/Linux) or **GoLand → Preferences** (macOS)
2. Navigate to **Go → GOROOT**
3. If not already set, add your Go installation:
   - Click the **+** button
   - Select **Download** to let GoLand download Go 1.24.5
   - Or browse to your existing Go installation
4. Click **OK**

### 3. Enable Docker Plugin

1. Go to **File → Settings → Plugins**
2. Search for "Docker"
3. Ensure the **Docker** plugin is installed and enabled
4. Restart GoLand if prompted

## Docker Integration

### Connect GoLand to Docker

1. Go to **File → Settings → Build, Execution, Deployment → Docker**
2. Click the **+** button to add a new Docker configuration
3. Select your platform:
   - **Windows/macOS**: Choose "Docker for Windows" or "Docker for Mac"
   - **Linux**: Choose "Unix socket" and use `unix:///var/run/docker.sock`
4. Click **OK**
5. GoLand will test the connection - you should see "Connection successful"

### Docker Tool Window

1. Open the Docker tool window: **View → Tool Windows → Services** (or Alt+8)
2. You should see your Docker connection listed
3. From here you can:
   - View running containers
   - View images
   - View volumes
   - View networks
   - Manage containers (start, stop, remove)

## Running Docker Compose from GoLand

### Method 1: Using the Docker Compose Configuration (Recommended)

1. **Locate docker-compose.yml**: In the Project view, find `docker-compose.yml` at the root
2. **Right-click** on `docker-compose.yml`
3. **Select** one of these options:
   - **Run 'docker-compose.yml'**: Starts all services in foreground
   - **Deploy 'docker-compose.yml'**: Starts all services in detached mode
4. The Services tool window will open, showing:
   - `db`: PostgreSQL database container
   - `app`: Nevarol application container

### Method 2: Create a Run Configuration

1. Go to **Run → Edit Configurations**
2. Click the **+** button and select **Docker → Docker Compose**
3. Configure:
   - **Name**: `Nevarol Docker Compose`
   - **Compose files**: Click **+** and select `docker-compose.yml`
   - **Options**: Leave default or add `-d` for detached mode
4. Click **OK**
5. Now you can run it from the toolbar using the run button

### Viewing Logs

1. In the **Services** tool window, expand your docker-compose deployment
2. Click on a service (e.g., `app` or `db`)
3. Select the **Log** tab to view container logs in real-time
4. You can filter logs, search, and even export them

### Stopping Services

To stop Docker Compose services:
1. In the **Services** tool window, right-click on the deployment
2. Select **Down** to stop and remove containers
3. Or use the stop button in the toolbar

## Debugging with Docker

### Option 1: Debug Locally (Recommended for Development)

For faster iteration, you can run the database in Docker and the Go app locally:

1. **Start only the database**:
   ```bash
   docker-compose up -d db
   ```

2. **Create a Go Run Configuration**:
   - Go to **Run → Edit Configurations**
   - Click **+** and select **Go Build**
   - Configure:
     - **Name**: `Nevarol Local`
     - **Run kind**: Directory
     - **Directory**: `./cmd/web`
     - **Environment**: Add these variables:
       ```
       DB_HOST=localhost
       DB_PORT=5432
       DB_USER=postgres
       DB_PASSWORD=postgres
       DB_NAME=nevarol
       IN_PRODUCTION=false
       ```
   - Click **OK**

3. **Set breakpoints** in your code by clicking in the gutter next to line numbers

4. **Click the Debug button** (Shift+F9) to start debugging

5. The application will stop at your breakpoints, allowing you to:
   - Inspect variables
   - Step through code
   - Evaluate expressions

### Option 2: Remote Debugging with Docker

For debugging the actual Docker container:

1. **Modify the Dockerfile** to use Delve (Go debugger):
   ```dockerfile
   FROM golang:1.24-alpine AS builder
   
   WORKDIR /app
   
   # Install Delve
   RUN go install github.com/go-delve/delve/cmd/dlv@latest
   
   # Copy and build
   COPY go.mod go.sum ./
   RUN go mod download
   COPY . .
   RUN go build -gcflags="all=-N -l" -o main ./cmd/web
   
   FROM alpine:latest
   WORKDIR /app
   RUN apk --no-cache add ca-certificates
   COPY --from=builder /go/bin/dlv .
   COPY --from=builder /app/main .
   COPY --from=builder /app/templates ./templates
   COPY --from=builder /app/static ./static
   COPY --from=builder /app/migrations ./migrations
   
   EXPOSE 8080 2345
   
   CMD ["./dlv", "exec", "./main", "--headless", "--listen=:2345", "--api-version=2", "--accept-multiclient"]
   ```

2. **Update docker-compose.yml** to expose the debug port:
   ```yaml
   app:
     ports:
       - "8080:8080"
       - "2345:2345"  # Debug port
   ```

3. **Create a Remote Debug Configuration** in GoLand:
   - Go to **Run → Edit Configurations**
   - Click **+** and select **Go Remote**
   - Configure:
     - **Name**: `Nevarol Docker Debug`
     - **Host**: `localhost`
     - **Port**: `2345`
   - Click **OK**

4. **Start the services** with `docker-compose up`

5. **Attach the debugger** by running the Remote Debug configuration

## Database Tools Integration

GoLand includes powerful database tools that work great with Docker PostgreSQL:

### Connect to the PostgreSQL Database

1. Open the **Database** tool window: **View → Tool Windows → Database**

2. Click the **+** button and select **Data Source → PostgreSQL**

3. Configure the connection:
   - **Host**: `localhost`
   - **Port**: `5432`
   - **Database**: `nevarol`
   - **User**: `postgres`
   - **Password**: `postgres`

4. Click **Test Connection** to verify

5. Click **OK**

### Using Database Tools

Once connected, you can:

- **Browse Tables**: Expand the database to see `users`, `products`, `cart_items`, `orders`, etc.
- **View Data**: Double-click a table to view and edit data
- **Run Queries**: Right-click the database and select **New → Query Console**
- **Manage Schema**: View and modify table structures
- **Import/Export**: Import or export data in various formats
- **Generate Diagrams**: Right-click the database and select **Diagrams → Show Visualization**

### Running SQL Queries

1. Right-click the `nevarol` database
2. Select **New → Query Console**
3. Type your SQL query, for example:
   ```sql
   SELECT * FROM products;
   SELECT * FROM users ORDER BY created_at DESC;
   ```
4. Press **Ctrl+Enter** to execute
5. View results in the output pane

## Common Workflows

### Workflow 1: Starting Development

1. **Start Docker services**:
   - Right-click `docker-compose.yml` → **Run**
   - Or use the Run Configuration you created

2. **Verify services are running**:
   - Check the Services tool window
   - Both `db` and `app` should show green status

3. **Access the application**:
   - Open browser to `http://localhost:8080`

### Workflow 2: Development with Hot Reload

For faster development with local changes:

1. **Start only the database**:
   ```bash
   docker-compose up -d db
   ```

2. **Run the Go app locally** using the debug configuration

3. **Make code changes** and GoLand will:
   - Highlight errors in real-time
   - Provide auto-completion
   - Suggest improvements

4. **Restart the debug session** to see changes (Ctrl+F5)

### Workflow 3: Rebuilding Docker Images

When you change dependencies or Dockerfile:

1. **Stop running containers**:
   - In Services tool window, right-click deployment → **Down**

2. **Rebuild images**:
   ```bash
   docker-compose build --no-cache
   ```
   Or in GoLand terminal: **View → Tool Windows → Terminal**

3. **Start services again**:
   - Right-click `docker-compose.yml` → **Run**

### Workflow 4: Database Management

1. **View database in Database tool**:
   - **View → Tool Windows → Database**

2. **Run migrations manually** (if needed):
   - The app runs migrations automatically on startup
   - Or connect to the container and run SQL scripts

3. **Export/Backup data**:
   - Right-click a table → **Export Data**
   - Choose format (SQL, CSV, JSON, etc.)

### Workflow 5: Viewing Container Logs

1. **Open Services tool window** (Alt+8)

2. **Expand docker-compose deployment**

3. **Click on a service** (`app` or `db`)

4. **View logs** in the Log tab:
   - Application logs show startup, requests, errors
   - Database logs show queries and connections

5. **Filter logs**:
   - Use the search box to find specific entries
   - Right-click to clear or export logs

## Troubleshooting

### Docker Connection Issues

**Problem**: GoLand cannot connect to Docker

**Solutions**:
- Ensure Docker Desktop is running
- Check **Settings → Docker** and test connection
- On Windows, ensure Docker is running in Linux container mode
- Try restarting GoLand and Docker Desktop

### Port Already in Use

**Problem**: Port 8080 or 5432 already in use

**Solutions**:
- Stop other applications using these ports
- Or modify `docker-compose.yml` to use different ports:
  ```yaml
  ports:
    - "8081:8080"  # Use 8081 instead of 8080
  ```

### Database Connection Refused

**Problem**: App cannot connect to database

**Solutions**:
- Wait for database to be ready (check healthcheck in logs)
- Verify database container is running in Services window
- Check environment variables in docker-compose.yml match

### Go Modules Not Recognized

**Problem**: GoLand shows "Cannot resolve symbol" errors

**Solutions**:
- Right-click `go.mod` → **Sync Dependencies**
- Go to **File → Invalidate Caches** and restart
- Ensure GOROOT is properly configured

### Container Builds Fail

**Problem**: Docker build fails with errors

**Solutions**:
- Check the Build log in Services tool window
- Ensure Dockerfile syntax is correct
- Try building manually: `docker-compose build --no-cache`
- Check you have enough disk space

### Changes Not Reflected in Container

**Problem**: Code changes don't appear when running in Docker

**Solutions**:
- Rebuild the Docker image: `docker-compose build`
- Restart containers: `docker-compose restart`
- For development, use local debugging instead (see above)

### Services Tool Window Not Showing

**Problem**: Can't find the Services tool window

**Solutions**:
- Go to **View → Tool Windows → Services**
- Or press **Alt+8**
- Ensure Docker plugin is enabled in Settings → Plugins

## Additional Resources

- [GoLand Documentation](https://www.jetbrains.com/help/go/)
- [Docker Documentation](https://docs.docker.com/)
- [Go Documentation](https://golang.org/doc/)
- [Project README](README.md) - General project setup
- [Email Setup Guide](EMAIL_SETUP.md) - Email configuration

## Quick Reference

### Useful GoLand Shortcuts

- **Alt+8**: Open Services tool window
- **Shift+F9**: Debug
- **Shift+F10**: Run
- **Ctrl+Shift+F10**: Run context configuration
- **Ctrl+F5**: Rerun
- **Ctrl+F2**: Stop
- **Alt+F12**: Open Terminal

### Useful Docker Commands (in GoLand Terminal)

```bash
# Start services
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop services
docker-compose down

# Rebuild and start
docker-compose up -d --build

# View running containers
docker ps

# Execute command in container
docker-compose exec app sh

# Access PostgreSQL shell
docker-compose exec db psql -U postgres -d nevarol
```

### Environment Variables Reference

For local development (in Run Configuration):
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=nevarol
IN_PRODUCTION=false
```

For Docker (in docker-compose.yml):
```yaml
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=nevarol
IN_PRODUCTION=false
```

---

**Happy coding!** If you encounter any issues not covered here, please check the main [README](README.md) or open an issue on GitHub.
