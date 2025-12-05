namespace Mainframe.Server.Auth.Features.Passwords;

/// <summary>
/// A password hash that has already been through the hashing process. 
/// To create a new Password hash from plaintext, use <see cref="IPasswordHasher"/> instead.
/// </summary>
/// <param name="Value">The hashed string value that represents a password.</param>
[ExcludeFromCodeCoverage]
public record PasswordHash(string Value);
