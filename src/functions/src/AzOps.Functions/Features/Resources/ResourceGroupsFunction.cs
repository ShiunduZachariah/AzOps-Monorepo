using AzOps.Core.Features.Resources;
using AzOps.Functions.Features.Shared;
using Microsoft.Azure.Functions.Worker;
using Microsoft.Azure.Functions.Worker.Http;
using Microsoft.Extensions.Logging;

namespace AzOps.Functions.Features.Resources;

public sealed class ResourceGroupsFunction(IResourceGroupsService resourceGroupsService, ILogger<ResourceGroupsFunction> logger)
{
    [Function("ResourceGroups")]
    public Task<HttpResponseData> Run(
        [HttpTrigger(AuthorizationLevel.Anonymous, "get", Route = "resources/resource-groups")] HttpRequestData request,
        FunctionContext context)
    {
        return FunctionExecution.ExecuteAsync(
            request,
            logger,
            operationName: "resource-groups",
            handler: cancellationToken => resourceGroupsService.GetSnapshotAsync(cancellationToken),
            cancellationToken: context.CancellationToken);
    }
}
