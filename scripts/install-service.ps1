# Worker Monitoring Agent - Windows Service Installer
# Run as Administrator

param(
    [string]$Action = "install"
)

$ServiceName = "WorkerMonitorAgent"
$DisplayName = "Worker Monitoring Agent"
$Description = "Monitors CPU, RAM, and GPU metrics and provides HTTP API"
$BinaryPath = "$PSScriptRoot\..\agent.exe"

function Install-Service {
    Write-Host "Installing $DisplayName..."
    
    # Check if service already exists
    $service = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
    if ($service) {
        Write-Host "Service already exists. Removing..."
        Uninstall-Service
    }

    # Create service
    sc.exe create $ServiceName binPath= $BinaryPath start= auto DisplayName= $DisplayName
    sc.exe description $ServiceName $Description
    
    # Set environment variables (optional - can also use registry)
    # For now, the service will use default values
    
    Write-Host "Service installed successfully!"
    Write-Host "To start the service, run: Start-Service $ServiceName"
}

function Uninstall-Service {
    Write-Host "Uninstalling $DisplayName..."
    
    # Stop service if running
    $service = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
    if ($service -and $service.Status -eq 'Running') {
        Write-Host "Stopping service..."
        Stop-Service -Name $ServiceName -Force
        Start-Sleep -Seconds 2
    }
    
    # Delete service
    sc.exe delete $ServiceName
    Write-Host "Service uninstalled successfully!"
}

function Start-AgentService {
    Write-Host "Starting $DisplayName..."
    Start-Service -Name $ServiceName
    Write-Host "Service started!"
    Get-Service -Name $ServiceName
}

function Stop-AgentService {
    Write-Host "Stopping $DisplayName..."
    Stop-Service -Name $ServiceName -Force
    Write-Host "Service stopped!"
}

# Main execution
switch ($Action.ToLower()) {
    "install" {
        Install-Service
    }
    "uninstall" {
        Uninstall-Service
    }
    "start" {
        Start-AgentService
    }
    "stop" {
        Stop-AgentService
    }
    default {
        Write-Host "Usage: .\install-service.ps1 [install|uninstall|start|stop]"
        Write-Host "  install   - Install the service"
        Write-Host "  uninstall - Uninstall the service"
        Write-Host "  start     - Start the service"
        Write-Host "  stop      - Stop the service"
    }
}
