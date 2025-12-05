using Mainframe.Server.Auth.Features.Sessions;

namespace Mainframe.Server.Auth.UnitTests.Features.Sessions;

public class SessionTests
{
    [Fact]
    public void CreatingSession_ShouldHaveTwoHourDefaultExpiration()
    {
        var session = new Session(1);
        var now = DateTime.UtcNow;
        var low = now.AddHours(2).AddMinutes(-1);
        var high = now.AddHours(2).AddMinutes(1);
        Assert.InRange(session.ExpiresAt, low, high);
        Assert.False(session.IsExpired);
        Assert.Equal(1, session.UserId);
    }

    [Fact]
    public void CreatingSession_ShouldNotHaveDefaultGuid()
    {
        var session = new Session(1);
        Assert.NotEqual(Guid.Empty, session.Id);
    }

    [Fact]
    public void Session_ShouldExtendExpiration_WhenProvided()
    {
        var session = new Session(1);
        var now = DateTime.UtcNow;
        var extraTime = TimeSpan.FromHours(1);

        session.Slide(extraTime);

        var low = now.Add(extraTime).AddMinutes(-1);
        var high = now.Add(extraTime).AddMinutes(1);

        Assert.InRange(session.ExpiresAt, low, high);
    }

    [Fact]
    public void Session_ShouldExtendExpiration_WhenNotProvided()
    {
        var session = new Session(1);
        var now = DateTime.UtcNow;

        session.Slide(); // Default is 2 hours.

        var low = now.AddHours(2).AddMinutes(-1);
        var high = now.AddHours(2).AddMinutes(1);

        Assert.InRange(session.ExpiresAt, low, high);
    }

    [Fact]
    public void Session_IsExpired_WhenExpiresAtIsPast()
    {
        var session = new Session(1);
        session.Slide(TimeSpan.FromHours(-3));

        Assert.True(session.IsExpired);
    }
}
