# Overview

Simple windows service that reports via HTTP status and body if the SQL Server running on localhost has all mirroring databases in `principal`, `mirroring` or `mixed` modes. 

## Build the service

```
./make.sh build-windows
```
^ this creates a file called `CheckMirror.exe` in the directory. Copy this to you target SQL Server and run it standalone or as a service

## Run the standalone API in Powershell:
```
PS C:\> C:\Software\CheckMirror.exe 
```

## Install it as Service

```
PS C:\Users\azureuser> cd \

PS C:\> .\Software\CheckMirror.exe --service.do=install

PS C:\> Get-Service CheckMirror

# Manually or programatically set the "Log On" for the service to be a user that has SQL connect & query rights ... 

PS C:\> Start-Service CheckMirror

PS C:\> Invoke-WebRequest http://127.0.0.1:8282
```

## Remove the Service
```
PS C:\> .\Software\CheckMirror.exe --service.do=remove

PS C:\> Get-Service CheckMirror
```

