using AzOps.Core.Features.Health;
using AzOps.Core.Features.Ping;
using AzOps.Core.Features.Resources;
using AzOps.Core.Features.Secrets;
using AzOps.Infrastructure.Configuration;
using AzOps.Infrastructure.DependencyInjection;
using Microsoft.Extensions.DependencyInjection;

namespace AzOps.IntegrationTests.DependencyInjection;

public sealed class InfrastructureCompositionTests
{
    [Fact]
    public void AddAzOpsInfrastructure_RegistersSprint3Services()
    {
        var services = new ServiceCollection();

        services.AddAzOpsInfrastructure(new AzOpsFunctionsOptions
        {
            SubscriptionId = "00000000-0000-0000-0000-000000000000",
            KeyVaultName = "azops-demo"
        });

        using var provider = services.BuildServiceProvider();

        Assert.IsType<HealthService>(provider.GetRequiredService<IHealthService>());
        Assert.IsType<PingService>(provider.GetRequiredService<IPingService>());
        Assert.NotNull(provider.GetRequiredService<IResourceGroupsService>());
        Assert.NotNull(provider.GetRequiredService<ISecretInspectionService>());
    }
}
