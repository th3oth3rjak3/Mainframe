using Isopoh.Cryptography.Argon2;

namespace Mainframe.Server.Auth.Features.Passwords;

public class Argon2IdPasswordHasher : IPasswordHasher
{
    public Result<PasswordHash, Exception> Hash(string plaintext) =>
        Try(() =>
        {
            if (string.IsNullOrWhiteSpace(plaintext))
            {
                throw new ArgumentException("Password cannot be null or whitespace", nameof(plaintext));
            }

            // Argon2.Hash is not nullable aware, do null check to be sure the return hash is not null.
            var hash = Argon2.Hash(plaintext) ?? throw new InvalidOperationException("Password hash cannot be null");

            return new PasswordHash(hash);
        });


    public bool Verify(PasswordHash hash, string attempt)
    {
        if (string.IsNullOrWhiteSpace(attempt)) return false;
        return Argon2.Verify(hash.Value, attempt);
    }
}
