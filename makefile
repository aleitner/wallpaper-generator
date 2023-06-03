.PHONY: build bat install uninstall log config

# Define the username and the Windows Startup folder path
USERNAME := $(shell powershell -command "(Get-ChildItem Env:UserName).Value")
APPDATA := C:\Users\$(USERNAME)\AppData\Roaming\WallpaperGenerator

# Build the WallpaperGenerator.exe
build:
	go build -o WallpaperGenerator.exe main.go

# Install: copy wallpaper_generator.bat to the Windows Startup folder
install: build
	copy config.json $(APPDATA)

# Target to uninstall the wallpaper_generator.bat from the Startup folder
uninstall:
	del $(APPDATA)\logfile.log
	del $(APPDATA)\config.json

# Target to edit the config file
config:
	notepad $(APPDATA)\config.json

# Target to display the log content
log:
	notepad $(APPDATA)\logfile.log
