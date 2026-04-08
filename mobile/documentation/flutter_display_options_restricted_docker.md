# Flutter Display Options Under Restricted Docker Setup

This table summarizes **all available ways to display a Flutter app** given a Docker-based environment with **no USB passthrough and no GUI access**, indicating whether each method is possible and why.

## Display Methods

| Method | Command(s) | Where it displays | Possible? | Why / Notes |
|---|---|---|---|---|
| Web dev server | `flutter run -d web-server --web-port 8080 --web-hostname 0.0.0.0` | Host browser (`http://localhost:8080`) | Yes | No USB or GPU required; best for development and hot reload. |
| Build static web | `flutter build web` | Host browser (served statically) | Yes | Generates static files; no hot reload. |
| Android APK sideload (manual) | `flutter build apk` → copy APK → install | Real phone screen | Yes | No ADB needed; manual install; no hot reload/log streaming. |
| Android App Bundle (Play Store) | `flutter build appbundle` | Phones via Play Store | Yes | Requires Play Console account, signing, and release workflow. |
||
| Linux desktop (native) | `flutter run -d linux` | Linux window | No | Container has no GUI/X11/Wayland forwarding. |
| Android on real phone (live) | `flutter run -d <device>` | Real phone screen | No | USB passthrough blocked; ADB can’t see device. |
| Android emulator inside container | `flutter emulators` | Emulator window | No | Needs GUI + hardware acceleration (KVM). |
| Android emulator on host | Host SDK + emulator | Emulator window | No / Unlikely | Host restrictions usually prevent SDK/emulator install. |
| ADB over Wi‑Fi | `adb connect <ip>` | Real phone screen | No | Requires initial USB setup, which is blocked. |
| iOS simulator / device | `flutter run -d ios` | Simulator / iPhone | No | Requires macOS + Xcode. |
