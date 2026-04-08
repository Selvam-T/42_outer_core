# Flutter Reference Notes

## 1. Running a Flutter app on a local web server

Command:

```bash
flutter run -d web-server --web-port 8080 --web-hostname 0.0.0.0
```

### What this command does

- `flutter run` compiles your Dart code and starts a development session.
- `-d web-server` targets the lightweight web server device instead of opening Chrome automatically.
- `--web-port 8080` forces the server to use port 8080.
- `--web-hostname 0.0.0.0` makes the server listen on all available network interfaces.

### Is it on localhost?

Yes, on your own machine you can usually open:

- `http://localhost:8080`
- `http://127.0.0.1:8080`

Because `0.0.0.0` is used, other devices on the same local network may also be able to access it through your computer's LAN IP, for example:

- `http://192.168.1.15:8080`

### Typical use cases

This mode is useful when:
- running Flutter inside Docker
- working over SSH
- testing from another device on the same network
- manually opening the app in a browser

---

## 2. `flutter run -d chrome` vs `flutter run -d web-server`

### `flutter run -d chrome`

- opens Chrome automatically
- attaches debugger more directly
- convenient for local development on one machine

### `flutter run -d web-server`

- starts a web server only
- does not open a browser automatically
- you open the URL manually
- useful for Docker, remote access, and LAN testing

### Comparison

| Feature | `flutter run -d chrome` | `flutter run -d web-server` |
|---|---|---|
| Auto-open browser | Yes | No |
| Manual URL open needed | No | Yes |
| Good for Docker / SSH | Less suitable | Very suitable |
| LAN access | Usually more limited | Better with `0.0.0.0` |

---

## 3. How Flutter knows where the Dart code is

By convention, Flutter expects:

- a valid Flutter project in the current directory
- a `pubspec.yaml` file in the project root
- an entry file at `lib/main.dart`
- a `main()` function inside that file

So when you run `flutter run`, Flutter looks in the current project directory and uses the standard entry point unless you override it.

### Custom entry file

You can specify another entry file using `-t`:

```bash
flutter run -d web-server -t lib/login_screen.dart
```

---

## 4. Does `flutter run` compile the code?

Yes. The code is not assumed to be precompiled.

When you run `flutter run`, Flutter performs compilation as part of the run process.

### For Web

Browsers do not run Dart directly in the normal Flutter web flow, so Flutter compiles or transpiles Dart into JavaScript, then serves the generated files.

In debug mode, this output is optimized for development rather than final performance.

### For Android

In debug mode, Flutter uses a development-oriented process that supports fast iteration and hot reload.

In release mode, Flutter builds a production version optimized for performance.

---

## 5. Lifecycle of `flutter run` for Web

### Summary flow

1. **Check**
   - verifies you are in a Flutter project
   - checks for `pubspec.yaml`

2. **Locate**
   - finds the entry file, usually `lib/main.dart`

3. **Compile**
   - processes Dart code for the web target

4. **Serve**
   - starts the local web server on the requested host and port

5. **Stay resident**
   - remains active for logs and development commands

### Practical result

For a command like:

```bash
flutter run -d web-server --web-port 8080 --web-hostname 0.0.0.0
```

Flutter:
- compiles the web app
- serves it on port 8080
- keeps the session alive in the terminal

---

## 6. Lifecycle of `flutter run` on Android phone

When running on an Android device, the flow is different from web.

### Summary flow

1. **Preparation**
   - validates the Flutter project
   - checks that an Android device is connected
   - checks toolchain requirements

2. **Dependency / project sync**
   - ensures packages and plugins are available

3. **Compilation phase**
   - prepares Dart code in development form suitable for debugging

4. **Android build phase**
   - Gradle performs Android-side build work
   - native/plugin code is built if needed
   - assets and resources are packaged

5. **Deployment**
   - APK is pushed to the device
   - app is installed

6. **Execution**
   - app launches on the phone
   - Flutter engine starts
   - Dart runtime starts
   - app begins from `main()`

7. **Resident development loop**
   - terminal remains connected
   - logs are streamed back
   - hot reload / restart can occur during development

---

## 7. `flutter run` vs `flutter build`

This distinction is about **purpose**, not platform.

### Use `flutter run` when:

- you are developing
- you want logs
- you want debugging
- you want rapid iteration

### Use `flutter build` when:

- you want a final deployable app
- you want optimized output
- you want to distribute to users

### Web example

#### Development
```bash
flutter run -d chrome
```
or
```bash
flutter run -d web-server --web-port 8080 --web-hostname 0.0.0.0
```

#### Deployment
```bash
flutter build web
```

This produces deployable static web files.

### Android example

#### Development
```bash
flutter run
```

#### Deployment
```bash
flutter build apk
```
or
```bash
flutter build appbundle
```

### Comparison

| Feature | `flutter run` | `flutter build` |
|---|---|---|
| Main purpose | Development | Deployment |
| Hot reload | Yes | No |
| Debugging | Yes | No |
| Optimized output | No | Yes |
| Stays attached in terminal | Yes | No |

---

## 8. Do you need JDK or Android SDK for Web?

If you are only targeting web, typically:

- Flutter SDK: needed
- Dart SDK: needed indirectly through Flutter
- JDK: not needed for web-only workflow
- Android SDK: not needed for web-only workflow

### Why not?

- Web build does not require Android build tools
- Chrome/web-server targets do not use Android packaging or ADB
- Java/Gradle tools are mainly needed for Android workflows

### What you need for web

- Flutter SDK
- browser support such as Chrome or another modern browser

### About `flutter doctor`

You may still see warnings related to Android toolchain if Android is not installed. Those can usually be ignored for a web-only workflow.

---

## 9. Dart SDK vs Flutter SDK vs JDK vs Android SDK

A simple breakdown:

### Dart SDK

What it does:
- supports the Dart language
- provides compiler/runtime/tooling for Dart

Why it matters:
- your app logic is written in Dart

### Flutter SDK

What it does:
- provides Flutter framework and widgets
- provides Flutter engine integration
- provides the `flutter` command

Why it matters:
- used to build UI and run Flutter apps

### JDK

What it does:
- provides Java tooling required by Android build systems such as Gradle

Why it matters:
- needed for Android app building
- not needed for web-only work

### Android SDK

What it does:
- provides Android platform tools
- includes tools for device communication and Android builds

Why it matters:
- needed to run/build on Android devices or emulators
- not needed for web-only work

### Simple table

| Tool | Main role | Needed for Web | Needed for Android |
|---|---|---|---|
| Dart SDK | language tooling | Yes | Yes |
| Flutter SDK | Flutter framework/tooling | Yes | Yes |
| JDK | Java-based Android build support | No | Yes |
| Android SDK | Android platform/device/build tools | No | Yes |

---

## 10. Python / PyQt analogy

This was explained using a comparison with Python and PyQt desktop development.

### Dart vs Python

- Python is the language in a PyQt app
- Dart is the language in a Flutter app

Both are where your logic lives.

### Flutter vs PyQt

- PyQt gives you desktop UI classes and widgets
- Flutter gives you its own widget system and UI framework

Both provide the visual building blocks.

### JDK analogy

JDK is not where you write app logic.
It is more like a required supporting build environment used by Android's packaging/build pipeline.

A rough analogy:
- not equal to PyQt itself
- closer to extra platform build tooling you need in order to produce a platform-specific packaged app

### Android SDK analogy

Android SDK is the platform bridge/toolset for Android.

Rough analogy:
- similar in spirit to OS-specific development tools required to package and target a specific platform

### Practical summary using your PyQt style thinking

Running:

```bash
flutter run -d web-server
```

is somewhat like running a development version of your GUI app directly, except that:
- the UI is rendered in a browser
- Flutter compiles/serves the app for web
- you are not yet producing a final packaged Android app

---

## 11. One-line practical summary

### If you are developing
Use:

```bash
flutter run
```

### If you are deploying
Use:

```bash
flutter build ...
```

### If you are targeting web only
You only need:
- Flutter SDK
- browser support

### If you are targeting Android
You also need:
- JDK
- Android SDK

---

## 12. Final takeaway

For this command:

```bash
flutter run -d web-server --web-port 8080 --web-hostname 0.0.0.0
```

the main idea is:

- Flutter finds your project entry point
- compiles the app for web
- starts a local web server
- serves the app on port 8080
- keeps the session active for development use

You can then open it manually in a browser, usually through:

- `http://localhost:8080`
