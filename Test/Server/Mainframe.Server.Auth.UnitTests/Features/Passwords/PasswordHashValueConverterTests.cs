using Mainframe.Server.Auth.Features.Passwords;

namespace Mainframe.Server.Auth.UnitTests.Features.Passwords;

public class PasswordHashValueConverterTests
{
    [Fact]
    public void ValueConverter_ShouldConvertFromString()
    {
        var converter = new PasswordHashValueConverter();
        var fakeHash = "some hash";

        var actual = converter.ConvertFromProvider(fakeHash);

        Assert.IsType<PasswordHash>(actual);
        Assert.Equal(fakeHash, ((PasswordHash)actual).Value);
    }

    [Fact]
    public void ValueConverter_ShouldConvertToString()
    {
        var converter = new PasswordHashValueConverter();
        var fakeHash = "some hash";
        var pwHash = new PasswordHash(fakeHash);

        var actual = converter.ConvertToProvider(pwHash);

        Assert.IsType<string>(actual);
        Assert.Equal(fakeHash, actual);
    }
}
