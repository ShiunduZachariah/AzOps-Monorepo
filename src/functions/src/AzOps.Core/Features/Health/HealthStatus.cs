namespace AzOps.Core.Features.Health;

public sealed record HealthStatus(string Status, string Service, string Message);
