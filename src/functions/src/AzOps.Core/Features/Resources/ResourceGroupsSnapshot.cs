namespace AzOps.Core.Features.Resources;

public sealed record ResourceGroupsSnapshot(int Count, IReadOnlyList<ResourceGroupSummary> Items);
