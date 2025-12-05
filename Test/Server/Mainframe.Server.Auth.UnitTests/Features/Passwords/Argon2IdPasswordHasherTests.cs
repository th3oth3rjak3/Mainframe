using Mainframe.Server.Auth.Features.Passwords;

namespace Mainframe.Server.Auth.UnitTests.Features.Passwords;

public class Argon2IdPasswordHasherTests
{
    [Fact]
    public void PasswordHasher_ShouldError_WhenInputsInvalid()
    {
        var hasher = new Argon2IdPasswordHasher();
        // Should be err variant.
        var err = hasher.Hash("").UnwrapError();
        Assert.IsType<ArgumentException>(err);

        err = hasher.Hash(null!).UnwrapError();
        Assert.IsType<ArgumentException>(err);
    }

    [Fact]
    public void PasswordHasher_ShouldHash_WhenInputsOk()
    {
        var hasher = new Argon2IdPasswordHasher();
        var rawPassword = "plaintext password";
        var hash = hasher.Hash(rawPassword).Unwrap();

        Assert.NotEqual(rawPassword, hash.Value);
        Assert.StartsWith("$argon2id", hash.Value);
    }

    [Fact]
    public void PasswordHasher_ShouldNotProduceIdenticalHashes_WhenInputsSame()
    {
        var hasher = new Argon2IdPasswordHasher();
        var rawPassword = "plaintext password";
        var hash1 = hasher.Hash(rawPassword).Unwrap();
        var hash2 = hasher.Hash(rawPassword).Unwrap();

        Assert.NotEqual(hash1.Value, hash2.Value);
    }

    [Fact]
    public void PasswordHasher_ShouldVerify_WhenPasswordsValid()
    {
        var hasher = new Argon2IdPasswordHasher();
        var rawPassword = "plaintext password";
        var hash = hasher.Hash(rawPassword).Unwrap();
        Assert.True(hasher.Verify(hash, rawPassword));
    }

    [Fact]
    public void PasswordHasher_ShouldNotVerify_WhenPasswordsInvalid()
    {
        var hasher = new Argon2IdPasswordHasher();
        var rawPassword = "plaintext password";
        var hash = hasher.Hash(rawPassword).Unwrap();

        Assert.False(hasher.Verify(hash, ""));
        Assert.False(hasher.Verify(hash, null!));
        Assert.False(hasher.Verify(hash, rawPassword.ToUpper()));
    }
}
