using AzOps.Core.Features.Secrets;
using AzOps.Functions.Features.Shared;
using Microsoft.Azure.Functions.Worker;
using Microsoft.Azure.Functions.Worker.Http;
using Microsoft.Extensions.Logging;

namespace AzOps.Functions.Features.Secrets;

public sealed class SecretInspectionFunction(ISecretInspectionService secretInspectionService, ILogger<SecretInspectionFunction> logger)
{
    [Function("SecretInspection")]
    public Task<HttpResponseData> Run(
        [HttpTrigger(AuthorizationLevel.Anonymous, "get", Route = "secrets/{secretName}")] HttpRequestData request,
        string secretName,
        FunctionContext context)
    {
        return FunctionExecution.ExecuteAsync(
            request,
            logger,
            operationName: "secret-inspection",
            handler: cancellationToken => secretInspectionService.InspectAsync(secretName, cancellationToken),
            cancellationToken: context.CancellationToken);
    }
}
