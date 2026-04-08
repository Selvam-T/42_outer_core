## 1. Dockerfile Installations Explained


**The Logic: Layer by Layer**  

Think of building a Docker image like making a sandwich. You have to lay the ingredients down in a specific order.  

***The Bread (Base OS & Basics):***  
- Plan: "git, curl, unzip, ca-certificates"  

- Docker: We start with a lightweight Linux (like Ubuntu or Debian). We install ca-certificates immediately so we can talk to the internet securely.  

***The Spread (Java/JDK 17):***  
- Plan: "JDK 17"  

- Docker: Android tools are written in Java. Without this layer, the Android SDK cannot run. We install OpenJDK 17.  

***The Meat (Android SDK):***  
- Plan: "cmdline-tools, platform-tools, build-tools..."  

- Docker: This is the trickiest part.  
 We download the Command Line Tools zip file from Google.  
 We unzip it.  
 We use the sdkmanager (inside that zip) to download the specific parts you listed (platforms;android-36, build-tools;36.0.0).  
 To prevent the container from trying to download the NDK at runtime, we "bake" it into the image. 
 Critical Step: We must accept the Google Licenses automatically using a script (using yes), otherwise the build fails waiting for a human to type "y".  

***The Topping (Flutter):***  

- Plan: "Flutter SDK"  

- Docker: We git clone the Flutter repository. We add the flutter/bin folder to the system PATH so you can type flutter anywhere.  

**Dockerfile necessary software components**  
***To compile the code required for Exercises 00–03.***  

- Flutter SDK: Required to create the projects (ex00, ex01, etc.) and widgets.  
- Android SDK & Gradle: Required to build the APKs for the calculator logic.  
- JDK 17: Required to run the Android build tools.  
- Internet Dependencies (curl/git/ca-certificates): Required because Exercise 03 (Page 11) explicitly allows/mentions using the math_expressions library. Flutter will need to download this package from pub.dev when you build.  

**Docker command broken down**  
```
docker run -it --privileged -v /dev/bus/usb:/dev/bus/usb -v $(pwd):/app my-flutter-image bash
```

1. docker run -it  
- What it does: Runs the container in Interactive Mode.  
- Why you need it: Without this, the container would start, finish its job, and close immediately. You want to stay "inside" the container (get a command prompt) so you can run multiple commands like flutter clean or flutter run repeatedly.  
2. --privileged  
- What it does: This is "God Mode" for the container. It removes standard security isolation.  
- Why you need it: By default, Docker containers are not allowed to access hardware devices (like USB ports, Webcams, or GPUs). This flag tells Docker: "Trust this container with full access to my host hardware."  
3. -v /dev/bus/usb:/dev/bus/usb  
- What it does: This mounts your Linux USB device directory into the container.  
- The Logic: On Linux, everything is a file. Your physical USB ports live at /dev/bus/usb. By mapping this folder, when you plug a phone into your laptop, the container thinks the phone is plugged into it.  
4. -v $(pwd):/app  
- What it does: This mirrors your current project folder (where your code is) to the /app folder inside the container.  
- Why you need it:  
<> You write code on your Laptop (Host) using VS Code.  
<> The Container reads the code from /app.    
<> Because it is mirrored, changes happen instantly. You don't need to rebuild the image when you change a line of code.  
5. bash  
- What it does: This is the specific command to run upon starting.  
- Why you need it: It opens a standard Linux shell (Terminal) inside the container so you can type commands.  

## 2. Tell Make to ignore the exit code

***The make: *** [Makefile:14: run] Error 1 occurs because Make is designed to stop and complain whenever a command returns a "non-zero" exit code.***

If you want to exit your container and not have make scream at you, you can tell Makefile to ignore the return code by adding a - (minus sign) before the docker command.

-@docker run --rm -it \
		-p 8080:8080 \
		-v $(PWD)/mobilePiscine:/app \
		$(NAME) bash

## 3. Improve the chmod target (Dynamic Container ID)  

The chown in your Dockerfile is ineffective the moment you use a Volume Mount (-v).  

***1. Why the Dockerfile chown fails for /app***  
- Build-Time vs. Runtime: The chown command in your Dockerfile happens during the Build Phase. At that moment, it correctly sets permissions for the empty folder inside the image.  
- The "Masking" Effect: When you run the container with -v $(PWD)/mobilePiscine:/app, Docker "mounts" your host folder over the internal one. It’s like placing a sticker over a drawing; the drawing (the Dockerfile permissions) is still there, but you can only see and interact with the sticker (the Host permissions).  
- Ownership Inheritance: Docker does not change the ownership of your host files to match the container user. If your host files are owned by "You" (Host User), the container user "developer" is treated as a "Stranger" with no right to delete or modify them.  

***2. Why the PathAccessException happens***  
- Flutter is very aggressive with its build/ folder. Every time you run the app, it tries to delete "Stamps" (temporary files) to ensure a fresh build.  
- If those Stamps were created by your Host Machine (if you ever ran flutter there) or a previous Container session that had different permissions, the current "developer" user crashes because it doesn't have the "Write/Delete" authority over those specific files.  

***3. Why make chmod is the solution***  
- Runtime Fix: make chmod doesn't happen during the image build; it happens while the "sticker" (the host volume) is already applied.  
- Root Authority: By using docker exec -u 0 (Root), you are telling the container to use "Superuser" powers to reach out and change the permissions of those host files from the inside.
The Result: It changes the permissions to a+rwx (All users can Read, Write, and Execute). Now, the "developer" user is no longer a "Stranger" and is allowed to delete the build stamps.  

***Comparison Table for Notes:***  
 
Dockerfile:  
RUN chown developer /app  
Sets permissions for the internal (empty) folder.  

Docker Run:  
-v host_folder:/app  
Overrides internal permissions with Host permissions.   

Make chmod:  
docker exec -u 0 chmod ...  
Forces new permissions onto the Host files while they are inside the container.  

## 4. What is the Flutter build/ folder?  
The build/ folder is the workspace where Flutter does all its heavy lifting. It is not part of your source code; it is the output of the compilation process.  

***What’s inside it:***  
- Platform-specific code: Subfolders for web, android, ios, etc., containing the actual files the OS needs to run the app.  
- Compiled Dart: Your Dart code translated into JavaScript (for web) or Bytecode (for mobile).  
- Assets: Optimized versions of your images, fonts, and icons.
Incremental "Stamps": Small files (like the .stamp files you saw) that tell Flutter: "I already compiled this part, don't do it again."  

## 5. Does the Dockerfile need a "Dart" section?
No. Because you have $FLUTTER_HOME/bin in your PATH, you can already run Dart commands by typing flutter dart.  

## 6. What display/target am I prepared for?

**From the Dockerfile:**

I installed -     
- Flutter  
- Android SDK command-line tools  
- platform-tools / build-tools / platform 34  

I do not have -  
- Chrome installation  
- Android emulator setup  
- Play Store signing/release setup  

**So I am mainly prepared for:**

Best-supported by this Dockerfile

- 1. Android APK sideload (manual) → Yes  
- 2. Web dev server → Possible  (Not permitted in the project)
- 3. Build static web → Possible (Not permitted in the project)
- 4. Android App Bundle (Play Store) → possible to build later, but not really what this Dockerfile is specifically prepared for  

**Most accurate summary:**

Prepared strongest for:  
- Android build tooling  

Also usable for:  
- Flutter web development/build  

Not clearly prepared for:  
- full Play Store release pipeline  

### 7. Development vs. Deployment: A Technical Comparison  

***1. Compilation & Target***  
Web-Server: Uses the Dart Dev Compiler (dartdevc). It transpiles your Dart code into readable JavaScript. It is optimized for "Hot Restart" speed, not execution performance.  
APK Build: Uses Ahead-of-Time (AOT) compilation. It compiles Dart into native ARM Machine Code. It then uses Gradle (Java) to bundle the code, the Flutter Engine, and your assets into a single compressed file.  

***2. Resource Requirements (The "Heavyweight" Factor)***    
Web-Server (Lightweight): Requires almost no extra tools. It uses the Dart SDK already inside Flutter. It has a very small disk footprint.  
APK Build (Heavyweight): Requires the Full Android Toolchain.  
JDK: To run the Gradle build system.  
Android SDK: To package the resources.  
NDK: To compile native C/C++ components (often 1GB+).  
Disk Space: Needs significantly more storage for caches (stored in .gradle).  

***3. Docker Networking & Ports***    
Web-Server (Needs Ports): Requires mapping ports (e.g., -p 8080:8080) because it hosts a live web server. It also uses -p 5555:5555 (VM Service) to allow the Host debugger to talk to the Container.  
APK Build (No Ports): Does not need any port mapping. It doesn't "listen" for connections; it just performs a calculation and writes a file to the disk.  

***4. File Output & Persistence***  
Web-Server: Its output is volatile. The compiled JavaScript is usually served from the build/ folder but is intended to be viewed in a browser.  
APK Build: Its output is a physical file.
Path: build/app/outputs/flutter-apk/app-release.apk.
Persistence: Because of your Docker Volume (-v), this file appears on your Host machine automatically so you can move it to a phone or submit it for grading.  

***5. The "Permission" Conflict***  
Web-Server: Very sensitive to permission errors in the build/ folder because it constantly tries to delete and rewrite "Stamps" (.stamp files) every time you refresh.  
Fix: Often requires the Symlink trick to keep the build/ folder inside the container's private memory.  
APK Build: Less sensitive to live permission locks but prone to "No space left on device" errors due to the massive size of the Android NDK and Gradle caches.  

## Original failed attempts - read for interest

### 1. My initial attempt at Physical Phone Workflow (Failed)
(It was suggested I Pivot to using web dev server.)
Additionally, I must not that my dockerfile installatio is not set up for physical phone workflow.

This is the exact sequence of actions you will take to solve Exercise 00:  

1. Physical Setup: Enable "USB Debugging" on your real Android phone and plug it into your computer via USB.  
2. Kill Host ADB: (Crucial Step) On your laptop terminal, run adb kill-server.  
- Reason: Only one thing can control the USB connection at a time. You want the Container to control the phone, not your laptop.  
3. Enter the Matrix: Run the long Docker command above.  
- Your terminal prompt will change (e.g., root@a1b2c3d4:/app#).  
4. Verify: Type adb devices inside the container.  
- You should see your phone listed (e.g., ZF523... device).  
- If you see "unauthorized," check your phone screen and click "Allow".  
5. Develop:  
- Type: flutter run  
- The app installs on your phone.  
- The app opens on your phone screen.  
- The Logs: You look at your Computer Terminal. It will say: "Waiting for connection..."  
6. The Test (Ex00):  
- You tap the button on your Phone.  
- You see "Button pressed" appear on your Computer Terminal.  

**Why this is the "Winner" for your Subject**  
- Satisfies Log Requirements: flutter run streams logs directly to the terminal (unlike flutter build apk).  
- Fast: It uses your phone's CPU, not a slow emulator inside Docker.  
- No GUI Complexity: You don't need X11, Wayland, or remote desktops. You look at your phone for UI, and your terminal for code/logs.


### 2. I decided to pivot to using Web dev server as target
  
### Why we stopped:  
The combination of Hardware Restrictions (USB Passthrough denied), Network Isolation (WiFi blocked), Admin Restrictions (Tethering blocked), and Hardware Detection Issues (Cable/Driver issues) created too many points of failure.

To set up a Dockerized Flutter environment where the code resides on the Host (42 iMacs) but compiles inside a Docker container, and runs on a physical Android phone connected via USB.

### Phase 1: Infrastructure Setup (Docker & Makefile)  

***Attempt:*** We created a Dockerfile to install Android SDK/Flutter and a Makefile to automate the build/run commands.  

**Error 1 (Dockerfile):** Cannot change ownership to uid... during build.  

**Cause:** Running flutter precache as root inside Docker caused permission issues with downloaded archives.  

**Solution:** Modified Dockerfile to create a non-root user (developer) and properly set permissions on SDK directories.   

**Error 2 (Makefile):** mkdir /mobilePiscine: permission denied.  

**Cause:** Used ${pwd} in the Makefile. Make interpreted this as an empty variable, resulting in Docker trying to mount the root directory / instead of the project path.  

**Solution:** Changed syntax to $(shell pwd) to correctly resolve the current working directory.  


### Phase 2: The USB Passthrough Approach  

***Attempt:*** We tried to mount the physical USB bus into the container using --privileged -v /dev/bus/usb:/dev/bus/usb.  

**Error (Makefile):** error creating device nodes: ... /dev/bus/usb/001/001: permission denied.  

**Cause:** The 42 School environment (/goinfre) enforces strict security policies (cgroups/AppArmor) that prevent standard users from passing physical hardware devices into containers, even with privileged flags.  

**Result:** Direct USB control from inside Docker was impossible.  


### Phase 3: The Network Approach (ADB Wireless)  

***Attempt:*** Since physical USB failed, we tried to connect Docker to the phone over the Network (WiFi).  

***Attempt 3.1 (School WiFi):***  

**Observation:** Host IP was 10.12.x.x. Phone IP was 10.37.x.x.  

**Result:** Ping failed. The school network isolates the wired network (iMacs) from the wireless network (Phones). They cannot talk to each other.  

***Attempt 3.2 (USB Tethering):***  

**Action:** Enabled USB Tethering to create a private network between Phone and iMac.  

**Observation:** hostname -I showed no new IP address.  

**Result:** The school iMacs are configured to ignore/block new network interfaces initiated by USB devices. We could not establish a network link.  


### Phase 4: The Tunneling Approach (Port Forwarding)  

***Attempt:*** We tried to use the Host machine as a "Bridge" by forwarding traffic from the Host to the Phone via USB (adb forward tcp:9999 tcp:5555).  

**Prerequisite:** Needed adb installed locally on the Host.
Action: Successfully downloaded and installed platform-tools in the home directory.  

**Error:** error: no devices/emulators found.  

**Context:** While trying to set up the tunnel, the Host machine itself failed to detect the Android phone via USB.  

**Possible Causes:** Defective USB cable (Charge-only), USB Debugging not authorized on the device, or restricted USB ports on the iMac.  


### The Solution:  
We are switching to Flutter Web (web-server).  

**Workflow:** Run the app as a website inside Docker (flutter run -d web-server).  

**Access:** View the app via Chrome on the Host machine (localhost:8080).  

**Submission:** We will only use the Android build tools at the very end to generate the required APK (flutter build apk) for grading, while doing all development and logic testing in the browser.  
