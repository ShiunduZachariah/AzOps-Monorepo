namespace AzOps.Core.Features.Health;

public sealed class HealthService : IHealthService
{
    public HealthStatus GetStatus()
    {
        return new HealthStatus(
            Status: "Healthy",
            Service: "AzOps.Functions",
            Message: "The Functions host is ready.");
    }
}
