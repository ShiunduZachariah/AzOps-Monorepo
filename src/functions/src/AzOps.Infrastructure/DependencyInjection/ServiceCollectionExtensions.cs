using AzOps.Core.Features.Health;
using AzOps.Core.Features.Ping;
using AzOps.Core.Features.Resources;
using AzOps.Core.Features.Secrets;
using AzOps.Infrastructure.Configuration;
using AzOps.Infrastructure.Features.Resources;
using AzOps.Infrastructure.Features.Secrets;
using Azure.Core;
using Azure.Identity;
using Microsoft.Extensions.DependencyInjection;

namespace AzOps.Infrastructure.DependencyInjection;

public static class ServiceCollectionExtensions
{
    public static IServiceCollection AddAzOpsInfrastructure(this IServiceCollection services, AzOpsFunctionsOptions options)
    {
        services.AddSingleton(options);
        services.AddSingleton(TimeProvider.System);
        services.AddSingleton<TokenCredential>(_ => new DefaultAzureCredential());
        services.AddSingleton<IHealthService, HealthService>();
        services.AddSingleton<IPingService, PingService>();
        services.AddSingleton<IResourceGroupReader, AzureResourceGroupReader>();
        services.AddSingleton<IResourceGroupsService, ResourceGroupsService>();
        services.AddSingleton<ISecretValueReader, KeyVaultSecretValueReader>();
        services.AddSingleton<ISecretInspectionService, SecretInspectionService>();

        return services;
    }
}
