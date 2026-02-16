# Frontend

SwiftUI iOS/macOS app for Clientele Manager.

## Structure

| File | Description |
|------|--------------|
| `Clientele Manager/Clientele_ManagerApp.swift` | App entry point; shows `LoginView` |
| `Clientele Manager/LoginView.swift` | Username/password login form |
| `Clientele Manager/AuthService.swift` | POST login to configurable URL |
| `Clientele Manager/AuthConfig.swift` | Placeholder `loginURL` – replace with real backend URL |
| `Clientele Manager/Assets.xcassets` | App assets including AppIcon |

## How to Build

1. Open `Clientele Manager/Clientele Manager.xcodeproj` in Xcode
2. Select your target device or simulator
3. Build and run (⌘R)

## Current State

- Login UI and auth service are implemented
- `AuthConfig.loginURL` is a placeholder (`https://your-api.com/auth/login`) and must be updated to point at your backend auth endpoint before login will work
