namespace AzOps.Core.Features.Health;

public interface IHealthService
{
    HealthStatus GetStatus();
}
