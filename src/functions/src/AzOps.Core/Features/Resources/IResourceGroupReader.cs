namespace AzOps.Core.Features.Resources;

public interface IResourceGroupReader
{
    Task<IReadOnlyList<ResourceGroupSummary>> ListAsync(CancellationToken cancellationToken);
}
