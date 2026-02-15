//
//  AuthService.swift
//  Clientele Manager
//
//  Created by Evan Haeick on 2/15/26.
//

import Foundation

enum AuthError: Error {
    case invalidURL
    case loginFailed
}

struct AuthService {
    static func signIn(username: String, password: String) async throws -> Bool {
        guard let url = URL(string: AuthConfig.loginURL) else {
            throw AuthError.invalidURL
        }

        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.httpBody = try JSONEncoder().encode([
            "username": username,
            "password": password
        ])

        let (_, response) = try await URLSession.shared.data(for: request)

        guard let httpResponse = response as? HTTPURLResponse,
              (200...299).contains(httpResponse.statusCode) else {
            throw AuthError.loginFailed
        }

        return true
    }
}
