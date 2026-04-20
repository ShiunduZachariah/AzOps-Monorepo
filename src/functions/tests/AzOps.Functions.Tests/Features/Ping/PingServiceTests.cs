using AzOps.Core.Features.Ping;

namespace AzOps.Functions.Tests.Features.Ping;

public sealed class PingServiceTests
{
    [Fact]
    public void GetStatus_ReturnsPongWithUtcTimestamp()
    {
        var fixedTime = new DateTimeOffset(2026, 4, 20, 8, 30, 0, TimeSpan.Zero);
        var service = new PingService(new FakeTimeProvider(fixedTime));

        var status = service.GetStatus();

        Assert.Equal("Healthy", status.Status);
        Assert.Equal("AzOps.Functions", status.Service);
        Assert.Equal("pong", status.Message);
        Assert.Equal(fixedTime, status.TimestampUtc);
    }

    private sealed class FakeTimeProvider(DateTimeOffset utcNow) : TimeProvider
    {
        public override DateTimeOffset GetUtcNow() => utcNow;
    }
}
