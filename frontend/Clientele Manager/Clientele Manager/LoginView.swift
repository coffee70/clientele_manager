//
//  LoginView.swift
//  Clientele Manager
//
//  Created by Evan Haeick on 2/15/26.
//

import SwiftUI

struct LoginView: View {
    @State private var username = ""
    @State private var password = ""
    @State private var isLoading = false
    @State private var errorMessage: String?

    private let inputCornerRadius: CGFloat = 14
    private let inputHeight: CGFloat = 52

    var body: some View {
        ZStack {
            Color(.systemGroupedBackground)
                .ignoresSafeArea()

            VStack(spacing: 0) {
                Spacer()

                VStack(spacing: 24) {
                    VStack(spacing: 8) {
                        Text("Clientele")
                            .font(.title)
                            .fontWeight(.semibold)
                        Text("Login with your Clientbook credentials.")
                            .font(.subheadline)
                            .foregroundStyle(.secondary)
                            .multilineTextAlignment(.center)
                    }
                    .padding(.bottom, 8)

                    TextField("Username", text: $username)
                        .textContentType(.username)
                        .textInputAutocapitalization(.never)
                        .autocorrectionDisabled()
                        .padding(.horizontal, 20)
                        .frame(height: inputHeight)
                        .background(Color(.secondarySystemGroupedBackground))
                        .clipShape(RoundedRectangle(cornerRadius: inputCornerRadius))

                    SecureField("Password", text: $password)
                        .textContentType(.password)
                        .padding(.horizontal, 20)
                        .frame(height: inputHeight)
                        .background(Color(.secondarySystemGroupedBackground))
                        .clipShape(RoundedRectangle(cornerRadius: inputCornerRadius))

                    Button {
                        signIn()
                    } label: {
                        if isLoading {
                            ProgressView()
                                .progressViewStyle(CircularProgressViewStyle(tint: .white))
                                .frame(maxWidth: .infinity)
                        } else {
                            Text("Sign In")
                                .fontWeight(.medium)
                                .frame(maxWidth: .infinity)
                        }
                    }
                    .frame(height: inputHeight)
                    .background(isLoading || username.isEmpty || password.isEmpty
                        ? Color.gray.opacity(0.3)
                        : Color.accentColor)
                    .foregroundStyle(.white)
                    .clipShape(RoundedRectangle(cornerRadius: inputCornerRadius))
                    .disabled(isLoading || username.isEmpty || password.isEmpty)
                }
                .padding(32)
                .background(Color(.systemBackground))
                .clipShape(RoundedRectangle(cornerRadius: 20))
                .shadow(color: .black.opacity(0.08), radius: 20, x: 0, y: 8)
                .shadow(color: .black.opacity(0.04), radius: 4, x: 0, y: 2)
                .padding(.horizontal, 24)

                if let errorMessage {
                    Text(errorMessage)
                        .foregroundStyle(.red.opacity(0.9))
                        .font(.subheadline)
                        .padding(.top, 16)
                }

                Spacer()
            }
            .frame(maxWidth: 400)
        }
    }

    private func signIn() {
        errorMessage = nil
        isLoading = true

        Task { @MainActor in
            do {
                _ = try await AuthService.signIn(username: username, password: password)
                // TODO: Navigate to main content on success
            } catch AuthError.invalidURL {
                errorMessage = "Invalid server URL"
            } catch AuthError.loginFailed {
                errorMessage = "Login failed. Please check your credentials."
            } catch {
                errorMessage = error.localizedDescription
            }
            isLoading = false
        }
    }
}

#Preview {
    LoginView()
}
