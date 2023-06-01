.PHONY: build create_wallpaper_generator_bat install uninstall log

# Define the username and the Windows Startup folder path
USERNAME := $(shell powershell -command "(Get-ChildItem Env:UserName).Value")
STARTUP_FOLDER := "C:\Users\$(USERNAME)\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup"

# Build the WallpaperGenerator.exe
build:
	go build -o WallpaperGenerator.exe main.go

# Create the .bat file content with the correct path to the generated exe
create_wallpaper_generator_bat:
	@echo start "" "%cd%\WallpaperGenerator.exe" > wallpaper_generator.bat

# Install: copy wallpaper_generator.bat to the Windows Startup folder
install: build create_wallpaper_generator_bat
	copy wallpaper_generator.bat $(STARTUP_FOLDER)

# Target to display the log content
log:
	type "%LOCALAPPDATA%\WallpaperGenerator\Logs\logfile.log"

# Target to uninstall the wallpaper_generator.bat from the Startup folder
uninstall:
	del $(STARTUP_FOLDER)\wallpaper_generator.bat