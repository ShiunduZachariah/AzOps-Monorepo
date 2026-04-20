namespace AzOps.Core.Features.Resources;

public interface IResourceGroupsService
{
    Task<ResourceGroupsSnapshot> GetSnapshotAsync(CancellationToken cancellationToken);
}
