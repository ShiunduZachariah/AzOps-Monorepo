namespace AzOps.Core.Features.Ping;

public sealed class PingService(TimeProvider timeProvider) : IPingService
{
    public PingStatus GetStatus()
    {
        return new PingStatus(
            Status: "Healthy",
            Service: "AzOps.Functions",
            Message: "pong",
            TimestampUtc: timeProvider.GetUtcNow());
    }
}
