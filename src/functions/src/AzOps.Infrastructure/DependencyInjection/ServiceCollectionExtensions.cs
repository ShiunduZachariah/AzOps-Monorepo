using AzOps.Core.Features.Health;
using AzOps.Infrastructure.Features.Health;
using Microsoft.Extensions.DependencyInjection;

namespace AzOps.Infrastructure.DependencyInjection;

public static class ServiceCollectionExtensions
{
    public static IServiceCollection AddAzOpsInfrastructure(this IServiceCollection services)
    {
        services.AddSingleton<IHealthService, HealthService>();

        return services;
    }
}
