using AzOps.Core.Features.Resources;

namespace AzOps.Functions.Tests.Features.Resources;

public sealed class ResourceGroupsServiceTests
{
    [Fact]
    public async Task GetSnapshotAsync_OrdersResultsByName()
    {
        var service = new ResourceGroupsService(new FakeResourceGroupReader(
            new ResourceGroupSummary("/subscriptions/1/resourceGroups/zeta", "zeta", "eastus"),
            new ResourceGroupSummary("/subscriptions/1/resourceGroups/alpha", "alpha", "westus")));

        var snapshot = await service.GetSnapshotAsync(CancellationToken.None);

        Assert.Equal(2, snapshot.Count);
        Assert.Collection(
            snapshot.Items,
            item => Assert.Equal("alpha", item.Name),
            item => Assert.Equal("zeta", item.Name));
    }

    private sealed class FakeResourceGroupReader(params ResourceGroupSummary[] items) : IResourceGroupReader
    {
        public Task<IReadOnlyList<ResourceGroupSummary>> ListAsync(CancellationToken cancellationToken)
        {
            return Task.FromResult<IReadOnlyList<ResourceGroupSummary>>(items);
        }
    }
}
