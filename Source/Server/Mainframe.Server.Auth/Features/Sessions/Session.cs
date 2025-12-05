namespace Mainframe.Server.Auth.Features.Sessions;

/// <summary>
/// A session represents an authorized period of time that a user is to remain logged in.
/// </summary>
public class Session
{
    public Guid Id { get; } = Guid.NewGuid();
    public int UserId { get; private set; }
    public DateTime ExpiresAt { get; private set; } = DateTime.UtcNow.AddHours(2);
    public bool IsExpired => ExpiresAt < DateTime.UtcNow;

    public Session(int userId)
    {
        UserId = userId;
    }

    [ExcludeFromCodeCoverage]
    private Session() { }

    /// <summary>
    /// Update the expiration to be UtcNow plus the session duration provided by the caller.
    /// </summary>
    /// <param name="sessionDuration">The length of the session duration.</param>
    public void Slide(TimeSpan sessionDuration = default)
    {
        var duration = sessionDuration == default ? TimeSpan.FromHours(2) : sessionDuration;
        ExpiresAt = DateTime.UtcNow.Add(duration);
    }
}
