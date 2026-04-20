using AzOps.Core.Common;
using AzOps.Core.Features.Resources;
using AzOps.Infrastructure.Configuration;
using Azure.Core;
using Azure.ResourceManager;
using Azure.ResourceManager.Resources;
using System.Net;

namespace AzOps.Infrastructure.Features.Resources;

public sealed class AzureResourceGroupReader : IResourceGroupReader
{
    private readonly ArmClient armClient;
    private readonly AzOpsFunctionsOptions options;

    public AzureResourceGroupReader(TokenCredential credential, AzOpsFunctionsOptions options)
    {
        armClient = new ArmClient(credential);
        this.options = options;
    }

    public async Task<IReadOnlyList<ResourceGroupSummary>> ListAsync(CancellationToken cancellationToken)
    {
        var subscriptionId = options.SubscriptionId?.Trim();
        if (string.IsNullOrWhiteSpace(subscriptionId))
        {
            throw new AzOpsException(
                statusCode: HttpStatusCode.ServiceUnavailable,
                code: "subscription_not_configured",
                message: "AZOPS_SUBSCRIPTION_ID is required for the resource-groups endpoint.");
        }

        var subscriptionResource = armClient.GetSubscriptionResource(new ResourceIdentifier($"/subscriptions/{subscriptionId}"));
        var items = new List<ResourceGroupSummary>();

        await foreach (ResourceGroupResource resourceGroup in subscriptionResource.GetResourceGroups().GetAllAsync(cancellationToken: cancellationToken))
        {
            items.Add(new ResourceGroupSummary(
                Id: resourceGroup.Id.ToString(),
                Name: resourceGroup.Data.Name,
                Location: resourceGroup.Data.Location.ToString()));
        }

        return items;
    }
}
