namespace AzOps.Core.Features.Resources;

public sealed class ResourceGroupsService(IResourceGroupReader resourceGroupReader) : IResourceGroupsService
{
    public async Task<ResourceGroupsSnapshot> GetSnapshotAsync(CancellationToken cancellationToken)
    {
        var items = await resourceGroupReader.ListAsync(cancellationToken);
        var orderedItems = items
            .OrderBy(item => item.Name, StringComparer.OrdinalIgnoreCase)
            .ToArray();

        return new ResourceGroupsSnapshot(
            Count: orderedItems.Length,
            Items: orderedItems);
    }
}
