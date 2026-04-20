using AzOps.Core.Features.Health;

namespace AzOps.Functions.Tests.Features.Health;

public sealed class HealthServiceTests
{
    [Fact]
    public void GetStatus_ReturnsHealthySnapshot()
    {
        var service = new HealthService();

        var status = service.GetStatus();

        Assert.Equal("Healthy", status.Status);
        Assert.Equal("AzOps.Functions", status.Service);
        Assert.Equal("The Functions host is ready.", status.Message);
    }
}
