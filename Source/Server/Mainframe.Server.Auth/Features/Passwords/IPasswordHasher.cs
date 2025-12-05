namespace Mainframe.Server.Auth.Features.Passwords;

/// <summary>
/// A password hasher used for handling plaintext to password hash conversion and verification.
/// </summary>
public interface IPasswordHasher
{
    /// <summary>
    /// Convert a plaintext password into a password hash suitable for storing in a database.
    /// </summary>
    /// <param name="plaintext">The raw password text to hash.</param>
    /// <returns>A password hash or an error.</returns>
    Result<PasswordHash, Exception> Hash(string plaintext);

    /// <summary>
    /// Verify that the password attempt is the correct plaintext password when compared against the
    /// password hash.
    /// </summary>
    /// <param name="hash">The existing password hash.</param>
    /// <param name="attempt">The plaintext password to hash and compare against the existing hash.</param>
    /// <returns>True when they match, otherwise false.</returns>
    bool Verify(PasswordHash hash, string attempt);
}
